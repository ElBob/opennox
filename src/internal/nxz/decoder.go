package nxz

import (
	"encoding/binary"
	"errors"
)

// group of symbols in the decoder table.
type group struct {
	bits int32
	off  int32
}

// decLenTable is used to decode long match lengths for symbols 264-271: {extra bits, base}.
var decLenTable = [8]group{
	{1, 0x8}, {2, 0xa}, {3, 0xe}, {4, 0x16}, {5, 0x26}, {6, 0x46}, {7, 0x86}, {8, 0x106},
}

// decDistTable is used to decode the high bits of a match distance: {extra bits, base}.
// It is indexed with the 3 bit code + 1; the first entry is used when the input is exhausted (code -1).
var decDistTable = [9]group{
	{8, 0x106},
	{0, 0x0}, {0, 0x1}, {1, 0x2}, {2, 0x4}, {3, 0x8}, {4, 0x10}, {5, 0x20}, {6, 0x40},
}

var errTruncated = errors.New("nxz: unexpected end of data")

// bitReader reads bits like the original: it loads at most one byte per read and
// when the input is exhausted it sets the bit count to 0 and returns -1.
type bitReader struct {
	data []byte
	buf  uint32
	n    int32
}

func (r *bitReader) read(n int32) int32 {
	if r.n < n {
		if len(r.data) == 0 {
			r.n = 0
			return -1
		}
		r.buf |= uint32(r.data[0]) << uint32(24-r.n)
		r.data = r.data[1:]
		r.n += 8
	}
	v := r.buf >> (32 - n)
	r.buf <<= n
	r.n -= n
	return int32(v)
}

// Decompress the whole file. Input must include the 4 byte header with the decompressed size.
//
// Like the original, problems in the compressed stream are not reported: decoding stops
// and the rest of the output (which always has the size from the header) is left zeroed.
func Decompress(data []byte) ([]byte, error) {
	if len(data) < 4 {
		return nil, errTruncated
	}
	sz := int(binary.LittleEndian.Uint32(data))
	out := make([]byte, sz)
	d := &decoder{order: initOrder, groups: initGroups}
	r := &bitReader{data: data[4:]}
	o := 0
	for len(r.data) > 0 && o < sz {
		// each chunk starts at a byte boundary; remaining bits in buf are zero padding and the original keeps them
		r.n = 0
		var ok bool
		o, ok = d.decodeChunk(r, out, o)
		if !ok {
			break
		}
	}
	return out, nil
}

type decoder struct {
	win    [windowSize]byte
	pos    uint32
	m      model
	order  [numSyms]uint16
	groups [16]group
}

// decodeChunk decodes one chunk into out starting at o. It returns the new output offset and true
// if the end of the chunk was reached, or false if decoding stopped (as nxz_decompress returning 0).
func (d *decoder) decodeChunk(r *bitReader, out []byte, o int) (int, bool) {
	for {
		// input is never exhausted here (see the check at the end of the loop),
		// the original would index the group table with -1 otherwise
		g := r.read(4)
		if g < 0 {
			return o, false
		}
		gr := d.groups[g]
		ind := gr.off
		if gr.bits != 0 {
			ind += r.read(gr.bits)
		}
		if ind >= numSyms {
			return o, false
		}
		if ind < 0 {
			// the original reads out of bounds of the table here
			return o, false
		}
		sym := d.order[ind]
		d.m.freq[sym]++
		switch {
		case sym < 256:
			if o >= len(out) {
				return o, false
			}
			out[o] = byte(sym)
			o++
			d.win[uint16(d.pos)] = byte(sym)
			d.pos++
		case sym == symRebuild:
			pairs, _ := d.m.sorted()
			for i, p := range pairs {
				d.order[i] = p.sym
			}
			bits, off := int32(0), int32(0)
			for i := range d.groups {
				// exhausted input (-1) ends the run like a set bit
				for r.read(1) == 0 {
					bits++
				}
				d.groups[i] = group{bits: bits, off: off}
				off += 1 << bits
			}
			if bits > 16 {
				// the original would shift by invalid amounts when reading such groups
				return o, false
			}
		case sym == symEnd:
			return o, true
		default:
			n := int32(sym) - 256
			if sym >= 264 {
				t := decLenTable[sym-264]
				n = r.read(t.bits) + t.off
			}
			n += 4
			t := decDistTable[r.read(3)+1]
			// high 8 bits are read separately, like in the original
			hi := r.read(8)
			lo := r.read(t.bits + 1)
			dist := t.off<<9 + (hi<<(t.bits+1) | lo)
			if dist < 0 {
				// the original copies a negative amount of memory here
				return o, false
			}
			if o+int(n) > len(out) {
				return o, false
			}
			for i := int32(0); i < n; i++ {
				b := d.win[uint16(d.pos-uint32(dist))]
				if dist == 0 {
					// the original copies the output onto itself here, leaving it unchanged
					b = out[o]
				}
				out[o] = b
				o++
				d.win[uint16(d.pos)] = b
				d.pos++
			}
		}
		if len(r.data) == 0 {
			return o, false
		}
	}
}
