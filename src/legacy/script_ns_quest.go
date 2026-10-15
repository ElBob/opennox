package legacy

/*
#include "defs.h"
void nox_xxx_cliBuildJournalString_469BC0();
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/libs/object"

	"github.com/opennox/opennox/v1/server"
)

var Nox_xxx_comJournalEntryAdd_427500 func(u *server.Object, name string, typ uint16)

//export nox_xxx_comJournalEntryAdd_427500
func nox_xxx_comJournalEntryAdd_427500(u *nox_object_t, name *C.char, typ C.short) {
	Nox_xxx_comJournalEntryAdd_427500(asObjectS(u), GoString(name), uint16(typ))
}

//export nox_xxx_journalEntryAdd_427490
func nox_xxx_journalEntryAdd_427490(pl *nox_playerInfo, name *C.char, typ C.short) *C.nox_playerInfo_journal {
	e := asPlayerS(pl).JournalAdd(GoString(name), uint16(typ))
	return (*C.nox_playerInfo_journal)(unsafe.Pointer(e))
}

//export nox_xxx_journalEntryRemove_427590
func nox_xxx_journalEntryRemove_427590(pl *nox_playerInfo, name *C.char) C.int {
	if asPlayerS(pl).JournalRemove(GoString(name)) {
		return 1
	}
	return 0
}

//export nox_xxx_journalUpdateEntry_4276B0
func nox_xxx_journalUpdateEntry_4276B0(pl *nox_playerInfo, name *C.char, typ C.short) {
	asPlayerS(pl).JournalUpdate(GoString(name), uint16(typ))
}

//export sub_4277B0
func sub_4277B0(u *nox_object_t, typ C.ushort) {
	obj := asObjectS(u)
	if obj == nil || !obj.Class().Has(object.ClassPlayer) {
		return
	}
	obj.UpdateDataPlayer().Player.JournalRemoveByType(uint16(typ))
}

func Nox_xxx_cliBuildJournalString_469BC0() {
	C.nox_xxx_cliBuildJournalString_469BC0()
}
