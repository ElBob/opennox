package opennox

import (
	"fmt"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/sound"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func init() {
	legacy.PlayerGiveExp = func(u *server.Object, xp float32) {
		noxServer.playerGiveExp(u, xp)
	}
	legacy.PlayerCheckLevelUp = func(u *server.Object) {
		noxServer.playerCheckLevelUp(u)
	}
}

func (s *Server) playerGiveExp(u *server.Object, xp float32) { // nox_xxx_plyrGiveExp_4EF3A0_exp_level
	pl := u.UpdateDataPlayer().Player
	u.Experience += xp
	legacy.Sub_56FA40(pl.ProtUnitExperience, xp)
	s.NetReportExperience(u)
	str := s.Strings().GetStringInFile("health.c:gainpoints", "explevel.c")
	legacy.Nox_xxx_netSendLineMessage_4D9EB0(u, fmt.Sprintf(str, int32(int64(xp))))
	s.playerCheckLevelUp(u)
}

func (s *Server) playerCheckLevelUp(u *server.Object) { // sub_4EF2E0_exp_level
	pl := u.UpdateDataPlayer().Player
	if legacy.Nox_xxx_gameGet_4DB1B0() && legacy.Sub_4DB1C0() != nil {
		return
	}
	if s.Balance.FloatInd("XPTable", int(int8(pl.Level))+1) > float64(u.Experience) {
		return
	}
	pl.Level++
	legacy.Sub_56F980(pl.ProtPlayerLevel, 1)
	legacy.Nox_xxx_plrReadVals_4EEDC0(u, 1)
	if noxflags.HasGame(noxflags.GameModeCoop) {
		legacy.Sub_57AF30(u, 0)
	} else {
		s.Audio.EventObj(sound.SoundLevelUp, u, 2, u.NetCode)
		legacy.Nox_xxx_netSendLineMessage_4D9EB0(u, s.Strings().GetStringInFile("LevelUP", "explevel.c"))
	}
}
