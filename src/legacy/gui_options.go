package legacy

/*
#include <math.h>
#include "client__gui__window.h"
#include "GAME1_2.h"
#include "GAME1_3.h"
#include "GAME2.h"
#include "GAME2_1.h"
#include "GAME2_2.h"
#include "GAME3.h"
#include "GAME3_1.h"
extern void* dword_5d4594_1309720;
extern nox_gui_animation* nox_wnd_xxx_1309740;
extern uint32_t dword_5d4594_1309728;
extern uint32_t dword_5d4594_1309732;
extern uint32_t dword_5d4594_1309736;
extern uint32_t dword_5d4594_1522604;
extern nox_gui_animation* nox_wnd_xxx_1522608;
extern uint32_t dword_5d4594_1522612;
extern nox_window* dword_5d4594_1522616;
extern nox_window* dword_5d4594_1522620;
extern nox_window* dword_5d4594_1522624;
extern nox_window* dword_5d4594_1522628;
extern uint32_t dword_5d4594_1522632;
extern uint32_t nox_xxx_normalWndBits_587000_172880;
void sub_4AA650();

static float nox_options_sensitivity(int a4) { return powf(10.0, a4 / 50.0f - 1.0); }
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/libs/strman"

	"github.com/opennox/opennox/v1/client/gui"
)

var (
	Nox_video_setMenuOptions    func(win *gui.Window)
	Sub_4A19F0                  func(name strman.ID)
	Sub_4AAA10                  func() int
	Nox_game_showOptions_4AA6B0 func() int
	Sub_4AB0C0                  func() int
	Sub_4CB880                  func() int
	Sub_4CBB70                  func() int
	Sub_4CBBB0                  func() int
	Sub_4C3A90                  func(a1, a2 int, a3 unsafe.Pointer, a4 int) int
)

//export nox_video_setMenuOptions
func nox_video_setMenuOptions(cwin *nox_window) { Nox_video_setMenuOptions(asWindow(cwin)) }

//export sub_4A19F0
func sub_4A19F0(name *C.char) { Sub_4A19F0(strman.ID(GoString(name))) }

//export sub_4AAA10
func sub_4AAA10() int { return Sub_4AAA10() }

//export sub_4C3A90
func sub_4C3A90(a1, a2 int, a3 *C.int, a4 int) int { return Sub_4C3A90(a1, a2, unsafe.Pointer(a3), a4) }

func Sub_430AA0(v int) {
	C.sub_430AA0(C.int(v))
}

func Sub_4C35B0(v int) {
	C.sub_4C35B0(C.int(v))
}

func Get_nox_wnd_xxx_1309740() *gui.Anim {
	return asGUIAnim(C.nox_wnd_xxx_1309740)
}

func Get_dword_5d4594_1309720() *gui.Window {
	return AsWindowP(C.dword_5d4594_1309720)
}

func Set_dword_5d4594_1309720(win *gui.Window) {
	C.dword_5d4594_1309720 = unsafe.Pointer(win)
}

func Set_nox_wnd_xxx_1309740(a *gui.Anim) {
	C.nox_wnd_xxx_1309740 = (*C.nox_gui_animation)(unsafe.Pointer(a))
}

// The options screen keeps its windows in uint32 C globals.
func Get_dword_5d4594_1309728() *gui.Window {
	return AsWindowP(unsafe.Pointer(uintptr(C.dword_5d4594_1309728)))
}
func Set_dword_5d4594_1309728(win *gui.Window) {
	C.dword_5d4594_1309728 = C.uint32_t(uintptr(unsafe.Pointer(win)))
}
func Get_dword_5d4594_1309732() *gui.Window {
	return AsWindowP(unsafe.Pointer(uintptr(C.dword_5d4594_1309732)))
}
func Set_dword_5d4594_1309732(win *gui.Window) {
	C.dword_5d4594_1309732 = C.uint32_t(uintptr(unsafe.Pointer(win)))
}
func Get_dword_5d4594_1309736() *gui.Window {
	return AsWindowP(unsafe.Pointer(uintptr(C.dword_5d4594_1309736)))
}
func Set_dword_5d4594_1309736(win *gui.Window) {
	C.dword_5d4594_1309736 = C.uint32_t(uintptr(unsafe.Pointer(win)))
}

// Input configuration screen globals (GAME3_1.c).
func Get_dword_5d4594_1522604() *gui.Window {
	return AsWindowP(unsafe.Pointer(uintptr(C.dword_5d4594_1522604)))
}
func Set_dword_5d4594_1522604(win *gui.Window) {
	C.dword_5d4594_1522604 = C.uint32_t(uintptr(unsafe.Pointer(win)))
}
func Get_nox_wnd_xxx_1522608() *gui.Anim {
	return asGUIAnim(C.nox_wnd_xxx_1522608)
}
func Set_nox_wnd_xxx_1522608(a *gui.Anim) {
	C.nox_wnd_xxx_1522608 = (*C.nox_gui_animation)(unsafe.Pointer(a))
}
func Get_dword_5d4594_1522612() *gui.Window {
	return AsWindowP(unsafe.Pointer(uintptr(C.dword_5d4594_1522612)))
}
func Set_dword_5d4594_1522612(win *gui.Window) {
	C.dword_5d4594_1522612 = C.uint32_t(uintptr(unsafe.Pointer(win)))
}
func Set_dword_5d4594_1522616(win *gui.Window) {
	C.dword_5d4594_1522616 = (*C.nox_window)(win.C())
}
func Set_dword_5d4594_1522620(win *gui.Window) {
	C.dword_5d4594_1522620 = (*C.nox_window)(win.C())
}
func Set_dword_5d4594_1522624(win *gui.Window) {
	C.dword_5d4594_1522624 = (*C.nox_window)(win.C())
}
func Set_dword_5d4594_1522628(win *gui.Window) {
	C.dword_5d4594_1522628 = (*C.nox_window)(win.C())
}
func Get_dword_5d4594_1522632() *gui.Window {
	return AsWindowP(unsafe.Pointer(uintptr(C.dword_5d4594_1522632)))
}
func Set_dword_5d4594_1522632(win *gui.Window) {
	C.dword_5d4594_1522632 = C.uint32_t(uintptr(unsafe.Pointer(win)))
}

func Set_nox_xxx_normalWndBits_587000_172880(v int) {
	C.nox_xxx_normalWndBits_587000_172880 = C.uint32_t(v)
}

// SliderSensitivity is powf(10.0, a4 / 50.0f - 1.0), computed by the C library as in the original.
func SliderSensitivity(a4 int) float32 {
	return float32(C.nox_options_sensitivity(C.int(a4)))
}

// Function pointers of C-callable exports, stored in the animation records.
func Get_sub_4AA9C0() unsafe.Pointer { return C.sub_4AA9C0 }
func Get_sub_4AAA10() unsafe.Pointer { return C.sub_4AAA10 }
func Get_sub_4AB0C0() unsafe.Pointer { return C.sub_4AB0C0 }
func Get_sub_4CB880() unsafe.Pointer { return C.sub_4CB880 }
func Get_sub_4CBB70() unsafe.Pointer { return C.sub_4CBB70 }
func Get_sub_4CBBB0() unsafe.Pointer { return C.sub_4CBBB0 }

func Sub_43DC10()                          { C.sub_43DC10() }
func Sub_44D970()                          { C.sub_44D970() }
func Nox_xxx____setargv_9_453060()         { C.nox_xxx____setargv_9_453060() }
func Sub_4AA650()                          { C.sub_4AA650() }
func Sub_47DBC0() int                      { return int(C.sub_47DBC0()) }
func Nox_xxx_wndRetNULL_46A8A0()           { C.nox_xxx_wndRetNULL_46A8A0() }
func Sub_4A4800(widget unsafe.Pointer) int { return int(C.sub_4A4800(C.int(uintptr(widget)))) }

func Nox_xxx_wndListboxProcPre_4A30D0(win *gui.Window, ev uint32, a3 uint32, a4 int) int {
	return int(C.nox_xxx_wndListboxProcPre_4A30D0((*nox_window)(win.C()), C.uint(ev), C.uint32_t(a3), C.int(a4)))
}

func Nox_xxx_wndListboxProcWithoutData10_4A28E0(win *gui.Window, ev int, a3 uint32, a4 int) int {
	return int(C.nox_xxx_wndListboxProcWithoutData10_4A28E0((*C.uint32_t)(win.C()), C.int(ev), C.uint(a3), C.int(a4)))
}

//export nox_game_showOptions_4AA6B0
func nox_game_showOptions_4AA6B0() C.int { return C.int(Nox_game_showOptions_4AA6B0()) }

//export sub_4AB0C0
func sub_4AB0C0() int { return Sub_4AB0C0() }

//export sub_4CB880
func sub_4CB880() int { return Sub_4CB880() }

//export sub_4CBB70
func sub_4CBB70() int { return Sub_4CBB70() }

//export sub_4CBBB0
func sub_4CBBB0() int { return Sub_4CBBB0() }
