package legacy

/*
#include <stdlib.h>
#include "defs.h"
#include "noxstring.h"
#include "client__gui__guicon.h"
#include "client__gui__guimsg.h"
#include "client__network__inform.h"
extern uint32_t dword_5d4594_1203836;
extern uint32_t dword_5d4594_1203840;

// The original passes a variable number of arguments to nox_swprintf (format is loaded from the string table).
// Extra arguments are ignored by cdecl callee, so a fixed number of slots is enough.
static void inform_swprintf(wchar2_t* dst, const wchar2_t* fmt, unsigned int a1, unsigned int a2) {
	nox_swprintf(dst, fmt, a1, a2);
}
static int inform_console_printf(unsigned char cl, wchar2_t* fmt, wchar2_t* a1, wchar2_t* a2) {
	return nox_gui_console_Printf_450C00(cl, fmt, a1, a2);
}
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
)

// Port of client/Network/inform.c (MSG_INFORM handler), client/Network/deathmsg.c and the two
// "get error/success" helpers from plyrspel.c and Ability.c that are only used by the handler.
// Decompiler artefacts are kept: packet and object pointers are passed around as 32-bit ints.

const (
	informSrcFile    = "C:\\NoxPost\\src\\client\\Network\\inform.c"
	deathmsgSrcFile  = "C:\\NoxPost\\src\\client\\Network\\deathmsg.c"
	plyrspelSrcFile  = "C:\\NoxPost\\src\\Server\\Magic\\plyrspel.c"
	abilitySrcFile   = "C:\\NoxPost\\src\\Server\\Ability\\Ability.c"
	informPlayerName = 4704 // offset of the player name in nox_playerInfo
)

func informLoad(name string, file string, line int) *wchar2_t {
	return nox_strman_loadString_40F1D0(internCStr(name), nil, internCStr(file), line)
}

func informLoadP(name *C.char, file string, line int) *wchar2_t {
	return nox_strman_loadString_40F1D0(name, nil, internCStr(file), line)
}

func informPtr(p unsafe.Pointer) uint32 { return uint32(uintptr(p)) }

func informPrint(s *wchar2_t) { C.nox_xxx_printCentered_445490(s) }

// informPrintf is nox_swprintf into buf + nox_xxx_printCentered_445490(buf).
func informSwprintf(buf unsafe.Pointer, format *wchar2_t, a1, a2 uint32) {
	C.inform_swprintf((*C.wchar2_t)(buf), format, C.uint(a1), C.uint(a2))
}

func informPlayerNamePtr(p *nox_playerInfo) uint32 {
	return informPtr(unsafe.Add(unsafe.Pointer(p), informPlayerName))
}

// Nox_xxx_abilGetError_4FB0B0_magic_plyrspel prints the spell error message.
func nox_xxx_abilGetError_4FB0B0_magic_plyrspel(a1 int) {
	v1 := informLoadP(*(**C.char)(memmap.PtrOff(0x587000, uintptr(216380+4*a1))), plyrspelSrcFile, 86)
	informPrint(v1)
}

func nox_xxx_abilGetSuccess_4FB960_ability(a1 int) {
	v1 := informLoadP(*(**C.char)(memmap.PtrOff(0x587000, uintptr(217092+4*a1))), abilitySrcFile, 75)
	informPrint(v1)
}

func Get_dword_5d4594_1203836() uint32 { return uint32(C.dword_5d4594_1203836) }
func Set_dword_5d4594_1203836(v uint32) { C.dword_5d4594_1203836 = C.uint(v) }
func Get_dword_5d4594_1203840() uint32 { return uint32(C.dword_5d4594_1203840) }
func Set_dword_5d4594_1203840(v uint32) { C.dword_5d4594_1203840 = C.uint(v) }

// sub_4952E0 prints the death message to the console.
func sub_4952E0(a1 unsafe.Pointer) int {
	u16 := func(i int) uint16 { return *(*uint16)(unsafe.Add(a1, 2*i)) }
	// wchar2_t v11[32], v12[32], v13[64]
	// NOTE: v11 stays uninitialized in the original when the victim is not found (stack garbage).
	v11 := C.malloc(32 * 2)
	v12 := C.malloc(32 * 2)
	v13 := C.malloc(64 * 2)
	defer C.free(v11)
	defer C.free(v12)
	defer C.free(v13)

	if u16(1) != 0 {
		if v1 := nox_common_playerInfoGetByID_417040(int(u16(1))); v1 != nil {
			v8 := informPlayerNamePtr(v1)
			v2 := informLoad("die.c:LocalizeAttacker", deathmsgSrcFile, 48)
			informSwprintf(v13, v2, v8, 0)
			if u16(2) != 0 {
				v3 := nox_common_playerInfoGetByID_417040(int(u16(2)))
				if v3 != nil {
					fmt := [...]uint16{' ', '+', ' ', '%', 's', 0} // L" + %s"
					informSwprintf(v12, (*wchar2_t)(unsafe.Pointer(&fmt[0])), informPlayerNamePtr(v3), 0)
					C.nox_wcscat((*C.wchar2_t)(v13), (*C.wchar2_t)(v12))
				}
			}
			goto victim
		}
	}
	{
		v9 := informLoad("die.c:AttackerNasty", deathmsgSrcFile, 65)
		v4 := informLoad("die.c:LocalizeAttacker", deathmsgSrcFile, 64)
		informSwprintf(v13, v4, informPtr(unsafe.Pointer(v9)), 0)
	}
victim:
	if u16(3) != 0 {
		v5 := nox_common_playerInfoGetByID_417040(int(u16(3)))
		if v5 != nil {
			v10 := informPlayerNamePtr(v5)
			v6 := informLoad("die.c:LocalizeVictim", deathmsgSrcFile, 72)
			informSwprintf(v11, v6, v10, 0)
		}
	}
	return int(C.inform_console_printf(4, (*wchar2_t)(memmap.PtrOff(0x587000, 161668)), (*C.wchar2_t)(v11), (*C.wchar2_t)(v13))) // NOX_CONSOLE_WHITE
}

func sub_495210(a1 unsafe.Pointer) int {
	v1 := Get_dword_5d4594_1203836()
	if (Get_dword_5d4594_1203836()+1)%100 == Get_dword_5d4594_1203840() {
		Set_dword_5d4594_1203840((Get_dword_5d4594_1203840() + 1) % 100)
	}
	if *(*uint8)(unsafe.Add(a1, 10)) == 1 {
		v2 := int32(*(*uint16)(unsafe.Add(a1, 8)))
		if uint32(uint16(v2)) == *memmap.PtrUint32(0x5D4594, 1203844) { // (unsigned short)v2 == dword
			*(*uint16)(unsafe.Add(a1, 8)) = *(*uint16)(memmap.PtrOff(0x5D4594, 1203856))
			v1 = Get_dword_5d4594_1203836()
		} else if uint32(v2) == *memmap.PtrUint32(0x5D4594, 1203848) {
			*(*uint16)(unsafe.Add(a1, 8)) = *(*uint16)(memmap.PtrOff(0x5D4594, 1203852))
			v1 = Get_dword_5d4594_1203836()
		}
	}
	v3 := uintptr(24 * v1)
	*memmap.PtrUint32(0x5D4594, 1201428+v3) = uint32(*(*uint16)(unsafe.Add(a1, 2)))
	*memmap.PtrUint32(0x5D4594, 1201432+v3) = uint32(*(*uint16)(unsafe.Add(a1, 4)))
	*memmap.PtrUint32(0x5D4594, 1201436+v3) = uint32(*(*uint16)(unsafe.Add(a1, 6)))
	*memmap.PtrUint32(0x5D4594, 1201440+v3) = uint32(*(*uint16)(unsafe.Add(a1, 8)))
	*memmap.PtrUint32(0x5D4594, 1201444+v3) = uint32(*(*uint8)(unsafe.Add(a1, 10)))
	*memmap.PtrUint32(0x5D4594, 1201448+v3) = gameFrame()
	Set_dword_5d4594_1203836((v1 + 1) % 100)
	return sub_4952E0(a1)
}

//export nox_client_handlePacketInform_4C9BF0
func nox_client_handlePacketInform_4C9BF0(a1p int) int {
	a1 := unsafe.Pointer(uintptr(a1p))
	u32 := func(off int) int { return int(*(*uint32)(unsafe.Add(a1, off))) }
	// wchar2_t v62[256], v63[256]
	v62 := C.malloc(256 * 2)
	v63 := C.malloc(256 * 2)
	defer C.free(v62)
	defer C.free(v63)
	sf := informSrcFile

	switch *(*uint8)(unsafe.Add(a1, 1)) {
	case 0:
		nox_xxx_abilGetError_4FB0B0_magic_plyrspel(u32(2))
		return 6
	case 1:
		v50 := nox_xxx_spellTitle_424930(u32(2))
		v10 := informLoad("plyrspel.c:SpellCastSuccess", sf, 117)
		informSwprintf(v62, v10, informPtr(unsafe.Pointer(v50)), 0)
		informPrint((*wchar2_t)(v62))
		return 6
	case 2:
		nox_xxx_abilGetSuccess_4FB960_ability(u32(2))
		return 6
	case 3:
		v13 := nox_common_playerInfoGetByID_417040(u32(2))
		if v13 == nil {
			return 6
		}
		v51 := informPlayerNamePtr(v13)
		v14 := informLoad("netserv.c:PlayerTimeout", sf, 147)
		informSwprintf(v62, v14, v51, 0)
		informPrint((*wchar2_t)(v62))
		return 6
	case 4:
		v15 := nox_common_playerInfoGetByID_417040(u32(2))
		if v15 == nil {
			return 6
		}
		v52 := informPlayerNamePtr(v15)
		v16 := informLoad("objcoll.c:FlagRetrieveNotice", sf, 164)
		informSwprintf(v62, v16, v52, 0)
		informPrint((*wchar2_t)(v62))
		Nox_xxx_clientPlaySoundSpecial_452D80(305, 100)
		return 6
	case 5:
		v18 := nox_common_playerInfoGetByID_417040(u32(2))
		if v18 == nil {
			return 10
		}
		v54 := nox_server_teamTitle_418C20(u32(6))
		v19 := informLoad("objcoll.c:FlagCaptureNotice", sf, 196)
		informSwprintf(v62, v19, informPlayerNamePtr(v18), informPtr(unsafe.Pointer(v54)))
		informPrint((*wchar2_t)(v62))
		Nox_xxx_clientPlaySoundSpecial_452D80(306, 100)
		return 10
	case 6:
		v20 := nox_common_playerInfoGetByID_417040(u32(2))
		if v20 == nil {
			return 10
		}
		v55 := nox_server_teamTitle_418C20(u32(6))
		v21 := informLoad("objcoll.c:FlagPickupNotice", sf, 213)
		informSwprintf(v62, v21, informPlayerNamePtr(v20), informPtr(unsafe.Pointer(v55)))
		informPrint((*wchar2_t)(v62))
		Nox_xxx_clientPlaySoundSpecial_452D80(303, 100)
		return 10
	case 7:
		v22 := nox_common_playerInfoGetByID_417040(u32(2))
		if v22 != nil {
			v56 := nox_server_teamTitle_418C20(u32(6))
			v23 := informLoad("drop.c:FlagDropNotice", sf, 230)
			informSwprintf(v62, v23, informPlayerNamePtr(v22), informPtr(unsafe.Pointer(v56)))
			informPrint((*wchar2_t)(v62))
			Nox_xxx_clientPlaySoundSpecial_452D80(304, 100)
		}
		return 10
	case 8:
		v53 := nox_server_teamTitle_418C20(u32(2))
		v17 := informLoad("update.c:FlagRespawnNotice", sf, 179)
		informSwprintf(v62, v17, informPtr(unsafe.Pointer(v53)), 0)
		informPrint((*wchar2_t)(v62))
		Nox_xxx_clientPlaySoundSpecial_452D80(305, 100)
		return 6
	case 9:
		v38 := nox_xxx_getTeamByID_418AB0(u32(6))
		v39 := nox_common_playerInfoGetByID_417040(u32(2))
		if v39 == nil {
			if v38 != nil {
				v41 := informLoad("objcoll.c:FlagBallUnknownNotice", sf, 374)
				informSwprintf(v62, v41, informPtr(unsafe.Pointer(v38)), 0)
				informPrint((*wchar2_t)(v62))
			}
			return 10
		}
		if v38 == nil {
			return 10
		}
		v46 := informPlayerNamePtr(v39)
		v40 := informLoad("objcoll.c:FlagBallNotice", sf, 365)
		informSwprintf(v62, v40, v46, informPtr(unsafe.Pointer(v38)))
		informPrint((*wchar2_t)(v62))
		return 10
	case 0xA:
		v34 := nox_common_playerInfoGetByID_417040(u32(2))
		v35 := nox_xxx_getTeamByID_418AB0(u32(6))
		if v35 != nil {
			if v34 == nil {
				return 10
			}
			v36 := informLoad("pickup.c:PickUpTeamCrown", sf, 334)
			informSwprintf(v62, v36, informPlayerNamePtr(v34), informPtr(unsafe.Pointer(v35)))
			informPrint((*wchar2_t)(v62))
			return 10
		}
		if v34 == nil {
			return 10
		}
		v37 := informLoad("pickup.c:PickUpCrown", sf, 343)
		informSwprintf(v62, v37, informPlayerNamePtr(v34), 0)
		informPrint((*wchar2_t)(v62))
		return 10
	case 0xB:
		v30 := nox_common_playerInfoGetByID_417040(u32(2))
		v31 := nox_xxx_getTeamByID_418AB0(u32(6))
		if v31 != nil {
			if v30 != nil {
				v32 := informLoad("drop.c:DropTeamCrown", sf, 303)
				informSwprintf(v62, v32, informPlayerNamePtr(v30), informPtr(unsafe.Pointer(v31)))
				informPrint((*wchar2_t)(v62))
				return 10
			}
		} else if v30 != nil {
			v33 := informLoad("drop.c:DropCrown", sf, 312)
			informSwprintf(v62, v33, informPlayerNamePtr(v30), 0)
			informPrint((*wchar2_t)(v62))
			return 10
		}
		return 10
	case 0xC:
		v11 := informLoad("Netserv.c:InObservationMode", sf, 128)
		informPrint(v11)
		if u32(2) == 0 {
			return 6
		}
		v12 := informLoad("Netserv.c:PressJump", sf, 132)
		informPrint(v12)
		return 6
	case 0xD:
		switch u32(2) {
		case 0:
			informPrint(informLoad("atckexec.c:PlayerStunned", sf, 27))
			return 6
		case 1:
			informPrint(informLoad("atckexec.c:PlayerConfused", sf, 31))
			return 6
		case 2:
			informPrint(informLoad("atckexec.c:PlayerPoisoned", sf, 35))
			return 6
		case 3:
			informPrint(informLoad("player.c:TooHeavy", sf, 39))
			return 6
		default:
			return 6
		}
	case 0xE:
		sub_495210(a1)
		return 11
	case 0xF:
		// strlen((const char*)(a1 + 3)) + 1
		v42 := uint32(0)
		for *(*uint8)(unsafe.Add(a1, 3+int(v42))) != 0 {
			v42++
		}
		v42++
		if nox_xxx_gameGetPlayState_4356B0() != 3 {
			return int(v42 + 3)
		}
		if *(*uint8)(unsafe.Add(a1, 2)) != 0 {
			v47 := informLoadP((*C.char)(unsafe.Add(a1, 3)), sf, 397)
			v43 := informLoad("use.c:SignSays", sf, 396)
			informSwprintf(v63, v43, informPtr(unsafe.Pointer(v47)), 0)
			informPrint((*wchar2_t)(v63))
			return int(v42 + 3)
		}
		v44 := informLoadP((*C.char)(unsafe.Add(a1, 3)), sf, 402)
		informPrint(v44)
		return int(v42 + 3)
	case 0x10:
		v49 := nox_server_teamTitle_418C20(u32(2))
		v9 := informLoad("pickup.c:WrongTeam", sf, 105)
		informSwprintf(v62, v9, informPtr(unsafe.Pointer(v49)), 0)
		informPrint((*wchar2_t)(v62))
		return 6
	case 0x11:
		v45 := informLoad("Noxworld.c:ErrChangedClass", sf, 51)
		v6 := informLoad("guiserv.c:Notice", sf, 50)
		nox_xxx_dialogMsgBoxCreate_449A10(nil, v6, v45, 33, nil, nil)
		return 2
	case 0x12:
		v24 := nox_common_playerInfoGetByID_417040(u32(2))
		if v24 == nil {
			return 6
		}
		v57 := informPlayerNamePtr(v24)
		v25 := informLoad("objcoll.c:PlayerExited", sf, 247)
		informSwprintf(v62, v25, v57, 0)
		informPrint((*wchar2_t)(v62))
		return 6
	case 0x13:
		v26 := nox_common_playerInfoGetByID_417040(u32(2))
		if v26 != nil {
			v58 := informPlayerNamePtr(v26)
			v27 := informLoad("objcoll.c:PlayerExitedWarp", sf, 265)
			informSwprintf(v62, v27, v58, 0)
			informPrint((*wchar2_t)(v62))
		}
		return 6
	case 0x14:
		v7 := nox_common_playerInfoGetByID_417040(u32(2))
		if v7 == nil {
			return 6
		}
		v48 := informPlayerNamePtr(v7)
		v8 := informLoad("GeneralPrint:SecretFoundOther", sf, 90)
		informSwprintf(v62, v8, v48, 0)
		informPrint((*wchar2_t)(v62))
		return 6
	case 0x15:
		v59 := uint32(u32(2))
		v28 := informLoad("GeneralPrint:AdvanceToStage1", sf, 281)
		informSwprintf(v62, v28, v59, 0)
		v29 := informLoad("use.c:SignSays", sf, 282)
		informSwprintf(v63, v29, informPtr(v62), 0)
		informPrint((*wchar2_t)(v63))
		return 6
	default:
		return 0
	}
}
