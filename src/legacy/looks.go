package legacy

/*
#include "client__system__parsecmd.h"
#include "common__object__armrlook.h"
#include "common__object__modifier.h"
#include "common__object__weaplook.h"
#include "common__system__settings.h"
#include "server__dbase__objdb.h"

static void nox_sub_440A20_int(wchar2_t* format, int v) { sub_440A20(format, v); }
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

var (
	LoadWeaponLookNames func()
	LoadArmorLookNames  func()
	ModifierEffectName  func(bit int8) (string, bool)
	SetGameLessons      func(mode uint16, v uint16)
	NPCDisplayName      func(obj *server.Object) string
)

//export nox_xxx_loadModifyers_4158C0
func nox_xxx_loadModifyers_4158C0() { LoadWeaponLookNames() }

//export nox_xxx_loadLook_415D50
func nox_xxx_loadLook_415D50() { LoadArmorLookNames() }

//export sub_413480
func sub_413480(bit C.char) *C.wchar2_t {
	s, ok := ModifierEffectName(int8(bit))
	if !ok {
		return nil
	}
	return (*C.wchar2_t)(unsafe.Pointer(alloc.InternCString16(s)))
}

//export sub_409FB0_settings
func sub_409FB0_settings(mode C.short, v C.ushort) { SetGameLessons(uint16(mode), uint16(v)) }

//export sub_4E39F0_obj_db
func sub_4E39F0_obj_db(obj *nox_object_t) *C.wchar2_t {
	return (*C.wchar2_t)(unsafe.Pointer(alloc.InternCString16(NPCDisplayName(asObjectS(obj)))))
}

func Sub_440A20_Int(format string, v int) {
	cstr, free := CWString(format)
	defer free()
	C.nox_sub_440A20_int(cstr, C.int(v))
}
