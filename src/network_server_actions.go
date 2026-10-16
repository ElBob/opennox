package opennox

import (
	"encoding/binary"

	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/spell"
	"github.com/opennox/libs/things"
	"github.com/opennox/libs/types"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func (s *Server) playerSetCustomWaypoint(u *server.Object, pos types.Pointf) { // nox_xxx_playerSetCustomWP_4F79A0
	ud := u.UpdateDataPlayer()
	if ud.Player.Field3680&0x3 != 0 {
		return
	}
	i := int(byte(ud.Field45))
	if i >= len(ud.CustomWaypoints) {
		// The original would access memory past the waypoint array here.
		return
	}
	if wp := ud.CustomWaypoints[i]; wp != nil {
		asObjectS(wp).SetPos(pos)
		return
	}
	wp := s.NewObjectByTypeID("PlayerWaypoint")
	ud.CustomWaypoints[i] = wp
	s.CreateObjectAt(wp, u, pos)
}

func (s *Server) netTryCreatureCommand(u *server.Object, code uint16, order byte) { // MSG_TRY_CREATURE_COMMAND
	netcode := s.Objs.DynamicNetCode(code)
	netDebugTestHighBit(code)
	if u.UpdateDataPlayer().Player.Field3680&0x1 != 0 {
		return
	}
	if code == 0 {
		legacy.Nox_xxx_orderUnit_533900(u, nil, uint32(order))
		return
	}
	if obj := s.getObjectFromNetCode(netcode); obj != nil {
		legacy.Nox_xxx_orderUnit_533900(u, obj, uint32(order))
	}
}

func (s *Server) netTrySpell(u *server.Object, data []byte) { // MSG_TRY_SPELL
	pl := u.UpdateDataPlayer().Player
	ok := true
	if pl.Field3680&0x1 != 0 {
		s.NetPriMsgToPlayer(u, "GeneralPrint:NoSpellWarningGeneral", 0)
		ok = false
	}
	if pl.Field3680&0x2 != 0 {
		s.NetPriMsgToPlayer(u, "GeneralPrint:ConjureNoSpellWarning1", 0)
		ok = false
	}
	if !noxflags.HasGame(noxflags.GameModeCoop) && u.Flags().Has(object.FlagAirborne) {
		ok = false
	}
	if noxflags.HasGame(noxflags.GameModeChat) || !ok {
		return
	}
	var spells [5]spell.ID
	cnt := 0
	for i := range spells {
		spells[i] = spell.ID(int32(binary.LittleEndian.Uint32(data[1+4*i:])))
		if spells[i] != 0 {
			cnt++
		}
	}
	// Offensive spells are not cast on a friendly unit under the cursor.
	if cnt == 1 && s.Spells.HasFlags(spells[0], things.SpellOffensive) {
		if targ := u.UpdateDataPlayer().CursorObj; targ != nil && !u.IsEnemyTo(targ) && !noxflags.HasGame(noxflags.GameModeQuest) {
			return
		}
	}
	if legacy.Nox_xxx_spellByBookInsert_4FE340(u, spells, cnt, 3, data[21]) || cnt != 1 {
		return
	}
	for _, sp := range spells {
		if sp != 0 {
			s.NetReportSpellStat(pl.Index(), sp, 0)
		}
	}
}

func (s *Server) tradeItemByNetCode(u *server.Object, code int) *server.Object { // sub_510DE0
	ts := u.UpdateDataPlayer().Trade70
	if ts == nil {
		return nil
	}
	for it := ts.Field20; it != nil; it = it.Field8 {
		if int(it.Item0.NetCode) == code {
			return it.Item0
		}
	}
	return nil
}

func (s *Server) netInfoBookData(u *server.Object, code uint16, typ byte) { // MSG_INFO_BOOK_DATA
	netcode := s.Objs.DynamicNetCode(code)
	netDebugTestHighBit(code)
	it := u.InventoryItemByNetCode(netcode)
	if it == nil {
		it = s.tradeItemByNetCode(u, netcode)
	}
	if it == nil {
		it = s.getObjectFromNetCode(netcode)
	}
	if it == nil || it.UseData.Ptr == nil {
		// The original dereferenced NULL use data here.
		return
	}
	var buf [4]byte
	buf[0] = byte(netmsg.MSG_INFO_BOOK_DATA)
	binary.LittleEndian.PutUint16(buf[1:], uint16(s.GetUnitNetCode(it)))
	if typ == 2 {
		buf[3] = byte(legacy.Nox_xxx_guide_427010(it.UseData.Ptr))
	} else {
		buf[3] = *(*byte)(it.UseData.Ptr)
	}
	s.NetSendPacketXxx0(u.UpdateDataPlayer().Player.Index(), buf[:], 0, 1)
}

func (s *Server) netNewAlias(pl *server.Player, data []byte) { // MSG_NEW_ALIAS
	i := int(data[1])
	if i >= len(pl.NetData16) {
		// The original code would write over the player unit pointer here.
		return
	}
	pl.NetData16[i] = server.PlayerNetData{
		Field0: binary.LittleEndian.Uint16(data[2:]),
		Field2: binary.LittleEndian.Uint16(data[4:]),
		Frame4: binary.LittleEndian.Uint32(data[6:]),
	}
}
