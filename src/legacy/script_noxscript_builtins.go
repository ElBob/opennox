package legacy

/*
#include "defs.h"
*/
import "C"
import (
	"github.com/opennox/noxscript/ns/asm"
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
