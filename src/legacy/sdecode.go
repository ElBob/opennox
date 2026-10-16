package legacy

/*
#include "defs.h"
#include "GAME3_2.h"
#include "GAME4.h"
#include "GAME4_3.h"
#include "GAME5_2.h"
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/server"
)

//export nox_xxx_packetDynamicUnitCode_578B40
func nox_xxx_packetDynamicUnitCode_578B40(a1 C.int) C.int {
	return C.int(GetServer().S().Objs.DynamicNetCode(uint16(a1)))
}

//export nox_xxx_equipedItemByCode_4F7920
func nox_xxx_equipedItemByCode_4F7920(a1, a2 C.int) C.int {
	it := asObjectS((*nox_object_t)(unsafe.Pointer(uintptr(a1)))).InventoryItemByNetCode(int(a2))
	return C.int(uintptr(unsafe.Pointer(asObjectC(it))))
}

//export nox_xxx_useByNetCode_53F8E0
func nox_xxx_useByNetCode_53F8E0(a1, a2 C.int) C.int {
	return C.int(bool2int(asObjectS((*nox_object_t)(unsafe.Pointer(uintptr(a2)))).CallUse(asObjectS((*nox_object_t)(unsafe.Pointer(uintptr(a1)))))))
}

func Nox_xxx_netSendSecondaryWeapon_4D9670(pli ntype.PlayerInd, item *server.Object, a3 byte) {
	C.nox_xxx_netSendSecondaryWeapon_4D9670(C.int(pli), (*C.uint32_t)(unsafe.Pointer(asObjectC(item))), C.char(a3))
}
