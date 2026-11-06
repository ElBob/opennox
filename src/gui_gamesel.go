package opennox

import (
	"fmt"
	"image"
	"unsafe"

	"github.com/opennox/libs/noxnet"
	"github.com/opennox/libs/strman"

	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/client/noxrender"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/sound"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

// Multiplayer game select screen (server list, filter window): nox_game_showGameSel_4379F0 and friends.
// Windows and state stay in the original C globals (reached via legacy getters/setters) and in the 0x5D4594 blob.

const gameSelSrcFile = "C:\\NoxPost\\src\\client\\shell\\noxworld.c"

// The decompiled Window is accessed as a 404-byte blob (memcpy in sub_438C80).
var _ = [1]struct{}{}[404-unsafe.Sizeof(gui.Window{})]

func init() {
	legacy.Nox_game_showGameSel_4379F0 = func() int {
		if nox_game_showGameSel_4379F0() {
			return 1
		}
		return 0
	}
}

// gameSelTicks is nox_platform_get_ticks (returns a 32-bit unsigned int).
func gameSelTicks() uint64 { return uint64(uint32(platformTicks())) }

// wndIsShown is wndIsShown_nox_xxx_wndIsShown_46ACC0. Quirk: despite the name it returns 1 for a hidden (or NULL) window.
func wndIsShown(win *gui.Window) bool {
	if win == nil {
		return true
	}
	return win.GetFlags().IsHidden()
}

// gameSelRootProcFnc is nox_xxx_windowMultiplayerSub_439E70 (still in C).
func gameSelRootProcFnc() unsafe.Pointer { return legacy.Get_nox_xxx_windowMultiplayerSub_439E70() }

func gameSelCall(win *gui.Window, code int, a1, a2 uintptr) gui.WindowEventResp {
	return win.Func94(gui.AsWindowEvent(code, a1, a2))
}

func gameSelPtr(win *gui.Window) uintptr { return uintptr(unsafe.Pointer(win)) }

func nox_game_showGameSel_4379F0() bool {
	c := noxClient
	noxServer.nox_xxx_setQuest_4D6F60(0)
	legacy.Sub_4D6F80(0)
	if !noxflags.HasGame(noxflags.GameFlag26) && !sub4D6F30() {
		legacy.Set_nox_game_createOrJoin_815048(0)
		legacy.Set_dword_587000_87412(-1)
	}
	*memmap.PtrUint32(0x5D4594, 815076) = 0
	*memmap.PtrUint32(0x5D4594, 815080) = 0
	*memmap.PtrUint32(0x5D4594, 815084) = 0
	c.GameAddStateCode(client.StateServerList)
	sub_4A24C0(1)
	sub_4A1BE0(0)
	if legacy.Get_dword_587000_87412() != -1 {
		legacy.Sub_49FDB0(0)
	}
	if world := legacy.Get_nox_wol_wnd_world_814980(); world != nil {
		legacy.Set_dword_5d4594_815044(0)
		anim := legacy.Get_nox_wnd_xxx_815040()
		anim.SetState(gui.AnimIn)
		anim.FncDoneOutPtr = legacy.Get_sub_438330()
		sub_43BE40(3)
		clientPlaySoundSpecial(sound.ID(922), 100)
		world.SetHidden(false)
		legacy.Get_dword_5d4594_815000().SetHidden(false)
		legacy.Nox_xxx_wnd_46ABB0(legacy.Get_dword_5d4594_814984(), 1)
		if legacy.Get_nox_game_createOrJoin_815048() != 0 {
			legacy.Sub_4375C0(0)
		}
		legacy.Sub_46ACE0(world, 10047, 10051, bool2int(legacy.Get_dword_587000_87408() == 0))
		if sub_4D6FA0() == 2 {
			v25 := world.ChildByID(10010)
			gameSelCall(world, 16391, gameSelPtr(v25), 0)
			v26 := world.ChildByID(10010)
			gameSelCall(world, 16391, gameSelPtr(v26), 0)
		}
		return true
	}
	world := nox_new_window_from_file("noxworld.wnd", gui.WrapFuncC(gameSelRootProcFnc()))
	legacy.Set_nox_wol_wnd_world_814980(world)
	if world == nil {
		return false
	}
	sub_49FF20()
	world.ShowModal()
	world.SetAllFuncs(gameSelRootProc, nil, nil)
	anim := nox_gui_makeAnimation(world, 0, 0, 0, -480, 0, 20, 0, -40)
	legacy.Set_nox_wnd_xxx_815040(anim)
	if anim == nil {
		return false
	}
	anim.StateID = 10000
	anim.Func12Ptr = legacy.Get_sub_438370()
	anim.FncDoneOutPtr = legacy.Get_sub_438330()
	w814984 := world.ChildByID(10020)
	legacy.Set_dword_5d4594_814984(w814984)
	w814988 := world.ChildByID(10021)
	legacy.Set_dword_5d4594_814988(w814988)
	w814984.SetFunc93(gameSelRootProc)
	w814988.SetFunc93(gameSelRootProc)
	w814988.SetFunc94(gui.WrapFuncC(gameSelRootProcFnc()))
	w814996 := world.ChildByID(10011)
	legacy.Set_dword_5d4594_814996(w814996)
	gameList := world.ChildByID(10037)
	legacy.Set_nox_wol_wnd_gameList_815012(gameList)
	w815004 := world.ChildByID(10034)
	legacy.Set_dword_5d4594_815004(w815004)
	w815000 := world.ChildByID(10033)
	legacy.Set_dword_5d4594_815000(w815000)
	legacy.Set_dword_5d4594_814992(sub_489B80(world))
	*memmap.PtrT[*gui.Window](0x5D4594, 815008) = world.ChildByID(10001)
	legacy.SetGameSelColumn(0, world.ChildByID(10038))
	legacy.SetGameSelColumn(1, world.ChildByID(10039))
	legacy.SetGameSelColumn(2, world.ChildByID(10040))
	legacy.SetGameSelColumn(3, world.ChildByID(10041))
	legacy.SetGameSelColumn(4, world.ChildByID(10042))
	v2 := world.ChildByID(10002)
	if noxflags.HasGame(noxflags.GameFlag25) {
		legacy.Nox_xxx_wnd_46ABB0(v2, 0)
	}
	v3 := world.ChildByID(10003)
	v3.DrawData().Field0 &= 0xFFFFFFFB
	v4 := world.ChildByID(10046)
	v4.SetDraw(sub_438C80)
	gameList.SetFunc94(gui.WrapFuncC(gameSelRootProcFnc()))
	gameList.SetHidden(true)
	w815000.SetHidden(true)
	legacy.Nox_xxx_wnd_46B280(*memmap.PtrT[*gui.Window](0x5D4594, 815008), world)

	legacy.Set_dword_587000_87412(0)
	w814988.SetHidden(true)
	v5 := c.Strings().GetStringInFile("JoinServer", gameSelSrcFile)
	w814996.Func94(&gui.StaticTextSetText{Str: v5})
	legacy.Sub_46ACE0(w814984, 10620, 10631, 1)
	legacy.Sub_49FDB0(0)

	*memmap.PtrT[noxrender.ImageHandle](0x5D4594, 814556) = nox_xxx_gLoadImg("NWGameIconLargeGreen").C()
	*memmap.PtrT[noxrender.ImageHandle](0x5D4594, 814560) = nox_xxx_gLoadImg("NWGameIconLargeGreenLit").C()
	*memmap.PtrT[noxrender.ImageHandle](0x5D4594, 814564) = nox_xxx_gLoadImg("NWGameIconSmallGreen").C()
	*memmap.PtrT[noxrender.ImageHandle](0x5D4594, 814568) = nox_xxx_gLoadImg("NWGameIconSmallGreenLit").C()
	*memmap.PtrT[noxrender.ImageHandle](0x5D4594, 814572) = nox_xxx_gLoadImg("NWGameIconLargeYellow").C()
	*memmap.PtrT[noxrender.ImageHandle](0x5D4594, 814576) = nox_xxx_gLoadImg("NWGameIconLargeYellowLit").C()
	*memmap.PtrT[noxrender.ImageHandle](0x5D4594, 814580) = nox_xxx_gLoadImg("NWGameIconSmallYellow").C()
	*memmap.PtrT[noxrender.ImageHandle](0x5D4594, 814584) = nox_xxx_gLoadImg("NWGameIconSmallYellowLit").C()
	*memmap.PtrT[noxrender.ImageHandle](0x5D4594, 814588) = nox_xxx_gLoadImg("NWGameIconLargeRed").C()
	*memmap.PtrT[noxrender.ImageHandle](0x5D4594, 814592) = nox_xxx_gLoadImg("NWGameIconLargeRedLit").C()
	*memmap.PtrT[noxrender.ImageHandle](0x5D4594, 814596) = nox_xxx_gLoadImg("NWGameIconSmallRed").C()
	*memmap.PtrT[noxrender.ImageHandle](0x5D4594, 814600) = nox_xxx_gLoadImg("NWGameIconSmallRedLit").C()
	*memmap.PtrT[noxrender.ImageHandle](0x5D4594, 814900) = nox_xxx_gLoadImg("NWMapULLg").C()
	*memmap.PtrT[noxrender.ImageHandle](0x5D4594, 814904) = nox_xxx_gLoadImg("NWMapURLg").C()
	*memmap.PtrT[noxrender.ImageHandle](0x5D4594, 814908) = nox_xxx_gLoadImg("NWMapLLLg").C()
	*memmap.PtrT[noxrender.ImageHandle](0x5D4594, 814912) = nox_xxx_gLoadImg("NWMapLRLg").C()
	v8 := gameList.ChildByID(10053)
	v9 := gameList.ChildByID(10043)
	v10 := gameList.ChildByID(10044)
	v11 := (*gui.ScrollListBoxData)(gameList.WidgetData)
	legacy.Nox_xxx_wnd_46B280(v8, gameList)
	legacy.Nox_xxx_wnd_46B280(v9, gameList)
	legacy.Nox_xxx_wnd_46B280(v10, gameList)
	v11.Field_9 = unsafe.Pointer(v8)
	v11.Field_7 = unsafe.Pointer(v9)
	v11.Field_8 = unsafe.Pointer(v10)
	v8.Field100Ptr.SizeVal.X = 16
	v8.Field100Ptr.SizeVal.Y = 12
	v8.Field100Ptr.DrawData().SetImagePoint(image.Pt(0, -15)) // nox_xxx_wndSetOffsetMB_46AE40
	v12 := w815000.ChildByID(10032)
	v13 := w815000.ChildByID(10035)
	v14 := w815000.ChildByID(10036)
	v15 := (*gui.ScrollListBoxData)(w815004.WidgetData)
	legacy.Nox_xxx_wnd_46B280(v12, w815004)
	legacy.Nox_xxx_wnd_46B280(v13, w815004)
	legacy.Nox_xxx_wnd_46B280(v14, w815004)
	v15.Field_9 = unsafe.Pointer(v12)
	v15.Field_7 = unsafe.Pointer(v13)
	v15.Field_8 = unsafe.Pointer(v14)
	v12.Field100Ptr.SizeVal.X = 16
	v12.Field100Ptr.SizeVal.Y = 10
	sub_438480()
	v16 := noxServer.ServerPort()
	legacy.Sub_40A3C0(v16)
	legacy.Set_nox_game_createOrJoin_815048(0)
	legacy.Set_dword_5d4594_815044(0)
	legacy.Set_dword_5d4594_815052(0)
	legacy.Sub_43AF90(0) // dword_5d4594_814548
	legacy.Set_qword_5d4594_815068(uint64(uint32(platformTicks()) + 1000))
	if legacy.Get_dword_587000_87408() == 1 {
		legacy.Sub_4383A0()
	}
	_ = nox_xxx_createSocketLocal(0)
	legacy.Set_nox_wol_server_result_cnt_815088(0)
	nox_xxx_loadModifyers_4158C0()
	nox_xxx_loadLook_415D50()
	sub_430C30_set_video_max(noxMaxWidth, noxMaxHeight)
	c.Nox_client_setCursorType(gui.CursorSelect)
	if legacy.Get_dword_5d4594_815096() != 0 {
		v27 := c.Strings().GetStringInFile("Kicked", gameSelSrcFile)
		v17 := c.Strings().GetStringInFile("Notification", gameSelSrcFile)
		NewDialogWindow(nil, v17, v27, 33, nil, nil)
		sub_44A360(1)
		legacy.Set_dword_5d4594_815096(0)
	} else if legacy.Get_dword_5d4594_815100() != 0 {
		v28 := c.Strings().GetStringInFile("Timeout", gameSelSrcFile)
		v18 := c.Strings().GetStringInFile("Notification", gameSelSrcFile)
		NewDialogWindow(nil, v18, v28, 33, nil, nil)
		sub_44A360(1)
		legacy.Set_dword_5d4594_815100(0)
	}
	noxServer.SetUpdateFunc2(sub_438770) // sub_43DE40
	if sub44A4A0() {
		sub_44A4B0()
	}
	w814984.SetDraw(sub_438E30)
	v19 := world.ChildByID(10054)
	v19.SetDraw(sub_438E30)
	v20 := world.ChildByID(10055)
	v20.SetDraw(sub_438E30)
	v21 := world.ChildByID(10056)
	v21.SetDraw(sub_438E30)
	v22 := world.ChildByID(10057)
	v22.SetDraw(sub_438E30)
	if legacy.Get_dword_587000_87408() == 0 {
		legacy.Sub_46ACE0(world, 10047, 10051, 1)
	}
	if sub_4D6FA0() == 1 {
		v23 := world.ChildByID(10002)
		gameSelCall(world, 16391, gameSelPtr(v23), 0)
		noxServer.nox_xxx_setQuest_4D6F60(1)
		v24 := world.ChildByID(10020)
		v24.Func93(gui.AsWindowEvent(5, 15663512, 0))
	}
	return true
}

func sub_438770() bool {
	switch legacy.Sub_43AF80() { // dword_5d4594_814548
	case 0:
		if (!wndIsShown(legacy.Get_dword_5d4594_814984()) ||
			!wndIsShown(legacy.Get_dword_5d4594_814988()) ||
			!wndIsShown(legacy.Get_nox_wol_wnd_gameList_815012())) &&
			legacy.Get_nox_game_createOrJoin_815048() == 0 && legacy.Get_dword_5d4594_815044() == 0 &&
			legacy.Get_dword_5d4594_815052() == 0 &&
			wndIsShown(legacy.Get_dword_5d4594_815000()) {
			if gameSelTicks() > legacy.Get_qword_5d4594_815068() {
				nox_client_refreshServerList_4378B0()
				return true
			}
			sub_438770_waitList()
		}
		return true
	case 2:
		sub_438BD0()
		legacy.Sub_43AF90(1)
		return true
	case 3:
		if gameSelTicks() < legacy.Get_qword_5d4594_814956() {
			return true
		}
		nox_client_setConnError_43AFA0(noxnet.ConnectError(8))
		return true
	case 4:
		legacy.Sub_43AF90(3)
		v2 := noxClient.Strings().GetStringInFile("TestCon", gameSelSrcFile)
		sub449E30(v2)
		legacy.Set_qword_5d4594_814956(uint64(uint32(platformTicks()) + 20000))
		return true
	case 5:
		v3 := noxClient.Strings().GetStringInFile("Password", gameSelSrcFile)
		sub449E00(v3)
		v4 := noxClient.Strings().GetStringInFile("PasswordRequired", gameSelSrcFile)
		sub449E30(v4)
		sub449EA0(gui.DialogFlags(7))
		sub_44A360(0)
		legacy.Sub_43AF90(6)
		sub_4A24C0(1)
		return true
	case 7:
		sub_44A360(1)
		v5 := noxClient.Strings().GetStringInFile("Connected", gameSelSrcFile)
		sub449E30(v5)
		sub449EA0(0)
		noxClient.SetDrawFunc(func() bool { // nox_game_SetCliDrawFunc
			return nox_xxx_cliDrawConnectedLoop_43B360() != 0
		})
		legacy.Sub_43AF90(1)
		return true
	case 8:
		v1 := gameSelTicks() // long long
		legacy.Sub_43AF90(9)
		*memmap.PtrUint64(0x5D4594, 814972) = v1 + 1000
		return true
	case 9:
		if gameSelTicks() <= *memmap.PtrUint64(0x5D4594, 814972) {
			return true
		}
		nox_client_joinGame_438A90()
		return true
	case 0xA:
		v6 := (*uint16)(gui.EventRespPtr(sub449E60(4)))
		v7 := (*gui.Window)(nil).ChildByID(4001) // nox_xxx_wndGetChildByID_46B0C0(0, 4001)
		if v6 != nil && *v6 != 0 {
			legacy.Nox_xxx_wnd_46ABB0(v7, 1)
		} else {
			legacy.Nox_xxx_wnd_46ABB0(v7, 0)
		}
		return true
	default:
		return true
	}
}

func sub_438BD0() int {
	if !sub44A4A0() {
		NewDialogWindow(legacy.Get_nox_wol_wnd_world_814980(), "", "", 0, nil, nil)
	}
	v0 := uint32(legacy.Get_nox_client_connError_814552())
	if v0 != 8 && v0 != 9 && v0 != 10 {
		v1 := noxClient.Strings().GetStringInFile("ConnError", gameSelSrcFile)
		sub449E00(v1)
		v0 = uint32(legacy.Get_nox_client_connError_814552())
	}
	name := alloc.GoString(*(**byte)(memmap.PtrOff(0x587000, 87416+4*uintptr(v0))))
	v2 := noxClient.Strings().GetStringInFile(strman.ID(name), gameSelSrcFile)
	sub449E30(v2)
	legacy.Set_dword_5d4594_815044(0)
	sub449EA0(1)
	sub_44A360(1)
	return sub_4A24C0(1)
}

func sub_438480() int {
	gl := legacy.Get_nox_wol_wnd_gameList_815012()
	gl.SetFunc94(gameSelListProc)
	gl.SetFunc93(sub_438EF0)
	cols := legacy.GameSelColumns()
	for _, col := range cols {
		col.SetParent(gl) // sub_46B120
		col.SetFunc93(sub_438EF0)
	}
	w815004 := legacy.Get_dword_5d4594_815004()
	data := (*gui.ScrollListBoxData)(w815004.WidgetData)
	*(*uint32)(data.Field_7) = 10035
	data = (*gui.ScrollListBoxData)(w815004.WidgetData)
	*(*uint32)(data.Field_8) = 10036
	data = (*gui.ScrollListBoxData)(w815004.WidgetData)
	*(*uint32)(data.Field_9) = 10032
	ldata := (*gui.ScrollListBoxData)(gl.WidgetData)
	for _, col := range cols {
		col.Func94(gui.AsWindowEvent(16408, uintptr(ldata.Field_7), 0))
	}
	for _, col := range cols {
		col.Func94(gui.AsWindowEvent(16409, uintptr(ldata.Field_8), 0))
	}
	for _, col := range cols[:4] {
		col.Func94(gui.AsWindowEvent(16410, uintptr(ldata.Field_9), 0))
	}
	return gui.EventRespInt(cols[4].Func94(gui.AsWindowEvent(16410, uintptr(ldata.Field_9), 0)))
}

func sub_4A28B0() bool { return legacy.Get_dword_5d4594_1307716() != 0 }

func sub_438C80(win *gui.Window, _ *gui.WindowData) int {
	mpos := noxClient.GetMousePos()
	w815000 := legacy.Get_dword_5d4594_815000()
	world := legacy.Get_nox_wol_wnd_world_814980()
	if !wndIsShown(w815000) {
		v2 := *w815000 // memcpy(v2, window, 404)
		v2.Off.X -= 32
		v2.Off.Y -= 32
		v2.SizeVal.X += 64
		v2.SizeVal.Y += 64
		if legacy.Get_dword_5d4594_815044() == 0 && !legacy.Nox_xxx_wndPointInWnd_46AAB0(unsafe.Pointer(&v2), mpos.X, mpos.Y) {
			w815000.SetHidden(true)
			gameSelCall(legacy.Get_nox_wol_wnd_gameList_815012(), 16403, ^uintptr(0), 0)
			legacy.Set_dword_5d4594_815056(0)
			w815000.StackPop()
			noxClient.GUI.Focus(world)
		}
	}
	if sub_4A28B0() && !legacy.Nox_xxx_wndPointInWnd_46AAB0(unsafe.Pointer(*memmap.PtrT[*gui.Window](0x5D4594, 815036)), mpos.X, mpos.Y) {
		legacy.Sub_4A2890()
		noxClient.GUI.Focus(world)
	}
	if legacy.Get_nox_game_createOrJoin_815048() != 0 && legacy.Sub_438DD0(mpos.X, mpos.Y) != 0 {
		noxClient.Nox_client_setCursorType(gui.CursorCreateGame)
	} else if !sub44A4A0() {
		noxClient.Nox_client_setCursorType(gui.CursorSelect)
	}
	return 1
}

func sub_438E30(win *gui.Window, _ *gui.WindowData) int {
	v6 := win.GlobalPos()
	d := win.DrawData()
	v2 := d.ImagePoint().Y
	if d.Field0&6 != 0 {
		noxClient.r.DrawImageAt(noxClient.r.GetBag().AsImage(d.HlImageHnd), image.Pt(v6.X+d.ImagePoint().X, v6.Y+v2))
	} else {
		noxClient.r.DrawImageAt(noxClient.r.GetBag().AsImage(d.BgImageHnd), image.Pt(v6.X+d.ImagePoint().X, v6.Y+v2))
	}
	for v3 := win.Field100Ptr; v3 != nil; v3 = v3.Prev() {
		if !v3.Flags.Has(gui.StatusHidden) && v3.DrawData().Style == gui.StyleStaticText {
			v4 := (*gui.StaticTextData)(v3.WidgetData)
			noxClient.r.Data().SetTextColor(nox_color_white_2523948)
			noxClient.r.DrawStringStyle(noxClient.r.GetFonts().AsFont(v3.DrawData().FontC()), alloc.GoString16(v4.Text),
				image.Pt(v6.X+v3.Off.X, v6.Y+v3.Off.Y))
		}
	}
	return 1
}

func sub_438EF0(win *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	a2 := ev.EventCode()
	var arg uintptr
	switch a2 {
	case 19:
		arg = uintptr((*gui.ScrollListBoxData)(win.WidgetData).Field_7)
	case 20:
		arg = uintptr((*gui.ScrollListBoxData)(win.WidgetData).Field_8)
	default:
		return legacy.Nox_xxx_wndListboxProcWithoutData10_4A28E0(win, ev)
	}
	gameSelCall(legacy.Get_nox_wol_wnd_gameList_815012(), 16391, arg, 0)
	for _, col := range legacy.GameSelColumns() {
		gameSelCall(col, 16391, arg, 0)
	}
	return nil
}

func sub_4A04C0(a1 int) unsafe.Pointer {
	v1 := 0
	result := nox_common_list_getFirstSafe_425890(legacy.Sub_4A0020())
	if result == nil {
		return nil
	}
	for a1 != v1 {
		v1++
		result = nox_common_list_getNextSafe_4258A0(result)
		if result == nil {
			return nil
		}
	}
	return result
}

func gameSelListProc(win *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	a2 := uint(ev.EventCode())
	a3, a4 := ev.EventArgsC()
	cols := legacy.GameSelColumns()
	if a2 > 0x400F {
		if a2 == 16400 {
			v7 := legacy.AsWindowP(unsafe.Pointer(a3)).ID()
			if v7 >= 10038 && v7 <= 10042 {
				for _, col := range cols {
					gameSelCall(col, 16403, a4, 0)
				}
				if uint32(a4) < uint32(legacy.Get_nox_wol_server_result_cnt_815088()) {
					pos := noxClient.GetMousePos()
					legacy.Set_dword_5d4594_814624(sub_4A04C0(int(a4)))
					legacy.Nox_client_gui_serverInfoBlockCheckExp_439370(pos, legacy.Get_dword_5d4594_814624())
				}
			}
		} else if a2 == 16403 || a2 == 16412 {
			for _, col := range cols {
				gameSelCall(col, int(a2), a3, 0)
			}
		}
	} else if a2 >= 0x400E {
		for _, col := range cols {
			gameSelCall(col, int(a2), a3, a4)
		}
	} else {
		switch a2 {
		case 0x17:
			return gui.RawEventResp(1)
		case 0x4000:
			gl := legacy.Get_nox_wol_wnd_gameList_815012()
			if a3 == gameSelPtr(gl.ChildByID(10043)) || a3 == gameSelPtr(gl.ChildByID(10044)) {
				for _, col := range cols {
					gameSelCall(col, 0x4000, a3, 0)
				}
			}
		case 0x4009:
			v4 := win.WidgetData
			legacy.Nox_xxx_wndListboxProcPre_4A30D0(win, ev)
			v5 := legacy.Sub_4A4800(v4)
			for _, col := range cols {
				gameSelCall(col, 16412, uintptr(v5), 0)
			}
		}
	}
	return legacy.Nox_xxx_wndListboxProcPre_4A30D0(win, ev)
}

func gameSelRootProc(win *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	a2 := ev.EventCode()
	a3, a4 := ev.EventArgsC()
	if a2 == 5 {
		if win.ID() == 10020 && legacy.Get_nox_game_createOrJoin_815048() == 1 {
			legacy.Sub_439D90(int(uint16(a3)), int(uint32(a3)>>16))
			return gui.RawEventResp(1)
		}
		return nil
	}
	if a2 != 21 {
		return nil
	}
	if a3 != 1 {
		if a3 != 28 && a3 == 57 {
			mpos := noxClient.GetMousePos()
			win.Func93(gui.AsWindowEvent(5, uintptr(mpos.X|(mpos.Y<<16)), 0))
		}
		return nil
	}
	if a4 == 2 {
		sub_4373A0()
	}
	return gui.RawEventResp(1)
}

func sub_4D6FA0() uint32 { return memmap.Uint32(0x5D4594, 1556104) }

func sub_4D70B0() unsafe.Pointer {
	result := *memmap.PtrT[unsafe.Pointer](0x5D4594, 1556152)
	if memmap.Uint32(0x5D4594, 1556152) == 0 {
		serverSetOpen() // sub_4169F0
		result = unsafe.Pointer(getServerSettings())
	}
	return result
}

// sub_489B80 creates the server list filter window.
func sub_489B80(a1 *gui.Window) *gui.Window {
	result := nox_new_window_from_file("filter.wnd", filterWindowProc)
	legacy.Set_dword_5d4594_1193380(result)
	if result != nil {
		legacy.Set_dword_5d4594_1193384(result.ChildByID(10012))
		v2 := uintptr(0)
		v3 := memmap.PtrOff(0x5D4594, 1193388+44*v2)
		win := legacy.Get_dword_5d4594_1193380()
		win.SetParent(a1)
		sub := legacy.Get_dword_5d4594_1193384()
		sub.SetParent(win)
		sub.SetFunc94(filterWindowProc)
		if *(*uint32)(v3) != 0 {
			v4 := win.ChildByID(10028)
			v4.DrawData().Field0 |= 4
		}
		v13, free13 := alloc.Make([]uint16(nil), 16) // wchar2_t v13[16], on the C heap
		for i, ch := range fmt.Sprintf("%d", int32(*(*uint32)(unsafe.Add(v3, 16)))) {
			v13[i] = uint16(ch)
		}
		v5 := win.ChildByID(10031)
		v5.Func94(gui.AsWindowEvent(16414, uintptr(unsafe.Pointer(&v13[0])), ^uintptr(0)))
		free13()
		if *(*uint32)(unsafe.Add(v3, 4)) != 0 {
			v6 := win.ChildByID(10029)
			v6.DrawData().Field0 |= 4
		}
		if *(*uint32)(unsafe.Add(v3, 8)) != 0 {
			v7 := win.ChildByID(10030)
			v7.DrawData().Field0 |= 4
		}
		if *(*uint32)(unsafe.Add(v3, 12)) != 0 {
			v8 := win.ChildByID(10015)
			v8.DrawData().Field0 |= 4
		}
		var v9 *gui.Window
		if *(*uint8)(unsafe.Add(v3, 12))&2 != 0 {
			v9 = win.ChildByID(10017)
		} else {
			v9 = win.ChildByID(10016)
		}
		v10 := v9
		v9.DrawData().Field0 |= 4
		if *(*uint32)(unsafe.Add(v3, 20)) != 0 {
			v10 = win.ChildByID(10014)
			v10.DrawData().Field0 |= 4
		}
		if *(*uint32)(unsafe.Add(v3, 40)) != 0 {
			v10 = win.ChildByID(10018)
			v10.DrawData().Field0 |= 4
		}
		v11 := *memmap.PtrUint32(0x5D4594, 1193372+4*v2)
		if v11 != 0 {
			v12 := v11 - 1
			if v12 != 0 {
				if v12 == 1 {
					v10 = win.ChildByID(10026)
					sub_489DC0()
				}
			} else {
				v10 = win.ChildByID(10025)
				sub.SetHidden(true)
			}
		} else {
			v10 = win.ChildByID(10024)
			sub.SetHidden(true)
		}
		v10.DrawData().Field0 |= 4
		result = legacy.Get_dword_5d4594_1193380()
	}
	return result
}

func sub_489DC0() {
	sub := legacy.Get_dword_5d4594_1193384()
	sub.SetHidden(false)
	if sub.ChildByID(10028).DrawData().Field0&4 != 0 {
		v0 := sub.ChildByID(10031)
		legacy.Nox_xxx_wnd_46ABB0(v0, 1)
	} else {
		v1 := sub.ChildByID(10031)
		legacy.Nox_xxx_wnd_46ABB0(v1, 0)
	}
	if sub.ChildByID(10015).DrawData().Field0&4 != 0 {
		legacy.Sub_46AD20(sub, 10016, 10017, 1)
	} else {
		legacy.Sub_46AD20(sub, 10016, 10017, 0)
	}
}

// filterWindowProc is nox_xxx_windowMplayFilterProc_489E70.
func filterWindowProc(win *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	a2 := ev.EventCode()
	a3, _ := ev.EventArgsC()
	v3 := uintptr(0)
	if a2 == 23 {
		return gui.RawEventResp(1)
	}
	if a2 != 16391 {
		return nil
	}
	a3win := legacy.AsWindowP(unsafe.Pointer(a3))
	v4 := a3win.ID()
	clientPlaySoundSpecial(sound.ID(766), 100)
	sub := legacy.Get_dword_5d4594_1193384()
	switch v4 {
	case 10015:
		legacy.Sub_46AD20(legacy.Get_dword_5d4594_1193380(), 10016, 10017, int((^a3win.DrawData().Field0>>2)&1))
		return nil
	case 10024:
		*memmap.PtrUint32(0x5D4594, 1193372+4*v3) = 0
		sub.SetHidden(true)
	case 10025:
		*memmap.PtrUint32(0x5D4594, 1193372+4*v3) = 1
		sub.SetHidden(true)
	case 10026:
		*memmap.PtrUint32(0x5D4594, 1193372+4*v3) = 2
		sub_489DC0()
	case 10028:
		v6 := legacy.Get_dword_5d4594_1193380().ChildByID(10031)
		legacy.Nox_xxx_wnd_46ABB0(v6, int((^a3win.DrawData().Field0>>2)&1))
	default:
		return nil
	}
	return nil
}

// sub_49FF20 resets the polygon data once, if it was initialized.
func sub_49FF20() unsafe.Pointer {
	result := legacy.Get_dword_5d4594_1305788()
	if result != 0 {
		r := sub_421B10()
		legacy.Set_dword_5d4594_1305788(0)
		return r
	}
	return nil
}

func sub_421B10() unsafe.Pointer {
	sub_421430()
	sub_421010()
	result := sub_420C70()
	*memmap.PtrUint32(0x5D4594, 588076) = 1
	*memmap.PtrUint16(0x5D4594, 588072) = 0
	return result
}

func sub_420C70() unsafe.Pointer {
	result := legacy.Get_dword_5d4594_588068()
	if result != nil {
		for {
			v1 := *(*unsafe.Pointer)(unsafe.Add(result, 8))
			legacy.CFree(result)
			result = v1
			if v1 == nil {
				break
			}
		}
		legacy.Set_dword_5d4594_588068(nil)
	} else {
		legacy.Set_dword_5d4594_588068(nil)
	}
	return result
}

func sub_421010() unsafe.Pointer {
	result := memmap.PtrOff(0x5D4594, 535872)
	for { // do-while
		*(*uint32)(result) = 0
		result = unsafe.Add(result, 16)
		if uintptr(result) >= uintptr(memmap.PtrOff(0x5D4594, 552240)) {
			break
		}
	}
	legacy.Set_nox_xxx_polygonNextAngle_587000_60356(1)
	return result
}

func sub_421430() unsafe.Pointer {
	v0 := memmap.PtrOff(0x5D4594, 552476)
	var result unsafe.Pointer
	for i := 0; i < 255; i++ {
		if p := unsafe.Add(v0, -27*4); *(*uint32)(p) != 0 {
			if memmap.Uint32(0x5D4594, 588076) != 0 {
				legacy.CFree(*(*unsafe.Pointer)(p))
			}
			*(*uint32)(p) = 0
		}
		result = *(*unsafe.Pointer)(v0)
		if *(*uint32)(v0) != 0 {
			if memmap.Uint32(0x5D4594, 588076) != 0 {
				legacy.CFree(*(*unsafe.Pointer)(v0))
			}
			*(*uint32)(v0) = 0
		}
		*(*uint16)(unsafe.Add(v0, 10*2)) = 0
		*(*uint32)(unsafe.Add(v0, -7*4)) = 0
		*(*int32)(unsafe.Add(v0, 2*4)) = -1
		*(*int32)(unsafe.Add(v0, 4*4)) = -1
		*(*uint32)(unsafe.Add(v0, -6*4)) = 0
		v0 = unsafe.Add(v0, 140)
	}
	legacy.Set_nox_xxx_polygonNextIdx_587000_60352(1)
	return result
}

func nox_xxx_loadLook_415D50() unsafe.Pointer {
	result := *memmap.PtrT[unsafe.Pointer](0x5D4594, 371256)
	if memmap.Uint32(0x5D4594, 371256) == 0 {
		result = *memmap.PtrT[unsafe.Pointer](0x587000, 35500)
		if memmap.Uint32(0x587000, 35500) != 0 {
			result = memmap.PtrOff(0x587000, 35500)
			v1 := memmap.PtrOff(0x587000, 35500)
			for {
				name := alloc.GoString(*(**byte)(result))
				str := noxServer.Strings().GetStringInFile(strman.ID(name), "C:\\NoxPost\\src\\common\\Object\\ArmrLook.c")
				*(*unsafe.Pointer)(unsafe.Add(v1, -4)) = unsafe.Pointer(alloc.InternCString16(str))
				v2 := *(*uint32)(unsafe.Add(v1, 12))
				v1 = unsafe.Add(v1, 12)
				result = v1
				if v2 == 0 {
					break
				}
			}
		}
		*memmap.PtrUint32(0x5D4594, 371256) = 1
	}
	return result
}

func nox_xxx_loadModifyers_4158C0() unsafe.Pointer {
	result := *memmap.PtrT[unsafe.Pointer](0x5D4594, 371248)
	if memmap.Uint32(0x5D4594, 371248) == 0 {
		result = *memmap.PtrT[unsafe.Pointer](0x587000, 33396)
		if memmap.Uint32(0x587000, 33396) != 0 {
			result = memmap.PtrOff(0x587000, 33396)
			v1 := memmap.PtrOff(0x587000, 33396)
			for {
				name := alloc.GoString(*(**byte)(result))
				str := noxServer.Strings().GetStringInFile(strman.ID(name), "C:\\NoxPost\\src\\common\\Object\\WeapLook.c")
				*(*unsafe.Pointer)(unsafe.Add(v1, -4)) = unsafe.Pointer(alloc.InternCString16(str))
				v2 := *(*uint32)(unsafe.Add(v1, 12))
				v1 = unsafe.Add(v1, 12)
				result = v1
				if v2 == 0 {
					break
				}
			}
		}
		*memmap.PtrUint32(0x5D4594, 371248) = 1
	}
	return result
}
