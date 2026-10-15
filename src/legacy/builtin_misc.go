package legacy

/*
#include "GAME1_1.h"
#include "GAME2_1.h"
#include "GAME3_3.h"
#include "GAME4_1.h"
#include "GAME4_3.h"
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/libs/types"

	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func Sub_467B00(typ int, a2 int) int {
	return int(C.sub_467B00(C.int(typ), C.int(a2)))
}

func Nox_xxx_ItemIsDroppable_53EBF0(item *server.Object) bool {
	return C.nox_xxx_ItemIsDroppable_53EBF0(C.int(uintptr(unsafe.Pointer(asObjectC(item))))) != 0
}

func Nox_xxx_shopGetItemCost_50E3D0(a1, a2 int, item *server.Object) int {
	// The original code passes the object pointer in a float argument.
	p := uint32(uintptr(unsafe.Pointer(asObjectC(item))))
	return int(C.nox_xxx_shopGetItemCost_50E3D0(C.int(a1), C.int(a2), *(*C.float)(unsafe.Pointer(&p))))
}

func Sub_4ED970(dist float32, pos types.Pointf) types.Pointf {
	cpos, free := alloc.New(types.Pointf{})
	defer free()
	cout, free2 := alloc.New(types.Pointf{})
	defer free2()
	*cpos = pos
	C.sub_4ED970(C.float(dist), (*C.float2)(unsafe.Pointer(cpos)), (*C.float2)(unsafe.Pointer(cout)))
	return *cout
}

func Sub_4277B0(u *server.Object, typ uint16) {
	C.sub_4277B0(asObjectC(u), C.ushort(typ))
}
