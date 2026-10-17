package server

import (
	"encoding/binary"

	"github.com/opennox/libs/noxnet/netmsg"

	"github.com/opennox/opennox/v1/common/ntype"
)

// NetMonitorCreature tells the player to show the creature in the creature monitor.
func (s *Server) NetMonitorCreature(pind ntype.PlayerInd, obj *Object) int { // nox_xxx_netMonitorCreature_4D9250
	var buf [5]byte
	buf[0] = byte(netmsg.MSG_REPORT_MONITOR_CREATURE)
	binary.LittleEndian.PutUint16(buf[1:], uint16(s.GetUnitNetCode(obj)))
	binary.LittleEndian.PutUint16(buf[3:], obj.TypeInd)
	s.NetSendPacketXxx1(int(pind), buf[:5], 0, 1)
	return s.NetReportTotalHealth(pind, obj)
}

// NetUnmonitorCreature tells the player to remove the creature from the creature monitor.
func (s *Server) NetUnmonitorCreature(pind ntype.PlayerInd, obj *Object) int { // nox_xxx_netSendUnMonitorCrea_4D92A0
	var buf [3]byte
	buf[0] = byte(netmsg.MSG_REPORT_UNMONITOR_CREATURE)
	binary.LittleEndian.PutUint16(buf[1:], uint16(s.GetUnitNetCode(obj)))
	return s.NetSendPacketXxx1(int(pind), buf[:3], 0, 1)
}

// NetReportTotalHealth sends current and max health of the object to the player.
func (s *Server) NetReportTotalHealth(pind ntype.PlayerInd, obj *Object) int { // nox_xxx_netReportTotalHealth_4D85C0
	h := obj.HealthData
	if h == nil {
		return 0
	}
	var buf [7]byte
	buf[0] = byte(netmsg.MSG_REPORT_TOTAL_HEALTH)
	binary.LittleEndian.PutUint16(buf[1:], uint16(s.GetUnitNetCode(obj)))
	binary.LittleEndian.PutUint16(buf[3:], h.Cur)
	binary.LittleEndian.PutUint16(buf[5:], h.Max)
	return s.NetSendPacketXxx1(int(pind), buf[:7], 0, 1)
}

// ObjBecomePet makes the creature a pet of the player unit.
func (s *Server) ObjBecomePet(owner, obj *Object) { // nox_xxx_unitBecomePet_4E7B00
	if owner == nil || obj == nil {
		return
	}
	obj.ObjSubClass |= 0x80
	pind := owner.UpdateDataPlayer().Player.PlayerIndex()
	s.NetMonitorCreature(pind, obj)
	s.Players.Nox_xxx_netMarkMinimapObject_417190(pind, obj, 1)
	s.ObjSetOwner(owner, obj)
}

// ObjRemoveMonitors reverts ObjBecomePet.
func (s *Server) ObjRemoveMonitors(owner, obj *Object) { // nox_xxx_monsterRemoveMonitors_4E7B60
	if owner == nil || obj == nil {
		return
	}
	obj.ObjSubClass &^= 0x80
	pind := owner.UpdateDataPlayer().Player.PlayerIndex()
	s.NetUnmonitorCreature(pind, obj)
	s.Players.Nox_xxx_netUnmarkMinimapObj_417300(pind, obj, 1)
	s.ObjClearOwner(obj)
}
