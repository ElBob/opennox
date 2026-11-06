package opennox

import (
	"image"
	"unsafe"

	noxcolor "github.com/opennox/libs/color"
	"github.com/opennox/libs/strman"

	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
)

// Quest (Solo Quest) briefing splash screen. Port of sub_44F300 from client/Gui/GUIBrief.c.
// The drawables and the font stay in the original C globals (dword_5d4594_8324xx..8325xx);
// the blink state stays in the 0x587000+123012 blob slot.
const briefSrcFile = "C:\\NoxPost\\src\\client\\Gui\\GUIBrief.c"

func init() {
	legacy.Sub_44F300 = sub_44F300
}

// sub_44F300(a1 = window, a2 = window draw data pointer; the decompiler passes pointers as ints).
func sub_44F300(a1, a2 int) int {
	c := noxClient
	r := c.r
	wd := (*gui.WindowData)(unsafe.Pointer(uintptr(a2)))

	load := func(id string) string { // nox_strman_loadString_40F1D0(id, 0, briefSrcFile, line)
		return c.Strings().GetStringInFile(strman.ID(id), briefSrcFile)
	}
	setColor := func(cl uint32) { // nox_xxx_drawSetTextColor_434390
		r.Data().SetTextColor(noxcolor.RGBA5551(cl))
	}
	size := func(font unsafe.Pointer, s string) int { // nox_xxx_drawGetStringSize_43F840(font, s, &w, 0, 0)
		return r.GetStringSizeWrapped(r.GetFonts().AsFont(font), s, 0).X
	}
	draw := func(font unsafe.Pointer, s string, x, y int) int { // nox_xxx_drawString_43F6E0
		return r.DrawString(r.GetFonts().AsFont(font), s, image.Pt(x, y))
	}
	wrap := func(font unsafe.Pointer, s string, x, y, maxW int) { // nox_xxx_drawStringWrap_43FAF0(..., maxW, 0)
		r.DrawStringWrapped(r.GetFonts().AsFont(font), s, image.Rect(x, y, x+maxW, y))
	}
	black := legacy.Get_nox_color_black_2650656
	white := legacy.Get_nox_color_white_2523948
	orange := legacy.Get_nox_color_orange_2614256
	font := legacy.Get_dword_5d4594_832484

	vp := c.Viewport() // nox_draw_getViewport_437250
	// place is sub_473A10 + the drawable draw call (*(draw_func*)(d + 300))(vp, d).
	place := func(d *client.Drawable, x, y int) {
		d.PosVec = image.Pt(x+vp.World.Min.X-vp.Screen.Min.X, y+vp.World.Min.Y-vp.Screen.Min.Y)
		d.CallDraw(vp)
	}

	legacy.Sub_44E110()
	nox_client_resetScreenParticles_431510()
	legacy.Nox_xxx_bookHideMB_45ACA0(1)
	legacy.Sub_446780()
	v3 := (nox_win_width - noxDefaultWidth) / 2
	v4 := (nox_win_height - noxDefaultHeight) / 2

	// Title with outline.
	v5 := load("GeneralPrint:QuestSplash1")
	v31 := size(wd.FontPtr, v5)
	v29 := v3 - v31/2 + 320
	setColor(black())
	draw(wd.FontPtr, v5, v29-1, v4+19)
	draw(wd.FontPtr, v5, v29+1, v4+19)
	draw(wd.FontPtr, v5, v29-1, v4+21)
	draw(wd.FontPtr, v5, v29+1, v4+21)
	setColor(orange())
	draw(wd.FontPtr, v5, v29, v4+20)

	place(legacy.Get_dword_5d4594_832492(), v3+73, v4+123)

	// 2a/2b
	x, y := v3+109, v4+76
	v6 := load("GeneralPrint:QuestSplash2a")
	setColor(orange())
	draw(font(), v6, x, y)
	v27 := size(font(), v6)
	x += v27 + 4
	setColor(white())
	v7 := load("GeneralPrint:QuestSplash2b")
	wrap(font(), v7, x, y, v3-x+520)

	place(legacy.Get_dword_5d4594_832496(), v3+565, v4+117)

	// pair draws a label/value pair either right-aligned at right (if both fit in maxFit) or
	// at (xWrap, y) with the value wrapped.
	pair := func(idA, idB string, right, maxFit, y, xWrap int) {
		a := load(idA)
		wa := size(font(), a)
		b := load(idB)
		wb := size(font(), b)
		if wb+wa <= maxFit {
			x := right - wb - wa - 4
			setColor(orange())
			draw(font(), a, x, y)
			x = right - wb
			setColor(white())
			draw(font(), b, x, y)
		} else {
			x := xWrap
			setColor(orange())
			draw(font(), a, x, y)
			x += wa + 4
			setColor(white())
			wrap(font(), b, x, y, right-x)
		}
	}
	// labelWrap draws label (orange) at (x, y) then the value wrapped to the right.
	labelWrap := func(idA, idB string, x, y, wrapRight int) {
		a := load(idA)
		setColor(orange())
		draw(font(), a, x, y)
		w := size(font(), a)
		x += w + 4
		setColor(white())
		b := load(idB)
		wrap(font(), b, x, y, wrapRight-x)
	}

	// 3a/3b
	pair("GeneralPrint:QuestSplash3a", "GeneralPrint:QuestSplash3b", v3+520, 390, v4+115, v3+199)

	place(legacy.Get_dword_5d4594_832504(), v3+133, v4+192)

	// 4a/4b. Quirk (kept from the original): the wrap right edge is v4+630, i.e. relative to the y offset, not v3.
	labelWrap("GeneralPrint:QuestSplash4a", "GeneralPrint:QuestSplash4b", v3+157, v4+156, v4+630)

	place(legacy.Get_dword_5d4594_832500(), v3+525, v4+222)

	// 7a/7b
	pair("GeneralPrint:QuestSplash7a", "GeneralPrint:QuestSplash7b", v3+500, 215, v4+198, v3+250)

	place(legacy.Get_dword_5d4594_832528(), v3+182, v4+262)
	place(legacy.Get_dword_5d4594_832536(), v3+201, v4+251)
	place(legacy.Get_dword_5d4594_832532(), v3+185, v4+234)

	// 5a/5b
	labelWrap("GeneralPrint:QuestSplash5a", "GeneralPrint:QuestSplash5b", v3+221, v4+240, v3+470)

	place(legacy.Get_dword_5d4594_832516(), v3+484, v4+278)
	place(legacy.Get_dword_5d4594_832520(), v3+503, v4+303)

	// 6a/6b
	pair("GeneralPrint:QuestSplash6a", "GeneralPrint:QuestSplash6b", v3+462, 350, v4+286, v3+113)

	place(legacy.Get_dword_5d4594_832512(), v3+186, v4+333)
	place(legacy.Get_dword_5d4594_832508(), v3+219, v4+345)
	place(legacy.Get_dword_5d4594_832524(), v3+220, v4+322)

	// 8a/8b
	labelWrap("GeneralPrint:QuestSplash8a", "GeneralPrint:QuestSplash8b", v3+241, v4+330, v3+550)

	// 9, 10, 11: centered pairs.
	centered := func(idA, idB string, y int) {
		a := load(idA)
		b := load(idB)
		wa := size(font(), a)
		wb := size(font(), b)
		x := v3 - (wa+wb)/2 + 320
		setColor(orange())
		draw(font(), a, x, y)
		x += wa + 4
		setColor(white())
		draw(font(), b, x, y)
	}
	centered("GeneralPrint:QuestSplash9a", "GeneralPrint:QuestSplash9b", v4+370)
	centered("GeneralPrint:QuestSplash10a", "GeneralPrint:QuestSplash10b", v4+395)
	centered("GeneralPrint:QuestSplash11a", "GeneralPrint:QuestSplash11b", v4+420)

	// Blinking prompt (the blink flag lives in the 0x587000 blob).
	frame := noxServer.Frame()
	flag := memmap.PtrUint32(0x587000, 123012)
	result := int(frame / 0x1E)
	if frame%0x1E != 0 {
		if *flag != 1 {
			return result
		}
	} else {
		result = 1
		if *flag == 1 {
			*flag = frame % 0x1E
			return result
		}
		*flag = 1
	}
	v24 := white()
	v25 := load("GeneralPrint:QuestSplash12")
	v27 = size(font(), v25)
	setColor(v24)
	return draw(wd.FontPtr, v25, v3-v27/2+320, v4+450)
}
