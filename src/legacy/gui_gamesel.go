package legacy

/*
#include <stdlib.h>
#include "defs.h"
#include "GAME1.h"
#include "GAME1_2.h"
#include "GAME2_3.h"
#include "GAME3.h"
#include "client__gui__window.h"
#include "client__shell__noxworld.h"
*/
import "C"
import (
	"image"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
)

var Nox_game_showGameSel_4379F0 func() int

//export nox_game_showGameSel_4379F0
func nox_game_showGameSel_4379F0() C.int { return C.int(Nox_game_showGameSel_4379F0()) }

func Get_sub_438330() unsafe.Pointer { return C.sub_438330 }
func Get_sub_438370() unsafe.Pointer { return C.sub_438370 }
func Get_nox_xxx_windowMultiplayerSub_439E70() unsafe.Pointer {
	return C.nox_xxx_windowMultiplayerSub_439E70
}

// CFree is C free().
func CFree(p unsafe.Pointer) { C.free(p) }

func Nox_xxx_wndPointInWnd_46AAB0(win unsafe.Pointer, x, y int) bool {
	return bool(C.nox_xxx_wndPointInWnd_46AAB0((*C.uint)(win), C.int(x), C.int(y)))
}

func Sub_46ACE0(win *gui.Window, a2, a3, a4 int) {
	C.sub_46ACE0((*C.uint)(win.C()), C.int(a2), C.int(a3), C.int(a4))
}

func Sub_46AD20(win *gui.Window, a2, a3, a4 int) {
	C.sub_46AD20((*C.uint)(win.C()), C.int(a2), C.int(a3), C.int(a4))
}

func Nox_xxx_wnd_46B280(a1, a2 *gui.Window) int {
	return int(C.nox_xxx_wnd_46B280(C.int(uintptr(unsafe.Pointer(a1))), C.int(uintptr(unsafe.Pointer(a2)))))
}

func Sub_49FDB0(a1 int) { C.sub_49FDB0(C.int(a1)) }

func Sub_40A3C0(a1 int) { C.sub_40A3C0(C.uint(a1)) }

func Sub_4383A0() { C.sub_4383A0() }

func Sub_438DD0(a1, a2 int) int { return int(C.sub_438DD0(C.uint(a1), C.uint(a2))) }

func Sub_439D90(a1, a2 int) { C.sub_439D90(C.uint(a1), C.uint(a2)) }

func Sub_4A2890() { C.sub_4A2890() }

func Sub_4A4800(a1 unsafe.Pointer) int { return int(C.sub_4A4800(C.int(uintptr(a1)))) }

// Sub_4A0020 returns the head of the server list (nox_gui_wol_servers_list).
func Sub_4A0020() unsafe.Pointer { return unsafe.Pointer(C.sub_4A0020()) }

func Nox_client_gui_serverInfoBlockCheckExp_439370(pos image.Point, srv unsafe.Pointer) {
	p := C.int2{field_0: C.int(pos.X), field_4: C.int(pos.Y)}
	C.nox_client_gui_serverInfoBlockCheckExp_439370(&p, C.int(uintptr(srv)))
}

func listboxResp(r C.int) gui.WindowEventResp {
	if r == 0 {
		return nil
	}
	return gui.RawEventResp(r)
}

func Nox_xxx_wndListboxProcWithoutData10_4A28E0(win *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	a3, a4 := ev.EventArgsC()
	return listboxResp(C.nox_xxx_wndListboxProcWithoutData10_4A28E0((*C.uint32_t)(win.C()), C.int(ev.EventCode()), C.uint(a3), C.int(a4)))
}

func Nox_xxx_wndListboxProcPre_4A30D0(win *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	a3, a4 := ev.EventArgsC()
	return listboxResp(C.nox_xxx_wndListboxProcPre_4A30D0((*C.nox_window)(win.C()), C.uint(ev.EventCode()), C.uint32_t(a3), C.int(a4)))
}
