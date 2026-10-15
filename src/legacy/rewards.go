package legacy

/*
#include "server__object__health.h"
#include "server__ability__ability.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/server"
)

var (
	SoloMonsterKillReward func(obj *server.Object)
	AbilityReward         func(u *server.Object, abil int, notify int) bool
	AbilityPrintSuccess   func(code int)
)

//export nox_xxx_soloMonsterKillReward_4EE500_obj_health
func nox_xxx_soloMonsterKillReward_4EE500_obj_health(a1 C.int) {
	SoloMonsterKillReward(asObjectInt(a1))
}

//export nox_xxx_abilityRewardServ_4FB9C0_ability
func nox_xxx_abilityRewardServ_4FB9C0_ability(a1, abil, notify C.int) C.int {
	return C.int(bool2int(AbilityReward(asObjectInt(a1), int(abil), int(notify))))
}

//export nox_xxx_abilGetSuccess_4FB960_ability
func nox_xxx_abilGetSuccess_4FB960_ability(code C.int) {
	AbilityPrintSuccess(int(code))
}
