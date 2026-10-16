package opennox

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

// playerCanTouchItems checks if the player is allowed to drop or pick up items.
func playerCanTouchItems(pl *server.Player, u *server.Object) bool {
	ud := u.UpdateDataPlayer()
	return pl.Field3680&0x3 == 0 && ud.Trade70 == nil && ud.DialogWith == nil && !u.Flags().Has(object.FlagNoUpdate)
}

func (s *Server) netTryDropItem(pl *server.Player, u *server.Object, code uint16, pos types.Pointf) { // MSG_TRY_DROP
	netcode := s.Objs.DynamicNetCode(code)
	if !playerCanTouchItems(pl, u) {
		return
	}
	if it := u.InventoryItemByNetCode(netcode); it != nil {
		legacy.Nox_xxx_drop_4ED810(u, it, pos)
	}
}

func (s *Server) netTryGetItem(pl *server.Player, u *server.Object, code uint16) { // MSG_TRY_GET
	netcode := s.Objs.DynamicNetCode(code)
	if nox_xxx_gameGet_4DB1B0() || !playerCanTouchItems(pl, u) {
		return
	}
	it := s.getObjectFromNetCode(netcode)
	if it == nil {
		return
	}
	weight := 0
	for inv := u.FirstItem(); inv != nil; inv = inv.NextItem() {
		weight += int(inv.Weight)
	}
	if weight+int(it.Weight) > int(u.CarryCapacity) {
		s.NetPriMsgToPlayer(u, "pickup.c:CarryingTooMuch", 0)
		return
	}
	gameexPickupTrap(u, it)
}

func (s *Server) netTryUseItem(pl *server.Player, u *server.Object, code uint16) { // MSG_TRY_USE
	netcode := s.Objs.DynamicNetCode(code)
	if nox_xxx_gameGet_4DB1B0() || pl.Field3680&0x3 != 0 {
		return
	}
	if it := u.InventoryItemByNetCode(netcode); it != nil {
		it.CallUse(u)
	}
}

func (s *Server) netTryEquipItem(pl *server.Player, u *server.Object, code uint16) { // MSG_TRY_EQUIP
	netcode := s.Objs.DynamicNetCode(code)
	if nox_xxx_gameGet_4DB1B0() || pl.Field3680&0x3 != 0 {
		return
	}
	if it := u.InventoryItemByNetCode(netcode); it != nil {
		legacy.Nox_xxx_playerTryEquip_4F2F70(u, it)
	}
}

func (s *Server) netTryDequipItem(pl *server.Player, u *server.Object, code uint16) { // MSG_TRY_DEQUIP
	netcode := s.Objs.DynamicNetCode(code)
	if pl.Field3680&0x3 != 0 {
		return
	}
	it := u.InventoryItemByNetCode(netcode)
	if it == nil {
		return
	}
	// Crossbow cannot be removed while it's being used.
	if u.UpdateDataPlayer().State == server.PlayerState1 && it.WeaponClass().Has(object.WeaponCrossbow) {
		return
	}
	legacy.Nox_xxx_playerTryDequip_4F2FB0(u, it)
}

func (s *Server) netTryCollide(pl *server.Player, u *server.Object, code uint16) { // MSG_TRY_COLLIDE
	netcode := s.Objs.DynamicNetCode(code)
	ud := u.UpdateDataPlayer()
	if pl.Field3680&0x3 != 0 || ud.Trade70 != nil || ud.DialogWith != nil {
		return
	}
	if obj := s.getObjectFromNetCode(netcode); obj != nil {
		obj.CallCollide(int(uintptr(u.CObj())), 0)
	}
}

func (s *Server) setSecondaryWeapon(u, it *server.Object) { // sub_53AB90
	if u == nil {
		return
	}
	ud := u.UpdateDataPlayer()
	if it != nil {
		if !nox_xxx_playerClassCanUseItem_57B3D0(it, ud.Player.PlayerClass()) || !legacy.Nox_xxx_playerCheckStrength_4F3180(u, it) {
			legacy.Nox_xxx_netSendSecondaryWeapon_4D9670(ud.Player.PlayerIndex(), nil, 1)
		}
	}
	ud.SecondaryWeapon = it
}

func (s *Server) netInventoryFail(u *server.Object, netcode uint16) { // MSG_INVENTORY_FAIL
	if it := u.InventoryItemByNetCode(int(netcode)); it != nil {
		asObjectS(u).forceDropAt(it, u.Pos())
		s.NetPriMsgToPlayer(u, "pickup.c:CarryingTooMuch", 0)
	}
}
