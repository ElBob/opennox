package legacy

/*
#include "server__object__die__die.h"
#include "server__object__objutil.h"
#include "GAME4_3.h"
#include "GAME5.h"

static int nox_call_objectType_parseDeath_go(int (*fnc)(char*, void*), char* arg1, void* arg2) { return fnc(arg1, arg2); }
*/
import "C"
import (
	"fmt"
	"strings"
	"unsafe"

	"github.com/opennox/opennox/v1/server"
)

var (
	Nox_xxx_dieWeapon_54E370 func(obj *server.Object)
	Nox_xxx_dieArmor_54E170  func(obj *server.Object)
)

func init() {
	server.RegisterObjectDeath("PlayerDie", C.nox_xxx_diePlayer_54D2B0, 0)
	server.RegisterObjectDeath("PotionDie", C.nox_xxx_diePotion_54CBB0, 0)
	server.RegisterObjectDeath("ImpEggDie", C.nox_xxx_dieImpEgg_54CAE0, 0)
	server.RegisterObjectDeath("GlyphDie", C.nox_xxx_dieGlyph_54DF30, 0)
	server.RegisterObjectDeath("BarrelDie", C.nox_xxx_dieBarrel_54DFA0, 0)
	server.RegisterObjectDeath("CreateObjectDie", C.nox_xxx_dieCreateObject_54E010, 132)
	server.RegisterObjectDeath("SpawnObjectDie", C.nox_xxx_dieSpawnObject_54E070, 132)
	server.RegisterObjectDeath("PolypDie", C.nox_xxx_diePolyp_54CB10, 0)
	server.RegisterObjectDeath("MarkerDie", C.nox_xxx_dieMarker_54E460, 0)
	server.RegisterObjectDeath("WeaponDie", C.nox_xxx_dieWeapon_54E370_obj_die, 0)
	server.RegisterObjectDeath("ArmorDie", C.nox_xxx_dieArmor_54E170_obj_die, 0)
	server.RegisterObjectDeath("BoulderDie", C.nox_xxx_dieBoulder_54E4B0, 0)
	server.RegisterObjectDeath("GameBallDie", C.nox_xxx_dieGameBall_54E620, 0)
	server.RegisterObjectDeath("MonsterGeneratorDie", C.nox_xxx_dieMonsterGen_54E630, 0)

	server.RegisterObjectDeathParse("CreateObjectDie", wrapObjectDeathParseC(C.sub_536B40))
	server.RegisterObjectDeathParse("SpawnObjectDie", wrapObjectDeathParseC(C.sub_536B40))
}

func wrapObjectDeathParseC(ptr unsafe.Pointer) server.ObjectParseFunc {
	return func(objt *server.ObjectType, args []string) error {
		if Nox_call_objectType_parseDeath_go(ptr, strings.Join(args, " "), objt.DeathData) == 0 {
			return fmt.Errorf("cannot parse death data for %q", objt.ID())
		}
		return nil
	}
}

func Nox_call_objectType_parseDeath_go(a1 unsafe.Pointer, a2 string, a3 unsafe.Pointer) int {
	cstr := CString(a2)
	defer StrFree(cstr)
	return int(C.nox_call_objectType_parseDeath_go((*[0]byte)(a1), cstr, a3))
}

//export nox_xxx_dieWeapon_54E370_obj_die
func nox_xxx_dieWeapon_54E370_obj_die(a1 C.int) {
	Nox_xxx_dieWeapon_54E370(asObjectS((*nox_object_t)(unsafe.Pointer(uintptr(a1)))))
}

//export nox_xxx_dieArmor_54E170_obj_die
func nox_xxx_dieArmor_54E170_obj_die(a1 C.int) {
	Nox_xxx_dieArmor_54E170(asObjectS((*nox_object_t)(unsafe.Pointer(uintptr(a1)))))
}

func Nox_xxx_itemGetName_4E77E0(obj *server.Object) string {
	return GoWString(C.nox_xxx_itemGetName_4E77E0_obj_util(C.int(uintptr(unsafe.Pointer(asObjectC(obj))))))
}
