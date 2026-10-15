package legacy

/*
#include "defs.h"
int nox_script_RetreatLevel_515DF0();
int nox_script_RetreatLevelGroup_515E50();
int nox_script_SetResumeLevel_515E80();
int nox_script_SetResumeLevelGroup_515EE0();
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/noxscript/ns/asm"

	"github.com/opennox/opennox/v1/legacy/common/ccall"
	"github.com/opennox/opennox/v1/server/noxscript"
)

var (
	Nox_script_shouldReadMoreXxx     func(fi asm.Builtin) bool
	Nox_script_shouldReadEvenMoreXxx func(fi asm.Builtin) bool
)

//export nox_script_shouldReadMoreXxx
func nox_script_shouldReadMoreXxx(fi int) C.bool {
	return C.bool(Nox_script_shouldReadMoreXxx(asm.Builtin(fi)))
}

//export nox_script_shouldReadEvenMoreXxx
func nox_script_shouldReadEvenMoreXxx(fi int) C.bool {
	return C.bool(Nox_script_shouldReadEvenMoreXxx(asm.Builtin(fi)))
}

func wrapScriptC(fnc unsafe.Pointer) noxscript.Builtin {
	return func(_ noxscript.VM) int {
		return ccall.CallIntVoid(fnc)
	}
}

func CallScriptBuiltin(fi asm.Builtin) (int, bool) {
	if fi < 0 || int(fi) >= len(noxScriptBuiltins) {
		return 0, false
	}
	fnc := noxScriptBuiltins[fi]
	if fnc == nil {
		return 0, false
	}
	res := fnc(GetServer().NoxScriptC())
	return res, true
}

var noxScriptBuiltins = [asm.BuiltinGetScore + 1]noxscript.Builtin{}
