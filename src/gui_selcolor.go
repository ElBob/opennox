package opennox

import (
	"fmt"
	"image"
	"os"
	"unsafe"

	noxcolor "github.com/opennox/libs/color"
	"github.com/opennox/libs/datapath"
	"github.com/opennox/libs/ifs"
	"github.com/opennox/libs/types"

	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/client/input"
	"github.com/opennox/opennox/v1/client/noxrender"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/binfile"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/ccall"
	"github.com/opennox/opennox/v1/server"
)

const selColorSrcFile = "C:\\NoxPost\\src\\client\\shell\\SelColor.c"

func init() {
	legacy.Nox_game_showSelColor_4A5D00 = func() int {
		if nox_game_showSelColor_4A5D00() {
			return 1
		}
		return 0
	}
	legacy.Sub_4A6890 = sub_4A6890
	legacy.Sub_4A6C90 = sub_4A6C90
}

// selColorWidget returns the widget data of the window as the original uint32 value.
// The color selection windows keep (colorIndex<<16 | colorTable) there instead of a real pointer.
func selColorWidget(win *gui.Window) *uint32 {
	return (*uint32)(unsafe.Pointer(&win.WidgetData))
}

// selColorRGB returns the color picked by a color selector window (original: table lookup at 0x5D4594+1307796).
func selColorRGB(win *gui.Window) types.RGB {
	w := *selColorWidget(win)
	idx := uintptr(uint16(w>>16)) + 32*uintptr(uint16(w))
	return *memmap.PtrT[types.RGB](0x5D4594, 1307796+3*idx)
}

// selColorMaterial sets a render material from a picked color.
func selColorMaterial(ind int, c types.RGB) {
	noxClient.r.Data().SetMaterial(ind, noxcolor.RGBA5551(nox_color_rgb_4344A0(int(c.R), int(c.G), int(c.B))))
}

func wcsLen(p *uint16) int {
	n := 0
	for *(*uint16)(unsafe.Add(unsafe.Pointer(p), 2*n)) != 0 {
		n++
	}
	return n
}

// wcsCopy is nox_wcscpy.
func wcsCopy(dst, src *uint16) {
	n := 0
	for ; *(*uint16)(unsafe.Add(unsafe.Pointer(src), 2*n)) != 0; n++ {
		*(*uint16)(unsafe.Add(unsafe.Pointer(dst), 2*n)) = *(*uint16)(unsafe.Add(unsafe.Pointer(src), 2*n))
	}
	*(*uint16)(unsafe.Add(unsafe.Pointer(dst), 2*n)) = 0
}

func selColorDefaultName() *uint16 {
	str := noxClient.Strings().GetStringInFile("DefaultName", selColorSrcFile)
	return alloc.InternCString16(str)
}

// selColorEntryText returns the text buffer of the name entry field (event 0x401d).
func selColorEntryText() *uint16 {
	r := legacy.Get_dword_5d4594_1308152().Func94(gui.AsWindowEvent(0x401d, 0, 0))
	return (*uint16)(gui.EventRespPtr(r))
}

func nox_game_showSelColor_4A5D00() bool {
	c := noxClient
	c.GameAddStateCode(client.StateColorSelect)
	host := legacy.Nox_xxx_getHostInfoPtr_431770()
	legacy.Set_dword_5d4594_1307784(host)
	host.SetIsFemale(0)
	win := newWindowFromFile(c.GUI, "SelColor.wnd", selColorProc)
	legacy.Set_dword_5d4594_1308084(win)
	if win == nil {
		return false
	}
	win.SetFunc93(sub4A18E0)
	anim := nox_gui_makeAnimation(win, 0, 0, 0, -440, 0, 20, 0, -40)
	legacy.Set_nox_wnd_xxx_1308092(anim)
	if anim == nil {
		return false
	}
	anim.StateID = client.StateColorSelect
	anim.Func12Ptr = legacy.Get_sub_4A6890()
	anim.FncDoneOutPtr = legacy.Get_sub_4A6C90()
	sub_4A5E90()
	for i := uint(720); i <= 729; i++ {
		win.ChildByID(i).SetDraw(selColorDrawSwatch)
	}
	for j := uint(761); j <= 792; j++ {
		win.ChildByID(j).SetDraw(selColorDrawSwatch)
	}
	if legacy.Get_dword_587000_171388() != 0 {
		v6 := selColorDefaultName()
		legacy.Get_dword_5d4594_1308152().Func94(gui.AsWindowEvent(16414, uintptr(unsafe.Pointer(v6)), 0))
	}
	legacy.Nox_xxx_wndRetNULL_46A8A0()
	popup := win.ChildByID(760)
	legacy.Set_dword_5d4594_1308088(popup)
	popup.SetFunc94(selColorProc)
	popup.SetFunc93(selColorPopupProc)
	popup.SetParent(nil)
	sub_4BFAD0()
	win.ChildByID(740).SetDraw(selColorDrawPreview)
	return true
}

// sub_4A5E90 binds the color selector windows and restores the previously selected colors.
func sub_4A5E90() {
	sub_4A5E90_A()
	root := legacy.Get_dword_5d4594_1308084()
	host := legacy.Get_dword_5d4594_1307784()

	v1 := root.ChildByID(720)
	legacy.Set_dword_5d4594_1308096(v1)
	*selColorWidget(v1) = 131074
	v2 := root.ChildByID(721)
	legacy.Set_dword_5d4594_1308100(v2)
	*selColorWidget(v2) = 589825
	v3 := root.ChildByID(722)
	legacy.Set_dword_5d4594_1308104(v3)
	*selColorWidget(v3) = 589825
	v4 := root.ChildByID(723)
	legacy.Set_dword_5d4594_1308108(v4)
	*selColorWidget(v4) = 589825
	v5 := root.ChildByID(724)
	legacy.Set_dword_5d4594_1308112(v5)
	*selColorWidget(v5) = 589825
	v6 := root.ChildByID(725)
	legacy.Set_dword_5d4594_1308116(v6)
	*selColorWidget(v6) = uint32(*memmap.PtrUint16(0x587000, 171372)) << 16
	v7 := root.ChildByID(726)
	legacy.Set_dword_5d4594_1308120(v7)
	*selColorWidget(v7) = uint32(*memmap.PtrUint16(0x587000, 171374)) << 16
	v8 := root.ChildByID(727)
	legacy.Set_dword_5d4594_1308124(v8)
	*selColorWidget(v8) = uint32(*memmap.PtrUint16(0x587000, 171376)) << 16
	v9 := root.ChildByID(728)
	legacy.Set_dword_5d4594_1308128(v9)
	*selColorWidget(v9) = uint32(*memmap.PtrUint16(0x587000, 171378)) << 16
	v10 := root.ChildByID(729)
	legacy.Set_dword_5d4594_1308132(v10)
	*selColorWidget(v10) = uint32(*memmap.PtrUint16(0x587000, 171380)) << 16
	legacy.Set_dword_5d4594_1308136(root.ChildByID(711))
	legacy.Set_dword_5d4594_1308140(root.ChildByID(712))
	legacy.Set_dword_5d4594_1308144(root.ChildByID(713))
	legacy.Set_dword_5d4594_1308148(root.ChildByID(714))
	legacy.Set_dword_5d4594_1308152(root.ChildByID(751))
	if legacy.Get_dword_587000_171388() == 0 {
		legacy.Get_dword_5d4594_1308152().Func94(gui.AsWindowEvent(16414, uintptr(host.C()), 0))
		sub_4A61E0(legacy.Get_dword_5d4594_1308096(), 2, &host.Colors.Skin)
		sub_4A61E0(legacy.Get_dword_5d4594_1308100(), 1, &host.Colors.Hair)
		sub_4A61E0(legacy.Get_dword_5d4594_1308104(), 1, &host.Colors.Mustache)
		sub_4A61E0(legacy.Get_dword_5d4594_1308108(), 1, &host.Colors.Goatee)
		sub_4A61E0(legacy.Get_dword_5d4594_1308112(), 1, &host.Colors.Beard)
		*selColorWidget(legacy.Get_dword_5d4594_1308116()) = uint32(host.Colors.Pants) << 16
		*selColorWidget(legacy.Get_dword_5d4594_1308120()) = uint32(host.Colors.Shirt1) << 16
		*selColorWidget(legacy.Get_dword_5d4594_1308124()) = uint32(host.Colors.Shirt2) << 16
		*selColorWidget(legacy.Get_dword_5d4594_1308128()) = uint32(host.Colors.Shoes1) << 16
		*selColorWidget(legacy.Get_dword_5d4594_1308132()) = uint32(host.Colors.Shoes2) << 16
	}
}

// sub_4A61E0 finds the color in a color table (a2 is the table index) and selects it.
func sub_4A61E0(win *gui.Window, a2 int, a3 *types.RGB) {
	root := legacy.Get_dword_5d4594_1308084()
	tbl := memmap.PtrT[[32]types.RGB](0x5D4594, uintptr(1307796+96*a2))
	v3 := uint32(0)
	found := false
	for ; v3 < 0x20; v3++ {
		if tbl[v3] == *a3 {
			found = true
			break
		}
	}
	if found {
		if a2 == 1 {
			v5 := root.ChildByID(win.ID() - 10)
			legacy.Nox_xxx_wnd_46ABB0(v5, 1)
			res := root.ChildByID(win.ID() + 10)
			res.DrawData().Field0 |= 6
		}
	}
	if v3 == 32 && a2 == 1 {
		v6 := root.ChildByID(win.ID() - 10)
		legacy.Nox_xxx_wnd_46ABB0(v6, 0)
		v3 = 9 // the original sets only the low word of v3
	}
	*selColorWidget(win) = uint32(uint16(a2)) | uint32(uint16(v3))<<16
}

func sub_4A6890() int {
	legacy.Get_nox_wnd_xxx_1308092().SetState(gui.AnimOut)
	sub_43BE40(2)
	clientPlaySoundSpecial(923, 100)
	sub_4A68C0()
	return 1
}

// sub_4A68C0 stores the selected name and colors into the host player info.
func sub_4A68C0() *server.PlayerInfo {
	v0 := selColorEntryText()
	if *v0 == 0 {
		wcsCopy(v0, selColorDefaultName())
	}
	host := legacy.Get_dword_5d4594_1307784()
	name := (*uint16)(host.C())
	wcsCopy(name, v0)
	if sub_4A6B50(name) == 0 {
		wcsCopy(name, selColorDefaultName())
	}
	host.Colors.Skin = selColorRGB(legacy.Get_dword_5d4594_1308096())
	if legacy.Get_dword_5d4594_1308136().Flags.Has(gui.StatusEnabled) {
		host.Colors.Hair = selColorRGB(legacy.Get_dword_5d4594_1308100())
	} else {
		host.Colors.Hair = host.Colors.Skin
	}
	if legacy.Get_dword_5d4594_1308140().Flags.Has(gui.StatusEnabled) {
		host.Colors.Mustache = selColorRGB(legacy.Get_dword_5d4594_1308104())
	} else {
		host.Colors.Mustache = host.Colors.Skin
	}
	if legacy.Get_dword_5d4594_1308144().Flags.Has(gui.StatusEnabled) {
		host.Colors.Goatee = selColorRGB(legacy.Get_dword_5d4594_1308108())
	} else {
		host.Colors.Goatee = host.Colors.Skin
	}
	if legacy.Get_dword_5d4594_1308148().Flags.Has(gui.StatusEnabled) {
		host.Colors.Beard = selColorRGB(legacy.Get_dword_5d4594_1308112())
	} else {
		host.Colors.Beard = host.Colors.Skin
	}
	host.Colors.Pants = byte(*selColorWidget(legacy.Get_dword_5d4594_1308116()) >> 16)
	host.Colors.Shirt1 = byte(*selColorWidget(legacy.Get_dword_5d4594_1308120()) >> 16)
	host.Colors.Shirt2 = byte(*selColorWidget(legacy.Get_dword_5d4594_1308124()) >> 16)
	host.Colors.Shoes1 = byte(*selColorWidget(legacy.Get_dword_5d4594_1308128()) >> 16)
	host.Colors.Shoes2 = byte(*selColorWidget(legacy.Get_dword_5d4594_1308132()) >> 16)
	return host
}

// sub_4A6B50 trims the name and replaces file name characters. Returns 1 if the name is not empty.
func sub_4A6B50(a1 *uint16) int {
	at := func(p *uint16, i int) *uint16 { return (*uint16)(unsafe.Add(unsafe.Pointer(p), 2*i)) }
	v2 := 0
	v3 := 1
	v10 := 0
	var v12 [26]uint16 // original has the same fixed size; the name entry is limited to fit it
	v4 := wcsLen(a1)
	if v4 >= 1 {
		v1 := 0
		v5 := 0
		v11 := v4
		for {
			if legacy.Iswspace(*at(a1, v1)) {
				if v3 == 0 {
					v12[v5] = *at(a1, v1)
					v5++
					v10++
				}
			} else {
				if v3 == 1 {
					v6 := *at(a1, v1)
					if v6 == 42 || v6 == 63 || v6 == 60 || v6 == 62 || v6 == 92 || v6 == 47 || v6 == 58 || v6 == 34 || v6 == 124 {
						*at(a1, v1) = 45
					}
				}
				v12[v5] = *at(a1, v1)
				v5++
				v10++
				v3 = 0
				v2 = 1
			}
			v1++
			v11--
			if v11 == 0 {
				break
			}
		}
		v12[v10] = 0
		if v2 != 0 {
			wcsCopy(a1, &v12[0])
			v7 := wcsLen(&v12[0]) - 1
			if v7 >= 0 {
				for iswsp := legacy.Iswspace(*at(a1, v7)); iswsp; iswsp = legacy.Iswspace(*at(a1, v7)) {
					v7--
					if v7 < 0 {
						return v2
					}
				}
				*at(a1, v7+1) = 0
			}
		}
	}
	return v2
}

func sub_4A6C90() int {
	anim := legacy.Get_nox_wnd_xxx_1308092()
	v0 := anim.Func13Ptr
	anim.Free()
	legacy.Get_dword_5d4594_1308084().Destroy()
	legacy.Get_dword_5d4594_1308088().Destroy()
	if v0 != nil {
		ccall.CallIntVoid(v0)
	} else {
		nox_client_resetScreenParticles_431510()
		noxClient.GUI.Draw()
		if !noxflags.HasGame(noxflags.GameOnline) {
			nox_client_guiXxxDestroy_4A24A0()
			return 1
		}
		nox_xxx_serverHost_43B4D0()
		return 1
	}
	return 1
}

// selColorDrawSwatch draws a color swatch (sub_4A6D20).
func selColorDrawSwatch(win *gui.Window, _ *gui.WindowData) int {
	pos := win.GlobalPos()
	h := win.EndPos.Y - win.Off.Y
	w := win.EndPos.X - win.Off.X
	c := selColorRGB(win)
	r := noxClient.r
	r.Data().SetColor2(noxcolor.RGBA5551(nox_color_rgb_4344A0(int(c.R), int(c.G), int(c.B))))
	r.DrawRectFilledOpaque(pos.X, pos.Y, w, h, r.Data().Color2())
	return 1
}

func selColorBlobImage(off uintptr) *noxrender.Image {
	return noxClient.r.Bag.AsImage(noxrender.ImageHandle(*memmap.PtrPtr(0x973A20, off)))
}

// selColorDrawPreview draws the character preview (sub_4A6DC0).
func selColorDrawPreview(win *gui.Window, _ *gui.WindowData) int {
	r := noxClient.r
	pos := win.GlobalPos()
	host := legacy.Get_dword_5d4594_1307784()
	w1308096 := legacy.Get_dword_5d4594_1308096()
	w1308100 := legacy.Get_dword_5d4594_1308100()
	w1308104 := legacy.Get_dword_5d4594_1308104()
	w1308108 := legacy.Get_dword_5d4594_1308108()
	w1308112 := legacy.Get_dword_5d4594_1308112()

	v4 := selColorRGB(w1308096)
	selColorMaterial(6, v4)
	selColorMaterial(1, selColorRGB(w1308100))
	if legacy.Get_dword_5d4594_1308144().Flags.Has(gui.StatusEnabled) {
		selColorMaterial(2, selColorRGB(w1308108))
	} else {
		selColorMaterial(2, v4)
	}
	selColorMaterial(3, v4)
	if legacy.Get_dword_5d4594_1308148().Flags.Has(gui.StatusEnabled) {
		selColorMaterial(4, selColorRGB(w1308112))
	} else {
		selColorMaterial(4, v4)
	}
	if legacy.Get_dword_5d4594_1308140().Flags.Has(gui.StatusEnabled) {
		selColorMaterial(5, selColorRGB(w1308104))
	} else {
		selColorMaterial(5, v4)
	}
	female := uintptr(*(*byte)(unsafe.Add(host.C(), 67)))
	if legacy.Get_dword_5d4594_1308136().Flags.Has(gui.StatusEnabled) {
		r.DrawImageAt(selColorBlobImage(16+4*female), pos)
	} else {
		r.DrawImageAt(selColorBlobImage(24+4*female), pos)
	}
	imgs := 32 + 104*female
	for v21 := 0; v21 < 26; v21++ {
		var mod unsafe.Pointer
		switch 1 << v21 {
		case 4:
			mod = legacy.Get_dword_5d4594_1308156()
		case 1024:
			mod = legacy.Get_dword_5d4594_1308160()
		case 1:
			mod = legacy.Get_dword_5d4594_1308164()
		default:
			continue // LABEL_27
		}
		v27 := 1
		for i := 3; i < 21; i += 3 {
			r.Data().SetMaterialRGB(v27, int(*(*uint8)(unsafe.Add(mod, i+12))), int(*(*uint8)(unsafe.Add(mod, i+13))), int(*(*uint8)(unsafe.Add(mod, i+14))))
			v27++
		}
		switch 1 << v21 {
		case 4:
			selColorMaterial(int(*(*uint32)(unsafe.Add(mod, 40))), selColorRGB(legacy.Get_dword_5d4594_1308116()))
		case 1024:
			selColorMaterial(int(*(*uint32)(unsafe.Add(mod, 40))), selColorRGB(legacy.Get_dword_5d4594_1308120()))
			selColorMaterial(int(*(*uint32)(unsafe.Add(mod, 44))), selColorRGB(legacy.Get_dword_5d4594_1308124()))
		case 1:
			selColorMaterial(int(*(*uint32)(unsafe.Add(mod, 40))), selColorRGB(legacy.Get_dword_5d4594_1308128()))
			selColorMaterial(int(*(*uint32)(unsafe.Add(mod, 36))), selColorRGB(legacy.Get_dword_5d4594_1308132()))
		}
		r.DrawImageAt(selColorBlobImage(imgs+4*uintptr(v21)), pos)
	}
	return 1
}

// selColorPopupProc is the window proc of the color picker popup (sub_4A7270).
func selColorPopupProc(win *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	m, ok := ev.(*gui.WindowMouseState)
	if !ok || m.State != input.NOX_MOUSE_LEFT_DOWN {
		return nil
	}
	if win.PointIn(m.Pos) {
		return nil
	}
	sub_4A72D0(0xDEAD)
	return gui.RawEventResp(1)
}

// sub_4A72D0 hides the color picker popup and applies the selected color.
func sub_4A72D0(a1 uint16) {
	popup := legacy.Get_dword_5d4594_1308088()
	popup.StackPop()
	popup.SetHidden(true)
	if a1 < 0x20 {
		res := popup.ChildByID(uint(legacy.Get_dword_5d4594_1307792()))
		if res != nil {
			*selColorWidget(res) = uint32(*memmap.PtrUint16(0x5D4594, 1307788)) | uint32(a1)<<16
		}
	}
}

func selColorProc(win *gui.Window, ev gui.WindowEvent) gui.WindowEventResp { // sub_4A7330
	switch ev := ev.(type) {
	case *WindowEvent0x4005:
		clientPlaySoundSpecial(920, 100)
		return gui.RawEventResp(1)
	case *WindowEvent0x4007:
		a4 := uint32(ev.Val)
		v5 := int(a4 >> 16)
		v6 := -2
		if ev.Win != nil {
			v6 = int(ev.Win.ID())
		}
		switch {
		case v6 == 720:
			legacy.Set_dword_5d4594_1307792(uint32(v6))
			sub_4A7530(2)
			sub_4A7580(int(uint16(a4)), v5)
		case v6 >= 721 && v6 <= 724:
			legacy.Set_dword_5d4594_1307792(uint32(v6))
			sub_4A7530(1)
			sub_4A7580(int(uint16(a4)), v5)
		case v6 >= 725 && v6 <= 729:
			legacy.Set_dword_5d4594_1307792(uint32(v6))
			sub_4A7530(0)
			sub_4A7580(int(uint16(a4)), v5)
		case v6 >= 731 && v6 <= 734:
			v7 := legacy.Get_dword_5d4594_1308084().ChildByID(uint(v6 - 20))
			if v7 != nil {
				legacy.Nox_xxx_wnd_46ABB0(v7, int((^uint32(v7.Flags)>>3)&1))
			}
		case v6 >= 761 && v6 <= 792:
			sub_4A72D0(uint16(v6 - 761))
		case v6 == 799:
			if *memmap.PtrUint32(0x5D4594, 1308168) == 1 {
				noxClient.GamePopState()
			}
			noxClient.GamePopState()
			noxClient.GamePopState()
			legacy.Set_dword_587000_171388(1)
			if sub_4A75C0() != 0 {
				switch legacy.Get_dword_5d4594_1307784().PlayerClass() {
				case 0:
					legacy.Nox_xxx_gameSetMapPath_409D70("war01a.map")
				case 1:
					legacy.Nox_xxx_gameSetMapPath_409D70("wiz01a.map")
				case 2:
					legacy.Nox_xxx_gameSetMapPath_409D70("con01a.map")
				}
				sub_4A24C0(0)
				sub_4A6890()
				legacy.Get_nox_wnd_xxx_1308092().Func13Ptr = nil
			}
		}
		clientPlaySoundSpecial(921, 100)
		return gui.RawEventResp(1)
	}
	return nil
}

// sub_4A7530 sets the color table of all swatches of the popup.
func sub_4A7530(a1 uint16) {
	*memmap.PtrUint32(0x5D4594, 1307788) = uint32(a1)
	popup := legacy.Get_dword_5d4594_1308088()
	for v1 := 761; v1 <= 792; v1++ {
		res := popup.ChildByID(uint(v1))
		if res != nil {
			*selColorWidget(res) = uint32(a1) | uint32(uint16(v1-761))<<16
		}
	}
}

// sub_4A7580 shows the color picker popup at the position.
func sub_4A7580(a1, a2 int) {
	popup := legacy.Get_dword_5d4594_1308088()
	popup.ShowModal()
	popup.StackPush()
	popup.SetPos(image.Pt(a1-popup.SizeVal.X, a2-int(uint32(popup.SizeVal.Y)/2)))
}

func cstrLen(b []byte) int {
	for i, c := range b {
		if c == 0 {
			return i
		}
	}
	return len(b)
}

// sub_4A75C0 creates the player file for the selected character.
func sub_4A75C0() int {
	if noxflags.HasGame(noxflags.GameModeCoop) {
		deleteSaveDir("WORKING", false)
	}
	var sg server.SaveGameInfo // original: zeroed 0x4FC bytes + word at 1276
	v0 := selColorEntryText()
	if *v0 == 0 {
		wcsCopy(v0, selColorDefaultName())
	}
	name := &sg.Player.NameBuf[0]
	wcsCopy(name, v0)
	if sub_4A6B50(name) == 0 {
		wcsCopy(name, selColorDefaultName())
	}
	host := legacy.Get_dword_5d4594_1307784()
	sg.Field1276 = 1
	sg.Player.PlayerClassVal = byte(host.PlayerClass())
	sg.Player.SkinColor = selColorRGB(legacy.Get_dword_5d4594_1308096())
	skin := sg.Player.SkinColor
	if legacy.Get_dword_5d4594_1308136().Flags.Has(gui.StatusEnabled) {
		sg.Player.HairColor = selColorRGB(legacy.Get_dword_5d4594_1308100())
	} else {
		sg.Player.HairColor = skin
	}
	if legacy.Get_dword_5d4594_1308140().Flags.Has(gui.StatusEnabled) {
		sg.Player.MustacheColor = selColorRGB(legacy.Get_dword_5d4594_1308104())
	} else {
		sg.Player.MustacheColor = skin
	}
	if legacy.Get_dword_5d4594_1308144().Flags.Has(gui.StatusEnabled) {
		sg.Player.GoateeColor = selColorRGB(legacy.Get_dword_5d4594_1308108())
	} else {
		sg.Player.GoateeColor = skin
	}
	if legacy.Get_dword_5d4594_1308148().Flags.Has(gui.StatusEnabled) {
		sg.Player.BeardColor = selColorRGB(legacy.Get_dword_5d4594_1308112())
	} else {
		sg.Player.BeardColor = skin
	}
	sg.Player.PantsColor = byte(*selColorWidget(legacy.Get_dword_5d4594_1308116()) >> 16)
	sg.Player.Shirt1Color = byte(*selColorWidget(legacy.Get_dword_5d4594_1308120()) >> 16)
	sg.Player.Shirt2Color = byte(*selColorWidget(legacy.Get_dword_5d4594_1308124()) >> 16)
	sg.Player.Shoes1Color = byte(*selColorWidget(legacy.Get_dword_5d4594_1308128()) >> 16)
	sg.Player.Shoes2Color = byte(*selColorWidget(legacy.Get_dword_5d4594_1308132()) >> 16)

	path := sg.PathBuf[:]
	root := datapath.Data()
	copy(path, root)
	path[len(root)] = 0
	l := cstrLen(path)
	copy(path[l:], memmap.Slice(0x587000, 171764)[:4])
	copy(path[l+4:], memmap.Slice(0x587000, 171768)[:2])
	path[l+6] = memmap.Slice(0x587000, 171770)[0]
	if noxflags.HasGame(noxflags.GameModeCoop) {
		l = cstrLen(path)
		copy(path[l:], memmap.Slice(0x587000, 171772)[:8])
		l = cstrLen(path)
		copy(path[l:], memmap.Slice(0x587000, 171780)[:2])
	}
	ifs.Mkdir(string(path[:cstrLen(path)]))
	ifs.Chdir(string(path[:cstrLen(path)]))
	i := 0
	var v24 string
	if noxflags.HasGame(noxflags.GameModeCoop) {
		v24 = "Player.plr"
	} else {
		// The original passes the UTF-16 name to "%.6s", which stops at the first zero byte,
		// so only the low byte of the first character ends up in the file name.
		nb := unsafe.Slice((*byte)(unsafe.Pointer(name)), 2*len(sg.Player.NameBuf))
		pfx := nb[:cstrLen(nb)]
		if len(pfx) > 6 {
			pfx = pfx[:6]
		}
		for i = 0; i < 100; i++ {
			v24 = fmt.Sprintf("%s%2.2d.plr", pfx, i)
			f, err := ifs.Open(v24)
			if err != nil {
				break
			}
			f.Close()
		}
	}
	ifs.Chdir(root)
	if i > 99 {
		return 0
	}
	l = cstrLen(path)
	copy(path[l:], v24)
	path[l+len(v24)] = 0
	if noxflags.HasGame(noxflags.GameModeCoop) {
		switch sg.Player.PlayerClassVal {
		case 0:
			legacy.Nox_xxx_gameSetMapPath_409D70("war01a.map")
		case 1:
			legacy.Nox_xxx_gameSetMapPath_409D70("wiz01a.map")
		case 2:
			legacy.Nox_xxx_gameSetMapPath_409D70("con01a.map")
		}
	}
	return sub_41CEE0(&sg, 1)
}

// sub_41CEE0 writes the player file for the given save game info.
func sub_41CEE0(a1p *server.SaveGameInfo, a2 int) int {
	a1 := int(uintptr(unsafe.Pointer(a1p))) // decompiler artefact: pointer reused as int
	// sizeof(nox_savegame_xxx) in C is 1278 (the Go struct is 1280).
	copy(unsafe.Slice((*byte)(unsafe.Pointer(clientCurSave())), 1278), unsafe.Slice((*byte)(unsafe.Pointer(a1p)), 1278))
	cur := clientCurSave()
	if err := cryptfile.OpenGlobal(cur.Path(), cryptfile.WriteOnly, 27); err != nil {
		if !os.IsNotExist(err) {
			binfile.Log.Println(err)
		}
		return 0
	}
	v3 := a2
	if a2 == 0 {
		a1 = 1
	}
	cf := cryptfile.Global()
	if *memmap.PtrUint32(0x587000, 55936) != 0 {
		off := uintptr(55940)
		for {
			id := memmap.PtrUint32(0x587000, off)
			if v3 != 0 || int(*id) == a1 {
				cf.ReadWrite(unsafe.Slice((*byte)(unsafe.Pointer(id)), 4))
				cf.SectionStart()
				fnc := *memmap.PtrPtr(0x587000, off+4)
				v6 := ccall.CallIntInt(fnc, 0)
				cf.SectionEnd()
				if v6 == 0 {
					cryptfile.Close()
					return 0
				}
			}
			v8 := *memmap.PtrUint32(0x587000, off+8)
			off += 12
			if v8 == 0 {
				break
			}
		}
	}
	var zero [4]byte
	cf.ReadWrite(zero[:])
	cryptfile.Close()
	return 1
}

// sub_4BFAD0 loads images used by the character preview.
func sub_4BFAD0() int {
	load := func(strOff uintptr, dst uintptr) {
		name := *memmap.PtrPtr(0x587000, strOff)
		img := nox_xxx_gLoadImg(alloc.GoString((*byte)(name)))
		*memmap.PtrPtr(0x973A20, dst) = unsafe.Pointer(img.C())
	}
	v0, v1 := uintptr(0), uintptr(0)
	for i := uintptr(0); i < 8; i += 4 {
		load(180960+i, 16+i)
		load(180968+i, 24+i)
		for v5 := 26; v5 != 0; v5-- {
			load(180976+v1, 32+v1)
			v1 += 4
		}
		for v6 := 27; v6 != 0; v6-- {
			load(181184+v0, 256+v0)
			v0 += 4
		}
	}
	*memmap.PtrPtr(0x5D4594, 1319052) = unsafe.Pointer(nox_xxx_gLoadImg("MaleMedievalCloakTop").C())
	return 1
}
