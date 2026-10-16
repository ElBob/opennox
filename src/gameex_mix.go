package opennox

import (
	"unsafe"

	"github.com/opennox/libs/strman"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/sound"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func gameexCanUseWeapon(u, it *server.Object) bool {
	pl := u.UpdateDataPlayer().Player
	if f := noxServer.Weapons.Nox_xxx_weaponInventoryEquipFlags_415820(it); f == 0 || f == 2 {
		return false
	}
	return nox_xxx_playerClassCanUseItem_57B3D0(it, pl.PlayerClass()) && legacy.Nox_xxx_playerCheckStrength_4F3180(u, it)
}

// gameexWeaponRoll switches to the next or previous weapon in the inventory.
func gameexWeaponRoll(u *server.Object, next bool) bool { // mix_MouseKeyboardWeaponRoll
	ud := u.UpdateDataPlayer()
	if ud.Player.Field3680&3 != 0 || ud.State == 1 {
		return false
	}
	step := func(it *server.Object) *server.Object {
		if next {
			return it.InvNextItem
		}
		return it.Field125
	}
	if cur := (*server.Object)(unsafe.Pointer(uintptr(ud.Field26))); cur != nil {
		for it := step(cur); it != nil; it = step(it) {
			if gameexCanUseWeapon(u, it) {
				return legacy.Nox_xxx_playerTryDequip_4F2FB0(u, cur) && legacy.Nox_xxx_playerTryEquip_4F2F70(u, it)
			}
		}
		return false
	}
	for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
		if gameexCanUseWeapon(u, it) {
			return legacy.Nox_xxx_playerTryEquip_4F2F70(u, it)
		}
	}
	return false
}

// isGameexTrap checks the class byte that the original code uses to detect traps.
func isGameexTrap(it *server.Object) bool {
	return byte(uint32(it.ObjClass)>>16) == 17
}

// gameexDropTrapFor drops the first trap from the player's inventory.
func gameexDropTrapFor(u *server.Object) bool { // playerDropATrap
	if u == nil {
		return false
	}
	ud := u.UpdateDataPlayer()
	pos := ud.Player.Pos3632()
	if ud.Player.Field3680&3 != 0 || ud.State == 1 {
		return false
	}
	for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
		if isGameexTrap(it) {
			legacy.Nox_xxx_drop_4ED810(u, it, pos)
			return true
		}
	}
	return false
}

// gameexPickupTrap only allows warriors to pick up certain traps.
func gameexPickupTrap(u, item *server.Object) { // OnLibraryNotice_420
	cl := byte(u.UpdateDataPlayer().Player.PlayerClass())
	if !isGameexTrap(item) {
		nox_xxx_inventoryServPlace_4F36F0(u, item, 1, 1)
		return
	}
	switch t := byte(item.TypeInd); {
	case t == 0x6A && cl == 1, t != 0x6A && !((t == 0x6B || t == 0x6D) && cl != 0):
		nox_xxx_inventoryServPlace_4F36F0(u, item, 1, 1)
		return
	}
	id := alloc.GoString((*byte)(memmap.PtrOff(0x587000, 215732))) // pickup.c:ObjectEquipClassFail
	noxServer.NetPriMsgToPlayer(u, strman.ID(id), 0)
	noxServer.Audio.EventObj(sound.ID(925), u, 2, u.NetCode)
}

// playerNameOEM converts the player name the same way as the original code did with "%S".
func playerNameOEM(pl *server.Player) string {
	var b []byte
	for _, c := range pl.NameFinal {
		if c == 0 {
			break
		}
		b = append(b, byte(c))
	}
	return string(b)
}

// gameexFindPlayer finds a player by name stored at buf[2:] and writes its class and team to buf[0:2].
func gameexFindPlayer(buf []byte) (*server.ObjectTeam, bool) { // playerInfoStructParser_0, playerInfoStructParser_1
	name := alloc.GoStringS(buf[2:])
	for _, pl := range noxServer.Players.List() {
		if playerNameOEM(pl) != name {
			continue
		}
		t := nox_xxx_objGetTeamByNetCode_418C80(int(pl.NetCodeVal))
		if t != nil {
			buf[1] = byte(t.ID)
		}
		buf[0] = byte(pl.PlayerClass())
		return t, true
	}
	return nil, false
}

// gameexFlagValue returns 2 to the power of i, using integer division for negative powers.
func gameexFlagValue(i int) uint32 { // getFlagValueFromFlagIndex
	v := uint32(i)
	if i < 0 {
		v = uint32(-i)
	}
	r, b := int32(1), int32(2)
	for {
		if v&1 != 0 {
			r *= b
		}
		v >>= 1
		if v == 0 {
			break
		}
		b *= b
	}
	if i >= 0 {
		return uint32(r)
	} else if r == 0 {
		return 0 // the original code divides by zero here
	}
	return uint32(1 / r)
}
