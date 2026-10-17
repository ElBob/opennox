package legacy

/*
#include "defs.h"
int sub_435F60();
int sub_40AA60(int a1);
int* sub_416770(int a1, wchar2_t* a2, const char* a3);
int* sub_48D7B0();
int* sub_455800();
nox_map_list_item* nox_common_maplist_first_4D09B0();
nox_map_list_item* nox_common_maplist_next_4D09C0(nox_map_list_item* list);
void sub_4EF410(int a1, unsigned char a2);
int sub_4ADA40();
char* nox_server_currentMapGetFilename_409B30();
void nox_xxx_motd_4467F0();
int sub_40A710(int a1);
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/server"
)

func Sub_435F60() { C.sub_435F60() }

func Sub_40AA60(a1 int) { C.sub_40AA60(C.int(a1)) }

// Sub_416770 adds a name (and optionally a serial) to the ban list.
func Sub_416770(a1 int, name string, serial *byte) {
	cname, free := CWString(name)
	defer free()
	C.sub_416770(C.int(a1), cname, (*C.char)(unsafe.Pointer(serial)))
}

func Sub_48D7B0() { C.sub_48D7B0() }

func Sub_455800() { C.sub_455800() }

func Nox_common_maplist_first_4D09B0() *Nox_map_list_item {
	return (*Nox_map_list_item)(unsafe.Pointer(C.nox_common_maplist_first_4D09B0()))
}

func Nox_common_maplist_next_4D09C0(it *Nox_map_list_item) *Nox_map_list_item {
	return (*Nox_map_list_item)(unsafe.Pointer(C.nox_common_maplist_next_4D09C0((*nox_map_list_item)(unsafe.Pointer(it)))))
}

func Sub_4EF410(u *server.Object, lvl byte) {
	C.sub_4EF410(C.int(uintptr(unsafe.Pointer(asObjectC(u)))), C.uchar(lvl))
}

func Sub_4ADA40() { C.sub_4ADA40() }

func Nox_server_currentMapGetFilename_409B30() { C.nox_server_currentMapGetFilename_409B30() }

// CmdShowMotd shows the message of the day (nox_xxx_motd_4467F0); available in all builds.
func CmdShowMotd() { C.nox_xxx_motd_4467F0() }

// CmdRateForQuality returns the network rate for a connection quality (sub_40A710); available in all builds.
func CmdRateForQuality(a1 int) int { return int(C.sub_40A710(C.int(a1))) }
