package legacy

/*
#include "defs.h"
#include "client__gui__window.h"
#include "GAME3_3.h"
#include "MixPatch.h"
int sub_4BDFD0();
int  modifyWndInputHandler(int a1, int a2, int a3, int a4);
int  nox_xxx_clientUpdateButtonRow_45E110(int a1);
unsigned int*  nox_xxx_objGetTeamByNetCode_418C80(int a1);
void  nox_xxx_printCentered_445490(wchar2_t* a1);
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/libs/types"

	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func Sub_4BDFD0() {
	C.sub_4BDFD0()
}

var GameexPickupTrap func(u, item *server.Object)

//export OnLibraryNotice_420
func OnLibraryNotice_420(a1, a2, a3, a4 C.uint32_t) {
	GameexPickupTrap(asObjectS((*nox_object_t)(unsafe.Pointer(uintptr(a1)))), asObjectS((*nox_object_t)(unsafe.Pointer(uintptr(a2)))))
}

func Nox_xxx_drop_4ED810(u, item *server.Object, pos types.Pointf) {
	cpos, free := alloc.New(types.Pointf{})
	defer free()
	*cpos = pos
	C.nox_xxx_drop_4ED810(C.int(uintptr(unsafe.Pointer(asObjectC(u)))), C.int(uintptr(unsafe.Pointer(asObjectC(item)))), (*C.float)(unsafe.Pointer(cpos)))
}
