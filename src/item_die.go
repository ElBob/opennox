package opennox

import (
	"fmt"
	"strings"

	"github.com/opennox/libs/strman"
	"github.com/opennox/libs/types"

	"github.com/opennox/opennox/v1/common/sound"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func itemDiePos(obj *server.Object) (*server.Object, types.Pointf) {
	if owner := obj.InvHolder; owner != nil {
		return owner, owner.PosVec
	}
	return nil, obj.PosVec
}

func (s *Server) dieWeapon(obj *server.Object) { // nox_xxx_dieWeapon_54E370_obj_die
	owner, pos := itemDiePos(obj)
	switch {
	case obj.Material&0x10 != 0:
		str := s.Strings().GetStringInFile("WeaponDieMetal", "Die.c")
		legacy.Nox_xxx_netSendLineMessage_4D9EB0(owner, fmt.Sprintf(str, legacy.Nox_xxx_itemGetName_4E77E0(obj)))
		s.Audio.EventPos(sound.SoundMetalWeaponBreak, pos, 0, 0)
	case obj.Material&0x8 != 0:
		str := s.Strings().GetStringInFile("WeaponDieWood", "Die.c")
		legacy.Nox_xxx_netSendLineMessage_4D9EB0(owner, fmt.Sprintf(str, legacy.Nox_xxx_itemGetName_4E77E0(obj)))
		s.Audio.EventPos(sound.SoundWoodWeaponBreak, pos, 0, 0)
	default:
		str := s.Strings().GetStringInFile("WeaponDieGeneric", "Die.c")
		legacy.Nox_xxx_netSendLineMessage_4D9EB0(owner, fmt.Sprintf(str, s.Armor.Sub_415B60(obj)))
	}
	s.DelayedDelete(obj)
}

func (s *Server) dieArmor(obj *server.Object) { // nox_xxx_dieArmor_54E170_obj_die
	plural := false
	if lang := s.Strings().Lang(); lang == 0 || lang == 1 {
		if def := s.Modif.Nox_xxx_equipClothFindDefByTT413270(int(obj.TypeInd)); def != nil {
			name := alloc.GoString16(def.Desc8)
			plural = strings.HasSuffix(name, "S") || strings.HasSuffix(name, "s")
		}
	}
	owner, pos := itemDiePos(obj)
	var (
		id  string
		snd sound.ID
	)
	switch {
	case obj.Material&0x10 != 0:
		id, snd = "ArmorDieMetal", sound.SoundMetalArmorBreak
	case obj.Material&0x8 != 0:
		id, snd = "ArmorDieWood", sound.SoundWoodenArmorBreak
	case obj.Material&0x4 != 0:
		id, snd = "ArmorDieHide", sound.SoundLeatherArmorBreak
	case obj.Material&0x2 != 0:
		id, snd = "ArmorDieCloth", sound.SoundClothArmorBreak
	default:
		// The original code passes an uninitialized value as a sound ID here.
		id = "ArmorDieGeneric"
	}
	if plural && snd != 0 {
		id += "Plural"
	}
	str := s.Strings().GetStringInFile(strman.ID(id), "Die.c")
	legacy.Nox_xxx_netSendLineMessage_4D9EB0(owner, fmt.Sprintf(str, legacy.Nox_xxx_itemGetName_4E77E0(obj)))
	if snd != 0 {
		s.Audio.EventPos(snd, pos, 0, 0)
	}
	s.DelayedDelete(obj)
}
