package legacy

/*
#include "defs.h"
#include "GAME1_1.h"
#include "GAME2.h"
#include "GAME2_1.h"
#include "client__gui__window.h"
#include "noxstring.h"

extern uint32_t dword_5d4594_1090040;
extern uint32_t dword_5d4594_1090044;
extern nox_window* dword_5d4594_1090048;
extern nox_window* dword_5d4594_1090100;
extern uint32_t dword_5d4594_1090108;
extern uint32_t dword_5d4594_1090112;
extern uint32_t dword_587000_145664;
extern uint32_t dword_587000_145668;
extern uint32_t dword_587000_145672;

// The original passes a variable number of arguments to nox_swprintf (format is loaded from the string table).
// Extra arguments are ignored by cdecl callee, so a fixed number of slots is enough.
static void rank_swprintf(wchar2_t* dst, const wchar2_t* fmt, unsigned int a1, unsigned int a2, unsigned int a3) {
	nox_swprintf(dst, fmt, a1, a2, a3);
}
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
)

// RankSwprintf is nox_swprintf with up to 3 untyped 32-bit arguments (pointers and ints).
func RankSwprintf(dst, format *uint16, a1, a2, a3 uint32) {
	C.rank_swprintf((*C.wchar2_t)(unsafe.Pointer(dst)), (*C.wchar2_t)(unsafe.Pointer(format)), C.uint(a1), C.uint(a2), C.uint(a3))
}

func Sub_470CC0() int {
	return int(C.sub_470CC0())
}

func Sub_470CD0() int {
	return int(C.sub_470CD0())
}

func Get_dword_5d4594_1090040() uint32 {
	return uint32(C.dword_5d4594_1090040)
}
func Set_dword_5d4594_1090040(v uint32) {
	C.dword_5d4594_1090040 = C.uint(v)
}
func Get_dword_5d4594_1090044() uint32 {
	return uint32(C.dword_5d4594_1090044)
}
func Set_dword_5d4594_1090044(v uint32) {
	C.dword_5d4594_1090044 = C.uint(v)
}
func Set_dword_5d4594_1090048(win *gui.Window) {
	C.dword_5d4594_1090048 = (*C.nox_window)(win.C())
}
func Set_dword_5d4594_1090100(win *gui.Window) {
	C.dword_5d4594_1090100 = (*C.nox_window)(win.C())
}
func Get_dword_5d4594_1090108() *gui.Window {
	return AsWindowP(unsafe.Pointer(uintptr(C.dword_5d4594_1090108)))
}
func Set_dword_5d4594_1090108(win *gui.Window) {
	C.dword_5d4594_1090108 = C.uint(uintptr(win.C()))
}
func Get_dword_5d4594_1090112() *gui.Window {
	return AsWindowP(unsafe.Pointer(uintptr(C.dword_5d4594_1090112)))
}
func Set_dword_5d4594_1090112(win *gui.Window) {
	C.dword_5d4594_1090112 = C.uint(uintptr(win.C()))
}
func Get_dword_587000_145664() uint32 {
	return uint32(C.dword_587000_145664)
}
func Get_dword_587000_145668() uint32 {
	return uint32(C.dword_587000_145668)
}
func Get_dword_587000_145672() int32 {
	return int32(C.dword_587000_145672)
}
func Set_dword_587000_145672(v int32) {
	C.dword_587000_145672 = C.uint(uint32(v))
}

func Nox_gui_getQuestStage_450B10() int {
	return int(C.nox_gui_getQuestStage_450B10())
}

// Nox_float2int is the C nox_float2int (float to int truncation).
func Nox_float2int(v float32) int {
	return int(C.nox_float2int(C.float(v)))
}
