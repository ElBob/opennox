package opennox

import (
	"github.com/opennox/libs/object"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

// monstersAllBelongToHost gives the host all the creatures that were saved together with the game.
func (s *Server) monstersAllBelongToHost() { // nox_xxx_monstersAllBelongToHost_4DB6A0
	pl := s.Players.ByInd(server.HostPlayerIndex)
	if pl == nil || pl.PlayerUnit == nil {
		return
	}
	if memmap.Uint32(0x5D4594, 1563124) == 0 {
		*memmap.PtrUint32(0x5D4594, 1563124) = uint32(s.Types.IndByID("SaveGameLocation"))
	}
	typ := int(memmap.Uint32(0x5D4594, 1563124))
	var loc *server.Object
	for it := s.Objs.First(); it != nil; it = it.Next() {
		if int(it.TypeInd) == typ {
			loc = it
			break
		}
	}
	if loc == nil {
		return
	}
	for it := loc.Field129; it != nil; {
		next := it.Field128
		s.ObjSetOwner(pl.PlayerUnit, it)
		if it.Class().Has(object.ClassMonster) && it.UpdateDataMonster().StatusFlags.Has(object.MonStatusSummoned) {
			it.ObjSubClass |= 0x80
			legacy.Nox_xxx_netReportAcquireCreature_4D91A0(int(pl.PlayerIndex()), it)
			s.Players.Nox_xxx_netMarkMinimapObject_417190(pl.PlayerIndex(), it, 1)
		}
		it = next
	}
	loc.ScriptIDVal = 0
	s.DelayedDelete(loc)
}
