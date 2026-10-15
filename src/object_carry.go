package opennox

import (
	"github.com/opennox/libs/object"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

// playerCanCarryItem drops the cheapest item from the player's inventory if there is no space for a new item.
func (s *Server) playerCanCarryItem(u, item *server.Object) { // nox_xxx_playerCanCarryItem_513B00
	if *memmap.PtrUint32(0x5D4594, 2386856) == 0 {
		*memmap.PtrUint32(0x5D4594, 2386856) = uint32(s.Types.IndByID("Glyph"))
	}
	if legacy.Sub_467B00(int(item.TypeInd), 1)-legacy.Get_dword_5d4594_2386848() > 0 {
		return
	}
	glyph := int(*memmap.PtrUint32(0x5D4594, 2386856))
	var drop *server.Object
	cost := 999999
	for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
		if it.Class().Has(object.ClassFood) || it.Flags().Has(object.FlagEquipped) ||
			int(it.TypeInd) == glyph || legacy.Nox_xxx_ItemIsDroppable_53EBF0(it) {
			continue
		}
		if c := legacy.Nox_xxx_shopGetItemCost_50E3D0(1, 0, it); c < cost {
			cost, drop = c, it
		}
	}
	if drop == nil {
		return
	}
	pos := legacy.Sub_4ED970(50.0, u.PosVec)
	legacy.Nox_xxx_drop_4ED790(u, drop, pos)
	if legacy.Get_dword_5d4594_2386852() == 0 {
		s.NetPriMsgToPlayer(u, "pickup.c:CarryingTooMuch", 0)
		legacy.Set_dword_5d4594_2386852(1)
	}
}
