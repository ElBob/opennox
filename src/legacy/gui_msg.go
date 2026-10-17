package legacy

/*
#include "defs.h"
#include "GAME2_1.h"
#include "GAME5_2.h"
#include "client__gui__gamewin__gamewin.h"
#include "client__gui__guimsg.h"
extern uint32_t dword_5d4594_825736;
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/client"
)

var (
	Nox_xxx_printCentered func(str string)
	Nox_xxx_clientPickup  func(dr *client.Drawable)
)

//export nox_xxx_printCentered_445490
func nox_xxx_printCentered_445490(str *wchar2_t) {
	if str == nil {
		return
	}
	Nox_xxx_printCentered(GoWString(str))
}

//export nox_xxx_clientPickup_46C140
func nox_xxx_clientPickup_46C140(dr *nox_drawable) {
	Nox_xxx_clientPickup(asDrawable(dr))
}

//export nox_xxx_netGetUnitCodeCli_578B00
func nox_xxx_netGetUnitCodeCli_578B00(dr C.int) C.uint {
	return C.uint(asDrawable((*nox_drawable)(unsafe.Pointer(uintptr(dr)))).NetCodeCli())
}

func Get_dword_5d4594_825736() uint32  { return uint32(C.dword_5d4594_825736) }
func Set_dword_5d4594_825736(v uint32) { C.dword_5d4594_825736 = C.uint32_t(v) }

func Sub_467B00(typ, a2 int) bool {
	return C.sub_467B00(C.int(typ), C.int(a2)) != 0
}
