package opennox

import (
	"unsafe"

	"github.com/opennox/libs/player"
	"github.com/opennox/libs/strman"

	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/sound"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/ccall"
	"github.com/opennox/opennox/v1/server"
)

func init() {
	legacy.Nox_game_showSelClass_4A4840 = func() int {
		if nox_game_showSelClass_4A4840() {
			return 1
		}
		return 0
	}
	legacy.Sub_4A4970 = selClassAnimOutStart
	legacy.Sub_4A49A0 = selClassAnimOutDone
	legacy.Sub_4AA650 = sub_4AA650
}

var (
	winSelClass       *gui.Window        // dword_5d4594_1307736
	winSelClassAnim   *gui.Anim          // nox_wnd_xxx_1307732
	winSelClassAccept *gui.Window        // 0x5D4594+1307728
	selClassInfo      *server.PlayerInfo // dword_5d4594_1307724
	selClassCurrent   int                // 0x5D4594+1307740
)

func nox_game_showSelClass_4A4840() bool {
	c := noxClient
	legacy.Sub_5007E0("*:*")
	sub_4A1BE0(1)
	selClassInfo = legacy.Nox_xxx_getHostInfoPtr_431770()
	c.GameAddStateCode(client.StateClassSelect)
	win := newWindowFromFile(c.GUI, "SelClass.wnd", selClassProc)
	winSelClass = win
	if win == nil {
		return false
	}
	win.SetFunc93(gui.WrapFuncC(legacy.Get_sub_4A18E0()))
	anim := nox_gui_makeAnimation(win, 0, 0, 0, -460, 0, 20, 0, -40)
	winSelClassAnim = anim
	if anim == nil {
		return false
	}
	anim.StateID = client.StateClassSelect
	anim.Func12Ptr = legacy.Get_sub_4A4970()
	anim.FncDoneOutPtr = legacy.Get_sub_4A49A0()
	for _, id := range []uint{601, 603, 602} {
		win.ChildByID(id).SetDraw(selClassButtonDraw)
	}
	winSelClassAccept = win.ChildByID(610)
	legacy.Nox_xxx_wndRetNULL_46A8A0()
	selClassCurrent = 0
	guiSetBackButtonText("OptsBack.wnd:Back")
	legacy.Sub_4602F0()
	return true
}

func selClassAnimOutStart() int { // sub_4A4970
	winSelClassAnim.SetState(gui.AnimOut)
	sub_43BE40(2)
	clientPlaySoundSpecial(sound.SoundShellSlideOut, 100)
	return 1
}

func selClassAnimOutDone() int { // sub_4A49A0
	next := winSelClassAnim.Func13Ptr
	winSelClassAnim.Free()
	winSelClassAnim = nil
	winSelClass.Destroy()
	winSelClass = nil
	ccall.CallIntVoid(next)
	return 1
}

// selClassButtonDraw darkens all class buttons except the selected one.
func selClassButtonDraw(win *gui.Window, _ *gui.WindowData) int { // sub_4A49D0
	if selClassCurrent != int(win.ID()) {
		pos := win.GlobalPos()
		size := win.EndPos.Sub(win.Off)
		noxClient.r.DrawRectFilledAlpha(pos.X, pos.Y, size.X, size.Y)
	}
	return 1
}

func selClassProc(win *gui.Window, ev gui.WindowEvent) gui.WindowEventResp { // sub_4A4A20
	c := noxClient
	switch ev := ev.(type) {
	case *WindowEvent0x4005:
		id := int(ev.Win.ID())
		if id >= 601 && id <= 603 {
			legacy.Nox_xxx_wnd_46ABB0(winSelClassAccept, 1)
			desc := winSelClass.ChildByID(605)
			cl := byte(id - 89)
			selClassInfo.SetPlayerClass(player.Class(cl))
			name := alloc.GoString(*(**byte)(memmap.PtrOff(0x587000, 170208+4*uintptr(cl))))
			str := c.Strings().GetStringInFile(strman.ID(name), "C:\\NoxPost\\src\\client\\shell\\SelClass.c")
			desc.Func94(&gui.StaticTextSetText{Str: str})
			selClassCurrent = id
		}
		clientPlaySoundSpecial(sound.ID(920), 100)
		return gui.RawEventResp(1)
	case *WindowEvent0x4007:
		id := int(ev.Win.ID())
		if id >= 601 {
			if id <= 603 {
				return gui.RawEventResp(1)
			}
			if id == 610 {
				if noxflags.HasGame(noxflags.GameOnline) && !noxflags.HasGame(noxflags.GameModeQuest) {
					sub_4A4B70(noxServer.nox_xxx_isQuest_4D6F50() || sub_4D6F70())
				}
				selClassAnimOutStart()
				winSelClassAnim.Func13Ptr = legacy.Get_nox_game_showSelColor_4A5D00()
			}
		}
		clientPlaySoundSpecial(sound.ID(921), 100)
		return gui.RawEventResp(1)
	}
	return nil
}

// sub_4A4B70 sets a random set of spells for the selected class.
func sub_4A4B70(empty bool) {
	cl := uintptr(byte(selClassInfo.PlayerClass()))
	sets := *(*unsafe.Pointer)(memmap.PtrOff(0x587000, 170156+4*cl))
	set := func(i, k int) byte { return *(*byte)(unsafe.Add(sets, 5*i+k)) }
	n := 0
	for set(n, 0) != 0 {
		n++
	}
	if n == 0 {
		return
	}
	ind := byte(noxServer.Rand.Other.Int(0, n-1))
	spell := func(k int) int {
		if empty {
			return 0
		}
		return int(set(int(ind), k))
	}
	if cl != 0 {
		for row := 0; row < 5; row++ {
			legacy.Nox_xxx_clientUpdateButtonRow_45E110(row)
			for k := 0; k < 5; k++ {
				legacy.Nox_xxx_book_45DBE0(2, spell(k), k)
			}
		}
		legacy.Nox_xxx_clientUpdateButtonRow_45E110(0)
	} else {
		legacy.Nox_xxx_clientUpdateButtonRow_45E110(0)
		for k := 0; k < 5; k++ {
			legacy.Nox_xxx_book_45DBE0(3, spell(k), k)
		}
	}
}

// sub_4AA650 plays one of the narrator lines on the options screen.
func sub_4AA650() {
	if legacy.Dialogs.Sub_44D930() {
		return
	}
	cnt := memmap.PtrUint32(0x5D4594, 1309744)
	name := alloc.GoString(*(**byte)(memmap.PtrOff(0x587000, 172892+4*uintptr(*cnt))))
	*cnt++
	v, _ := noxServer.Strings().GetVariantInFile(strman.ID(name), "C:\\NoxPost\\src\\client\\shell\\Options.c")
	*cnt %= 3
	// The original code checks the interned string pointer, not its contents.
	if alloc.InternCString(v.Str2) != nil {
		legacy.Dialogs.PlayFile(v.Str2, 100)
	}
}
