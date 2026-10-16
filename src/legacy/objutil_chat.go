package legacy

/*
#include "GAME1.h"
#include "server__object__objutil.h"
#include "server__script__script.h"
*/
import "C"
import (
	"unicode/utf16"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/server"
)

var ScriptCallbackReadWrite func(cf *cryptfile.CryptFile, cb *server.ScriptCallback, name *byte) bool

//export nox_xxx_itemGetName_4E77E0_obj_util
func nox_xxx_itemGetName_4E77E0_obj_util(a1 C.int) *C.wchar2_t {
	name := GetServer().S().ItemName(asObjectS((*nox_object_t)(unsafe.Pointer(uintptr(a1)))))
	// The original code returns a static buffer.
	buf := unsafe.Slice((*uint16)(memmap.PtrOff(0x5D4594, 1565660)), (1567732-1565660)/2)
	n := copy(buf[:len(buf)-1], utf16.Encode([]rune(name)))
	buf[n] = 0
	return (*C.wchar2_t)(unsafe.Pointer(&buf[0]))
}

//export nox_xxx_xferReadScriptHandler_4F5580
func nox_xxx_xferReadScriptHandler_4F5580(a1 C.int, name *C.char) C.int {
	cb := (*server.ScriptCallback)(unsafe.Pointer(uintptr(a1)))
	return C.int(bool2int(ScriptCallbackReadWrite(cryptfile.Global(), cb, (*byte)(unsafe.Pointer(name)))))
}

func Nox_xxx_cliCanTalkMB_4100F0(s string) bool {
	cstr, free := CWString(s)
	defer free()
	return C.nox_xxx_cliCanTalkMB_4100F0((*C.short)(unsafe.Pointer(cstr))) != 0
}
