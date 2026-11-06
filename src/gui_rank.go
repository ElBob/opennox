package opennox

import (
	"image"
	"unsafe"

	noxcolor "github.com/opennox/libs/color"
	"github.com/opennox/libs/strman"

	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

// Player rank (scoreboard) window. Port of client/Gui/guirank.c and the related part of GAME2_1.c.
//
// Memory layout is kept as in the original (mods may read it directly):
//   - 0x5D4594+1084036: column widths (5 dwords), 1084056: 3 class name pointers
//   - 0x5D4594+1084132: player table, 80 bytes per entry:
//     +0 name, +52 team (int), +56 class (byte), +60 net code, +64 score, +68 ping, +72 flag type, +76 flags
//   - 0x5D4594+1087204: team table, 56 bytes per entry:
//     +0 name, +44 team id, +48 score, +52 color (byte)
//   - 0x5D4594+1090060..1090096: list box windows (5 columns x 2 rows)
//   - 0x5D4594+1090116/17/18: team count, player count, players with score
const rankSrcFile = "C:\\NoxPost\\src\\client\\Gui\\guirank.c"

const (
	rankOffTable   = 1084132 // player table
	rankOffTeams   = 1087204 // team table
	rankOffWidths  = 1084036 // column widths
	rankOffClasses = 1084056 // class name pointers
	rankOffLists   = 1090060 // list box windows
	rankOffCount   = 1090116 // team count; +1 player count; +2 players with score
)

func rankU32(off uintptr) *uint32 { return memmap.PtrUint32(0x5D4594, off) }
func rankU8(off uintptr) *uint8   { return memmap.PtrUint8(0x5D4594, off) }
func rankW(off uintptr) *uint16   { return (*uint16)(memmap.PtrOff(0x5D4594, off)) }
func rankW587(off uintptr) *uint16 {
	return (*uint16)(memmap.PtrOff(0x587000, off))
}

// rankWin reads a window pointer stored in the 0x5D4594 blob.
func rankWin(off uintptr) *gui.Window {
	return (*gui.Window)(unsafe.Pointer(uintptr(*rankU32(off))))
}

// rankPtr32 converts a pointer to the 32-bit value that was passed in the original C code.
func rankPtr32(p unsafe.Pointer) uint32 { return uint32(uintptr(p)) }

func rankWcscpy(dst, src unsafe.Pointer) { // nox_wcscpy
	n := uintptr(0)
	for {
		c := *(*uint16)(unsafe.Add(src, 2*n))
		*(*uint16)(unsafe.Add(dst, 2*n)) = c
		if c == 0 {
			break
		}
		n++
	}
}

func rankWcslen(p *uint16) int { return alloc.StrLen(p) }

// rankLoad is nox_strman_loadString_40F1D0(id, 0, "guirank.c", line).
func rankLoad(id string) *uint16 {
	s := noxClient.Strings().GetStringInFile(strman.ID(id), rankSrcFile)
	return alloc.InternCString16(s)
}

func rankFmt(s string) *uint16 { return alloc.InternCString16(s) } // wide string literal

func rankPlU32(p *server.Player, off uintptr) *uint32 { return (*uint32)(unsafe.Add(p.C(), off)) }
func rankTmU32(t *server.Team, off uintptr) *uint32   { return (*uint32)(unsafe.Add(t.C(), off)) }

func rankStrSize(s *uint16, maxW int) (w, h int) { // nox_xxx_drawGetStringSize_43F840(0, s, &w, &h, maxW)
	sz := noxClient.r.GetStringSizeWrapped(nil, alloc.GoString16(s), maxW)
	return sz.X, sz.Y
}

func rankSetColor2(c uint32) { // nox_client_drawSetColor_434460
	noxClient.r.Data().SetColor2(noxcolor.RGBA5551(c))
}

func rankCall94(win *gui.Window, ev int, a1, a2 uintptr) int { // nox_window_call_field_94
	if win == nil {
		return 0
	}
	r := win.Func94(gui.AsWindowEvent(ev, a1, a2))
	if r == nil {
		return 0
	}
	return int(r.EventRespC())
}

// rankWndIsShown is wndIsShown_nox_xxx_wndIsShown_46ACC0 (the name is misleading: it returns 1 for hidden windows).
func rankWndIsShown(win *gui.Window) int {
	if win == nil {
		return 1
	}
	return bool2int(win.GetFlags().IsHidden())
}

func rankSetHidden(win *gui.Window, hidden int) int { // nox_window_set_hidden
	if win == nil {
		return -2
	}
	if hidden != 0 {
		win.Hide()
	} else {
		win.Show()
	}
	return 0
}

// rankSwprintf is nox_swprintf on a blob buffer.
func rankSwprintf(dst, format *uint16, args ...uint32) {
	var a [3]uint32
	copy(a[:], args)
	legacy.RankSwprintf(dst, format, a[0], a[1], a[2])
}

// rankLineHeight returns the line height of a list box (*(short*)(widget + 2)).
func rankLineHeight(win *gui.Window) int {
	return int(int16((*gui.ScrollListBoxData)(win.WidgetData).Line_height))
}

// rankLineInd returns *(short*)(widget + 46) of a list box.
func rankLineInd(win *gui.Window) int {
	return int(int16((*gui.ScrollListBoxData)(win.WidgetData).Field_11_1))
}

func rankServGamedata(flags noxflags.GameFlag) int { // (unsigned short)nox_xxx_servGamedataGet_40A020(flags)
	return int(uint16(legacy.Nox_xxx_servGamedataGet_40A020(uint16(int16(flags)))))
}

func nox_xxx_guiDrawRank_46E870() uint32 {
	legacy.Set_dword_587000_145668(6)
	sub_46F030()
	v0 := noxClient.r.FontHeight(nil)
	v46 := v0
	*rankU32(rankOffWidths) = 80

	var v41, v42, v43, v44 int
	v41, v44 = rankStrSize(rankLoad("Flag"), nox_win_width)
	v2 := v41
	if v41 < 18 {
		v2 = 18
		v41 = 18
	}
	*rankU32(1084040) = uint32(v2 + 14)
	v41, v44 = rankStrSize(rankLoad("Score"), nox_win_width)
	v42, v43 = rankStrSize(rankLoad("HealthHeading"), nox_win_width)
	v5 := v41
	if v42 > v41 {
		v5 = v42
		v41 = v42
	}
	*rankU32(1084048) = uint32(v5 + 7)
	v41, v44 = rankStrSize(rankLoad("Ping"), nox_win_width)
	v42, v43 = rankStrSize(rankW587(145972), nox_win_width)
	v7 := v41
	if v42 > v41 {
		v7 = v42
		v41 = v42
	}
	*rankU32(1084052) = uint32(v7 + 7)
	v41, v44 = rankStrSize(rankLoad("Class"), nox_win_width)
	v42, v43 = rankStrSize(rankLoad("Warrior"), nox_win_width)
	if v42 > v41 {
		v41 = v42
	}
	v42, v43 = rankStrSize(rankLoad("Wizard"), nox_win_width)
	if v42 > v41 {
		v41 = v42
	}
	v42, v43 = rankStrSize(rankLoad("Conjurer"), nox_win_width)
	if v42 > v41 {
		v41 = v42
	}
	v42, v43 = rankStrSize(rankLoad("LivesHeading"), nox_win_width)
	v13 := v41
	if v42 > v41 {
		v13 = v42
		v41 = v42
	}
	_, _, _ = v44, v43, v41
	*rankU32(1084044) = uint32(v13 + 7)
	v14 := uint32(0)
	for off := uintptr(rankOffWidths); off < 1084056; off += 4 {
		v14 += *rankU32(off)
	}
	legacy.Set_dword_5d4594_1090040(v14)
	legacy.Set_dword_5d4594_1090044(uint32(439 - v0))
	W := int(legacy.Get_dword_5d4594_1090040())
	H := int(legacy.Get_dword_5d4594_1090044())

	root := noxClient.GUI.NewWindowRaw(nil, 1560, 0, v0+40, 1, 1, nil)
	legacy.Set_dword_5d4594_1090048(root)
	root.SetAllFuncs(sub_46F060, sub_46F080, nil)
	// Original writes 0x80000000 to the window data colors by offsets +56, +64, +72, +80, +88.
	root.DrawData().BgColorVal = 0x80000000
	root.DrawData().EnColorVal = 0x80000000
	root.DrawData().HlColorVal = 0x80000000
	root.DrawData().DisColorVal = 0x80000000
	root.DrawData().SelColorVal = 0x80000000

	// v49 is a window data struct on the stack; v48 is the list box data (both zeroed first).
	v49, free49 := alloc.New(gui.WindowData{})
	defer free49()
	v48, free48 := alloc.New(gui.ScrollListBoxData{})
	defer free48()
	rankWcscpy(unsafe.Add(v49.C(), 72), unsafe.Pointer(rankW(1090136))) // text
	v49.TextColorVal = *memmap.PtrUint32(0x85B3FC, 940)
	v49.BgColorVal = 0x80000000
	// *(uint16*)&v48[6] = 0 is the high half of Field_1 (already zero)
	v48.Field_3 = 0
	v48.Field_4 = 0
	v49.EnColorVal = 0x80000000
	v49.HlColorVal = 0x80000000
	v49.DisColorVal = 0x80000000
	v49.SelColorVal = 0x80000000
	v47 := 3 * v0
	v18 := 3*v0 + 1
	v49.Style = 32
	// *(uint32*)&v48[2] = (unsigned short)(v0+1): line height, the high half overlaps Field_1 (stays zero)
	v48.Line_height = uint16(v0 + 1)
	v48.Field_1 = 0
	// strcpy(v48, "@"): Count (uint16 at offset 0) becomes 0x40
	v48.Count = 0x40
	v48.Field_2 = 1
	v19 := 2 * (v0 + 1)
	v20 := 2 * v0
	for v17 := 0; v17 < 2; v17++ {
		v21 := legacy.Nox_gui_newScrollListBox_4A4310(root, 1088, v17*W, v18, W, H-v19, v49, v48)
		v22 := int(*rankU32(rankOffWidths))
		v23 := H - v19
		*rankU32(1090052 + uintptr(4*v17)) = rankPtr32(v21.C())
		v24 := legacy.Nox_gui_newScrollListBox_4A4310(v21, 1088, 0, v20, v22, v23, v49, v48)
		v25 := H
		*rankU32(rankOffLists + uintptr(4*v17)) = rankPtr32(v24.C())
		v26 := legacy.Nox_gui_newScrollListBox_4A4310(rankWin(1090052+uintptr(4*v17)), 1088,
			int(*rankU32(rankOffWidths)), v20, int(*rankU32(1084040)), v25-v19, v49, v48)
		v27 := H
		*rankU32(1090068 + uintptr(4*v17)) = rankPtr32(v26.C())
		*rankU32(1090076 + uintptr(4*v17)) = rankPtr32(legacy.Nox_gui_newScrollListBox_4A4310(rankWin(1090052+uintptr(4*v17)), 1088,
			int(*rankU32(rankOffWidths)+*rankU32(1084040)), v20, int(*rankU32(1084044)), v27-v19, v49, v48).C())
		v28 := legacy.Nox_gui_newScrollListBox_4A4310(rankWin(1090052+uintptr(4*v17)), 1088,
			int(*rankU32(rankOffWidths)+*rankU32(1084040)+*rankU32(1084044)), v20, int(*rankU32(1084048)), H-v19, v49, v48)
		v29 := H
		*rankU32(1090084 + uintptr(4*v17)) = rankPtr32(v28.C())
		*rankU32(1090092 + uintptr(4*v17)) = rankPtr32(legacy.Nox_gui_newScrollListBox_4A4310(rankWin(1090052+uintptr(4*v17)), 1088,
			int(*rankU32(rankOffWidths)+*rankU32(1084040)+*rankU32(1084048)+*rankU32(1084044)), v20, int(*rankU32(1084052)), v29-v19, v49, v48).C())
		rankWin(1090052 + uintptr(4*v17)).SetFunc94(nox_xxx_Proc_46F070)
		rankWin(1090060 + uintptr(4*v17)).SetParent(rankWin(1090052 + uintptr(4*v17)))
		rankWin(1090068 + uintptr(4*v17)).SetParent(rankWin(1090052 + uintptr(4*v17)))
		rankWin(1090076 + uintptr(4*v17)).SetParent(rankWin(1090052 + uintptr(4*v17)))
		rankWin(1090084 + uintptr(4*v17)).SetParent(rankWin(1090052 + uintptr(4*v17)))
		rankWin(1090092 + uintptr(4*v17)).SetParent(rankWin(1090052 + uintptr(4*v17)))
	}
	v49.Style = 2048
	v49.TextColorVal = legacy.Get_nox_color_yellow_2589772()
	v45 := gui.StaticTextData{Text: rankLoad("yourrank")}
	v30 := v46 + 1
	st := noxClient.GUI.NewStaticTextRaw(root, 1088, 0, v46, W, v46+1, v49, &v45)
	legacy.Set_dword_5d4594_1090100(st)
	v49.TextColorVal = legacy.Get_nox_color_white_2523948()
	v45 = gui.StaticTextData{Text: rankLoad("WindowDir:Empty")}
	legacy.Set_dword_5d4594_1090112(noxClient.GUI.NewStaticTextRaw(root, 1088, 0, v20, W, v30, v49, &v45))
	v49.TextColorVal = *memmap.PtrUint32(0x85B3FC, 940)
	v45 = gui.StaticTextData{Text: rankLoad("WindowDir:Empty")}
	legacy.Set_dword_5d4594_1090108(noxClient.GUI.NewStaticTextRaw(root, 1088, 0, v47, W, v30, v49, &v45))
	v49.TextColorVal = legacy.Get_dword_8531A0_2572()
	v45 = gui.StaticTextData{Text: rankLoad("TeamPlayerRank")}
	*rankU32(1090104) = rankPtr32(noxClient.GUI.NewStaticTextRaw(root, 1088, 0, 0, W, v30, v49, &v45).C())
	result := rankPtr32(root.C())
	legacy.Set_dword_587000_145664(1)
	return result
}

// sub_46DC60 shows a string (or the "InternalError" string if it is NULL) in a list box.
func sub_46DC60(win *gui.Window, color uint8, text *uint16) uintptr {
	if text != nil {
		return uintptr(sub_46DC00(win, color, text))
	}
	result := rankLoad("InternalError")
	if result != nil {
		return uintptr(sub_46DC00(win, color, result))
	}
	return 0
}

func sub_46F030() *uint16 {
	var result *uint16
	for i := uintptr(0); i < 12; i += 4 {
		id := alloc.GoString((*byte)(*memmap.PtrPtr(0x587000, 145676+i)))
		result = rankLoad(id)
		*rankU32(rankOffClasses + i) = rankPtr32(unsafe.Pointer(result))
	}
	return result
}

func sub_46F060(_ *gui.Window, _ gui.WindowEvent) gui.WindowEventResp { return nil }

func nox_xxx_Proc_46F070(_ *gui.Window, _ gui.WindowEvent) gui.WindowEventResp { return nil }

func sub_46F080(win *gui.Window, draw *gui.WindowData) int {
	var v2 *server.Team
	var v35 *server.Team
	if !noxflags.HasGame(8) && legacy.Get_dword_587000_145668() != 6 {
		legacy.Set_dword_5d4594_1090120(int(legacy.Get_dword_587000_145668() - 1))
		if int32(legacy.Get_dword_587000_145668())-1 < 0 {
			legacy.Set_dword_5d4594_1090120(5)
		}
		sub_4703F0()
		legacy.Set_dword_587000_145668(6)
		if legacy.Get_dword_5d4594_1090120() == 0 {
			return 1
		}
	}
	pos := win.GlobalPos()
	xLeft, yTop := pos.X, pos.Y
	var v6 int
	if noxflags.HasGame(1) {
		v6 = rankServGamedata(noxflags.GetGame())
	} else {
		v6 = int(getSettings2ByInd(0).Field54)
	}
	v34 := v6
	if !win.Flags.Has(gui.StatusImage) { // (signed char)flags >= 0
		if draw.BgColorVal != 0x80000000 {
			noxClient.r.DrawRectFilledAlpha(xLeft, yTop, win.SizeVal.X, win.SizeVal.Y)
		}
	} else {
		noxClient.r.DrawImageAt(draw.BackgroundImage(), image.Pt(xLeft, yTop))
	}
	// Decompiler artefact: a2 (window data pointer) is reused as a byte counter below, and as a string in the default branch.
	a2 := uint32(uintptr(unsafe.Pointer(draw)))
	if legacy.Get_dword_587000_145664() != 0 || noxServer.Frame() > *rankU32(1090124)+noxServer.TickRate() {
		noxClient.GUI.ValYYY = 1
		*rankU32(1090124) = noxServer.Frame()
		legacy.Set_dword_587000_145672(-1)
		v33 := 0
		sub_46DB80()
		sub_46DCC0()
		legacy.Set_dword_587000_145664(0)
		if v7 := nox_xxx_objGetTeamByNetCode_418C80(legacy.ClientPlayerNetCode()); v7 != nil {
			v2 = noxServer.Teams.ByID(v7.ID)
			v35 = v2
		}
		if pl := legacy.Get_dword_8531A0_2576(); pl != nil {
			v8 := *rankPlU32(pl, 3680)
			if v8&1 == 0 || v8&0x20 != 0 {
				v33 = 1
			}
		}
		v9 := legacy.Get_dword_5d4594_1090120()
		if *rankU8(rankOffCount) != 0 && (v9 == 2 || v9 == 3) {
			v10 := rankLoad("team")
			sub_46DC60(rankWin(1090060), 9, v10)
			sub_46DC60(rankWin(1090068), 9, rankW587(146512))
			sub_46DC60(rankWin(1090076), 9, rankW587(146516))
			v11 := rankLoad("score")
			sub_46DC60(rankWin(1090084), 9, v11)
			sub_46DC60(rankWin(1090092), 9, rankW587(146564))
			v12 := uint8(0)
			v31 := uint8(0)
			var v13 uint8
			if *rankU8(rankOffCount) != 0 {
				for {
					v13 = sub_46FEB0(v31)
					v14 := uintptr(56) * uintptr(v31)
					sub_46DC30(rankWin(1090060), v13, rankW587(146568), rankPtr32(memmap.PtrOff(0x5D4594, 1087204+v14)))
					sub_46DC30(rankWin(1090068), v13, rankW587(146576))
					sub_46DC30(rankWin(1090076), v13, rankW587(146580))
					sub_46DC30(rankWin(1090084), v13, rankW587(146584), *rankU32(1087252 + v14))
					sub_46DC30(rankWin(1090092), v13, rankW587(146592))
					v12++
					v31 = v12
					if !(v12 < *rankU8(rankOffCount)) {
						break
					}
				}
			} else {
				v13 = uint8(uintptr(unsafe.Pointer(win))) // unreachable (guarded by the condition above)
			}
			sub_46DC30(rankWin(1090060), v13, rankW587(146596))
			sub_46DC30(rankWin(1090068), v13, rankW587(146600))
			sub_46DC30(rankWin(1090076), v13, rankW587(146604))
			sub_46DC30(rankWin(1090084), v13, rankW587(146608))
			sub_46DC30(rankWin(1090092), v13, rankW587(146612))
			v9 = legacy.Get_dword_5d4594_1090120()
		}
		v36 := rankLineInd(rankWin(1090060))
		if *rankU8(rankOffCount+1) != 0 && (v9 == 2 || v9 == 4 || v9 == 5) {
			sub_46F8F0(0, 0)
			v9 = legacy.Get_dword_5d4594_1090120()
			if legacy.Get_dword_5d4594_1090120() == 4 {
				if *rankU8(rankOffCount+1) <= 3 {
					a2 = a2&^0xFF | uint32(*rankU8(rankOffCount+1))
				} else {
					a2 = a2&^0xFF | 3
				}
			} else {
				a2 = a2&^0xFF | uint32(*rankU8(rankOffCount+1))
			}
			v32 := uint8(0)
			if uint8(a2) != 0 {
				for {
					v15 := 0
					v16 := uintptr(v32 >> 4)
					if v32 == 16 {
						sub_46F8F0(1, v36)
					}
					e := uintptr(80) * uintptr(v32)
					if noxflags.HasGame(1024) && v34 > 0 && int32(*rankU32(1084196 + e)) >= int32(v34) {
						v15 = 1
					}
					v17 := e
					v18 := *rankU32(1084208 + v17)
					var v40 uint8
					if v18&1 == 0 || v18&0x20 != 0 || v15 != 0 {
						if int32(*rankU32(1084184 + e)) == -1 {
							v40 = 3
							if v18&0x20 != 0 || v15 != 0 {
								v40 = 2
							}
						} else {
							v19 := sub_46FE60(int(*rankU32(1084184 + e)))
							v40 = sub_46FEB0(v19)
							if *rankU8(1084208+e)&0x20 != 0 || v15 != 0 {
								v40 -= 2
							}
						}
					} else {
						v40 = 9
					}
					if *rankU32(1084192 + e) == uint32(legacy.ClientPlayerNetCode()) {
						legacy.Set_dword_587000_145672(int32(rankLineInd(rankWin(1090060 + 4*v16))))
						*rankU32(1088996) = uint32(v32 >> 4)
					}
					sub_46DC30(rankWin(1090060+4*v16), v40, rankW587(146616), rankPtr32(memmap.PtrOff(0x5D4594, rankOffTable+v17)))
					sub_46DC30(rankWin(1090076+4*v16), v40, rankW587(146624), *rankU32(rankOffClasses + 4*uintptr(*rankU8(1084188 + e))))
					if legacy.Get_dword_5d4594_1090120() != 5 || *rankU32(1084196 + e) > 0 {
						sub_46DC30(rankWin(1090084+4*v16), v40, rankW587(146632), *rankU32(1084196 + e))
					} else {
						sub_46DC30(rankWin(1090084+4*v16), v40, rankW587(146640))
					}
					sub_46DC30(rankWin(1090092+4*v16), v40, rankW587(146648), *rankU32(1084200 + e))
					if legacy.Get_dword_5d4594_1090120() == 5 {
						v20 := noxServer.Players.ByID(int(*rankU32(1084192 + e)))
						if v20 != nil {
							sub_46DC30(rankWin(1090068+4*v16), v40, rankW587(146656), rankPtr32(unsafe.Add(v20.C(), 2096)))
						}
					} else {
						var v37 uint8
						v21 := sub_46FB50(int(*rankU32(1084204 + e)), &v37)
						sub_46DC60(rankWin(1090068+4*v16), v37, v21)
					}
					v32++
					if !(v32 < uint8(a2)) {
						break
					}
				}
				v9 = legacy.Get_dword_5d4594_1090120()
				v2 = v35
			}
		} else if v9 == 1 {
			sub_46FFD0()
			v9 = legacy.Get_dword_5d4594_1090120()
		}
		var v22 uintptr
		switch v9 {
		case 1, 2, 4, 5:
			switch v9 {
			case 2:
				v22 = uintptr(unsafe.Pointer(rankLoad("TeamPlayerRank")))
			case 4:
				v22 = uintptr(unsafe.Pointer(rankLoad("Top3")))
			case 5:
				v22 = uintptr(unsafe.Pointer(rankLoad("WolRank")))
			case 1:
				v22 = uintptr(unsafe.Pointer(rankLoad("Noxworld.c:Quest")))
			default:
				v22 = uintptr(a2)
			}
			if !noxflags.HasGame(1) || !noxflags.HasEngine(noxflags.EngineNoRendering) {
				if legacy.Get_dword_5d4594_1090120() == 1 {
					// quest on-screen stats
					v28 := legacy.Nox_gui_getQuestStage_450B10()
					v23 := rankLoad("Noxworld.c:Stage")
					rankSwprintf(rankW(1086692), rankFmt("%s %d"), rankPtr32(unsafe.Pointer(v23)), uint32(v28))
				} else {
					v29 := *rankU8(1090118)
					v27 := uint8(sub_46FEE0())
					v24 := rankLoad("yourrank")
					rankSwprintf(rankW(1086692), rankFmt("%s %d / %d"), rankPtr32(unsafe.Pointer(v24)), uint32(v27), uint32(v29))
				}
			}
		case 3:
			var v25 uint8
			if v2 != nil {
				v33 = 1
				v25 = uint8(sub_46FF70(int(*rankTmU32(v2, 52))))
			} else {
				v25 = 0
			}
			v22 = uintptr(unsafe.Pointer(rankLoad("Teams")))
			if !noxflags.HasGame(1) || !noxflags.HasEngine(noxflags.EngineNoRendering) {
				v30 := *rankU8(rankOffCount)
				v26 := rankLoad("yourteamrank")
				rankSwprintf(rankW(1086692), rankFmt("%s %d / %d"), rankPtr32(unsafe.Pointer(v26)), uint32(v25), uint32(v30))
			}
		default:
			v22 = uintptr(a2)
		}
		if v33 != 0 {
			if rankWndIsShown(legacy.Get_dword_5d4594_1090100()) != 0 {
				rankSetHidden(legacy.Get_dword_5d4594_1090100(), 0)
			}
			rankCall94(legacy.Get_dword_5d4594_1090100(), 16385, uintptr(unsafe.Pointer(rankW(1086692))), 0)
		} else if rankWndIsShown(legacy.Get_dword_5d4594_1090100()) == 0 {
			rankSetHidden(legacy.Get_dword_5d4594_1090100(), 1)
		}
		rankCall94(rankWin(1090104), 16385, v22, 0)
		sub_46FC50()
		sub_46FD80()
	}
	if int32(legacy.Get_dword_587000_145672()) >= 0 {
		sub_46FAE0()
	}
	return 1
}

func sub_46F8F0(a1 int, a2 int) uintptr {
	v2 := a2
	b := 4 * uintptr(a1)
	if a2 > 0 {
		for {
			sub_46DC60(rankWin(1090060+b), 9, rankW587(147124))
			sub_46DC60(rankWin(1090068+b), 9, rankW587(147128))
			sub_46DC60(rankWin(1090076+b), 9, rankW587(147132))
			sub_46DC60(rankWin(1090084+b), 9, rankW587(147136))
			sub_46DC60(rankWin(1090092+b), 9, rankW587(147140))
			v2--
			if v2 == 0 {
				break
			}
		}
	}
	v3 := rankLoad("player")
	sub_46DC60(rankWin(1090060+b), 9, v3)
	sub_46DC60(rankWin(1090068+b), 9, rankW587(147188))
	var v6 *uint16
	if legacy.Get_dword_5d4594_1090120() == 1 {
		v6 = rankLoad("LivesHeading")
	} else {
		v6 = rankLoad("class")
	}
	sub_46DC60(rankWin(1090076+b), 9, v6)
	if legacy.Get_dword_5d4594_1090120() == 5 {
		v7 := rankLoad("rank")
		sub_46DC60(rankWin(1090084+b), 9, v7)
	} else {
		var v8 *uint16
		if legacy.Get_dword_5d4594_1090120() == 1 {
			v8 = rankLoad("HealthHeading")
		} else {
			v8 = rankLoad("score")
		}
		sub_46DC60(rankWin(1090084+b), 9, v8)
	}
	var v4 *uint16
	if legacy.Get_dword_5d4594_1090120() == 1 {
		v4 = rankLoad("class")
	} else {
		v4 = rankLoad("ping")
	}
	return sub_46DC60(rankWin(1090092+b), 9, v4)
}

func sub_46FB50(a1 int, a2 *uint8) *uint16 {
	var result *uint16
	switch a1 {
	case 4:
		v2 := rankLoad("Ball")
		rankSwprintf(rankW(1090024), rankFmt("<%s"), rankPtr32(unsafe.Pointer(v2)))
		*a2 = 4
		result = rankW(1090024)
	case 1:
		v4 := rankLoad("King")
		rankSwprintf(rankW(1090024), rankFmt("<%s"), rankPtr32(unsafe.Pointer(v4)))
		result = rankW(1090024)
		*a2 = 4
	case 2:
		v5 := rankLoad("Flag")
		rankSwprintf(rankW(1090024), rankFmt("<%s"), rankPtr32(unsafe.Pointer(v5)))
		result = rankW(1090024)
		*a2 = 7
	case 3:
		v6 := rankLoad("Flag")
		rankSwprintf(rankW(1090024), rankFmt("<%s"), rankPtr32(unsafe.Pointer(v6)))
		*a2 = 13
		result = rankW(1090024)
	default:
		result = rankW587(147724)
		*a2 = 4
	}
	return result
}

func sub_46FC50() int8 {
	var v2 uint32 // decompiler artefact: partially assigned (LOBYTE) result
	cond := func() bool {
		return !noxflags.HasGame(1) || noxServer.GetFlag3592() || legacy.Sub_40A180(noxflags.GetGame()) != 0
	}
	if legacy.Sub_40A220() != 0 && cond() {
		if cond() {
			if rankWndIsShown(legacy.Get_dword_5d4594_1090108()) != 0 {
				rankSetHidden(legacy.Get_dword_5d4594_1090108(), 0)
			}
			v3 := int32(legacy.Sub_40A230())
			v7 := v3 % 60000 / 1000
			v6 := v3 / 60000
			v4 := rankLoad("TimeRemaining")
			rankSwprintf(rankW(1084068), v4, uint32(v6), uint32(v7))
			v2 = v2&^0xFF | uint32(uint8(rankCall94(legacy.Get_dword_5d4594_1090108(), 16385, uintptr(unsafe.Pointer(rankW(1084068))), 0)))
		}
	} else {
		v2 = uint32(rankWndIsShown(legacy.Get_dword_5d4594_1090108()))
		if v2 == 0 {
			v2 = v2&^0xFF | uint32(uint8(rankSetHidden(legacy.Get_dword_5d4594_1090108(), 1)))
		}
	}
	return int8(v2)
}

func sub_46FD80() int {
	var result int
	if noxflags.HasGame(4224) {
		result = rankWndIsShown(legacy.Get_dword_5d4594_1090112())
		if result == 0 {
			result = rankSetHidden(legacy.Get_dword_5d4594_1090112(), 1)
		}
	} else {
		if rankWndIsShown(legacy.Get_dword_5d4594_1090112()) != 0 {
			rankSetHidden(legacy.Get_dword_5d4594_1090112(), 0)
		}
		var v2 int
		if noxflags.HasGame(1) {
			v2 = rankServGamedata(noxflags.GetGame())
		} else {
			v2 = int(getSettings2ByInd(0).Field54)
		}
		v4 := v2
		v3 := rankLoad("LessonLimit")
		rankSwprintf(rankW(1083972), v3, uint32(v4))
		result = rankCall94(legacy.Get_dword_5d4594_1090112(), 16385, uintptr(unsafe.Pointer(rankW(1083972))), 0)
	}
	return result
}

func sub_46DB80() int {
	var result int
	for i := uintptr(0); i < 8; i += 4 {
		rankCall94(rankWin(1090060+i), 16399, 1, 0)
		rankCall94(rankWin(1090068+i), 16399, 1, 0)
		rankCall94(rankWin(1090076+i), 16399, 1, 0)
		rankCall94(rankWin(1090084+i), 16399, 1, 0)
		result = rankCall94(rankWin(1090092+i), 16399, 1, 0)
	}
	return result
}

func sub_46DC00(win *gui.Window, color uint8, text *uint16) int {
	rankCall94(win, 16397, uintptr(unsafe.Pointer(text)), uintptr(color))
	return 1
}

func sub_46DC30(win *gui.Window, color uint8, format *uint16, args ...uint32) int {
	rankSwprintf(rankW(1089000), format, args...)
	return sub_46DC00(win, color, rankW(1089000))
}

func sub_46DCC0() unsafe.Pointer {
	if legacy.Get_dword_5d4594_1090120() == 5 {
		v23 := uint32(noxServer.Players.Count())
		v43 := v23
		*rankU8(1090117) = 0
		*rankU8(1090118) = 0
		if noxflags.HasGame(1) && noxflags.HasEngine(noxflags.EngineNoRendering) {
			v23--
			v43 = v23
		}
		for i := noxServer.Players.First(); i != nil; i = noxServer.Players.Next(i) {
			v25 := *rankPlU32(i, 3680)
			if v25&1 == 0 || v25&0x20 != 0 {
				if *rankPlU32(i, 2108) == 0 {
					*rankPlU32(i, 2108) = 0x8000000
				}
			} else {
				*rankPlU32(i, 2108) |= 0x80000000
			}
		}
		v26 := uint32(0)
		if uint32(*rankU8(1090117)) < v23 {
			var v27 *server.Player // original starts with the count in this pointer variable (crashes if never replaced)
			for {
				j := uint32(0xFFFFFFFF)
				for v28 := noxServer.Players.First(); v28 != nil; v28 = noxServer.Players.Next(v28) {
					if (*(*int8)(unsafe.Add(v28.C(), 2064)) != 31 || !noxflags.HasEngine(noxflags.EngineNoRendering)) &&
						*rankPlU32(v28, 2108) >= v26 && sub_46E1E0(int(*rankPlU32(v28, 2060))) == 0 && *rankPlU32(v28, 2108) < j {
						j = *rankPlU32(v28, 2108)
						v27 = v28
					}
				}
				v30 := *rankU8(1090117)
				v31 := uintptr(80) * uintptr(v30)
				*rankU32(1084192 + v31) = *rankPlU32(v27, 2060)
				v32 := *rankPlU32(v27, 3680)
				if v32&1 == 0 || v32&0x20 != 0 {
					v33 := *rankPlU32(v27, 2108)
					if v33 == 0x8000000 {
						*rankU32(1084196 + v31) = 0
					} else {
						*rankU32(1084196 + v31) = v33
						*rankU8(1090118)++
					}
				} else {
					*rankU32(1084196 + v31) = *rankPlU32(v27, 2108) + 0x80000000
				}
				*rankU32(1084200 + v31) = uint32(*(*uint16)(unsafe.Add(v27.C(), 2148)))
				*rankU32(1084208 + v31) = *rankPlU32(v27, 3680)
				v34 := v31
				*rankU32(1084204 + v34) = uint32(sub_46E080(v27))
				rankWcscpy(memmap.PtrOff(0x5D4594, rankOffTable+v34), unsafe.Add(v27.C(), 4704))
				sub_46E170(rankW(rankOffTable + 80*uintptr(*rankU8(1090117))))
				*rankU8(1084188 + 80*uintptr(*rankU8(1090117))) = *(*uint8)(unsafe.Add(v27.C(), 2251))
				v39 := rankFillTeam(v27)
				*rankU8(1090117) = v39 + 1
				v26 = j
				if !(uint32(uint8(v39+1)) < v43) {
					break
				}
			}
		}
		for result := noxServer.Players.First(); result != nil; result = noxServer.Players.Next(result) {
			v40 := *rankPlU32(result, 3680)
			if v40&1 == 0 || v40&0x20 != 0 {
				if *rankPlU32(result, 2108) == 0x8000000 {
					*rankPlU32(result, 2108) = 0
				}
			} else {
				*rankPlU32(result, 2108) &= 0x7FFFFFFF
			}
		}
		return nil
	} else if noxflags.HasGame(1024) {
		v1 := int32(-0x80000000)
		v41 := uint32(noxServer.Players.Count())
		v2 := uint32(noxServer.Teams.Count()) // nox_xxx_getTeamCounter_417DD0 returns a byte
		v2 = uint32(uint8(v2))
		*rankU8(1090116) = 0
		v42 := v2
		*rankU8(1090117) = 0
		*rankU8(1090118) = 0
		if noxflags.HasGame(1) && noxflags.HasEngine(noxflags.EngineNoRendering) {
			v41--
		}
		if uint32(*rankU8(1090116)) < v2 {
			var v3 *server.Team
			for {
				v4 := int32(0x7FFFFFFF)
				for k := noxServer.Teams.First(); k != nil; k = noxServer.Teams.Next(k) {
					if int32(*rankTmU32(k, 52)) >= v1 && sub_46E130(int(*(*uint8)(unsafe.Add(k.C(), 57)))) == 0 && int32(*rankTmU32(k, 52)) < v4 {
						v4 = int32(*rankTmU32(k, 52))
						v3 = k
					}
				}
				rankCopyTeam(v3)
				*rankU8(1090116)++
				v1 = v4
				if !(uint32(*rankU8(1090116)) < v42) {
					break
				}
			}
		}
		for l := noxServer.Players.First(); l != nil; l = noxServer.Players.Next(l) {
			v8 := *rankPlU32(l, 3680)
			if v8&1 != 0 && v8&0x20 == 0 {
				*rankPlU32(l, 2140) += 0xFFFF
			}
		}
		v9 := int32(-1)
		if uint32(*rankU8(1090117)) < v41 {
			var v10 *server.Player
			for {
				m := int32(0x7FFFFFFF)
				for v11 := noxServer.Players.First(); v11 != nil; v11 = noxServer.Players.Next(v11) {
					if (*(*int8)(unsafe.Add(v11.C(), 2064)) != 31 || !noxflags.HasEngine(noxflags.EngineNoRendering)) &&
						int32(*rankPlU32(v11, 2140)) >= v9 && sub_46E1E0(int(*rankPlU32(v11, 2060))) == 0 && int32(*rankPlU32(v11, 2140)) < m {
						m = int32(*rankPlU32(v11, 2140))
						v10 = v11
					}
				}
				v13 := *rankU8(1090117)
				v14 := uintptr(80) * uintptr(v13)
				*rankU32(1084192 + v14) = *rankPlU32(v10, 2060)
				v15 := *rankPlU32(v10, 3680)
				if v15&1 == 0 || v15&0x20 != 0 {
					*rankU32(1084196 + v14) = *rankPlU32(v10, 2140)
					*rankU8(1090118)++
				} else {
					*rankU32(1084196 + v14) = *rankPlU32(v10, 2140) - 0xFFFF
				}
				*rankU32(1084200 + v14) = uint32(*(*uint16)(unsafe.Add(v10.C(), 2148)))
				*rankU32(1084208 + v14) = *rankPlU32(v10, 3680)
				v16 := uintptr(80) * uintptr(v13)
				*rankU32(1084204 + v16) = uint32(sub_46E080(v10))
				rankWcscpy(memmap.PtrOff(0x5D4594, rankOffTable+v16), unsafe.Add(v10.C(), 4704))
				sub_46E170(rankW(rankOffTable + 80*uintptr(*rankU8(1090117))))
				*rankU8(1084188 + 80*uintptr(*rankU8(1090117))) = *(*uint8)(unsafe.Add(v10.C(), 2251))
				v21 := rankFillTeam(v10)
				*rankU8(1090117) = v21 + 1
				v9 = m
				if !(uint32(uint8(v21+1)) < v41) {
					break
				}
			}
		}
		for result := noxServer.Players.First(); result != nil; result = noxServer.Players.Next(result) {
			v22 := *rankPlU32(result, 3680)
			if v22&1 != 0 {
				if v22&0x20 == 0 {
					*rankPlU32(result, 2140) -= 0xFFFF
				}
			}
		}
		return nil
	}
	return sub_46E4E0()
}

// rankFillTeam writes the team field of the current player entry (+52). It returns the entry index (byte 1090117).
// This is the part of sub_46DCC0 / sub_46E4E0 that is repeated verbatim in all branches.
func rankFillTeam(p *server.Player) uint8 {
	var v39 uint8
	v35 := nox_xxx_objGetTeamByNetCode_418C80(int(*rankPlU32(p, 2060)))
	if v35 != nil && v35.Has() {
		v37 := noxServer.Teams.ByID(v35.ID)
		if v37 != nil {
			v38 := uint32(*(*uint8)(unsafe.Add(v37.C(), 57)))
			v39 = *rankU8(1090117)
			*rankU32(1084184 + 80*uintptr(*rankU8(1090117))) = v38
		} else {
			v39 = *rankU8(1090117)
			*rankU32(1084184 + 80*uintptr(*rankU8(1090117))) = 0xFFFFFFFF
		}
	} else {
		v39 = *rankU8(1090117)
		*rankU32(1084184 + 80*uintptr(*rankU8(1090117))) = 0xFFFFFFFF
	}
	return v39
}

// rankCopyTeam fills the current team table entry from a team (part repeated in sub_46DCC0 and sub_46E4E0).
func rankCopyTeam(t *server.Team) {
	v6 := uintptr(56) * uintptr(*rankU8(1090116))
	*rankU32(1087252 + v6) = *rankTmU32(t, 52)
	*rankU32(1087248 + v6) = uint32(*(*uint8)(unsafe.Add(t.C(), 57)))
	*rankU8(1087256 + v6) = *(*uint8)(unsafe.Add(t.C(), 56))
	rankWcscpy(memmap.PtrOff(0x5D4594, rankOffTeams+v6), t.C())
	sub_46E170(rankW(rankOffTeams + 56*uintptr(*rankU8(1090116))))
}

func sub_46E080(p *server.Player) int {
	if noxflags.HasGame(32) {
		v1 := int32(*rankPlU32(p, 2060))
		if v1 == int32(*(*uint16)(memmap.PtrOff(0x5D4594, 1090128))) {
			return 2
		}
		if v1 == int32(*(*uint16)(memmap.PtrOff(0x5D4594, 1090130))) {
			return 3
		}
	} else if noxflags.HasGame(64) {
		if *rankPlU32(p, 2060) == uint32(*(*uint16)(memmap.PtrOff(0x5D4594, 1090132))) {
			return 4
		}
	} else if noxflags.HasGame(16) {
		v3 := noxClient.Objs.ByNetCodeDynamic(int(*rankPlU32(p, 2060)))
		if v3 != nil {
			if v3.HasEnchant(server.EnchantID(30)) {
				return 1
			}
		}
	}
	return 0
}

func sub_46E130(a1 int) int {
	v1 := 0
	if *rankU8(1090116) == 0 {
		return 0
	}
	for i := uintptr(1087248); *rankU32(i) != uint32(a1); i += 56 {
		v1++
		if v1 >= int(*rankU8(1090116)) {
			return 0
		}
	}
	return 1
}

func sub_46E170(a1 *uint16) *uint16 {
	if a1 == nil || int32(*rankU32(rankOffWidths)) == 0 {
		return nil
	}
	v1 := a1
	v2 := rankWcslen(a1)
	a1a, _ := rankStrSize(v1, 0)
	if a1a == 0 {
		return nil
	}
	if a1a+5 > int(int32(*rankU32(rankOffWidths))) {
		v3 := (*uint16)(unsafe.Add(unsafe.Pointer(v1), 2*v2))
		for {
			*v3 = 0
			v3 = (*uint16)(unsafe.Add(unsafe.Pointer(v3), -2))
			a1a, _ = rankStrSize(v1, 0)
			if !(a1a+5 > int(int32(*rankU32(rankOffWidths)))) {
				break
			}
		}
	}
	return v1
}

func sub_46E1E0(a1 int) int {
	v1 := 0
	if *rankU8(1090117) == 0 {
		return 0
	}
	for i := uintptr(1084192); *rankU32(i) != uint32(a1); i += 80 {
		v1++
		if v1 >= int(*rankU8(1090117)) {
			return 0
		}
	}
	return 1
}

func sub_46E4E0() unsafe.Pointer {
	v0 := int32(0x7FFFFFFF)
	v1 := uint32(uint8(noxServer.Teams.Count())) // nox_xxx_getTeamCounter_417DD0 returns a byte
	v24 := v1
	v23 := uint32(noxServer.Players.Count())
	*rankU8(1090116) = 0
	*rankU8(1090117) = 0
	*rankU8(1090118) = 0
	if noxflags.HasGame(1) && noxflags.HasEngine(noxflags.EngineNoRendering) {
		v23--
	}
	if uint32(*rankU8(1090116)) < v1 {
		var v2 *server.Team
		for {
			v3 := int32(-0x80000000)
			for i := noxServer.Teams.First(); i != nil; i = noxServer.Teams.Next(i) {
				if int32(*rankTmU32(i, 52)) <= v0 && sub_46E130(int(*(*uint8)(unsafe.Add(i.C(), 57)))) == 0 && int32(*rankTmU32(i, 52)) > v3 {
					v3 = int32(*rankTmU32(i, 52))
					v2 = i
				}
			}
			rankCopyTeam(v2)
			*rankU8(1090116)++
			v0 = v3
			if !(uint32(*rankU8(1090116)) < v24) {
				break
			}
		}
	}
	for j := noxServer.Players.First(); j != nil; j = noxServer.Players.Next(j) {
		v7 := *rankPlU32(j, 3680)
		if v7&1 != 0 && v7&0x20 == 0 {
			*rankPlU32(j, 2136) -= 0xFFFF
		}
	}
	v8 := int32(0x7FFFFFFF)
	if uint32(*rankU8(1090117)) < v23 {
		var v9 *server.Player
		for {
			k := int32(-0x80000000)
			for v10 := noxServer.Players.First(); v10 != nil; v10 = noxServer.Players.Next(v10) {
				if (*(*int8)(unsafe.Add(v10.C(), 2064)) != 31 || !noxflags.HasEngine(noxflags.EngineNoRendering)) &&
					int32(*rankPlU32(v10, 2136)) <= v8 && sub_46E1E0(int(*rankPlU32(v10, 2060))) == 0 && int32(*rankPlU32(v10, 2136)) > k {
					k = int32(*rankPlU32(v10, 2136))
					v9 = v10
				}
			}
			v12 := *rankU8(1090117)
			v13 := uintptr(80) * uintptr(v12)
			*rankU32(1084192 + v13) = *rankPlU32(v9, 2060)
			v14 := *rankPlU32(v9, 3680)
			if v14&1 == 0 || v14&0x20 != 0 {
				*rankU32(1084196 + v13) = *rankPlU32(v9, 2136)
				*rankU8(1090118)++
			} else {
				*rankU32(1084196 + v13) = *rankPlU32(v9, 2136) + 0xFFFF
			}
			*rankU32(1084200 + v13) = uint32(*(*uint16)(unsafe.Add(v9.C(), 2148)))
			v15 := uintptr(80) * uintptr(v12)
			*rankU32(1084204 + v15) = uint32(sub_46E080(v9))
			*rankU32(1084208 + v15) = *rankPlU32(v9, 3680)
			rankWcscpy(memmap.PtrOff(0x5D4594, rankOffTable+v15), unsafe.Add(v9.C(), 4704))
			sub_46E170(rankW(rankOffTable + 80*uintptr(*rankU8(1090117))))
			*rankU8(1084188 + 80*uintptr(*rankU8(1090117))) = *(*uint8)(unsafe.Add(v9.C(), 2251))
			v20 := rankFillTeam(v9)
			*rankU8(1090117) = v20 + 1
			v8 = k
			if !(uint32(uint8(v20+1)) < v23) {
				break
			}
		}
	}
	for result := noxServer.Players.First(); result != nil; result = noxServer.Players.Next(result) {
		v22 := *rankPlU32(result, 3680)
		if v22&1 != 0 {
			if v22&0x20 == 0 {
				*rankPlU32(result, 2136) += 0xFFFF
			}
		}
	}
	return nil
}

func sub_46FAE0() {
	idx := uintptr(*rankU32(1088996))
	lst := rankWin(1090060 + 4*idx)
	pos := lst.GlobalPos()
	v1 := pos.X
	yTop := pos.Y
	h := rankLineHeight(rankWin(1090060 + 4*idx))
	yTop += int(int32(legacy.Get_dword_587000_145672()))*h + h/2
	rankSetColor2(legacy.Get_nox_color_yellow_2589772())
	noxClient.r.DrawPoint(image.Pt(v1+1, yTop), 3, noxClient.r.Data().Color2())
}

func sub_46FE60(a1 int) uint8 {
	result := uint8(0)
	v2 := uint8(0)
	if *rankU8(1090116) == 0 {
		return 0
	}
	for *rankU32(1087248 + 56*uintptr(v2)) != uint32(a1) {
		result++
		v2 = result
		if result >= *rankU8(1090116) {
			return 0
		}
	}
	return result
}

func sub_46FEB0(a1 uint8) uint8 {
	return *memmap.PtrUint8(0x587000, 145584+8*uintptr(*rankU8(1087256 + 56*uintptr(a1))%10))
}

func sub_46FEE0() int8 {
	v0 := 0
	v1 := int8(1)
	if *rankU8(1090117) == 0 {
		return 0
	}
	for i := uintptr(1084192); *rankU32(i) != uint32(legacy.ClientPlayerNetCode()); i += 80 {
		v0++
		if v0 >= int(*rankU8(1090117)) {
			return 0
		}
	}
	v4 := *rankU32(1084196 + 80*uintptr(v0))
	if !noxflags.HasGame(1024) {
		v7 := int(*rankU8(1090117))
		if v7 != 0 {
			v8 := uintptr(1084196)
			for {
				if *rankU32(v8) > v4 {
					v1++
				}
				v8 += 80
				v7--
				if v7 == 0 {
					break
				}
			}
		}
		return v1
	}
	v5 := int(*rankU8(1090117))
	if v5 == 0 {
		return v1
	}
	v6 := uintptr(1084196)
	for {
		if *rankU32(v6) < v4 {
			v1++
		}
		v6 += 80
		v5--
		if v5 == 0 {
			break
		}
	}
	return v1
}

func sub_46FF70(a1 int) int8 {
	v1 := int8(1)
	if !noxflags.HasGame(1024) {
		v5 := int(*rankU8(1090116))
		if v5 != 0 {
			v6 := uintptr(1087252)
			for {
				if *rankU32(v6) > uint32(a1) {
					v1++
				}
				v6 += 56
				v5--
				if v5 == 0 {
					break
				}
			}
		}
		return v1
	}
	v2 := int(*rankU8(1090116))
	if v2 == 0 {
		return v1
	}
	v3 := uintptr(1087252)
	for {
		if *rankU32(v3) < uint32(a1) {
			v1++
		}
		v3 += 56
		v2--
		if v2 == 0 {
			break
		}
	}
	return v1
}

func sub_46FFD0() uint8 {
	sub_46DB80()
	sub_46F8F0(0, 0)
	result := *rankU8(1090117)
	v23 := 0
	if *rankU8(1090117) != 0 {
		v1 := 1
		v2 := uintptr(1084184)
		for {
			v3 := noxServer.Players.ByID(int(*rankU32(v2 + 8)))
			v4 := v3
			if v3 != nil && *rankPlU32(v3, 4792) != 0 {
				var v21 uint8
				if *rankU8(v2+24)&1 != 0 {
					v21 = 9
				} else if int32(*rankU32(v2)) == -1 {
					v21 = 4
				} else {
					v5 := sub_46FE60(int(*rankU32(v2)))
					v21 = sub_46FEB0(v5)
				}
				if *rankU32(v2 + 8) == uint32(legacy.ClientPlayerNetCode()) {
					legacy.Set_dword_587000_145672(int32(rankLineInd(rankWin(1090060))))
					*rankU32(1088996) = 0
				}
				v6 := v21
				sub_46DC30(rankWin(1090060), v21, rankW587(147828), rankPtr32(memmap.PtrOff(0x5D4594, v2-52)))
				sub_46DC30(rankWin(1090076), v6, rankW587(147836), uint32(*(*uint8)(unsafe.Add(v4.C(), 4816))))
				v7 := *rankU32(v2 + 8)
				v8 := int(*(*uint8)(unsafe.Add(v4.C(), 2282)))
				var v22 uint8 = 4
				if v7 == uint32(legacy.ClientPlayerNetCode()) {
					v9 := legacy.Sub_470CD0()
					v25 := float32(v9)
					v10 := legacy.Sub_470CC0()
					v16 := float32(float64(v10) / float64(v25) * 100.0)
					v8 = legacy.Nox_float2int(v16)
				}
				if v8 > 25 {
					if v8 <= 50 {
						v22 = 15
					}
				} else {
					v22 = 6
				}
				sub_46DC30(rankWin(1090084), v22, rankW587(147844), uint32(v8))
				sub_46DC30(rankWin(1090092), v6, rankW587(147856), *rankU32(rankOffClasses + 4*uintptr(*rankU8(v2 + 4))))
				var v26 uint8
				v11 := sub_46FB50(int(*rankU32(v2 + 20)), &v26)
				sub_46DC60(rankWin(1090068), v26, v11)
				if *(*uint8)(unsafe.Add(v4.C(), 4824)) != 0 {
					p := rankWin(1090068).GlobalPos()
					v19, v18 := p.X, p.Y
					_ = rankWin(1090068).Size()
					v12 := rankLineHeight(rankWin(1090068))
					v19 += 5
					v18 = v18 + v12/2 + v12*v1 - 1
					rankSetColor2(legacy.Get_nox_color_white_2523948())
					noxClient.r.DrawCircle(v19, v18, 2, noxcolor.RGBA5551(legacy.Get_nox_color_white_2523948()))
					rankDrawLines(v19, v18)
				}
				if *(*uint8)(unsafe.Add(v4.C(), 4825)) != 0 {
					p := rankWin(1090068).GlobalPos()
					v17, v20 := p.X, p.Y
					_ = rankWin(1090068).Size()
					v13 := v17 + 5
					v14 := rankLineHeight(rankWin(1090068))
					v15 := *(*uint8)(unsafe.Add(v4.C(), 4824))
					v17 += 5
					if v15 == 1 {
						v17 = v13 + 15
					}
					v20 = v20 + v14/2 + v14*v1 - 1
					rankSetColor2(legacy.Get_nox_color_yellow_2589772())
					noxClient.r.DrawCircle(v17, v20, 2, noxcolor.RGBA5551(legacy.Get_nox_color_yellow_2589772()))
					rankDrawLines(v17, v20)
				}
				v1++
			}
			result = uint8(v23 + 1)
			v2 += 80
			v23++
			if !(v23 < int(*rankU8(1090117))) {
				break
			}
		}
		legacy.Set_dword_587000_145664(1)
	} else {
		legacy.Set_dword_587000_145664(1)
	}
	return result
}

// rankDrawLines draws the three small lines that follow the circle of a player marker in sub_46FFD0.
func rankDrawLines(x, y int) {
	r := noxClient.r
	col := r.Data().Color2()
	_ = col
	r.AddPoint(image.Pt(x+2, y))
	r.AddPoint(image.Pt(x+9, y))
	r.DrawLineFromPoints(r.Data().Color2())
	r.AddPoint(image.Pt(x+9, y))
	r.AddPoint(image.Pt(x+9, y+3))
	r.DrawLineFromPoints(r.Data().Color2())
	r.AddPoint(image.Pt(x+7, y))
	r.AddPoint(image.Pt(x+7, y+2))
	r.DrawLineFromPoints(r.Data().Color2())
}
