package legacy

/*
#include <stdlib.h>
#include "defs.h"
#include "GAME3.h"
#include "client__gui__window.h"
#include "noxstring.h"
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
)

// Generic list box window proc. Port of GAME3.c (nox_xxx_wndListboxProcPre_4A30D0 and the helpers only it uses).
//
// The window data (nox_scrollListBox_data) lives on the C heap and is accessed in place. The decompiler stores
// pointers as ints in several fields (field_7/8/9 are windows, field_12 is either a selected index or a pointer to
// an int array terminated by a negative value), so those are accessed as raw 32-bit values, as in the original.

const (
	lbItemSize = 524 // sizeof(nox_scrollListBox_item)
)

func lbU32(p unsafe.Pointer, off uintptr) *uint32 { return (*uint32)(unsafe.Add(p, off)) }

// lbArr is the int array that field_12 points to (multi-select mode).
func lbArr(sdata *gui.ScrollListBoxData) *[1 << 20]int32 {
	return (*[1 << 20]int32)(unsafe.Pointer(uintptr(*lbU32(unsafe.Pointer(sdata), 48))))
}

func lbF12(sdata *gui.ScrollListBoxData) *int32 {
	return (*int32)(unsafe.Pointer(lbU32(unsafe.Pointer(sdata), 48)))
}

func lbItem(sdata *gui.ScrollListBoxData, i int) *gui.ScrollListBoxItem {
	return (*gui.ScrollListBoxItem)(unsafe.Add(unsafe.Pointer(sdata.Items), lbItemSize*i))
}

func lbRemoveSel(sdata *gui.ScrollListBoxData, i int) { // sub_4A3090
	C.sub_4A3090((*C.short)(unsafe.Pointer(sdata)), C.int(i))
}

func lbScroll(win *nox_window, a2, a3 int) int { // nox_xxx_wndListBox_4A2D10
	return int(C.nox_xxx_wndListBox_4A2D10(C.int(uintptr(unsafe.Pointer(win))), C.int(a2), C.int(a3)))
}

// lbHeight is the byte at offset 520 of the item (line height); the original writes only this byte.
func lbHeight(it *gui.ScrollListBoxItem) *uint8 { return (*uint8)(unsafe.Pointer(&it.Field_130)) }

func lbCall94(win unsafe.Pointer, a2, a3, a4 int) int {
	return nox_window_call_field_94_fnc((*nox_window)(win), a2, a3, a4, nil, 0)
}

func lbWcsncpy(dst, src unsafe.Pointer, n int) {
	C.nox_wcsncpy((*C.wchar2_t)(dst), (*C.wchar2_t)(src), C.size_t(n))
}

func lbWcslen(s unsafe.Pointer) int { return int(C.nox_wcslen((*C.wchar2_t)(s))) }

func lbMemset(p unsafe.Pointer, v byte, n int) {
	b := unsafe.Slice((*byte)(p), n)
	for i := range b {
		b[i] = v
	}
}

func lbMemmove(dst, src unsafe.Pointer, n int) {
	copy(unsafe.Slice((*byte)(dst), n), unsafe.Slice((*byte)(src), n))
}

//export nox_xxx_wndListboxProcPre_4A30D0
func nox_xxx_wndListboxProcPre_4A30D0(cwin *nox_window, ev C.uint, a3u C.uint32_t, a4i C.int) C.int {
	return C.int(listboxProcPre(cwin, uint(ev), uint32(a3u), int(a4i)))
}

func listboxProcPre(cwin *nox_window, ev uint, a3 uint32, a4 int) int {
	win := unsafe.Pointer(cwin)
	wn := (*gui.Window)(win)
	sdata := (*gui.ScrollListBoxData)(wn.WidgetData)
	draw := wn.DrawData()
	// draw_data.text is at offset 108 in the window
	dtext := (*[64]uint16)(unsafe.Add(win, 108))
	f7 := lbU32(unsafe.Pointer(sdata), 28)
	f8 := lbU32(unsafe.Pointer(sdata), 32)
	f9 := lbU32(unsafe.Pointer(sdata), 36)
	if ev > 16402 {
		switch ev {
		case 0x4013:
			ind := int(int32(a3))
			if ind < 0 || ind >= int(sdata.Count) {
				if sdata.Field_4 != 0 {
					lbMemset(unsafe.Pointer(uintptr(*lbU32(unsafe.Pointer(sdata), 48))), 0xFF, 4*int(sdata.Count))
					return 0
				}
				*lbF12(sdata) = -1
				return 0
			}
			it := lbItem(sdata, ind)
			if it.Text[0] == 0 {
				return 0
			}
			if sdata.Field_4 != 0 {
				arr := lbArr(sdata)
				arr[0] = int32(ind)
				arr[1] = -1
				return 0
			}
			v38 := int(sdata.Field_13_1)
			*lbF12(sdata) = int32(ind)
			if it.Field_0 < uint32(v38) {
				lbCall94(win, 16412, ind, 0)
				return 0
			}
			v39 := int(int16(sdata.Field_13_0))
			if it.Field_0 > uint32(v38+v39) {
				if ind <= 0 {
					sdata.Field_13_1 = 0
				} else {
					sdata.Field_13_1 = uint16(int(uint16(it.Field_0)) - v39)
				}
				lbScroll(cwin, 0, 1)
				return 0
			}
			return 0
		case 0x4014:
			return int(*lbF12(sdata))
		case 0x4015:
			ind := int(int32(a3))
			if ind < 0 || ind >= int(sdata.Count) {
				if sdata.Field_4 == 0 {
					return 0
				}
				lbMemset(unsafe.Pointer(uintptr(*lbU32(unsafe.Pointer(sdata), 48))), 0xFF, 4*int(sdata.Count))
				return 0
			}
			if lbItem(sdata, ind).Text[0] == 0 || sdata.Field_4 == 0 {
				return 0
			}
			arr := lbArr(sdata)
			v34 := 0
			v35 := arr[0]
			if v35 < 0 {
				arr[v34] = int32(ind)
				arr[v34+1] = -1
				return 0
			}
			for v35 != int32(ind) {
				v35 = arr[v34+1]
				v34++
				if v35 < 0 {
					arr[v34] = int32(ind)
					arr[v34+1] = -1
					return 0
				}
			}
			lbRemoveSel(sdata, v34)
			return 0
		case 0x4016:
			ind := int(int32(a3))
			if ind < 0 || ind >= int(sdata.Count) {
				return 0
			}
			return int(uintptr(unsafe.Pointer(&lbItem(sdata, ind).Text[0])))
		case 0x4017:
			ind := a4
			wstr := unsafe.Pointer(uintptr(a3))
			if ind < 0 || ind >= int(sdata.Count) {
				return 0
			}
			it := lbItem(sdata, ind)
			lbWcsncpy(unsafe.Pointer(&it.Text[0]), wstr, 255)
			it.Text[lbWcslen(wstr)] = 0 // TODO: potential overflow
			return 0
		case 0x4018:
			*f7 = a3
			return 0
		case 0x4019:
			*f8 = a3
			return 0
		case 0x401A:
			*f9 = a3
			return 0
		case 0x401B:
			ind := int(int32(a3))
			v40 := int(int16(sdata.Field_11_0))
			if v40 < ind {
				return 0
			}
			lbMemmove(unsafe.Pointer(sdata.Items), unsafe.Pointer(lbItem(sdata, ind)), lbItemSize*int(uint32(v40-ind)))
			sdata.Field_11_0 -= uint16(ind)
			sdata.Field_11_1 = sdata.Field_11_0
			if sdata.Field_4 != 0 {
				arr := lbArr(sdata)
				v42 := 0
				if arr[0] >= 0 {
					for {
						v43 := int(arr[v42])
						if ind < v43 {
							lbRemoveSel(sdata, v42)
							v42--
						} else {
							arr[v42] = int32(v43 - ind)
						}
						arr = lbArr(sdata)
						v42++
						if arr[v42] < 0 {
							break
						}
					}
				}
			} else {
				v44 := int(*lbF12(sdata))
				if v44 > 0 {
					*lbF12(sdata) = int32(v44 - ind)
				}
			}
			if int(sdata.Field_13_1) > 0 {
				lbScroll(cwin, -ind, 1)
			}
			listboxRecalc(cwin)
			return 1
		case 0x401C:
			ind := int(int32(a3))
			if ind-1 >= 0 && ind-1 < int(sdata.Count) {
				sdata.Field_13_1 = uint16(lbItem(sdata, ind-1).Field_0 + 1)
			} else {
				sdata.Field_13_1 = 0
			}
			v45 := int(int16(sdata.Field_13_0))
			v46 := int(int32(sdata.Field_10))
			if int(sdata.Field_13_1)+v45 >= v46 {
				sdata.Field_13_1 = uint16(v46 - v45)
			}
			lbScroll(cwin, 0, 1)
			return 0
		default:
			return 0
		}
	}
	if ev == 16402 {
		if int32(a3) >= 0 {
			v32 := int16(sdata.Field_11_0)
			if int32(a3) <= int32(v32) {
				sdata.Field_11_1 = uint16(a3)
			} else {
				sdata.Field_11_1 = uint16(v32)
			}
		} else {
			sdata.Field_11_1 = 0
		}
		return 0
	}
	if ev > 0x4007 {
		switch ev {
		case 0x4009:
			v29 := int(int16(sdata.Field_13_0))
			v30 := int16(int(*(*uint16)(unsafe.Add(unsafe.Pointer(uintptr(*(*uint32)(unsafe.Add(unsafe.Pointer(uintptr(*f9)), 32)))), 4))) - a4)
			v31 := int(int32(sdata.Field_10)) - v29 + 1
			sdata.Field_13_1 = uint16(v30)
			if int(v30) > v31 {
				sdata.Field_13_1 = uint16(int(int32(sdata.Field_10)) - v29 + 1)
			}
			if int(int32(sdata.Field_13_1)) < 0 { // always false: the field is unsigned
				sdata.Field_13_1 = 0
			}
			lbScroll(cwin, 0, 0)
			return 0
		case 16397:
			a3p := unsafe.Pointer(uintptr(a3))
			v18 := int(int16(sdata.Field_11_1))
			v19 := int(int16(sdata.Field_11_0))
			if v18 != v19 {
				if v19 < int(sdata.Count) {
					v20 := v19 - 1
					if v20 >= v18 {
						v21 := lbItemSize * v20
						for {
							v20--
							// Byte offsets, as in the original binary and as in the list-full branch below. The repo's C
							// had `(uint32_t*)sdata->items + v21` here, which scales the offset by 4 and runs out of bounds;
							// that was introduced by the "Refactor nox_xxx_wndListboxProcPre_4A30D0" commit (e421a6fe8),
							// before it the C was `*((uint32_t*)v5 + 6) + v21` (a byte offset).
							v22 := unsafe.Add(unsafe.Pointer(sdata.Items), v21)
							v21 -= lbItemSize
							lbMemmove(unsafe.Add(v22, lbItemSize), v22, 0x20C)
							if !(v20 >= int(int16(sdata.Field_11_1))) {
								break
							}
						}
					}
					listboxAddLine(a3p, a4, cwin)
				} else {
					if sdata.Field_2 == 0 {
						return 0
					}
					lbCall94(win, 16411, 1, 0)
					v23 := int(sdata.Field_11_0) - 1
					if v23 >= int(int16(sdata.Field_11_1)) {
						v24 := lbItemSize * v23
						for {
							v23--
							v25 := unsafe.Add(unsafe.Pointer(sdata.Items), v24)
							v24 -= lbItemSize
							lbMemmove(unsafe.Add(v25, lbItemSize), v25, 0x20C)
							if !(v23 >= int(int16(sdata.Field_11_1))) {
								break
							}
						}
					}
					listboxAddLine(a3p, a4, cwin)
				}
			} else if v18 >= int(sdata.Count) {
				if sdata.Field_2 == 0 {
					return 0
				}
				lbCall94(win, 16411, 1, 0)
				listboxAddLine(a3p, a4, cwin)
			} else {
				listboxAddLine(a3p, a4, cwin)
			}
			if sdata.Field_1 != 0 {
				for {
					it := (*gui.ScrollListBoxItem)(unsafe.Add(unsafe.Pointer(sdata.Items), lbItemSize*(int(sdata.Field_11_1)-1)))
					if !(it.Field_0 >= uint32(int(sdata.Field_13_1)+int(sdata.Field_13_0))) {
						break
					}
					lbScroll(cwin, 1, 1)
				}
			}
			if sdata.Field_4 != 0 {
				arr := lbArr(sdata)
				if arr[0] >= 0 {
					v27 := 0
					for {
						v28 := int(arr[v27])
						if int(sdata.Field_11_1) < v28 {
							arr[v27] = int32(v28 + 1)
						}
						arr = lbArr(sdata)
						v27++
						if arr[v27] < 0 {
							break
						}
					}
					return 1
				}
			} else if int(sdata.Field_11_1) < int(*lbF12(sdata)) {
				*lbF12(sdata) = *lbF12(sdata) + 1
			}
			return 1
		case 0x400E:
			ind := int(int32(a3))
			v12 := int(int16(sdata.Field_11_0))
			if ind < 0 || ind >= v12 {
				return 0
			}
			for i := ind; i < v12-1; i++ {
				lbMemmove(unsafe.Pointer(lbItem(sdata, i)), unsafe.Pointer(lbItem(sdata, i+1)), lbItemSize)
			}
			sdata.Field_11_0--
			v13 := int(int16(sdata.Field_11_0))
			lbMemset(unsafe.Pointer(lbItem(sdata, v13)), 0, lbItemSize)
			sdata.Field_11_1 = uint16(v13)
			lbItem(sdata, v13).Text[0] = 0
			if sdata.Field_4 != 0 {
				arr := lbArr(sdata)
				v15 := 0
				if arr[0] >= 0 {
					for {
						v16 := int(arr[v15])
						if ind >= v16 {
							if ind == v16 {
								lbRemoveSel(sdata, v15)
								v15--
							}
						} else {
							arr[v15] = int32(v16 - 1)
						}
						arr = lbArr(sdata)
						v15++
						if arr[v15] < 0 {
							break
						}
					}
					listboxRecalc(cwin)
					return 0
				}
			} else {
				v17 := int(*lbF12(sdata))
				if ind < v17 {
					*lbF12(sdata) = int32(v17 - 1)
					listboxRecalc(cwin)
					return 0
				}
				if ind == v17 {
					*lbF12(sdata) = -1
				}
			}
			listboxRecalc(cwin)
			return 0
		case 0x400F:
			lbMemset(unsafe.Pointer(sdata.Items), 0, lbItemSize*int(sdata.Count))
			if a3 != 1 {
				sdata.Field_13_1 = 0
			}
			if sdata.Field_4 != 0 {
				lbMemset(unsafe.Pointer(uintptr(*lbU32(unsafe.Pointer(sdata), 48))), 0xFF, 4*int(sdata.Count))
			} else {
				*lbF12(sdata) = -1
			}
			sdata.Field_11_1 = 0
			sdata.Field_11_0 = 0
			sdata.Field_10 = 0
			lbScroll(cwin, 0, 1)
			return 0
		default:
			return 0
		}
	}
	if ev == 16391 || ev == 0x4000 { // both events have identical handlers in the original
		if a3 == *f7 {
			if int(sdata.Field_13_1) > 0 {
				lbScroll(cwin, -1, 1)
				return 0
			}
		} else if a3 == *f8 && uint32(int(sdata.Field_13_0)+int(sdata.Field_13_1)) <= sdata.Field_10 {
			lbScroll(cwin, 1, 1)
			return 0
		}
		return 0
	}
	if ev == 16385 {
		wstr := unsafe.Pointer(uintptr(a3))
		lbWcsncpy(unsafe.Pointer(&dtext[0]), wstr, 63)
		dtext[lbWcslen(wstr)] = 0 // TODO: this may overflow
		return 0
	}
	if ev == 16388 {
		if *f7 != 0 {
			v9 := unsafe.Pointer(uintptr(*f7))
			nox_window_setPos_46A9B0((*nox_window)(v9), int(int32(a3))-int(int32(*lbU32(v9, 8))), 0)
		}
		if *f8 != 0 {
			v10 := unsafe.Pointer(uintptr(*f8))
			nox_window_setPos_46A9B0((*nox_window)(v10), int(int32(a3))-int(int32(*lbU32(v10, 8))), a4-int(int32(*lbU32(v10, 12))))
		}
		if *f9 != 0 {
			v11 := unsafe.Pointer(uintptr(*f9))
			// the original dereferences field_7 here without a check
			nox_window_setPos_46A9B0((*nox_window)(v11), int(int32(a3))-int(int32(*lbU32(v11, 8))), int(int32(*lbU32(unsafe.Pointer(uintptr(*f7)), 12))))
			w9 := unsafe.Pointer(uintptr(*f9))
			C.sub_46AB20((*C.uint)(w9), C.int(*lbU32(w9, 8)),
				C.int(a4-2*int(int32(*lbU32(unsafe.Pointer(uintptr(*lbU32(w9, 400))), 12)))))
		}
		sdata.Field_13_0 = uint16(a4)
		if dtext[0] != 0 {
			sdata.Field_13_0 -= uint16(nox_xxx_guiFontHeightMB_43F320(draw.FontPtr))
		}
		return 0
	}
	if ev > 0x4000 {
		return 0
	}
	if ev == 2 {
		if sdata != nil {
			C.free(unsafe.Pointer(sdata.Items))
			if sdata.Field_4 != 0 {
				C.free(unsafe.Pointer(uintptr(*lbU32(unsafe.Pointer(sdata), 48))))
			}
			C.free(unsafe.Pointer(sdata))
		}
		wn.WidgetData = nil
		return 0
	}
	if ev == 23 {
		if a3 != 0 {
			draw.Field0 |= 2
		} else {
			draw.Field0 &= 0xFFFFFFFD
		}
		v7 := nox_xxx_wndGetID_46B0A0(cwin)
		lbCall94(unsafe.Pointer(draw.Window), 16387, int(int32(a3)), v7)
		return 1
	}
	return 0
}

// listboxRecalc is nox_xxx_wndListBox_4A3A70.
func listboxRecalc(cwin *nox_window) int {
	sdata := (*gui.ScrollListBoxData)((*gui.Window)(unsafe.Pointer(cwin)).WidgetData)
	v1 := 0
	v2 := 0
	if int16(sdata.Field_11_0) > 0 {
		for {
			it := lbItem(sdata, v1)
			v1++
			v2 += int(*lbHeight(it)) + 1
			it.Field_0 = uint32(v2)
			if !(v1 < int(int16(sdata.Field_11_0))) {
				break
			}
		}
	}
	sdata.Field_10 = uint32(v2)
	return lbScroll(cwin, 0, 1)
}

// listboxAddLine is nox_xxx_wndListBoxAddLine_4A3AC0.
func listboxAddLine(a1 unsafe.Pointer, a2 int, cwin *nox_window) int {
	win := unsafe.Pointer(cwin)
	sdata := (*gui.ScrollListBoxData)((*gui.Window)(win).WidgetData)
	v5 := int(int32(*lbU32(win, 8))) - 7
	if sdata.Field_3 != 0 {
		v5 = int(int32(*lbU32(win, 8))) - 17
	}
	it := lbItem(sdata, int(int16(sdata.Field_11_1)))
	if a2 >= 17 || a2 < 0 {
		it.Field_129 = *lbU32(win, 104)
	} else {
		it.Field_129 = **(**uint32)(memmap.PtrOff(0x85B3FC, uintptr(132+4*a2)))
	}
	v6 := unsafe.Pointer(&it.Text[0])
	if a1 != nil {
		lbWcsncpy(v6, a1, 0xFF)
		*(*uint16)(unsafe.Add(v6, 2*255)) = 0
		v7 := (*uint16)(unsafe.Add(v6, 2*(lbWcslen(v6)-1)))
		if *v7 == 10 {
			*v7 = 0
		}
	} else {
		*(*uint16)(v6) = 32
		*(*uint16)(unsafe.Add(v6, 2)) = 0
	}
	font := unsafe.Pointer(uintptr(*lbU32(win, 236)))
	if *lbU32(win, 4)&0x4000 == 0x4000 {
		*lbHeight(it) = uint8(nox_xxx_guiFontHeightMB_43F320(font))
	} else {
		var h C.int
		nox_xxx_drawGetStringSize_43F840(font, (*wchar2_t)(v6), nil, &h, v5)
		*lbHeight(it) = uint8(h)
	}
	sdata.Field_11_1++
	sdata.Field_11_0++
	return listboxRecalc(cwin)
}
