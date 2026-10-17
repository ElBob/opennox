package legacy

/*
#include <stddef.h>
#include "nox_wchar.h"

extern int nox_cheat_allowall;
extern int nox_cheat_charmall;
*/
import "C"

func CheatEquipAll(v bool) {
	C.nox_cheat_allowall = C.int(bool2int(v))
}

func CheatCharmAll(v bool) {
	C.nox_cheat_charmall = C.int(bool2int(v))
}
