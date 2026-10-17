package legacy

/*
#include "defs.h"
#include "GAME2.h"
#include "GAME2_1.h"
#include "GAME3.h"
#include "GAME4.h"
#include "client__shell__selcolor.h"
void sub_4AA650();
extern nox_gui_animation* nox_wnd_xxx_1307732;
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
)

var (
	Nox_game_showSelClass_4A4840 func() int
	Sub_4A4970                   func() int
	Sub_4A49A0                   func() int
	Sub_4AA650                   func()
)

//export nox_game_showSelClass_4A4840
func nox_game_showSelClass_4A4840() C.int { return C.int(Nox_game_showSelClass_4A4840()) }

//export sub_4A4970
func sub_4A4970() C.int { return C.int(Sub_4A4970()) }

//export sub_4A49A0
func sub_4A49A0() C.int { return C.int(Sub_4A49A0()) }

//export sub_4AA650
func sub_4AA650() { Sub_4AA650() }

func Get_sub_4A4970() unsafe.Pointer                   { return C.sub_4A4970 }
func Get_sub_4A49A0() unsafe.Pointer                   { return C.sub_4A49A0 }
func Get_nox_game_showSelColor_4A5D00() unsafe.Pointer { return C.nox_game_showSelColor_4A5D00 }

func Sub_5007E0(s string) {
	cs := CString(s)
	defer StrFree(cs)
	C.sub_5007E0(cs)
}

func Nox_xxx_wndRetNULL_46A8A0() {
	C.nox_xxx_wndRetNULL_46A8A0()
}

func Nox_xxx_book_45DBE0(a1, a2, a3 int) {
	C.nox_xxx_book_45DBE0(unsafe.Pointer(uintptr(a1)), C.int(a2), C.int(a3))
}

func Get_nox_wnd_xxx_1307732() *gui.Anim {
	return asGUIAnim(C.nox_wnd_xxx_1307732)
}

func Set_nox_wnd_xxx_1307732(a *gui.Anim) {
	C.nox_wnd_xxx_1307732 = (*C.nox_gui_animation)(unsafe.Pointer(a))
}
