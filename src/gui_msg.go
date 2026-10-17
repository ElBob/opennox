package opennox

import (
	"github.com/opennox/libs/noxnet/netmsg"

	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/common/sound"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
)

func init() {
	legacy.Nox_xxx_printCentered = printCentered
	legacy.Nox_xxx_clientPickup = func(dr *client.Drawable) { noxClient.clientPickup(dr) }
}

// clientPickup asks the server to pick up the item, if the player can carry it.
func (c *Client) clientPickup(dr *client.Drawable) { // nox_xxx_clientPickup_46C140
	gold := memmap.PtrUint32(0x5D4594, 1064928)
	pile := memmap.PtrUint32(0x5D4594, 1064932)
	chest := memmap.PtrUint32(0x5D4594, 1064936)
	if *gold == 0 {
		*gold = uint32(c.Cli().Things.IndByID("Gold"))
		*pile = uint32(c.Cli().Things.IndByID("QuestGoldPile"))
		*chest = uint32(c.Cli().Things.IndByID("QuestGoldChest"))
	}
	if dr == nil {
		return
	}
	if t := dr.TypeIDVal; t == *gold || t == *pile || t == *chest || legacy.Sub_467B00(int(t), 1) {
		code := dr.NetCodeCli()
		buf := [3]byte{byte(netmsg.MSG_TRY_GET), byte(code), byte(code >> 8)}
		c.srv.NetList.AddToMsgListCli(ntype.PlayerInd(31), netlist.Kind(0), buf[:])
		return
	}
	clientPlaySoundSpecial(sound.ID(925), 100)
	str := c.Strings().GetStringInFile("pickup.c:CarryingTooMuch", "C:\\NoxPost\\src\\Client\\Gui\\GameWin\\gamewin.c")
	printCentered(str)
}
