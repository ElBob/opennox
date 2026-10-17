package opennox

import (
	"strings"
	"unsafe"

	"github.com/opennox/libs/strman"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func init() {
	legacy.LoadWeaponLookNames = loadWeaponLookNames
	legacy.LoadArmorLookNames = loadArmorLookNames
	legacy.ModifierEffectName = modifierEffectName
	legacy.SetGameLessons = setGameLessons
	legacy.NPCDisplayName = npcDisplayName
}

// loadLookNames loads localized names into a table of {name *uint16, id *byte, _ uint32} entries.
func loadLookNames(done, off uintptr, file string) {
	if *memmap.PtrUint32(0x5D4594, done) != 0 {
		return
	}
	for ; ; off += 12 {
		p := *memmap.PtrPtr(0x587000, off)
		if p == nil {
			break
		}
		id := alloc.GoString((*byte)(p))
		*memmap.PtrPtr(0x587000, off-4) = unsafe.Pointer(alloc.InternCString16(strMan.GetStringInFile(strman.ID(id), file)))
	}
	*memmap.PtrUint32(0x5D4594, done) = 1
}

func loadWeaponLookNames() { // nox_xxx_loadModifyers_4158C0
	loadLookNames(371248, 33396, "WeapLook.c")
}

func loadArmorLookNames() { // nox_xxx_loadLook_415D50
	loadLookNames(371256, 35500, "ArmrLook.c")
}

func modifierEffectName(bit int8) (string, bool) { // sub_413480
	for i := uintptr(0); i < 6; i++ {
		// The original compares a signed char with an unsigned byte, so values >= 0x80 never match.
		if int(bit) == int(*memmap.PtrUint8(0x587000, 27332+20*i)) {
			id := alloc.GoString((*byte)(*memmap.PtrPtr(0x587000, 27344+20*i)))
			return strMan.GetStringInFile(strman.ID(id), "Modifier.c"), true
		}
	}
	return "", false
}

func setGameLessons(mode uint16, v uint16) { // sub_409FB0_settings
	ptr := memmap.PtrUint16(0x5D4594, 3488+2*uintptr(legacy.Sub_409A70(int(int16(mode)))))
	if *ptr == v {
		return
	}
	v = min(v, 999)
	*ptr = v
	legacy.Nox_server_gameSettingsUpdated_40A670()
	if nox_client_isConnected() {
		legacy.Sub_440A20_Int(strMan.GetStringInFile("parsecmd.c:lessonsset", "settings.c"), int(v))
	}
}

// npcDisplayName returns a localized name of the NPC based on its script name or type.
func npcDisplayName(obj *server.Object) string { // sub_4E39F0_obj_db
	name := obj.ID()
	if obj.IDPtr == nil {
		name = noxServer.Types.ByInd(int(obj.TypeInd)).ID()
	}
	if i := strings.LastIndexByte(name, ':'); i >= 0 {
		name = name[i+1:]
	}
	// The original builds the key in a static buffer.
	buf := memmap.Slice(0x5D4594, 1563460)
	n := copy(buf, "NPC:")
	for i := 0; i < len(name); i++ {
		if c := name[i]; c != '_' {
			buf[n] = c
			n++
		}
	}
	buf[n] = 0
	return strMan.GetStringInFile(strman.ID(string(buf[:n])), "objdb.c")
}
