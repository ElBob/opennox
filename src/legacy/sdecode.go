package legacy

/*
#include "defs.h"
#include "GAME1_1.h"
#include "GAME3_2.h"
#include "GAME4.h"
#include "GAME4_3.h"
#include "GAME5_2.h"
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/libs/spell"

	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
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

func Nox_xxx_spellByBookInsert_4FE340(u *server.Object, spells [5]spell.ID, cnt int, a4 int, a5 byte) bool {
	cspells, free := alloc.Make([]C.int{}, len(spells))
	defer free()
	for i, sp := range spells {
		cspells[i] = C.int(sp)
	}
	return C.nox_xxx_spellByBookInsert_4FE340(C.int(uintptr(unsafe.Pointer(asObjectC(u)))), &cspells[0], C.int(cnt), C.int(a4), C.int(a5)) != 0
}

func Nox_xxx_guide_427010(name unsafe.Pointer) int {
	return int(C.nox_xxx_guide_427010((*C.char)(name)))
}
