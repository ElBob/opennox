package legacy

/*
#include <wctype.h>
#include "defs.h"
#include "GAME2_1.h"
#include "GAME3.h"
#include "client__shell__selcolor.h"
extern nox_gui_animation* nox_wnd_xxx_1308092;
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
)

var Nox_game_showSelColor_4A5D00 func() int

//export nox_game_showSelColor_4A5D00
func nox_game_showSelColor_4A5D00() C.int { return C.int(Nox_game_showSelColor_4A5D00()) }

var (
	Sub_4A6890 func() int
	Sub_4A6C90 func() int
)

//export sub_4A6890
func sub_4A6890() C.int { return C.int(Sub_4A6890()) }

//export sub_4A6C90
func sub_4A6C90() C.int { return C.int(Sub_4A6C90()) }

func Get_nox_game_showSelColor_4A5D00() unsafe.Pointer { return C.nox_game_showSelColor_4A5D00 }
func Get_sub_4A6890() unsafe.Pointer                   { return C.sub_4A6890 }
func Get_sub_4A6C90() unsafe.Pointer                   { return C.sub_4A6C90 }

func Get_nox_wnd_xxx_1308092() *gui.Anim {
	return asGUIAnim(C.nox_wnd_xxx_1308092)
}

func Set_nox_wnd_xxx_1308092(a *gui.Anim) {
	C.nox_wnd_xxx_1308092 = (*C.nox_gui_animation)(unsafe.Pointer(a))
}

func Nox_xxx_wndRetNULL_46A8A0() {
	C.nox_xxx_wndRetNULL_46A8A0()
}

// Iswspace calls the C iswspace (locale dependent), on a 16-bit wchar.
func Iswspace(c uint16) bool {
	return C.iswspace(C.wint_t(c)) != 0
}
