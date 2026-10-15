package legacy

/*
#include "GAME1.h"
#include "GAME4.h"
#include "GAME4_1.h"
#include "server__magic__plyrgide.h"
#include "server__magic__plyrspel.h"
#include "server__magic__spell__execdur.h"
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/libs/types"

	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

var (
	SpellTeleportToTargetStart func(sp *server.DurSpell) bool
	SpellGrantToPlayer         func(u *server.Object, sp, notify, a4, lvl int) bool
	AwardBeastGuide            func(u *server.Object, guide, notify int) bool
	SpellPrintError            func(code int)
)

func asObjectInt(a1 C.int) *server.Object {
	return asObjectS((*nox_object_t)(unsafe.Pointer(uintptr(a1))))
}

//export sub_530A30_spell_execdur
func sub_530A30_spell_execdur(a1 C.int) C.int {
	return C.int(bool2int(SpellTeleportToTargetStart((*server.DurSpell)(unsafe.Pointer(uintptr(a1))))))
}

//export nox_xxx_spellGrantToPlayer_4FB550
func nox_xxx_spellGrantToPlayer_4FB550(u *nox_object_t, sp, notify, a4, lvl C.int) C.int {
	return C.int(bool2int(SpellGrantToPlayer(asObjectS(u), int(sp), int(notify), int(a4), int(lvl))))
}

//export nox_xxx_awardBeastGuide_4FAE80_magic_plyrgide
func nox_xxx_awardBeastGuide_4FAE80_magic_plyrgide(u, guide, notify C.int) C.int {
	return C.int(bool2int(AwardBeastGuide(asObjectInt(u), int(guide), int(notify))))
}

//export nox_xxx_abilGetError_4FB0B0_magic_plyrspel
func nox_xxx_abilGetError_4FB0B0_magic_plyrspel(code C.int) {
	SpellPrintError(int(code))
}

func Nox_xxx_mapTileAllowTeleport_411A90(pos types.Pointf) bool {
	cpos, free := alloc.New(types.Pointf{})
	defer free()
	*cpos = pos
	return C.nox_xxx_mapTileAllowTeleport_411A90((*C.float2)(unsafe.Pointer(cpos))) != 0
}

func Nox_xxx_netSendRewardNotify_4FAD50(u *server.Object, kind int, who *server.Object, id int) {
	C.nox_xxx_netSendRewardNotify_4FAD50(C.int(uintptr(unsafe.Pointer(asObjectC(u)))), C.int(kind), C.int(uintptr(unsafe.Pointer(asObjectC(who)))), C.char(id))
}

func Nox_xxx_shopExit_50F4C0(t *server.TradeSession) {
	C.nox_xxx_shopExit_50F4C0((*C.uint32_t)(unsafe.Pointer(t)))
}
