package legacy

/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "defs.h"
#include "noxstring.h"
#include "GAME1_1.h"

extern void* nox_monsterBin_head_2386924;
*/
import "C"
import (
	"encoding/binary"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
)

var (
	cStrMonsterBin = C.CString("monster.bin")
	cStrEND        = C.CString("END")
	cStrARENA      = C.CString("ARENA")
	cStrSOLO       = C.CString("SOLO")
	cStrNULL       = C.CString("NULL")
	cStrPlus       = C.CString("+")
)

func mem587(off uintptr) unsafe.Pointer { return memmap.PtrOff(0x587000, off) }

// Nox_xxx_loadMonsterBin_517010 is the port of nox_xxx_loadMonsterBin_517010.
// The monster definition list head stays in the C global nox_monsterBin_head_2386924.
func Nox_xxx_loadMonsterBin_517010() int {
	var v2 [256]byte
	C.nox_monsterBin_head_2386924 = nil
	f := nox_binfile_open_408CC0(cStrMonsterBin, 0)
	result := int(uintptr(unsafe.Pointer(f)))
	if result != 0 {
		result = nox_binfile_cryptSet_408D40(f, 23)
		// Quirk: the file is not closed when the key cannot be set.
		if result != 0 {
			for nox_xxx_readStr_517090(f, unsafe.Pointer(&v2[0])) != 0 &&
				nox_xxx_servParseMonsterDef_517170(f, (*C.char)(unsafe.Pointer(&v2[0]))) != 0 {
			}
			nox_binfile_close_408D90(f)
			result = 1
		}
	}
	return result
}

func cIsSpace(c byte) bool {
	// isspace in the C locale
	return c == ' ' || (c >= 9 && c <= 13)
}

// nox_xxx_readStr_517090 reads a whitespace-separated token, treating // as a comment start.
// Quirk: the destination buffer size is never checked.
func nox_xxx_readStr_517090(a1 *C.FILE, a2 unsafe.Pointer) int {
	v2 := a2
	var v3 uint32
	var charType [4]byte // uint16 CharType[2] in C, zeroed
	v4 := 1
	for {
		for {
			v5 := v3
			nox_binfile_fread_408E40((*C.char)(unsafe.Pointer(&charType[0])), 1, 1, a1)
			if nox_binfile_lastErr_409370(a1) == -1 {
				return 0
			}
			v3 = binary.LittleEndian.Uint32(charType[:])
			if cIsSpace(charType[0]) {
				break
			}
			v4 = 0
			if v3 != 47 || v5 != 47 {
				*(*byte)(v2) = byte(v3)
				v2 = unsafe.Add(v2, 1)
			} else {
				sub_517140(a1)
				v2 = a2
				v3 = binary.LittleEndian.Uint32(charType[:])
				v4 = 1
			}
		}
		if v4 == 0 {
			break
		}
	}
	*(*byte)(v2) = 0
	return 1
}

// sub_517140 skips to the end of the line.
func sub_517140(a1 *C.FILE) int {
	var b byte
	var result int
	for {
		b = 0
		nox_binfile_fread_408E40((*C.char)(unsafe.Pointer(&b)), 1, 1, a1)
		result = nox_binfile_lastErr_409370(a1)
		if !(result != -1 && b != 10) {
			break
		}
	}
	return result
}

func nox_xxx_servParseMonsterDef_517170(a1 *C.FILE, a2 *C.char) int {
	var v10 [256]byte
	buf := (*C.char)(unsafe.Pointer(&v10[0]))
	v3 := (*server.MonsterDef)(C.calloc(1, 0xF8))
	if v3 == nil {
		return 0
	}
	C.strcpy((*C.char)(unsafe.Pointer(v3)), a2)
	for {
		if nox_xxx_readStr_517090(a1, unsafe.Pointer(buf)) == 0 || C.nox_strcmpi(cStrEND, buf) == 0 {
			v3.Next244 = (*server.MonsterDef)(C.nox_monsterBin_head_2386924)
			C.nox_monsterBin_head_2386924 = unsafe.Pointer(v3)
			return 1
		}
		if nox_common_gameFlags_check_40A5C0(2048) || nox_common_gameFlags_check_40A5C0(0x200000) {
			if C.nox_strcmpi(cStrARENA, buf) != 0 {
				if C.nox_strcmpi(cStrSOLO, buf) != 0 {
					goto LABEL_10
				}
				// Quirk: SOLO is skipped here without skipping the rest of the line.
			} else {
				sub_517140(a1)
			}
			continue
		}
	LABEL_10:
		if nox_common_gameFlags_check_40A5C0(0x2000) {
			if C.nox_strcmpi(cStrSOLO, buf) == 0 {
				sub_517140(a1)
				continue
			}
			if C.nox_strcmpi(cStrARENA, buf) == 0 {
				continue
			}
		}
		// LABEL_14
		v4 := mem587(248192)
		if *(*uint32)(mem587(248192)) != 0 {
			for {
				if C.nox_strcmpi(*(**C.char)(v4), buf) == 0 {
					break
				}
				v5 := *(*uint32)(unsafe.Add(v4, 12))
				v4 = unsafe.Add(v4, 12)
				if v5 == 0 {
					break
				}
			}
		}
		if *(*uint32)(v4) == 0 {
			C.free(unsafe.Pointer(v3))
			return 0
		}
		v6 := unsafe.Add(unsafe.Pointer(v3), *(*uint32)(unsafe.Add(v4, 8)))
		switch *(*uint32)(unsafe.Add(v4, 4)) {
		case 0:
			nox_xxx_readStr_517090(a1, unsafe.Pointer(buf))
			*(*uint32)(v6) = uint32(C.atoi(buf))
			continue
		case 1:
			nox_xxx_readStr_517090(a1, unsafe.Pointer(buf))
			*(*float32)(v6) = float32(C.atof(buf))
			continue
		case 2:
			nox_xxx_readStr_517090(a1, unsafe.Pointer(buf))
			*(*uint32)(v6) = uint32(nox_xxx_utilFindSound_40AF50(buf))
			continue
		case 3:
			nox_xxx_readStr_517090(a1, unsafe.Pointer(buf))
			if nox_xxx_monsterLoadStrikeFn_549040(v3, buf) != 0 {
				continue
			}
			return 0
		case 4:
			nox_xxx_readStr_517090(a1, unsafe.Pointer(buf))
			if nox_xxx_monsterLoadDieFn_5490E0(v3, buf) != 0 {
				continue
			}
			return 0
		case 5:
			nox_xxx_readStr_517090(a1, unsafe.Pointer(buf))
			if nox_xxx_monsterLoadDeadFn_549180(v3, buf) != 0 {
				continue
			}
			return 0
		case 6:
			var v9 uint32
			nox_xxx_readStr_517090(a1, unsafe.Pointer(buf))
			set_bitmask_flags_from_plus_separated_names_423930(buf, &v9, (**C.char)(mem587(247536)))
			*(*uint16)(v6) = uint16(v9)
			continue
		case 7:
			nox_xxx_readStr_517090(a1, v6)
			if C.strcmp(cStrNULL, (*C.char)(v6)) == 0 {
				*(*byte)(v6) = 0
			}
			continue
		case 8:
			nox_xxx_readStr_517090(a1, unsafe.Pointer(buf))
			v3.MeleeAttackDamageType124 = 18
			v7 := 0
			v8 := mem587(247464)
			found := true
			for C.nox_strcmpi(buf, (*C.char)(unsafe.Add(*(*unsafe.Pointer)(v8), 7))) != 0 {
				v8 = unsafe.Add(v8, 4)
				v7++
				if uintptr(v8) >= uintptr(mem587(247536)) {
					found = false
					break
				}
			}
			if found {
				v3.MeleeAttackDamageType124 = uint32(v7)
			}
			// LABEL_27
			if v3.MeleeAttackDamageType124 == 18 {
				return 0
			}
			continue
		default:
			continue
		}
	}
}

func set_bitmask_flags_from_plus_separated_names_423930(input *C.char, bitmask *uint32, allowed_names **C.char) {
	// strtok keeps pointers into the buffer, so it must live on the C heap (Go stacks can move).
	input_copy := (*C.char)(C.malloc(256))
	defer C.free(unsafe.Pointer(input_copy))
	C.strcpy(input_copy, input)
	if C.strncmp(input_copy, cStrNULL, 4) == 0 {
		return
	}
	cur_value := C.strtok(input_copy, cStrPlus)
	for cur_value != nil {
		set_one_bitmask_flag_by_name_4239C0(cur_value, bitmask, allowed_names)
		cur_value = C.strtok(nil, cStrPlus)
	}
}

func set_one_bitmask_flag_by_name_4239C0(name *C.char, bitmask *uint32, allowed_names **C.char) int {
	names := unsafe.Pointer(allowed_names)
	for i := int8(0); *(**C.char)(unsafe.Add(names, int(i)*4)) != nil; i++ {
		if C.nox_strcmpi(*(**C.char)(unsafe.Add(names, int(i)*4)), name) == 0 {
			*bitmask |= uint32(1) << uint(i)
			return 1
		}
	}
	return 0
}

func monsterLoadFn(a1 *server.MonsterDef, a2 *C.char, tbl uintptr, dst *unsafe.Pointer) int {
	if C.nox_strcmpi(cStrNULL, a2) != 0 {
		v3 := *(**C.char)(mem587(tbl))
		v4 := 0
		if *(*uint32)(mem587(tbl)) != 0 {
			v5 := mem587(tbl)
			for C.strcmp(a2, v3) != 0 {
				v3 = *(**C.char)(unsafe.Add(v5, 8))
				v5 = unsafe.Add(v5, 8)
				v4++
				if v3 == nil {
					return 0
				}
			}
			*dst = *(*unsafe.Pointer)(mem587(tbl + 4 + uintptr(8*v4)))
			return 1
		}
		return 0
	}
	*dst = nil
	return 1
}

func nox_xxx_monsterLoadStrikeFn_549040(a1 *server.MonsterDef, a2 *C.char) int {
	return monsterLoadFn(a1, a2, 287096, &a1.MeleeStrikeFunc236)
}

func nox_xxx_monsterLoadDieFn_5490E0(a1 *server.MonsterDef, a2 *C.char) int {
	return monsterLoadFn(a1, a2, 287280, &a1.DieFunc228)
}

func nox_xxx_monsterLoadDeadFn_549180(a1 *server.MonsterDef, a2 *C.char) int {
	return monsterLoadFn(a1, a2, 287192, &a1.DeadFunc232)
}
