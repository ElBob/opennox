package legacy

/*
#include "GAME2_1.h"
#include "client__system__parsecmd.h"
*/
import "C"
import "github.com/opennox/opennox/v1/server"

var ServerHandleClientConsole func(pl *server.Player, typ byte, cmd string) int

//export nox_xxx_serverHandleClientConsole_443E90
func nox_xxx_serverHandleClientConsole_443E90(pl *nox_playerInfo, typ C.char, cmd *wchar2_t) C.int {
	return C.int(ServerHandleClientConsole(asPlayerS(pl), byte(typ), GoWString(cmd)))
}
func Nox_xxx_cmdSayDo_46A4B0(text string, a2 int) {
	cstr, free := CWString(text)
	defer free()
	C.nox_xxx_cmdSayDo_46A4B0(cstr, C.int(a2))
}
