package opennox

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/spell"
	"github.com/opennox/libs/strman"
	"github.com/opennox/libs/things"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/sound"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func init() {
	legacy.SpellTeleportToTargetStart = func(sp *server.DurSpell) bool { return noxServer.spellTeleportToTargetStart(sp) }
	legacy.SpellGrantToPlayer = func(u *server.Object, sp, notify, a4, lvl int) bool {
		return noxServer.spellGrantToPlayer(u, spell.ID(sp), notify, a4, lvl)
	}
	legacy.AwardBeastGuide = func(u *server.Object, guide int, notify bool) bool {
		return noxServer.awardBeastGuide(u, guide, notify)
	}
	legacy.SpellPrintError = spellPrintError
}

// spellTeleportToTargetStart checks if the caster can teleport to the target location.
func (s *Server) spellTeleportToTargetStart(sp *server.DurSpell) bool { // sub_530A30_spell_execdur
	if sp.Target48 == nil {
		sp.Target48 = sp.Caster16
	}
	var ud *server.PlayerUpdateData
	if c := sp.Caster16; c.Class().Has(object.ClassPlayer) {
		ud = c.UpdateDataPlayer()
	}
	if legacy.Nox_xxx_mapTileAllowTeleport_411A90(sp.Pos2) {
		legacy.Nox_xxx_netSendLineMessage_4D9EB0(sp.Target48, s.Strings().GetStringInFile("UnseenTarget", "ExecDur.c"))
		s.Audio.EventObj(sound.ID(231), sp.Target48, 0, 0)
		return true
	}
	if sp.Caster16 == nil || s.MapTraceRay9(sp.Target48.PosVec, sp.Pos2) {
		if noxflags.HasGame(noxflags.GameModeCoop) {
			dt := float32(s.Balance.FloatInd("TeleportDelay", int(sp.Level)-1))
			sp.Frame68 = s.Frame() + uint32(int32(dt))
		} else {
			sp.Frame68 = s.Frame() + 1
		}
		return false
	}
	if ud != nil {
		s.NetInformTextMsg(ud.Player.PlayerIndex(), 0, 2)
	}
	return true
}

func (s *Server) netSendSpellAward(u *server.Object, sp spell.ID, notify int, a4 int) { // nox_xxx_netSendSpellAward_4D7F90
	if !u.Class().Has(object.ClassPlayer) {
		return
	}
	pl := u.UpdateDataPlayer().Player
	var buf [4]byte
	buf[0] = 111 // MSG_REPORT_SPELL_AWARD
	buf[1] = byte(sp)
	buf[2] = byte(pl.SpellLvl[sp])
	buf[3] = byte(notify)
	if a4 != 0 {
		buf[3] |= 0x80
	}
	s.NetSendPacketXxx1(int(pl.PlayerIndex()), buf[:], 0, 1)
}

func spellPrintError(code int) { // nox_xxx_abilGetError_4FB0B0_magic_plyrspel
	id := alloc.GoString((*byte)(*memmap.PtrPtr(0x587000, 216380+4*uintptr(code))))
	legacy.Nox_xxx_printCentered_445490(strMan.GetStringInFile(strman.ID(id), "plyrspel.c"))
}

func (s *Server) notifyReward(u *server.Object, kind int, id int) {
	legacy.Nox_xxx_netSendRewardNotify_4FAD50(u, kind, u, id)
	if s.Players.CheckXxx(u) {
		return
	}
	for _, it := range s.Players.ListUnits() {
		if it != u {
			legacy.Nox_xxx_netSendRewardNotify_4FAD50(it, kind, u, id)
		}
	}
}

func (s *Server) spellGrantToPlayer(u *server.Object, sp spell.ID, notify, a4, lvl int) bool { // nox_xxx_spellGrantToPlayer_4FB550
	if !u.Class().Has(object.ClassPlayer) {
		return false
	}
	if sp <= 0 || sp >= 137 {
		legacy.Nox_xxx_netSendLineMessage_4D9EB0(u, s.Strings().GetStringInFile("AwardSpellError", "plyrspel.c"))
		return false
	}
	ud := u.UpdateDataPlayer()
	pl := ud.Player
	lvls := &pl.SpellLvl
	quest := noxflags.HasGame(noxflags.GameModeQuest)
	if noxflags.HasGame(noxflags.GameModeCoop|noxflags.GameModeQuest) && lvls[sp] == 3 || lvls[sp] == 5 {
		legacy.Nox_xxx_netSendLineMessage_4D9EB0(u, s.Strings().GetStringInFile("MaxSpellLevel", "plyrspel.c"))
		return false
	}
	if quest && lvls[sp] != 0 {
		switch sp {
		case 34, 134, 45, 46, 47, 48, 49, 117, 118, 119, 120, 121, 122, 123, 124, 125, 19:
			legacy.Nox_xxx_netSendLineMessage_4D9EB0(u, s.Strings().GetStringInFile("MaxSpellLevel", "plyrspel.c"))
			return false
		}
	}
	lvls[sp] = min(lvls[sp]+1, 5)
	if quest {
		lvls[sp] = min(lvls[sp], 3)
	}
	if lvl != 0 {
		lvls[sp] = uint32(lvl)
	}
	legacy.Nox_xxx_playerAwardSpellProtectionCRC_56FCE0(pl.Prot4636, int(sp), int(lvls[sp]))
	var group things.SpellFlags
	switch {
	case s.Spells.HasFlags(sp, 0x1000):
		group = 0x2000
	case s.Spells.HasFlags(sp, 0x4000):
		group = 0x8000
	case s.Spells.HasFlags(sp, 0x10000):
		group = 0x20000
	}
	if group != 0 {
		for i := spell.ID(1); i < 137; i++ {
			if !s.Spells.HasFlags(i, group) || !s.Spells.DefByInd(i).IsValid() {
				continue
			}
			if lvl != 0 {
				lvls[i] = uint32(lvl)
			} else {
				lvls[i]++
			}
			lvls[i] = min(lvls[i], 5)
			// The original code caps the level of the awarded spell, not the one in the group.
			if quest {
				lvls[sp] = min(lvls[sp], 3)
			}
			legacy.Nox_xxx_playerAwardSpellProtectionCRC_56FCE0(pl.Prot4636, int(sp), int(lvls[i]))
		}
	}
	if notify != 0 {
		s.Audio.EventObj(sound.ID(226), u, 0, 0)
		show := !(noxflags.HasGame(noxflags.GameModeCoop) && (sp == 34 || s.Spells.HasFlags(sp, 86016)))
		if (!quest || pl.Field4792 != 0) && show {
			s.notifyReward(u, 0, int(sp))
		}
	}
	if noxflags.HasGame(noxflags.GameModeCoop) && notify == 1 && a4 == 1 {
		if t := ud.Trade70; t != nil {
			legacy.Nox_xxx_shopExit_50F4C0(t)
		}
	}
	s.netSendSpellAward(u, sp, notify, a4)
	return true
}

func (s *Server) awardBeastGuide(u *server.Object, guide int, notify bool) bool { // nox_xxx_awardBeastGuide_4FAE80_magic_plyrgide
	if !u.Class().Has(object.ClassPlayer) {
		return false
	}
	if guide <= 0 || guide >= 41 {
		legacy.Nox_xxx_netSendLineMessage_4D9EB0(u, s.Strings().GetStringInFile("AwardGuideError", "PlyrGide.c"))
		return false
	}
	pl := u.UpdateDataPlayer().Player
	if pl.BeastScrollLvl[guide] != 0 {
		return false
	}
	pl.BeastScrollLvl[guide] = 1
	legacy.Nox_xxx_playerAwardSpellProtectionCRC_56FCE0(pl.Prot4640, guide, 1)
	if notify {
		s.Audio.EventObj(sound.ID(227), u, 0, 0)
		legacy.Nox_xxx_netSendRewardNotify_4FAD50(u, 1, u, guide)
	}
	// Some guides also give access to related creatures. Each group starts with the guide ID.
	for i := uintptr(0); ; i++ {
		grp := (*[64]uint32)(*memmap.PtrPtr(0x587000, 216292+4*i))
		if grp == nil {
			break
		}
		if int(grp[0]) != guide {
			continue
		}
		for j := 1; grp[j] != 0; j++ {
			pl.BeastScrollLvl[grp[j]] = 1
			legacy.Nox_xxx_playerAwardSpellProtectionCRC_56FCE0(pl.Prot4640, int(grp[j]), 1)
		}
	}
	if notify {
		for _, it := range s.Players.ListUnits() {
			if it != u {
				legacy.Nox_xxx_netSendRewardNotify_4FAD50(it, 1, u, guide)
			}
		}
	}
	s.netReportGuideAward(u, guide, notify)
	return true
}

func (s *Server) netReportGuideAward(u *server.Object, guide int, notify bool) { // nox_xxx_netReportGuideAward_4D8000
	pl := u.UpdateDataPlayer().Player
	var buf [3]byte
	buf[0] = 0xD1 // MSG_REPORT_GUIDE_AWARD
	buf[1] = byte(guide)
	if notify {
		buf[2] = 1
	}
	s.NetSendPacketXxx1(int(pl.PlayerIndex()), buf[:], 0, 1)
}
