package legacy

/*
#include "GAME5_2.h"
#include "common__gamemech__pausefx.h"
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/server"
)

var (
	PlayerGiveExp      func(u *server.Object, xp float32)
	PlayerCheckLevelUp func(u *server.Object)
)

//export nox_xxx_plyrGiveExp_4EF3A0_exp_level
func nox_xxx_plyrGiveExp_4EF3A0_exp_level(a1 C.int, xp C.float) {
	PlayerGiveExp(asObjectS((*nox_object_t)(unsafe.Pointer(uintptr(a1)))), float32(xp))
}

//export sub_4EF2E0_exp_level
func sub_4EF2E0_exp_level(a1 C.int) {
	PlayerCheckLevelUp(asObjectS((*nox_object_t)(unsafe.Pointer(uintptr(a1)))))
}

func Sub_56FA40(h uint32, v float32) {
	C.sub_56FA40(C.int(h), C.float(v))
}

func Sub_56F980(h uint32, v byte) {
	C.sub_56F980(C.int(h), C.uchar(v))
}

func Sub_57AF30(u *server.Object, a2 int) {
	C.sub_57AF30(C.int(uintptr(unsafe.Pointer(asObjectC(u)))), C.int(a2))
}
