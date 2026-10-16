package opennox

import (
	"encoding/binary"

	"github.com/opennox/libs/object"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func (s *Server) netVote(u *server.Object, data []byte) (int, bool) { // MSG_VOTE
	if len(data) < 2 {
		return 0, false
	}
	kind := 0
	if noxflags.HasGame(noxflags.GameModeQuest) {
		kind = 3
	}
	switch data[1] {
	case 0, 1, 2, 3:
		if len(data) < 52 {
			return 0, false
		}
		text := data[2:52]
		if data[1]&1 != 0 {
			kind = 1
		}
		if data[1] < 2 {
			legacy.Sub_506870(kind, u, text)
		} else {
			legacy.Sub_506C90(kind, u, text)
		}
		return 52, true
	case 4:
		legacy.Sub_506870(2, u, nil)
		return 2, true
	case 5:
		legacy.Sub_506C90(2, u, nil)
		return 2, true
	}
	return 0, false
}

func (s *Server) netTrade(pl *server.Player, u *server.Object, data []byte) (int, bool) { // MSG_TRADE
	if len(data) < 2 {
		return 0, false
	}
	ud := u.UpdateDataPlayer()
	ts := ud.Trade70
	sub := data[1]
	switch sub {
	case 0x0E, 0x11, 0x12:
		if ts == nil {
			return 2, true
		}
		switch sub {
		case 0x0E:
			legacy.Sub_50F3A0(ts)
		case 0x11:
			legacy.Nox_xxx_tradeAccept_50F5A0(ts, u)
		case 0x12:
			legacy.Nox_xxx_shopExit_50F4C0(ts)
		}
		return 2, true
	}
	if len(data) < 4 {
		return 0, false
	}
	code := binary.LittleEndian.Uint16(data[2:])
	switch sub {
	case 0x0F:
		netcode := s.Objs.DynamicNetCode(code)
		if it := u.InventoryItemByNetCode(netcode); it != nil && ts != nil {
			if legacy.Nox_xxx_tradeP2PAddOffer2_50F820_trade(ts, u, it) {
				legacy.Sub_4ED0C0(u, it)
			}
		}
		return 4, true
	case 0x10:
		netcode := s.Objs.DynamicNetCode(code)
		if ts != nil {
			legacy.Nox_xxx_tradeP2PAddOfferMB_50FE20(ts, netcode)
		}
		return 4, true
	case 0x15:
		if nox_xxx_gameGet_4DB1B0() || pl.Field3680&0x3 != 0 {
			return 4, true
		}
		netcode := s.Objs.DynamicNetCode(code)
		if obj := s.getObjectFromNetCode(netcode); obj != nil && obj.SubClass().AsMonster().Has(object.MonsterShopkeeper) {
			legacy.Nox_xxx_servShopStart_50EF10_trade(u, obj)
		}
		return 4, true
	case 0x16, 0x18, 0x1A, 0x1C, 0x1E:
		if ts == nil {
			return 4, true
		}
		switch sub {
		case 0x16:
			legacy.Sub_5100C0_trade(u, ts, int(code))
		case 0x18:
			legacy.Sub_510BE0_trade(u, ts, int(code))
		case 0x1A:
			legacy.Sub_510AE0(u, ts, int(code))
		case 0x1C:
			legacy.Sub_5109C0_trade(u, ts, int(code))
		case 0x1E:
			legacy.Sub_5108D0(u, ts, int(code))
		}
		return 4, true
	case 0x17, 0x19:
		if len(data) < 5 {
			return 0, false
		}
		if ts == nil {
			return 5, true
		}
		if sub == 0x17 {
			legacy.Sub_510640_trade(u, ts, int(code), int(data[4]))
		} else {
			legacy.Sub_510D10(u, ts, int(code), int(data[4]))
		}
		return 5, true
	}
	return 0, false
}

func (s *Server) netGauntlet(u *server.Object, data []byte) (int, bool) { // MSG_GAUNTLET
	if len(data) < 2 {
		return 0, false
	}
	switch data[1] {
	case 3:
		ud := u.UpdateDataPlayer()
		if pu := ud.Player.PlayerUnit; pu != nil && pu.Flags().Has(object.FlagDead) {
			ud.Field137 = 0
			legacy.Nox_xxx_playerRespawn_4F7EF0(pu)
		}
		return 2, true
	case 27:
		sub_4DD0B0(u)
		return 2, true
	}
	return 0, false
}
