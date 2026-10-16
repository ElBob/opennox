package server

import (
	"fmt"
	"strings"

	"github.com/opennox/libs/object"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

// ItemName returns a full name of the item, including its modifiers.
func (s *Server) ItemName(obj *Object) string { // nox_xxx_itemGetName_4E77E0_obj_util
	cl := obj.Class()
	if !cl.HasAny(object.Class(0x13001000)) {
		return s.sm.GetStringInFile("NoDescription", "objutil.c")
	}
	var def *Modifier
	if cl.HasAny(object.Class(0x11001000)) {
		def = s.Modif.Nox_xxx_getProjectileClassById413250(int(obj.TypeInd))
	} else {
		def = s.Modif.Nox_xxx_equipClothFindDefByTT413270(int(obj.TypeInd))
	}
	if def == nil {
		format := strings.ReplaceAll(s.sm.GetStringInFile("NoInfo", "objutil.c"), "%S", "%s")
		return fmt.Sprintf(format, s.Types.ByInd(int(obj.TypeInd)).ID())
	}
	var mods [4]*ModifierEff
	if obj.InitData != nil {
		mods = *(*[4]*ModifierEff)(obj.InitData)
	}
	name := alloc.GoString16((*uint16)(memmap.PtrOff(0x5D4594, 1567732)))
	for _, m := range mods[:2] {
		if m != nil && m.desc8 != nil {
			name += alloc.GoString16(m.desc8) + " "
		}
	}
	if def.Desc8 != nil {
		name += alloc.GoString16(def.Desc8)
	}
	if m := mods[2]; m != nil && m.desc8 != nil {
		name += " " + alloc.GoString16(m.desc8)
	}
	if m := mods[3]; m != nil && m.secdesc12 != nil {
		name += " " + alloc.GoString16(m.secdesc12)
	}
	return name
}
