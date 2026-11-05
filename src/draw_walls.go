package opennox

import (
	"encoding/binary"
	"image"
	"unsafe"

	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func (c *Client) wallImage(h unsafe.Pointer) *noxrender.Image {
	return c.r.GetBag().AsImage(noxrender.ImageHandle(h))
}

// sub_47D380 clips a horizontal span against the current draw clip rect.
func (c *Client) sub_47D380(a1, a2 int) int {
	v2 := a1
	v3 := a2
	if a1 > a2 {
		v2 = a2
		v3 = a1
	}
	clip := c.r.Data().ClipRect()
	v4 := clip.Min.X
	if v2 >= v4 {
		if v2 >= clip.Max.X {
			return 0
		}
	} else {
		v2 = clip.Min.X
	}
	v5 := clip.Max.X
	if v3 < v5 {
		if v3 < v4 {
			return 0
		}
	} else {
		v3 = clip.Max.X
	}
	if v2 == v3 {
		return 0
	}
	if v2 != v4 || v3 != v5 {
		*memmap.PtrUint32(0x973F18, 52) = uint32(v2)
		*memmap.PtrUint32(0x973F18, 12) = uint32(v3)
		legacy.Set_dword_5d4594_3799452(1)
	}
	return 1
}

// nox_xxx_edgeDraw_480EF0 draws a type-3 (RLE) image with a per-pixel light gradient between a4 and a5 colors.
// The original takes a4 and a5 as int* (3 ints each); a1 is the image pointer.
func (c *Client) edgeDraw(a1 unsafe.Pointer, a2, a3 int, a4, a5 []uint32, a6, a7, a8, a9, a10 int) {
	_ = a10 // unused in the original
	v10 := 0
	v44 := 0
	if a1 == nil {
		return
	}
	img := c.wallImage(a1)
	var v45, v46 image.Point
	if byte(img.Type())&0x3F != 3 {
		return
	}
	pd := img.Pixdata()
	if pd == nil {
		return
	}
	rd := func(off int) int { return int(int32(binary.LittleEndian.Uint32(pd[off:]))) }
	// LOBYTE(a1) is used by the original as a scratch register holding the last RLE opcode.
	// Before any opcode is read it holds the low byte of the image pointer.
	op := byte(uintptr(a1))

	v12 := rd(0)
	v13 := rd(4) - a9
	v43 := v12
	v40 := v12
	v41 := v13
	if v13 <= 0 {
		return
	}
	// "result" in the original is an int* variable that is used as an int (x position).
	x := rd(8) + a2
	v15 := rd(12)
	pos := 12 + 5 // (char*)v14 + 5
	v17 := v15 + a3
	v53 := v17
	// The clip bounds are stored by sub_430B50 in C globals (left, top, right, bottom).
	clipR := legacy.Get_dword_5d4594_3807116()
	clipL := legacy.Get_dword_5d4594_3807140()
	clipT := legacy.Get_dword_5d4594_3807136()
	clipB := legacy.Get_dword_5d4594_3807152()
	if !(x <= clipR && v17 <= clipB) {
		return
	}
	if x < clipL {
		if x+v12 < clipL {
			return
		}
		v10 = clipL - x
		v40 = v12 - (clipL - x)
		x = clipL
	}
	if x+v40 > clipR {
		v40 = clipR - x
	}
	if v17 < clipT {
		if v13+v17 < clipT {
			return
		}
		v18 := clipT - v17
		v17 = clipT
		v13 -= v18
		v53 = clipT
		v41 = v13
		v44 = v18
	}
	if v13+v17 > clipB {
		v13 = clipB - v17
		v41 = clipB - v17
	}
	if !(v40 > 0 && v13 > 0) {
		return
	}
	v19 := a6
	if a6 < v17 {
		v19 = v17
	}
	if v13+v17 > v19 {
		v41 = v19 - v17
	}
	v20 := 0
	if a7 > x+v10 {
		v20 = a7 - v10 - x
		v10 = a7 - x
		v40 -= v20
	}
	if a8 < x+v40+v10 {
		v40 = a8 - v10 - x
	}
	if v40 <= 0 {
		return
	}
	v52 := unsafe.Add(unsafe.Pointer(noxPixBuffer.rows[v17]), 2*(x+v20))
	v22 := int32(a4[1])
	v23 := int32(a4[0])
	v49 := int32(a4[2]) << 8
	v23 <<= 8
	v24 := int32(a5[0])<<8 - v23
	v47 := v23
	v48 := v22 << 8
	v25 := int32(a5[1]) << 8
	v26 := int32(a5[2]) << 8
	var v50 [3]uint32
	v50[0] = uint32(v24 / int32(v43))
	v50[1] = uint32((v25 - v48) / int32(v43))
	v50[2] = uint32((v26 - v49) / int32(v43))
	skip := func(n int) {
		for i := n; i > 0; {
			v29 := pd[pos]
			v57 := int(pd[pos+1])
			pos += 2
			op = v29
			if v29 == 2 {
				pos += 2 * v57
			}
			i -= v57
		}
	}
	for v27 := v44; v27 > 0; v27-- {
		skip(v43)
	}
	v46.Y = v53
	v45 = sub_473970(v46)
	var v30 uint32
	var v31 unsafe.Pointer
	if v41 == 0 {
		return
	}
	v56 := v41
	bpitch := nox_getBackbufferPitch()
	var v51 [3]uint32
	setGrad := func(n int) {
		v51[0] = uint32(v47 + int32(n)*int32(v50[0]))
		v51[1] = uint32(v48 + int32(n)*int32(v50[1]))
		v51[2] = uint32(v49 + int32(n)*int32(v50[2]))
	}
	u16s := func(p unsafe.Pointer, n int) []uint16 {
		return unsafe.Slice((*uint16)(p), n)
	}
	for {
		v32 := 0
		v54 := 0
		v55 := 0
		if legacy.Get_nox_client_highResFrontWalls_80820() != 0 || v45.Y&1 == 0 {
			var v35 byte
			if v10 <= 0 {
				v35 = op
			} else {
				for {
					v35 = pd[pos]
					v59 := int(pd[pos+1])
					pos += 2
					op = v35
					if v35 == 2 {
						pos += 2 * v59
					}
					v32 += v59
					if !(v32 < v10) {
						break
					}
				}
			}
			var v36 int
			if v32 <= v10 {
				v36 = v40
			} else {
				v36 = v40
				v37 := v40
				if v32-v10 <= v40 {
					v37 = v32 - v10
				}
				if v35 == 2 {
					setGrad(v10)
					sub_480860(u16s(v52, v37), u16s(unsafe.Pointer(&pd[pos-2*(v32-v10)]), v37), v37, v51[:], v50[:])
					v36 = v40
					v54 = 2 * v37
				}
				v55 = 0
			}
			if v32 < v36+v10 {
				for {
					v38 := pd[pos+1]
					op = pd[pos]
					pos += 2
					v60 := v38
					if op == 2 {
						if int(v38) > v10+v36-v32 {
							v55 = v32 + int(v38) - v36 - v10
							v60 = byte(-int(int8(v32 + -int(int8(v36)) - v10)))
						}
						setGrad(v32)
						sub_480860(u16s(unsafe.Add(v52, 2*(v32-v10)), int(v60)), u16s(unsafe.Pointer(&pd[pos]), int(v60)), int(v60), v51[:], v50[:])
						v36 = v40
						pos += 2 * int(v60)
					}
					v32 += int(v60)
					v54 += 2 * int(v60)
					if !(v32 < v36+v10) {
						break
					}
				}
				if v55 != 0 {
					if op == 2 {
						pos += 2 * v55
					}
					v32 += v55
				}
			}
			for v32 < v43 {
				v39 := pd[pos]
				v61 := int(pd[pos+1])
				pos += 2
				op = v39
				if v39 == 2 {
					pos += 2 * v61
				}
				v32 += v61
			}
		} else {
			if v31 != nil && v30 != 0 {
				copy(unsafe.Slice((*byte)(v52), int(v30)), unsafe.Slice((*byte)(v31), int(v30)))
			}
			skip(v43)
		}
		v31 = v52
		v52 = unsafe.Add(v52, bpitch)
		v45.Y++
		v56--
		if v56 == 0 {
			break
		}
		v30 = uint32(v54)
	}
}

// sub_4831C0 copies a 46x46 diamond tile (2116 bytes of 16-bit pixels) from a1 into the tile buffer a2 with stride dword_5d4594_3798804.
// It returns the last 16-bit value of the source (the C caller uses it as a return value).
func (c *Client) sub_4831C0(a1, a2 unsafe.Pointer) int16 {
	stride := c.tiles.dword_5d4594_3798804
	src := 0
	row := func(r, start, w int) {
		dst := unsafe.Add(a2, r*stride+start)
		copy(unsafe.Slice((*byte)(dst), w), unsafe.Slice((*byte)(unsafe.Add(a1, src)), w))
		src += w
	}
	for r := 0; r <= 22; r++ {
		row(r, 46-2*r, 2+4*r)
	}
	for r := 23; r <= 45; r++ {
		row(r, 2*(r-22), 90-4*(r-23))
	}
	return *(*int16)(unsafe.Add(a1, 2112))
}

// nox_xxx_drawWalls_473C10 draws a single wall segment.
func (c *Client) drawWall(vp *noxrender.Viewport, wl *server.Wall) {
	a4 := nox_win_width
	v72 := 0
	a3 := 0
	if wl == nil {
		return
	}
	v4 := byte(wl.Flags4)
	if v4&1 == 0 {
		return
	}
	v6 := vp.Screen.Min.X + 23*int(wl.X5) - vp.World.Min.X
	v82 := v6
	v7 := vp.Screen.Min.Y + 23*int(wl.Y6) - vp.World.Min.Y
	v74 := int(memmap.Int32(0x587000, uintptr(149364+4*int(wl.Field3))))
	v8 := v74
	if v74 == -1 {
		v8 = int(wl.Dir0)
		v74 = v8
	}
	v84 := v8
	if v8 != 0 {
		if v8 == 1 && v4&0x40 != 0 {
			v84 = 12
		}
	} else if v4&0x40 != 0 {
		v84 = 11
	}
	tbl := func(off int) int { return int(memmap.Int32(0x587000, uintptr(off))) }
	if memmap.Uint32(0x587000, 80808) != 0 {
		v9 := 16 * v74
		v10 := tbl(85440 + v9)
		v11 := tbl(85448 + v9)
		a1a := image.Pt(v6+v10, v7+tbl(85444+v9))
		a2a := image.Pt(v6+v11, v7+tbl(85452+v9))
		var v19 int
		if v74 == 7 || v74 == 9 {
			var v22 int
			if c.Sub4C42A0(a1a, a2a, &a3, &a4) != 0 {
				v22 = 1
			} else {
				v22 = 0
				a3 = a2a.X
			}
			v23 := tbl(85508 + v9)
			a1a.X = v82 + tbl(85504+v9)
			v24 := v82 + tbl(85512+v9)
			v25 := tbl(85516 + v9)
			a1a.Y = v7 + v23
			a2a.X = v24
			a2a.Y = v7 + v25
			if c.Sub4C42A0(a1a, a2a, &a3, &a4) != 0 {
				v19 = a3
			} else {
				if v22 == 0 {
					return
				}
				if a4 > a1a.X {
					a4 = a1a.X
				}
				v19 = a3
			}
		} else {
			if v74 != 8 && v74 != 10 {
				if c.Sub4C42A0(a1a, a2a, &a3, &a4) == 0 {
					return
				}
				v19 = a3
			} else {
				v76 := a2a.X
				v75 := a1a.X
				var v73 bool
				if c.Sub4C42A0(a1a, a2a, &v75, &v76) != 0 {
					v73 = v76-v75 >= 3
				} else {
					v73 = false
				}
				v14 := tbl(85504 + 16*v74)
				v15 := tbl(85516 + 16*v74)
				var v80, v81 image.Point
				v80.Y = v7 + tbl(85508+16*v74)
				v16 := v6 + v14
				v17 := v6 + tbl(85512+16*v74)
				v80.X = v16
				a3 = v16
				v81.X = v17
				a4 = v17
				v81.Y = v7 + v15
				v18 := c.Sub4C42A0(v80, v81, &a3, &a4)
				v19 = a3
				var v21 bool
				if v18 == 0 {
					v21 = false
				} else {
					v21 = a4-a3 >= 3
				}
				if v73 {
					if v21 {
						if a3 > v75 {
							v19 = v75
							a3 = v75
						}
						if v19 <= v80.X {
							v19 = 0
							a3 = 0
						}
						if a4 < v76 {
							a4 = v76
						}
						if a4 >= v81.X {
							a4 = nox_win_width
						}
					} else {
						v19 = v75
						a3 = v75
						a4 = v76
						if v74 != 8 {
							v84 = 1
							if v19 == v80.X {
								v19 = 0
								a3 = 0
							}
						} else {
							v84 = 0
							if v76 == v81.X {
								a4 = nox_win_width
							}
						}
					}
				} else {
					if !v21 {
						return
					}
					v84 = 13
					if v74 != 8 {
						v84 = 14
					}
					if a4 == v81.X {
						a4 = nox_win_width
					}
					if v19 == v80.X {
						v19 = 0
						a3 = 0
					}
				}
			}
		}
		if v19 >= a4 {
			wl.Field3 = 0
			wl.Flags4 = wl.Flags4 & 0xFC
			return
		}
	}
	v27 := v4
	v28 := v4 & 2
	var v73 int
	if v28 == 0 {
		v73 = int(v4>>2) & 2
	} else {
		skip61 := false
		if memmap.Uint32(0x5D4594, 805848) != 0 && legacy.Get_nox_client_translucentFrontWalls_805844() != 0 {
			if legacy.Get_nox_client_highResFrontWalls_80820() == 0 && legacy.Get_nox_client_highResFloors_154952() != 0 {
				v72 |= 4
				skip61 = true
			} else {
				v72 = 8
			}
		}
		if !skip61 && legacy.Get_nox_client_highResFrontWalls_80820() == 0 {
			v72 |= 4
		}
		v73 = int(v27&8|4) >> 2
	}
	if v28 != 0 && legacy.Get_nox_client_translucentFrontWalls_805844() != 0 && c.srv.Walls.DefByInd(int(wl.Tile1)).Flags32&4 == 0 {
		v72 |= 2
	} else {
		v72 |= 1
	}
	tile := int(wl.Tile1)
	def := c.srv.Walls.DefByInd(tile)
	var final1 unsafe.Pointer
	var fx, fy int
	if memmap.Uint32(0x587000, 80816) != 0 {
		var pt image.Point
		var l1 [3]uint32
		var l2 []uint32
		get1 := func() {
			// The original copies the first light into a local buffer (v83) before the second call overwrites the shared one.
			l := c.Sub469920(pt)
			l1 = [3]uint32{l[0], l[1], l[2]}
		}
		X := 23 * int(wl.X5)
		Y := int(wl.Y6)
		switch v74 {
		case 0, 3:
			pt = image.Pt(X, 23*(Y+1))
			get1()
			pt.X += 23
			pt.Y -= 23
		case 1, 4:
			pt = image.Pt(X, 23*Y)
			get1()
			pt.X += 23
			pt.Y += 23
		case 7:
			pt = image.Pt(X, 23*Y)
			get1()
			pt.X += 23
		case 8:
			pt = image.Pt(X+11, 23*Y+11)
			get1()
			pt.X -= 34
			pt.Y -= 34
		case 10:
			pt = image.Pt(X, 23*(Y+1))
			get1()
			pt.X += 11
			pt.Y -= 11
		default:
			pt = image.Pt(X, 23*(Y+1))
			get1()
			pt.X += 23
		}
		// In the original v74 is overwritten with the pointer returned by the second call.
		l2 = c.Sub469920(pt)
		off := def.DrawOffset(v84, int(wl.Field2), v73)
		v46 := v82 + off.X - 51
		v47 := -73 - off.Y + v7
		c.r.Data().SetMultiply14(1)
		c.r.SetColorMultAndIntensityRGB(byte(l1[0]), byte(l1[1]), byte(l1[2]))
		if v72&2 == 0 {
			spr := def.Sprite(v84, int(wl.Field2), v73)
			c.edgeDraw(spr, v46, v47, l1[:], l2, nox_win_height, a3, a4, 0, v72)
			c.wallDone(wl)
			return
		}
		if c.sub_47D380(a3, a4) == 0 {
			c.wallDone(wl)
			return
		}
		c.r.Data().SetAlphaEnabled(true)
		c.r.Data().SetAlpha(0x80)
		c.r.SetInterlacing(legacy.Get_nox_client_highResFrontWalls_80820() == 0, vp.World.Min.Y)
		fx, fy = v46, v47
		final1 = def.Sprite(v84, int(wl.Field2), v73)
	} else {
		pt := image.Pt(23*int(wl.X5)+11, 23*int(wl.Y6)+11)
		l := c.Sub469920(pt)
		off := def.DrawOffset(v84, int(wl.Field2), v73)
		v56 := v82 + off.X - 50
		v57 := -72 - off.Y + v7
		c.r.Data().SetMultiply14(1)
		c.r.SetColorMultAndIntensityRGB(byte(l[0]), byte(l[1]), byte(l[2]))
		if v72&2 == 0 {
			if c.sub_47D380(a3, a4) != 0 {
				c.r.SetInterlacing(legacy.Get_nox_client_highResFrontWalls_80820() == 0, vp.World.Min.Y)
				spr := def.Sprite(v84, int(wl.Field2), v73)
				c.r.DrawImageAt(c.wallImage(spr), image.Pt(v56, v57))
				c.r.SetInterlacing(false, 0)
			}
			c.wallDone(wl)
			return
		}
		if c.sub_47D380(a3, a4) == 0 {
			c.wallDone(wl)
			return
		}
		c.r.Data().SetAlphaEnabled(true)
		c.r.Data().SetAlpha(0x80)
		c.r.SetInterlacing(legacy.Get_nox_client_highResFrontWalls_80820() == 0, vp.World.Min.Y)
		fx, fy = v56, v57
		final1 = def.Sprite(v84, int(wl.Field2), v73)
	}
	c.r.DrawImageAt(c.wallImage(final1), image.Pt(fx, fy))
	c.r.SetInterlacing(false, 0)
	c.r.Data().SetAlphaEnabled(false)
	c.wallDone(wl)
}

// LABEL_106 of nox_xxx_drawWalls_473C10.
func (c *Client) wallDone(wl *server.Wall) {
	c.r.Data().SetMultiply14(0)
	wl.Field3 = 0
	wl.Flags4 &= 0xFC
	wl.Field12 = 1
}
