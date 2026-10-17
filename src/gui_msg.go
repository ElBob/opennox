package opennox

import (
	"github.com/opennox/libs/noxnet/netmsg"

	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/common/sound"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
)

func init() {
	legacy.Nox_xxx_printCentered = printCentered
	legacy.Nox_xxx_clientPickup = func(dr *client.Drawable) { noxClient.clientPickup(dr) }
}

var pickupGoldTypes [3]int // Gold, QuestGoldPile, QuestGoldChest

// clientPickup asks the server to pick up the item, if the player can carry it.
func (c *Client) clientPickup(dr *client.Drawable) { // nox_xxx_clientPickup_46C140
	if pickupGoldTypes[0] == 0 {
		for i, id := range []string{"Gold", "QuestGoldPile", "QuestGoldChest"} {
			pickupGoldTypes[i] = c.Cli().Things.IndByID(id)
		}
	}
	gold, pile, chest := pickupGoldTypes[0], pickupGoldTypes[1], pickupGoldTypes[2]
	if dr == nil {
		return
	}
	if t := int(dr.TypeIDVal); t == gold || t == pile || t == chest || legacy.Sub_467B00(t, 1) {
		code := dr.NetCodeCli()
		buf := [3]byte{byte(netmsg.MSG_TRY_GET), byte(code), byte(code >> 8)}
		c.srv.NetList.AddToMsgListCli(ntype.PlayerInd(31), netlist.Kind(0), buf[:])
		return
	}
	clientPlaySoundSpecial(sound.ID(925), 100)
	str := c.Strings().GetStringInFile("pickup.c:CarryingTooMuch", "C:\\NoxPost\\src\\Client\\Gui\\GameWin\\gamewin.c")
	printCentered(str)
}
