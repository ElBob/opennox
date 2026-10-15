package legacy

/*
#include "server__gamemech__explevel.h"
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/server"
)

func Nox_xxx_plyrGiveExp_4EF3A0_exp_level(u *server.Object, xp float32) {
	C.nox_xxx_plyrGiveExp_4EF3A0_exp_level(C.int(uintptr(unsafe.Pointer(asObjectC(u)))), C.float(xp))
}
