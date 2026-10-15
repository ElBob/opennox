package opennox

import (
	"unsafe"

	"github.com/opennox/libs/strman"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

const creatureGuidesCnt = 40

// loadCreatureGuides fills the creature guide table that is used by the client and other legacy code.
func loadCreatureGuides() bool { // nox_xxx_loadGuides_427070
	for i := 0; i < creatureGuidesCnt; i++ {
		name := alloc.GoString((*byte)(*memmap.PtrPtr(0x587000, 70504+4*uintptr(i))))
		th := noxClient.Things.TypeByID(name)
		if th == nil {
			return false
		}
		off := 740104 + 28*uintptr(i)
		setPtr := func(d uintptr, p unsafe.Pointer) {
			*memmap.PtrPtr(0x5D4594, off+d) = p
		}
		setPtr(0, unsafe.Pointer(alloc.InternCString16(strMan.GetStringInFile(strman.ID("creature:"+name), "ComGuide.c"))))
		if alloc.GoString(th.Name) == "Bomber" {
			*memmap.PtrUint32(0x5D4594, off+4) = 0
		} else {
			*memmap.PtrUint32(0x5D4594, off+4) = uint32(noxClient.Things.IndByID(name))
		}
		setPtr(8, unsafe.Pointer(alloc.InternCString16(strMan.GetStringInFile(strman.ID("creature_desc:"+name), "ComGuide.c"))))
		setPtr(12, unsafe.Pointer(nox_xxx_gLoadImg("CreatureCage"+name).C()))
		setPtr(16, unsafe.Pointer(nox_xxx_gLoadImg("Spellbook"+name).C()))
		*memmap.PtrUint32(0x5D4594, off+20) = 0
		switch sub := th.ObjSubClass; {
		case sub&1 != 0:
			*memmap.PtrUint8(0x5D4594, off+24) = 1
		case sub&2 != 0:
			*memmap.PtrUint8(0x5D4594, off+24) = 2
		default:
			*memmap.PtrUint8(0x5D4594, off+24) = 4
		}
	}
	return true
}
