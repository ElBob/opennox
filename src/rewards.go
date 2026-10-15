package opennox

import (
	"fmt"

	"github.com/opennox/libs/object"
	"github.com/opennox/libs/strman"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func init() {
	legacy.SoloMonsterKillReward = func(obj *server.Object) { noxServer.soloMonsterKillReward(obj) }
	legacy.AbilityReward = func(u *server.Object, abil int, notify int) bool { return noxServer.abilityReward(u, abil, notify) }
	legacy.AbilityPrintSuccess = abilityPrintSuccess
}

// unitGiveXP gives experience to the player for killing a creature with a given amount of experience.
func (s *Server) unitGiveXP(u *server.Object, exp float32) float64 { // nox_xxx_unitGiveXP_4EF270
	if float64(u.Experience) >= float64(exp) {
		return 0
	}
	mult := *memmap.PtrFloat32(0x587000, 206148)
	v4 := float64(exp-u.Experience)*float64(mult) + 1.0
	xp := float32(v4)
	u.Experience = float32(v4 + float64(u.Experience))
	legacy.Sub_56FA40(u.UpdateDataPlayer().Player.ProtUnitExperience, xp)
	s.NetReportExperience(u)
	s.playerCheckLevelUp(u)
	return float64(xp)
}

func (s *Server) soloMonsterKillReward(obj *server.Object) { // nox_xxx_soloMonsterKillReward_4EE500_obj_health
	if obj == nil || !noxflags.HasGame(noxflags.GameModeCoop) {
		return
	}
	killer := obj.Obj130
	if killer == nil {
		return
	}
	pl := killer.FindOwnerChainPlayer()
	if pl == nil || !pl.Class().Has(object.ClassPlayer) {
		return
	}
	if killer != pl {
		ok := true
		for it := killer; it != pl && it != nil; it = it.ObjOwner {
			if it.Class().Has(object.ClassMonster) {
				ok = server.Nox_xxx_creatureIsMonitored_500CC0(pl, it)
				break
			}
		}
		if !ok {
			return
		}
	}
	if xp := s.unitGiveXP(pl, obj.Experience); xp > 0 {
		str := s.Strings().GetStringInFile("gainpoints", "health.c")
		legacy.Nox_xxx_netSendLineMessage_4D9EB0(pl, fmt.Sprintf(str, int32(int64(xp))))
	}
}

func abilityPrintSuccess(code int) { // nox_xxx_abilGetSuccess_4FB960_ability
	id := alloc.GoString((*byte)(*memmap.PtrPtr(0x587000, 217092+4*uintptr(code))))
	legacy.Nox_xxx_printCentered_445490(strMan.GetStringInFile(strman.ID(id), "Ability.c"))
}

func (s *Server) netAbilityReport(u *server.Object, abil int, notify int) { // nox_xxx_netAbilityReport_4D8060
	pl := u.UpdateDataPlayer().Player
	var buf [3]byte
	buf[0] = 0xCD // MSG_REPORT_ABILITY_AWARD
	buf[1] = byte(abil)
	buf[2] = byte(pl.SpellLvl[abil])
	if notify != 0 {
		buf[2] |= 0x80
	}
	s.NetSendPacketXxx1(int(pl.PlayerIndex()), buf[:], 0, 1)
}

func (s *Server) abilityReward(u *server.Object, abil int, notify int) bool { // nox_xxx_abilityRewardServ_4FB9C0_ability
	if !u.Class().Has(object.ClassPlayer) {
		return false
	}
	if abil <= 0 || abil >= 6 {
		legacy.Nox_xxx_netSendLineMessage_4D9EB0(u, s.Strings().GetStringInFile("AwardAbilityError", "Ability.c"))
		return false
	}
	pl := u.UpdateDataPlayer().Player
	if pl.SpellLvl[abil] != 0 {
		s.NetPriMsgToPlayer(u, "use.c:HadAbility", 0)
		return false
	}
	// Abilities use the same level array as spells.
	pl.SpellLvl[abil] = 5
	legacy.Nox_xxx_playerAwardSpellProtectionCRC_56FCE0(pl.Prot4636, abil, 5)
	s.netAbilityReport(u, abil, notify)
	if noxflags.HasGame(noxflags.GameModeQuest) {
		legacy.Nox_xxx_netSendRewardNotify_4FAD50(u, 2, u, abil)
		if !s.Players.CheckXxx(u) {
			for _, it := range s.Players.ListUnits() {
				if it != u {
					legacy.Nox_xxx_netSendRewardNotify_4FAD50(it, 2, u, abil)
				}
			}
		}
	}
	return true
}
