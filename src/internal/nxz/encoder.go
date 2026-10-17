// Package nxz implements the compressor for Nox .nxz map files.
//
// It is a direct port of the original encoder, so the output is byte-exact with the one produced by the game.
package nxz

import (
	"encoding/binary"
	"math/bits"
)

const (
	numSyms    = 274
	symRebuild = 272
	symEnd     = 273
	windowSize = 0x10000
	hashSize   = 0x8000
	maxMatch   = 521
	chunkSize  = 500000
	// readPad is the number of zero bytes appended to the input.
	// The original encoder reads a few bytes past the end of the input when updating the hash.
	readPad = 32
)

// model is an adaptive prefix code for the 274 symbols.
type model struct {
	freq    [numSyms]uint16
	counter int32
	// codes: low 16 bits are the code length, high 16 bits are the code.
	codes [numSyms]uint32
}

// Encoder keeps the compressor state. The state is shared between chunks.
type Encoder struct {
	win  [windowSize]byte
	pos  uint32
	m    model
	hash [hashSize]uint16
}

// NewEncoder creates a new encoder.
func NewEncoder() *Encoder {
	e := &Encoder{}
	e.m.counter = 4096
	e.m.codes = initCodes
	for i := range e.hash {
		e.hash[i] = 0xFFFF
	}
	return e
}

// Compress the whole file. Output includes the 4 byte header with the decompressed size.
func Compress(data []byte) []byte {
	src := make([]byte, len(data)+readPad)
	copy(src, data)
	out := make([]byte, 4, 4+len(data)/2)
	binary.LittleEndian.PutUint32(out, uint32(len(data)))
	e := NewEncoder()
	for i := 0; i < len(data); i += chunkSize {
		n := min(len(data)-i, chunkSize)
		out = e.compress(out, src, i, n)
	}
	return out
}

func rol(v uint32, n int) uint32 {
	return bits.RotateLeft32(v, n)
}

func hash5(b []byte) uint32 {
	var h uint32
	for _, c := range b[:5] {
		h = rol(h^uint32(c), 5)
	}
	return h
}

// hashNext advances the hash of b[0:5] to b[1:6].
func hashNext(h uint32, b []byte) uint32 {
	return rol(h^uint32(b[5])^rol(uint32(b[0]), 25), 5)
}

func hashSlot(h uint32) uint32 {
	return (214013*h + 2531011) >> 17
}

func (e *Encoder) winCopy(b []byte) {
	o := e.pos & 0xFFFF
	for i, c := range b {
		e.win[(o+uint32(i))&0xFFFF] = c
	}
	e.pos += uint32(len(b))
}

// compress chunk src[p:p+n] and append it to out. Bytes after the chunk must be readable (see readPad).
func (e *Encoder) compress(out []byte, src []byte, p, n int) []byte {
	w := bitWriter{m: &e.m, out: out}
	end := p + n
	rem := n
	var h uint32
	if n >= 5 {
		h = hash5(src[p:])
	}
	for rem >= 5 {
		var (
			cur      int // offset of the current byte from p
			pendLen  int // lazy match candidate
			pendCand uint32
			pendOff  int
			pendPos  uint32
			pendHash uint32
			back     bool
		)
		lim := min(rem-5, 64)
		matched := false
	scan:
		for {
			mlen := 0
			s := hashSlot(h)
			cand := uint32(e.hash[s])
			wpos := e.pos & 0xFFFF
			e.hash[s] = uint16(wpos)
			if cand != 0xFFFF && cand != wpos {
				d := int((wpos - cand) & 0xFFFF)
				left := end - (p + cur)
				mx := d
				if min(d, left) >= maxMatch {
					mx = maxMatch
				} else if d >= left {
					mx = left
				}
				for mlen < mx && e.win[(cand+uint32(mlen))&0xFFFF] == src[p+cur+mlen] {
					mlen++
				}
				if (cand+uint32(mlen))&0xFFFF == wpos {
					// match reached the current position: continue matching on the input itself
					ext := min(left-mlen, maxMatch-mlen)
					k := 0
					if ext > 0 {
						for i := p + cur; src[i] == src[i+mlen]; i++ {
							k++
							if k >= ext {
								break
							}
						}
					}
					mlen += k
				}
				if mlen >= 3 {
					// try to extend the match backward over pending literals
					bl := min(maxMatch-mlen, cur)
					dist := int((e.pos - cand) & 0xFFFF)
					bl = min(bl, dist-mlen)
					if bl >= 0x10000-dist {
						bl = 0x10000 - dist
					}
					if bl > 0 {
						k := 0
						wi := uint16(cand - 1)
						si := p + cur - 1
						for e.win[wi] == src[si] {
							k++
							wi--
							si--
							if k >= bl {
								break
							}
						}
						if k > 0 {
							cur -= k
							e.pos -= uint32(k)
							cand = (cand - uint32(k)) & 0xFFFF
							mlen += k
							wpos = e.pos & 0xFFFF
							h = hash5(src[p+cur:])
							back = true
						}
					}
				}
			}
			if pendLen >= 4 {
				if mlen <= pendLen {
					// previous match is better
					e.pos = pendPos
					base := p + pendOff
					h1 := hashNext(pendHash, src[base:])
					w.match(src[p:base], uint32(pendLen-4), (pendPos-pendCand)&0xFFFF)
					cnt := min(pendLen-2, end-base-2)
					h = e.rehash(src, base+1, h1, cnt, e.pos+2)
					e.winCopy(src[base : base+pendLen])
					p = base + pendLen
				} else {
					p = e.emitMatch(&w, src, p, cur, end, mlen, wpos, cand, &h)
				}
				matched = true
				break scan
			}
			if mlen >= 4 {
				if back {
					p = e.emitMatch(&w, src, p, cur, end, mlen, wpos, cand, &h)
					matched = true
					break scan
				}
				pendLen = mlen
				pendCand = cand
				pendOff = cur
				pendPos = e.pos
				pendHash = h
			}
			if cur+1 > lim {
				break
			}
			c := src[p+cur]
			h = hashNext(h, src[p+cur:])
			cur++
			e.win[uint16(e.pos)] = c
			e.pos++
		}
		if !matched {
			if pendLen < 4 {
				if cur+5 >= rem && rem <= 64 {
					e.winCopy(src[p+cur : p+rem])
					cur = rem
				}
				w.literals(src[p : p+cur])
				p += cur
			} else {
				base := p + cur
				w.match(src[p:base], uint32(pendLen-4), (e.pos-pendCand)&0xFFFF)
				cnt := min(pendLen-1, end-base-1)
				h = e.rehash(src, base, h, cnt, e.pos+1)
				e.winCopy(src[base : base+pendLen])
				p = base + pendLen
			}
		}
		rem = end - p
	}
	if rem > 0 {
		e.winCopy(src[p : p+rem])
		w.literals(src[p : p+rem])
	}
	return w.finish()
}

// emitMatch writes literals src[p:p+cur] followed by a match at cur and inserts the matched positions to the hash table.
func (e *Encoder) emitMatch(w *bitWriter, src []byte, p, cur, end, mlen int, wpos, cand uint32, h *uint32) int {
	base := p + cur
	w.match(src[p:base], uint32(mlen-4), (wpos-cand)&0xFFFF)
	cnt := min(mlen-1, end-base-1)
	*h = e.rehash(src, base, *h, cnt, e.pos+1)
	e.winCopy(src[base : base+mlen])
	return base + mlen
}

// rehash inserts cnt positions starting from i+1 to the hash table, given h is the hash at i.
// Window position for i+1 is wpos. Returns the hash at i+cnt+1.
func (e *Encoder) rehash(src []byte, i int, h uint32, cnt int, wpos uint32) uint32 {
	if cnt <= 0 {
		return 0
	}
	if cnt > 1024 {
		// unreachable with the current match length limit, kept as in the original
		return hash5(src[i+cnt+1:])
	}
	h = hashNext(h, src[i:])
	for k := 0; k < cnt; k++ {
		e.hash[hashSlot(h)] = uint16(uint32(k) + wpos)
		h = hashNext(h, src[i+1+k:])
	}
	return h
}
