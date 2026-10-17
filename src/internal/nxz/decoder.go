package nxz

import (
	"encoding/binary"
	"errors"
	"fmt"
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
var decDistTable = [8]group{
	{0, 0x0}, {0, 0x1}, {1, 0x2}, {2, 0x4}, {3, 0x8}, {4, 0x10}, {5, 0x20}, {6, 0x40},
}

var errTruncated = errors.New("nxz: unexpected end of data")

type bitReader struct {
	data []byte
	buf  uint32
	n    int32
}

func (r *bitReader) read(n int32) (uint32, error) {
	for r.n < n {
		if len(r.data) == 0 {
			return 0, errTruncated
		}
		r.buf |= uint32(r.data[0]) << (24 - r.n)
		r.data = r.data[1:]
		r.n += 8
	}
	v := r.buf >> (32 - n)
	r.buf <<= n
	r.n -= n
	return v, nil
}

// Decompress the whole file. Input must include the 4 byte header with the decompressed size.
func Decompress(data []byte) ([]byte, error) {
	if len(data) < 4 {
		return nil, errTruncated
	}
	sz := int(binary.LittleEndian.Uint32(data))
	out := make([]byte, 0, sz)
	d := &decoder{order: initOrder, groups: initGroups}
	r := &bitReader{data: data[4:]}
	for len(out) < sz && len(r.data) > 0 {
		var err error
		// each chunk starts at a byte boundary; remaining bits in buf are zero padding and the original keeps them
		r.n = 0
		out, err = d.decodeChunk(r, out, sz)
		if err != nil {
			return nil, err
		}
	}
	if len(out) != sz {
		return nil, fmt.Errorf("nxz: unexpected size: %d vs %d", len(out), sz)
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

func (d *decoder) decodeChunk(r *bitReader, out []byte, max int) ([]byte, error) {
	for {
		g, err := r.read(4)
		if err != nil {
			return out, err
		}
		gr := d.groups[g]
		ind := gr.off
		if gr.bits != 0 {
			v, err := r.read(gr.bits)
			if err != nil {
				return out, err
			}
			ind += int32(v)
		}
		if ind >= numSyms {
			return out, fmt.Errorf("nxz: invalid symbol index: %d", ind)
		}
		sym := d.order[ind]
		d.m.freq[sym]++
		switch {
		case sym < 256:
			if len(out) >= max {
				return out, errors.New("nxz: output overflow")
			}
			out = append(out, byte(sym))
			d.win[uint16(d.pos)] = byte(sym)
			d.pos++
			continue
		case sym == symRebuild:
			pairs, _ := d.m.sorted()
			for i, p := range pairs {
				d.order[i] = p.sym
			}
			bits, off := int32(0), int32(0)
			for i := range d.groups {
				for {
					b, err := r.read(1)
					if err != nil {
						return out, err
					}
					if b != 0 {
						break
					}
					bits++
				}
				if bits > 16 {
					return out, fmt.Errorf("nxz: invalid group size: %d", bits)
				}
				d.groups[i] = group{bits: bits, off: off}
				off += 1 << bits
			}
			continue
		case sym == symEnd:
			return out, nil
		}
		n := int(sym) - 256
		if sym >= 264 {
			t := decLenTable[sym-264]
			v, err := r.read(t.bits)
			if err != nil {
				return out, err
			}
			n = int(v) + int(t.off)
		}
		n += 4
		c, err := r.read(3)
		if err != nil {
			return out, err
		}
		t := decDistTable[c]
		v, err := r.read(t.bits + 9)
		if err != nil {
			return out, err
		}
		dist := uint32(t.off)<<9 + v
		if len(out)+n > max {
			return out, errors.New("nxz: output overflow")
		}
		for i := 0; i < n; i++ {
			b := d.win[uint16(d.pos-dist)]
			out = append(out, b)
			d.win[uint16(d.pos)] = b
			d.pos++
		}
	}
}
