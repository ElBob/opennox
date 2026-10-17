package nxz

// bitWriter writes MSB-first codes from the adaptive model.
type bitWriter struct {
	m    *model
	out  []byte
	buf  uint32
	bits int32
}

func (w *bitWriter) put(n int32, v uint32) {
	// shift count is masked the same way x86 does it
	w.buf |= v << (uint32(32-w.bits-n) & 31)
	w.bits += n
	if w.bits >= 16 {
		w.out = append(w.out, byte(w.buf>>24), byte(w.buf>>16))
		w.bits -= 16
		w.buf <<= 16
	}
}

// count a symbol for the model and rebuild it if necessary.
func (w *bitWriter) count(s uint32) uint32 {
	w.m.counter--
	if w.m.counter <= 0 {
		w.rebuild()
	}
	w.m.freq[s]++
	return w.m.codes[s]
}

func (w *bitWriter) putSym(s uint32) {
	c := w.count(s)
	w.put(int32(int16(c)), c>>16)
}

func (w *bitWriter) rebuild() {
	w.m.counter = 1
	// the rebuild marker is written with the old codes
	w.m.freq[symRebuild]++
	c := w.m.codes[symRebuild]
	w.put(int32(int16(c)), c>>16)
	var groups [16]int32
	w.m.build(&groups)
	prev := int32(0)
	for _, g := range groups {
		w.put(g-prev+1, 1)
		prev = g
	}
}

func (w *bitWriter) literals(b []byte) {
	for _, c := range b {
		w.putSym(uint32(c))
	}
}

// match writes literals followed by a match with length l+4 and distance d.
func (w *bitWriter) match(lit []byte, l uint32, d uint32) {
	w.literals(lit)
	switch {
	case l >= 38:
		t := &lenTable[(l-38)>>5]
		w.putSym(t[0] + 4)
		w.put(int32(t[1]+4), (uint32(uint8(l))-38)&0x1F|uint32(uint16(t[2]))<<4)
	case l >= 8:
		t := &lenTable[(l-8)>>1]
		nb := int32(t[1])
		c := w.count(t[0])
		w.put(nb+int32(int16(c)), (c>>16)<<nb|(l-8)&1|t[2])
	default:
		w.putSym(l + 256)
	}
	t := &distTable[d>>9]
	n := int32(t[0]) + 9
	v := d&0x1FF | t[1]<<9
	if n <= 16 {
		w.put(n, v)
	} else {
		w.put(int32(t[0])-7, v>>16)
		w.put(16, v&0xFFFF)
	}
}

// finish writes the end marker and flushes remaining bits.
func (w *bitWriter) finish() []byte {
	w.putSym(symEnd)
	for w.bits > 0 {
		w.out = append(w.out, byte(w.buf>>24))
		w.bits -= 8
		w.buf <<= 8
	}
	return w.out
}

type symFreq struct{ sym, freq uint16 }

// sorted returns symbols sorted by frequency (descending, then by symbol) and the total frequency.
// Frequencies in the model are halved.
func (m *model) sorted() (pairs [numSyms]symFreq, total int32) {
	for i := range pairs {
		f := int16(m.freq[i])
		pairs[i] = symFreq{sym: uint16(i), freq: uint16(f)}
		total += int32(f)
		m.freq[i] = uint16(f >> 1)
	}
	// shell sort
	for gap := 40; gap > 0; gap /= 3 {
		for j := gap; j < numSyms; j++ {
			key := pairs[j]
			k := j
			for k >= gap {
				p := pairs[k-gap]
				c := int32(int16(p.freq)) - int32(int16(key.freq))
				if c == 0 {
					c = int32(int16(p.sym)) - int32(int16(key.sym))
				}
				if c >= 0 {
					break
				}
				pairs[k] = p
				k -= gap
			}
			pairs[k] = key
		}
	}
	return pairs, total
}

// build assigns new codes based on symbol frequencies and returns code lengths for 16 symbol groups.
func (m *model) build(groups *[16]int32) {
	pairs, total := m.sorted()
	freq := func(i int32) int32 { return int32(int16(pairs[i].freq)) }
	m.counter = 4096

	// first 14 groups try to split the remaining frequency evenly
	var (
		nb     int32 // bits for the current group
		pos    int32 // first symbol of the current group
		used   int32 // frequency covered so far
		ngroup int32
		left   = int32(16)
	)
	for {
		var cnt, prevSum, sum int32
		target := (total - used) / left
		for {
			last := false
			sz := int32(1) << nb
			if pos+sz+cnt > numSyms {
				last = true
				sz = numSyms - pos
			}
			for ; cnt < sz; cnt++ {
				sum += freq(pos + cnt)
			}
			if last || nb >= 8 || sum > target {
				break
			}
			prevSum = sum
			nb++
		}
		if nb != 0 && abs(prevSum-target) <= abs(sum-target) {
			sum = prevSum
			nb--
		}
		i := ngroup
		if left < 16 {
			for k := ngroup - 1; ; k-- {
				if nb >= groups[k] {
					break
				}
				groups[k+1] = groups[k]
				i--
				if i <= 0 {
					break
				}
			}
		}
		groups[i] = nb
		used += sum
		pos += 1 << nb
		left--
		ngroup++
		if left <= 2 {
			break
		}
		nb = 0
	}

	// last 2 groups split the rest minimizing the total code length
	var cnt, sum, rest, b1, b2 int32
	best := int32(0x7FFFFFFF)
	if pos < numSyms {
		for k := pos; k < numSyms; k++ {
			rest += freq(k)
		}
	}
	if pos+1 <= numSyms {
		nb1 := int32(0)
		sz := int32(1)
		for {
			for ; cnt < sz; cnt++ {
				sum += freq(pos + cnt)
			}
			nb2 := int32(0)
			if n := numSyms - cnt - pos; n > 1 {
				for {
					nb2++
					if 1<<nb2 >= n {
						break
					}
				}
			}
			if nb1 <= 8 && nb2 <= 8 {
				cost := sum*nb1 + nb2*(rest-sum)
				if cost >= best {
					break
				}
				best = cost
				b1, b2 = nb1, nb2
			}
			nb1++
			sz = 1 << nb1
			if pos+sz+cnt > numSyms {
				break
			}
		}
	}
	insertSorted(groups, 14, b1)
	insertSorted(groups, 15, b2)

	pos = 0
	for g := int32(0); g < 16; g++ {
		b := groups[g]
		sz := int32(1) << b
		n := min(sz, numSyms-pos)
		for j := int32(0); j < n; j++ {
			s := pairs[pos+j].sym
			m.codes[s] = uint32(uint16(b+4)) | uint32(uint16(j|g<<b))<<16
		}
		pos += sz
	}
}

// insertSorted inserts v into sorted groups[:n].
func insertSorted(groups *[16]int32, n int32, v int32) {
	k := n
	for {
		if v >= groups[k-1] {
			break
		}
		groups[k] = groups[k-1]
		k--
		if k <= 0 {
			break
		}
	}
	groups[k] = v
}

func abs(v int32) int32 {
	if v < 0 {
		return -v
	}
	return v
}
