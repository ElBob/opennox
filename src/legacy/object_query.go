package legacy

/*
#include "defs.h"
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/libs/object"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
)

// C callers of small object/team queries that already have Go implementations.

func queryObject(p C.int) *server.Object {
	return AsObjectP(unsafe.Pointer(uintptr(p)))
}

func queryTeam(p C.int) *server.ObjectTeam {
	return (*server.ObjectTeam)(unsafe.Pointer(uintptr(p)))
}

//export nox_xxx_inventoryCountObjects_4E7D30
func nox_xxx_inventoryCountObjects_4E7D30(a1 C.int, a2 C.int) C.int {
	obj := queryObject(a1)
	if obj == nil {
		return 0
	}
	cnt := 0
	for it := obj.InvFirstItem; it != nil; it = it.InvNextItem {
		if a2 == 0 || C.int(it.TypeInd) == a2 && !it.Flags().Has(object.FlagDestroyed) {
			cnt++
		}
	}
	return C.int(cnt)
}

//export nox_xxx_inventoryGetFirst_4E7980
func nox_xxx_inventoryGetFirst_4E7980(a1 C.int) C.int {
	return C.int(uintptr(unsafe.Pointer(asObjectC(queryObject(a1).FirstItem()))))
}

//export nox_xxx_servCompareTeams_419150
func nox_xxx_servCompareTeams_419150(a1 C.int, a2 C.int) C.int {
	return C.int(bool2int(queryTeam(a1).SameAs(queryTeam(a2))))
}

//export nox_xxx_servObjectHasTeam_419130
func nox_xxx_servObjectHasTeam_419130(a1 C.int) C.int {
	return C.int(bool2int(queryTeam(a1).Has()))
}

//export nox_xxx_testUnitBuffs_4FF350
func nox_xxx_testUnitBuffs_4FF350(unit *nox_object_t, buff C.char) C.int {
	if buff < 0 {
		return 0
	}
	return C.int(bool2int(asObjectS(unit).HasEnchant(server.EnchantID(buff))))
}

//export nox_xxx_unitGetBuffTimer_4FF550
func nox_xxx_unitGetBuffTimer_4FF550(unit *nox_object_t, buff C.int) C.int {
	if buff < 0 {
		return 0
	}
	return C.int(asObjectS(unit).EnchantDur(server.EnchantID(buff)))
}

//export nox_xxx_buffGetPower_4FF570
func nox_xxx_buffGetPower_4FF570(unit *nox_object_t, buff C.int) C.char {
	obj := asObjectS(unit)
	if obj == nil || buff < 0 || buff >= 32 {
		return 0
	}
	return C.char(obj.BuffsPower[buff])
}

//export nox_xxx_unitIsZombie_534A40
func nox_xxx_unitIsZombie_534A40(a1 C.int) C.int {
	zombie := memmap.PtrUint32(0x5D4594, 2488532)
	vileZombie := memmap.PtrUint32(0x5D4594, 2488536)
	if *zombie == 0 {
		*zombie = uint32(GetServer().S().Types.IndByID("Zombie"))
		*vileZombie = uint32(GetServer().S().Types.IndByID("VileZombie"))
	}
	obj := queryObject(a1)
	if obj == nil {
		return 0
	}
	typ := uint32(obj.TypeInd)
	return C.int(bool2int(typ == *zombie || typ == *vileZombie))
}

//export nox_xxx_monsterIsActionScheduled_50A090
func nox_xxx_monsterIsActionScheduled_50A090(a1 C.int, a2 C.int) C.int {
	// Unlike Object.MonsterActionIsScheduled, the original skips the current (top) action.
	ud := queryObject(a1).UpdateDataMonster()
	for i := int(ud.AIStackInd) - 1; i >= 0; i-- {
		if int32(ud.AIStack[i].Action) == int32(a2) {
			return 1
		}
	}
	return 0
}

//export nox_xxx_monsterActionIsCondition_50A010
func nox_xxx_monsterActionIsCondition_50A010(a1 C.int) C.int {
	// The original accepts any action above ACTION_DEPENDENCY_OR-1, not only valid dependency actions.
	return C.int(bool2int(a1 > 39))
}

//export nox_xxx_unitNeedSync_4E44F0
func nox_xxx_unitNeedSync_4E44F0(a1 *nox_object_t) {
	asObjectS(a1).NeedSync()
}

//export nox_client_drawable_testBuff_4356C0
func nox_client_drawable_testBuff_4356C0(dr *nox_drawable, a2 C.char) C.bool {
	if a2 < 0 {
		return false
	}
	return C.bool(asDrawable(dr).HasEnchant(server.EnchantID(a2)))
}

//export nox_xxx_protectionStringCRC_56FAC0
func nox_xxx_protectionStringCRC_56FAC0(a1 *C.int, a2 C.uint) C.int {
	var res C.int
	if n := int(a2 >> 2); n > 0 {
		for _, v := range unsafe.Slice(a1, n) {
			res ^= v
		}
	}
	return res
}
