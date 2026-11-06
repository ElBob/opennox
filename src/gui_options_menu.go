package opennox

import (
	"fmt"
	"unicode/utf16"
	"unsafe"

	"github.com/opennox/libs/client/keybind"

	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/sound"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/ccall"
	"github.com/opennox/opennox/v1/legacy/timer"
)

const inputCfgSrcFile = "C:\\NoxPost\\src\\Client\\shell\\InputCfg\\inputcfg.c"

func init() {
	legacy.Nox_game_showOptions_4AA6B0 = func() int {
		if nox_game_showOptions_4AA6B0() {
			return 1
		}
		return 0
	}
	legacy.Sub_4AB0C0 = sub_4AB0C0
	legacy.Sub_4CB880 = sub_4CB880
	legacy.Sub_4CBB70 = sub_4CBB70
	legacy.Sub_4CBBB0 = sub_4CBBB0
}

// winRespInt converts a C-style int return value of a window function.
func winRespInt(r int) gui.WindowEventResp {
	if r == 0 {
		return nil
	}
	return gui.RawEventResp(r)
}

// wndID is nox_xxx_wndGetID_46B0A0.
func wndID(win *gui.Window) int {
	if win == nil {
		return -2
	}
	return int(win.ID())
}

// optionsCallKeyEnter calls the raw function pointer at window offset 372 (field 93) with
// (win, 21, 28, 2), as the original does for the sound/music checkboxes (key event, Enter, pressed).
// Quirk: this reads the raw C slot. Windows created by the Go GUI code keep their handlers in the Go
// extension instead, so this slot is NULL for them and the original call is skipped.
func optionsCallKeyEnter(win *gui.Window) {
	if fnc := win.Field93Raw(); fnc != nil {
		ccall.CallUPtrUPtr4(fnc, uintptr(win.C()), 21, 28, 2)
	}
}

// optionsTimerSet is sub_486320.
func optionsTimerSet(p unsafe.Pointer, v int) {
	(*timer.Timer)(p).SetRaw(uint32(v))
}

// optionsTimerValue reads the high word of the timer current value: *(uint32*)(p + 4) >> 16.
func optionsTimerValue(p unsafe.Pointer) int {
	return int(*(*uint32)(unsafe.Add(p, 4)) >> 16)
}

func optionsSetupSlider(root *gui.Window, id uint, timerGroup unsafe.Pointer) {
	v1 := root.ChildByID(id)
	v1.Field100Ptr.SizeVal.X = 24
	v1.Field100Ptr.SizeVal.Y = 20
	// the original loads the same image names in this order (hl, sel, enabled)
	v19 := nox_xxx_gLoadImg("OptionsVolumeSliderLit")
	v16 := nox_xxx_gLoadImg("OptionsVolumeSliderLit")
	v2 := nox_xxx_gLoadImg("OptionsVolumeSlider")
	gui.ButtonSetImage(v1, nil, nil, v2.C(), v16.C(), v19.C())
	v1.Func94(gui.AsWindowEvent(16395, 0, 0x4000))
	v1.Func94(gui.AsWindowEvent(16394, uintptr(optionsTimerValue(timerGroup)), 0))
}

func optionsSetupCheck(win *gui.Window, enabled int) {
	if enabled == 1 {
		win.DrawData().Field0 |= 4
	} else {
		win.DrawData().Field0 &= 0xFFFFFFFB
	}
}

// nox_game_showOptions_4AA6B0 is the options screen constructor.
// It keeps a C export, the C state code references it by address.
func nox_game_showOptions_4AA6B0() bool {
	c := noxClient
	c.GameAddStateCode(client.StateOptions)
	root := nox_new_window_from_file("Options.wnd", guiOptionsProc)
	legacy.Set_dword_5d4594_1309720(root)
	if root == nil {
		return false
	}
	// nox_client_advVideoOpts_New_4CB590 always returns 1, the original checks it for zero
	c.NewGUIAdvOptsOn(root)
	root.SetFunc93(sub4A18E0)
	c.r.SetTabWidth(15) // nox_draw_setTabWidth_43FE20, returns the old value
	anim := nox_gui_makeAnimation(root, 0, 0, 0, -480, 0, 20, 0, -40)
	legacy.Set_nox_wnd_xxx_1309740(anim)
	if anim == nil {
		return false
	}
	anim.StateID = client.StateOptions
	anim.Func12Ptr = legacy.Get_sub_4AA9C0()
	anim.FncDoneOutPtr = legacy.Get_sub_4AAA10()

	optionsSetupSlider(root, 351, legacy.Get_dword_587000_127004())
	v := root.ChildByID(361)
	legacy.Set_dword_5d4594_1309728(v)
	optionsSetupCheck(v, legacy.Sub_453070())

	optionsSetupSlider(root, 352, legacy.Get_dword_587000_122852())
	v = root.ChildByID(362)
	legacy.Set_dword_5d4594_1309732(v)
	optionsSetupCheck(v, legacy.Sub_44D990())

	optionsSetupSlider(root, 353, legacy.Get_dword_587000_93164())
	v = root.ChildByID(363)
	legacy.Set_dword_5d4594_1309736(v)
	optionsSetupCheck(v, legacy.Sub_43DC30())

	legacy.Nox_xxx_wndRetNULL_46A8A0()
	guiSetBackButtonText("OptsBack.wnd:Back") // sub_4A19F0
	sub_4A1A40(0)
	sub_4AAA70()
	return true
}

// sub_4AAA70 selects the radio buttons for the current video settings.
func sub_4AAA70() {
	root := legacy.Get_dword_5d4594_1309720()
	nox_video_setMenuOptions(root)
	// TODO (original): This is original handling of 8/16bit and windowed/fullscreen.
	//       However, windowed/fullscreen switch is not present in options.wnd, thus it's not visible in the menu,
	//       even though the code for it exists here.
	//       So currently we just pretend that 8/16bit switch actually corresponds to windowed/fullscreen.
	//       Later we should refactor it to inject those switches.
	// The original has two '#if 0' blocks here (8/16 bit selection via 331/332 + 333/334, and
	// nox_video_bagexists_4300D0 checks); they are not compiled in the original either.
	var v1 *gui.Window
	if nox_video_getFullScreen() == 0 {
		v1 = root.ChildByID(331)
	} else {
		v1 = root.ChildByID(332)
	}
	v1.Func94(gui.AsWindowEvent(16392, 1, 0))
	v4 := nox_video_getCutSize_4766D0()
	*memmap.PtrUint32(0x587000, 172884) = uint32(v4)
	var id uint
	if v4 > 69 {
		if v4 <= 79 {
			id = 312
		} else if v4 <= 89 {
			id = 313
		} else {
			id = 314
		}
	} else {
		id = 311
	}
	if res := root.ChildByID(id); res != nil {
		res.Func94(gui.AsWindowEvent(16392, 1, 0))
	}
}

// guiOptionsProc is sub_4AABE0, the options window procedure.
func guiOptionsProc(win *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	switch ev := ev.(type) {
	case *WindowEvent0x4005: // 16389
		clientPlaySoundSpecial(sound.SoundShellSelect, 100)
		return gui.RawEventResp(1)
	case *WindowEvent0x4007: // 16391
		v4 := wndID(ev.Win)
		if v4 >= 380 {
			return winRespInt(nox_gui_menu_proc_ext(v4))
		}
		if v4 > 341 {
			if v4 > 363 {
				if v4 == 2099 {
					sub_4AAA70()
				}
			} else if v4 == 363 {
				if legacy.Sub_43DC30() == 1 {
					legacy.Sub_43DC00()
				} else {
					legacy.Sub_43DC10()
					optionsTimerSet(legacy.Get_dword_587000_93164(), optionsTimerValue(legacy.Get_dword_587000_93164()))
				}
			} else {
				v5 := v4 - 361
				if v5 != 0 {
					if v5 == 1 {
						if legacy.Sub_44D990() == 1 {
							legacy.Sub_44D960()
						} else {
							legacy.Sub_44D970()
						}
					}
				} else if legacy.Sub_453070() == 1 {
					legacy.Sub_453050()
				} else {
					legacy.Nox_xxx____setargv_9_453060()
				}
			}
		} else if v4 == 341 {
			sub_4AA9C0()
			anim := legacy.Get_nox_wnd_xxx_1309740()
			anim.FncDoneOutPtr = legacy.Get_sub_4AB0C0()
			anim.Func13Ptr = legacy.Get_sub_4CB880()
		} else {
			switch v4 {
			case 311:
				nox_video_setCutSize_4766A0(65)
			case 312:
				nox_video_setCutSize_4766A0(75)
			case 313:
				nox_video_setCutSize_4766A0(85)
			case 314:
				nox_video_setCutSize_4766A0(100)
			case 332:
				legacy.Set_nox_xxx_normalWndBits_587000_172880(16)
				fallthrough // original falls through here
			case 333:
				nox_video_setFullScreen(1)
			case 331:
				legacy.Set_nox_xxx_normalWndBits_587000_172880(8)
				fallthrough // original falls through here
			case 334:
				nox_video_setFullScreen(0)
			}
		}
		clientPlaySoundSpecial(sound.SoundShellClick, 100)
		return gui.RawEventResp(1)
	case *WindowEvent0x4009: // 16393
		a4 := ev.Val
		v7 := wndID(ev.Win)
		v8 := noxClient.GUI.WinYYY // nox_client_getWin1064916_46C720
		switch v7 {
		case 351:
			optionsTimerSet(legacy.Get_dword_587000_127004(), a4)
			if v8 != nil && v8.Parent() == ev.Win {
				return nil
			}
			chk := legacy.Get_dword_5d4594_1309728()
			if a4 != 0 {
				if legacy.Sub_453070() != 0 {
					clientPlaySoundSpecial(sound.ID(768), 100)
					return nil
				}
				optionsCallKeyEnter(chk)
				clientPlaySoundSpecial(sound.ID(768), 100)
				return nil
			}
			if legacy.Sub_453070() != 1 {
				return nil
			}
			optionsCallKeyEnter(chk)
			return nil
		case 352:
			optionsTimerSet(legacy.Get_dword_587000_122852(), a4)
			chk := legacy.Get_dword_5d4594_1309732()
			if a4 != 0 {
				if legacy.Sub_44D990() != 0 {
					legacy.Sub_4AA650()
					return nil
				}
				optionsCallKeyEnter(chk)
				legacy.Sub_4AA650()
				return nil
			}
			if legacy.Sub_44D990() != 1 {
				return nil
			}
			optionsCallKeyEnter(chk)
			return nil
		case 353:
			optionsTimerSet(legacy.Get_dword_587000_93164(), a4)
			chk := legacy.Get_dword_5d4594_1309736()
			if a4 != 0 {
				if legacy.Sub_43DC30() != 0 {
					return nil
				}
				optionsCallKeyEnter(chk)
				return nil
			}
			if legacy.Sub_43DC30() != 1 {
				return nil
			}
			optionsCallKeyEnter(chk)
			return nil
		case 316: // gamma slider
			nox_video_setGammaSlider(a4)
		case 318: // input sensitivity slider
			noxClient.SetSensitivity(legacy.SliderSensitivity(a4))
		default:
			return nil
		}
		return nil
	case *WindowEvent0x400c: // 16396
		a4 := int(ev.Val)
		v17 := wndID(ev.Win) - 351
		if v17 != 0 {
			if v17 != 1 {
				return nil
			}
			optionsTimerSet(legacy.Get_dword_587000_122852(), a4)
			chk := legacy.Get_dword_5d4594_1309732()
			if a4 != 0 {
				if legacy.Sub_44D990() == 0 {
					optionsCallKeyEnter(chk)
				}
				legacy.Sub_4AA650()
				return nil
			}
			if legacy.Sub_44D990() != 1 {
				return nil
			}
			optionsCallKeyEnter(chk)
			return nil
		}
		optionsTimerSet(legacy.Get_dword_587000_127004(), a4)
		chk := legacy.Get_dword_5d4594_1309728()
		if a4 != 0 {
			if legacy.Sub_453070() == 0 {
				optionsCallKeyEnter(chk)
			}
			clientPlaySoundSpecial(sound.ID(768), 100)
			return nil
		}
		if legacy.Sub_453070() != 1 {
			return nil
		}
		optionsCallKeyEnter(chk)
		return nil
	}
	return nil
}

// sub_4AB0C0 is the options screen "done out" callback.
func sub_4AB0C0() int {
	anim := legacy.Get_nox_wnd_xxx_1309740()
	v0 := anim.Func13Ptr
	anim.Free()
	legacy.Get_dword_5d4594_1309720().Destroy()
	ccall.CallIntVoid(v0)
	return 1
}

// --- input configuration (InputCfg.wnd) ---

// inputCfgProc is sub_4CBE70 as a window function.
func inputCfgProc(win *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	a3, a4 := ev.EventArgsC()
	return winRespInt(sub_4CBE70(int(uintptr(win.C())), ev.EventCode(), unsafe.Pointer(a3), int(a4)))
}

func inputCfgItemText(win *gui.Window, i int) *uint16 {
	return (*uint16)(gui.EventRespPtr(win.Func94(gui.AsWindowEvent(16406, uintptr(i), 0))))
}

func inputCfgSetItemText(win *gui.Window, str unsafe.Pointer, i int) {
	win.Func94(gui.AsWindowEvent(16407, uintptr(str), uintptr(i)))
}

// wcsEqual is nox_wcscmp(a, b) == 0.
func wcsEqual(a, b *uint16) bool {
	for i := 0; ; i++ {
		ca := *(*uint16)(unsafe.Add(unsafe.Pointer(a), 2*i))
		cb := *(*uint16)(unsafe.Add(unsafe.Pointer(b), 2*i))
		if ca != cb {
			return false
		}
		if ca == 0 {
			return true
		}
	}
}

// inputCfgTitleByKey is nox_xxx_keybind_titleByKey_42EA00.
func inputCfgTitleByKey(key uintptr) *uint16 {
	k := keybind.Key(key)
	if !k.IsValid() {
		return alloc.InternCString16("")
	}
	return alloc.InternCString16(k.Title(noxClient.Strings()))
}

// inputCfgTitleByKeyZero is nox_xxx_keybind_titleByKeyZero_42EA00.
func inputCfgTitleByKeyZero(key uintptr) *uint16 {
	k := keybind.Key(key)
	if !k.IsValid() {
		return nil
	}
	return alloc.InternCString16(k.Title(noxClient.Strings()))
}

// sub_4CB880 creates the input configuration screen.
func sub_4CB880() int {
	c := noxClient
	c.GameAddStateCode(900)
	root := nox_new_window_from_file("InputCfg.wnd", inputCfgProc)
	legacy.Set_dword_5d4594_1522604(root)
	if root == nil {
		return 0
	}
	root.SetFunc93(sub4A18E0)
	anim := nox_gui_makeAnimation(root, 0, 0, 0, -480, 0, 20, 0, -40)
	legacy.Set_nox_wnd_xxx_1522608(anim)
	if anim == nil {
		return 0
	}
	anim.StateID = 900
	anim.Func12Ptr = legacy.Get_sub_4CBB70()
	anim.FncDoneOutPtr = legacy.Get_sub_4CBBB0()
	win16 := root.ChildByID(910)
	legacy.Set_dword_5d4594_1522616(win16)
	win20 := root.ChildByID(911)
	legacy.Set_dword_5d4594_1522620(win20)
	win24 := root.ChildByID(912)
	legacy.Set_dword_5d4594_1522624(win24)
	win28 := root.ChildByID(913)
	legacy.Set_dword_5d4594_1522628(win28)
	if win16 == nil {
		return 0
	}
	data := (*gui.ScrollListBoxData)(win16.WidgetData)
	*(*uint32)(data.Field_7) = 921
	*(*uint32)(data.Field_8) = 922
	*(*uint32)(data.Field_9) = 920
	win16.SetFunc94(inputCfgListProc)
	win20.SetParent(win16)
	win24.SetParent(win16)
	win28.SetParent(win16)
	win24.SetFunc93(inputCfgListNoScroll)
	win28.SetFunc93(inputCfgListNoScroll)
	v2 := root.ChildByID(921)
	win20.Func94(gui.AsWindowEvent(16408, uintptr(v2.C()), 0))
	win24.Func94(gui.AsWindowEvent(16408, uintptr(v2.C()), 0))
	win28.Func94(gui.AsWindowEvent(16408, uintptr(v2.C()), 0))
	v3 := root.ChildByID(922)
	win20.Func94(gui.AsWindowEvent(16409, uintptr(v3.C()), 0))
	win24.Func94(gui.AsWindowEvent(16409, uintptr(v3.C()), 0))
	win28.Func94(gui.AsWindowEvent(16409, uintptr(v3.C()), 0))
	sub_4CBBF0()
	v4 := uint(971)
	v5 := uint(legacy.Sub_47DBC0()) + 971
	if v5 > 971 {
		for {
			v6 := root.ChildByID(v4)
			legacy.Nox_xxx_wnd_46ABB0(v6, 1)
			v4++
			if v4 >= v5 {
				break
			}
		}
	}
	v7 := legacy.Nox_client_mousePriKey_430AF0()
	v8 := root.ChildByID(uint(v7 + 971))
	v8.Func94(gui.AsWindowEvent(16392, 1, 0))
	win980 := root.ChildByID(980)
	legacy.Set_dword_5d4594_1522612(win980)
	win980.SetParent(nil)
	win980.SetFunc94(inputCfgProc)
	win980.SetFunc93(inputCfgPopupProc)
	win980.Hide()
	v9 := win980.ChildByID(981)
	v9.Func94(gui.AsWindowEvent(16385, uintptr(memmap.PtrOff(0x5D4594, 1522636)), 0))
	return 1
}

// sub_4CBB70 applies the key bindings and starts the screen slide out.
func sub_4CBB70() int {
	sub_4CBD30()
	nox_common_writecfgfile("nox.cfg")
	legacy.Get_nox_wnd_xxx_1522608().SetState(gui.AnimOut)
	sub_43BE40(2)
	clientPlaySoundSpecial(sound.SoundShellSlideOut, 100)
	return 1
}

// sub_4CBBB0 is the "done out" callback of the input configuration screen.
func sub_4CBBB0() int {
	anim := legacy.Get_nox_wnd_xxx_1522608()
	v0 := anim.Func13Ptr
	anim.Free()
	legacy.Get_dword_5d4594_1522604().Destroy()
	legacy.Get_dword_5d4594_1522612().Destroy()
	ccall.CallIntVoid(v0)
	return 1
}

// inputCfgListNoScroll is sub_4CC140, the function of the key columns (events 19 and 20 are ignored).
func inputCfgListNoScroll(win *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	a2 := uint32(ev.EventCode())
	a3, a4 := ev.EventArgsC()
	if a2 < 0x13 || a2 > 0x14 {
		return winRespInt(legacy.Nox_xxx_wndListboxProcWithoutData10_4A28E0(win, int(a2), uint32(a3), int(a4)))
	}
	return nil
}

func inputCfgClosePopup() {
	popup := legacy.Get_dword_5d4594_1522612()
	noxClient.GUI.Focus(nil) // nox_xxx_windowFocus_46B500(0)
	popup.StackPop()
	popup.Hide()
}

// inputCfgPopupProc is sub_4CC170, the "press a key" popup function.
func inputCfgPopupProc(win *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	a2 := ev.EventCode()
	a3, a4 := ev.EventArgsC()
	switch a2 {
	case 6, 7:
		sub_4CC3C0(0x10000)
		inputCfgClosePopup()
		return gui.RawEventResp(1)
	case 10, 11:
		sub_4CC3C0(0x10002)
		inputCfgClosePopup()
		return gui.RawEventResp(1)
	case 14, 15:
		sub_4CC3C0(0x10001)
		inputCfgClosePopup()
		return gui.RawEventResp(1)
	case 19:
		sub_4CC3C0(0x10003)
		inputCfgClosePopup()
		return gui.RawEventResp(1)
	case 20:
		sub_4CC3C0(0x10004)
		inputCfgClosePopup()
		return gui.RawEventResp(1)
	case 21:
		if a3 == 1 {
			if a4 == 2 {
				inputCfgClosePopup()
				if sel := legacy.Get_dword_5d4594_1522632(); sel != nil {
					sel.Func94(gui.AsWindowEvent(16403, ^uintptr(0), 0))
				}
			}
			return gui.RawEventResp(1)
		} else if a4 == 1 && sub_4CC280(a3) != 0 {
			sub_4CC3C0(0x10000)
			inputCfgClosePopup()
			return gui.RawEventResp(1)
		}
		return nil
	}
	return nil
}

// sub_4CC280 binds a keyboard key to the selected action, clearing it from other actions.
func sub_4CC280(a1 uintptr) int {
	v6 := inputCfgTitleByKeyZero(a1)
	if v6 == nil {
		return 0
	}
	sel := legacy.Get_dword_5d4594_1522632()
	if sel == nil {
		return 1
	}
	inputCfgClearBinding(legacy.Get_dword_5d4594_1522624(), v6, memmap.PtrOff(0x587000, 187824))
	inputCfgClearBinding(legacy.Get_dword_5d4594_1522628(), v6, memmap.PtrOff(0x587000, 187828))
	inputCfgAssign(sel, v6)
	return 1
}

// sub_4CC3C0 binds a mouse button to the selected action, clearing it from other actions.
func sub_4CC3C0(a1 uintptr) int {
	v1 := inputCfgTitleByKey(a1)
	sel := legacy.Get_dword_5d4594_1522632()
	if sel == nil {
		return 1
	}
	inputCfgClearBinding(legacy.Get_dword_5d4594_1522624(), v1, memmap.PtrOff(0x587000, 187832))
	inputCfgClearBinding(legacy.Get_dword_5d4594_1522628(), v1, memmap.PtrOff(0x587000, 187836))
	inputCfgAssign(sel, v1)
	return 1
}

// inputCfgClearBinding replaces all items of the list that are equal to title by the string at repl.
func inputCfgClearBinding(list *gui.Window, title *uint16, repl unsafe.Pointer) {
	data := (*gui.ScrollListBoxData)(list.WidgetData)
	for i := 0; i < int(int16(data.Field_11_0)); i++ {
		if wcsEqual(inputCfgItemText(list, i), title) {
			inputCfgSetItemText(list, repl, i)
		}
	}
}

// inputCfgAssign sets the title for the selected row of the list, and resets the selection.
func inputCfgAssign(sel *gui.Window, title *uint16) {
	data := (*gui.ScrollListBoxData)(sel.WidgetData)
	row := *(*int32)(unsafe.Add(unsafe.Pointer(data), 48)) // Field_12 used as the selected row index
	inputCfgSetItemText(sel, unsafe.Pointer(title), int(row))
	sel.Func94(gui.AsWindowEvent(16403, ^uintptr(0), 0))
	legacy.Set_dword_5d4594_1522632(nil)
}

// sub_4CBD30 applies the key bindings from the lists into the control event handler.
func sub_4CBD30() {
	c := noxClient
	win20 := legacy.Get_dword_5d4594_1522620()
	win24 := legacy.Get_dword_5d4594_1522624()
	win28 := legacy.Get_dword_5d4594_1522628()
	v0 := (*gui.ScrollListBoxData)(win20.WidgetData)
	c.ctrl.Reset() // sub_42CD90
	for i := 0; i < int(int16(v0.Field_11_0)); i++ {
		v3 := nox_xxx_bindevent_bindNameByTitle_42EA40(alloc.GoString16(inputCfgItemText(win20, i)))
		v5 := nox_xxx_keybind_nameByTitle_42E960(alloc.GoString16(inputCfgItemText(win24, i)))
		v7 := nox_xxx_keybind_nameByTitle_42E960(alloc.GoString16(inputCfgItemText(win28, i)))
		// the original passes a NULL name to %s, which prints "(null)"
		evName := "(null)"
		if v3 != nil {
			evName = v3.Name
		}
		if v7 != 0 {
			nox_client_parseConfigHotkeysLine_42CF50(fmt.Sprintf("%s = %s", v7.String(), evName))
		}
		if v5 != 0 {
			nox_client_parseConfigHotkeysLine_42CF50(fmt.Sprintf("%s = %s", v5.String(), evName))
		}
	}
	v8 := c.Strings().GetStringInFile("bindevent:ToggleQuitMenu", inputCfgSrcFile)
	v9 := nox_xxx_bindevent_bindNameByTitle_42EA40(v8)
	v10 := c.Strings().GetStringInFile("keybind:Esc", inputCfgSrcFile)
	if res := nox_xxx_keybind_nameByTitle_42E960(v10); res != 0 {
		evName := "(null)"
		if v9 != nil {
			evName = v9.Name
		}
		nox_client_parseConfigHotkeysLine_42CF50(fmt.Sprintf("%s = %s", res.String(), evName))
	}
}

// inputCfgListProc is sub_4CBF60, the function of the key binding list.
func inputCfgListProc(win *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	a2 := uint32(ev.EventCode())
	a3, a4 := ev.EventArgsC()
	pre := func(code uint32) gui.WindowEventResp {
		return winRespInt(legacy.Nox_xxx_wndListboxProcPre_4A30D0(win, code, uint32(a3), int(a4)))
	}
	root := legacy.Get_dword_5d4594_1522604()
	win20 := legacy.Get_dword_5d4594_1522620()
	win24 := legacy.Get_dword_5d4594_1522624()
	win28 := legacy.Get_dword_5d4594_1522628()
	if a2 > 0x4007 {
		if a2 == 16393 {
			v7 := win.WidgetData
			pre(0x4009)
			v8 := uintptr(legacy.Sub_4A4800(v7))
			win20.Func94(gui.AsWindowEvent(16412, v8, 0))
			win24.Func94(gui.AsWindowEvent(16412, v8, 0))
			win28.Func94(gui.AsWindowEvent(16412, v8, 0))
		} else if a2 == 16400 {
			target := legacy.AsWindowP(unsafe.Pointer(a3))
			v5 := unsafe.Pointer(target.WidgetData)
			if *(*int32)(unsafe.Add(v5, 48)) >= 0 {
				legacy.Set_dword_5d4594_1522632(target)
				v9 := inputCfgItemText(win20, int(*(*int32)(unsafe.Add(v5, 48))))
				v6 := noxClient.Strings().GetStringInFile("InputCfg.wnd:PressKey", inputCfgSrcFile)
				// nox_swprintf(buf, L"%s\n'%s'", v6, v9) into the blob buffer
				inputCfgSwprintf(memmap.PtrOff(0x5D4594, 1522636), v6+"\n'"+alloc.GoString16(v9)+"'")
				popup := legacy.Get_dword_5d4594_1522612()
				popup.ShowModal()
				noxClient.GUI.Focus(popup)
				popup.StackPush()
				return pre(0x4010)
			}
		}
	} else {
		if a2 != 16391 {
			if a2 == 23 {
				return gui.RawEventResp(1)
			}
			if a2 != 0x4000 {
				return pre(a2)
			}
		}
		tw := unsafe.Pointer(a3)
		if tw == root.ChildByID(921).C() || tw == root.ChildByID(922).C() {
			win20.Func94(gui.AsWindowEvent(int(a2), a3, 0))
			win24.Func94(gui.AsWindowEvent(int(a2), a3, 0))
			win28.Func94(gui.AsWindowEvent(int(a2), a3, 0))
			return pre(a2)
		}
	}
	return pre(a2)
}

// inputCfgSwprintf writes a UTF-16 string with a terminator to an unbounded buffer (nox_swprintf).
func inputCfgSwprintf(dst unsafe.Pointer, str string) {
	u := utf16.Encode([]rune(str))
	for i, c := range u {
		*(*uint16)(unsafe.Add(dst, 2*i)) = c
	}
	*(*uint16)(unsafe.Add(dst, 2*len(u))) = 0
}
