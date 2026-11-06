package legacy

/*
#include "defs.h"
#include "GAME3_3.h"
#include "GAME4_1.h"
#include "GAME4_2.h"
#include "GAME4_3.h"
#include "server__script__script.h"
void* nox_xxx_monsterDefByTT_517560(int a1);
*/
import "C"
import (
	"unsafe"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/cryptfile"
)

// This file is a literal port of the monster/NPC/AI action save-file serialization code (GAME4_2.c).
// The original code works on raw object memory with offsets, and this port intentionally does the same:
// offsets, local variable widths and order of reads/writes are preserved (see CONTRIBUTING.md).

func xro() bool { return cryptfile.Global().ReadOnly() }

func xgame(f uint32) bool { return noxflags.HasGame(noxflags.GameFlag(f)) }

// xrw is nox_xxx_fileReadWrite_426AC0_file3_fread. The result is ignored by all callers.
func xrw(p unsafe.Pointer, n int) {
	cryptfile.Global().ReadWrite(unsafe.Slice((*byte)(p), n))
}

func xat(p unsafe.Pointer, off int) unsafe.Pointer { return unsafe.Add(p, off) }
func xu8(p unsafe.Pointer, off int) *uint8         { return (*uint8)(unsafe.Add(p, off)) }
func xu16(p unsafe.Pointer, off int) *uint16       { return (*uint16)(unsafe.Add(p, off)) }
func xi16(p unsafe.Pointer, off int) *int16        { return (*int16)(unsafe.Add(p, off)) }
func xu32(p unsafe.Pointer, off int) *uint32       { return (*uint32)(unsafe.Add(p, off)) }
func xi32(p unsafe.Pointer, off int) *int32        { return (*int32)(unsafe.Add(p, off)) }
func xpp(p unsafe.Pointer, off int) *unsafe.Pointer {
	return (*unsafe.Pointer)(unsafe.Add(p, off))
}

// xlobyte is LOBYTE(v) = b
func xlobyte(v *int32, b uint8) { *v = *v&^0xff | int32(b) }

func xstrlen(p unsafe.Pointer) int {
	n := 0
	for *xu8(p, n) != 0 {
		n++
	}
	return n
}

var xEmptyCStr [1]byte

// xname converts a possibly nil C string returned by a name lookup into a valid pointer.
// GUARD: the original passes NULL to strlen() (crash) when a spell index has no name.
func xname(p *C.char) unsafe.Pointer {
	if p == nil {
		return unsafe.Pointer(&xEmptyCStr[0])
	}
	return unsafe.Pointer(p)
}

func xcs(b *[256]byte) *C.char { return (*C.char)(unsafe.Pointer(&b[0])) }

func xint(p unsafe.Pointer) C.int { return C.int(uintptr(p)) }

func xscript(v1 unsafe.Pointer, off int, v3 unsafe.Pointer, soff int) {
	var s unsafe.Pointer
	if v3 != nil {
		s = xat(v3, soff)
	}
	C.nox_xxx_xferReadScriptHandler_4F5580(xint(xat(v1, off)), (*C.char)(s))
}

func xdirIndexed(a1 unsafe.Pointer, v *[2]uint32) {
	C.nox_xxx_xferIndexedDirection_509E20(C.int(*xi16(a1, 124)), (*C.int2)(unsafe.Pointer(v)))
}

func xdirAngle(v *[2]uint32) int {
	return int(C.nox_xxx_xferDirectionToAngle_509E00((*C.uint32_t)(unsafe.Pointer(v))))
}

func xglyphID() *uint32 { return memmap.PtrUint32(0x5D4594, 2487692) }

//export nox_xxx_XFerMonster_528DB0
func nox_xxx_XFerMonster_528DB0(a1p *nox_object_t) int {
	a1 := unsafe.Pointer(a1p)
	var (
		v44, v45, v46, v48, v49, v50, v47 int32
		v52                               [2]uint32
		v53                               [256]byte
	)
	v1 := *xpp(a1, 748)
	v51 := *xu32(a1, 136)
	if *xglyphID() == 0 {
		*xglyphID() = uint32(nox_xxx_getNameId_4E3AA0(internCStr("Glyph")))
	}
	v45 = 64
	xrw(unsafe.Pointer(&v45), 2)
	if int16(v45) > 64 {
		return 0
	}
	if Nox_xxx_mapReadWriteObjData_4F4530(asObjectS(a1p), int(int16(v45))) == 0 {
		return 0
	}
	if !xro() {
		xdirIndexed(a1, &v52)
	}
	xrw(unsafe.Pointer(&v52), 8)
	if int16(v45) >= 3 {
		v3 := *xpp(a1, 756)
		xscript(v1, 1232, v3, 640)
		xscript(v1, 1264, v3, 896)
		xrw(xat(v1, 1220), 2)
		xscript(v1, 1224, v3, 768)
		if int16(v45) >= 31 {
			xscript(v1, 1240, v3, 1024)
			xscript(v1, 1248, v3, 1152)
			xscript(v1, 1256, v3, 1280)
			xscript(v1, 1272, v3, 1408)
			xscript(v1, 1280, v3, 1536)
			xscript(v1, 1288, v3, 1664)
			if int16(v45) >= 52 {
				xscript(v1, 1296, v3, 1792)
			}
		}
	} else {
		C.sub_4F5540(xint(xat(v1, 1232)))
		C.sub_4F5540(xint(xat(v1, 1264)))
		xrw(xat(v1, 1220), 2)
		C.sub_4F5540(xint(xat(v1, 1224)))
	}
	if !xro() {
		// write mode: nox_crypt_IsReadOnly() == 0
		v46 = 0
		xrw(unsafe.Pointer(&v46), 4)
	} else {
		a := int16(xdirAngle(&v52))
		*xi16(a1, 126) = a
		*xi16(a1, 124) = a
		if int16(v45) >= 11 {
			v46 = 0
			xrw(unsafe.Pointer(&v46), 4)
		}
	}
	if int16(v45) >= 31 {
		xrw(xat(v1, 1332), 1)
		if int16(v45) < 51 {
			xrw(unsafe.Pointer(&v46), 2)
			*xu32(v1, 1440) = uint32(uint16(v46))
		} else {
			xrw(xat(v1, 1440), 4)
		}
		xrw(xat(v1, 1352), 4)
		xrw(xat(v1, 1336), 4)
		xrw(xat(v1, 1344), 4)
		xrw(xat(v1, 1312), 4)
		if int16(v45) < 33 {
			nox_xxx_cryptSeekCur_40E0A0(2)
		}
		xrw(xat(v1, 1304), 4)
		*xu32(v1, 1308) = *xu32(v1, 1304)
		if int16(v45) < 34 {
			xrw(xat(v1, 1360), 4)
		}
		xlobyte(&v48, uint8(xstrlen(xat(v1, 1364))))
		xrw(unsafe.Pointer(&v48), 1)
		xrw(xat(v1, 1364), int(uint8(v48)))
		*xu8(v1, int(uint8(v48))+1364) = 0
		if int16(v45) >= 34 {
			if xro() {
				C.memset(xat(v1, 1488), 0, 0x224)
				xrw(unsafe.Pointer(&v44), 4)
				for i := int32(0); i < v44; i++ {
					xrw(unsafe.Pointer(&v48), 1)
					xrw(unsafe.Pointer(&v53[0]), int(uint8(v48)))
					v53[uint8(v48)] = 0
					v23 := nox_xxx_spellNameToN_4243F0(xcs(&v53))
					xrw(xat(v1, 4*v23+1488), 4)
				}
			} else {
				v16 := int32(0)
				v44 = 0
				for k := 0; k < 137; k++ {
					if *xu32(v1, 1488+4*k) != 0 {
						v16++
					}
				}
				v44 = v16
				xrw(unsafe.Pointer(&v44), 4)
				v47p := xat(v1, 1488)
				for v20 := 0; v20 < 137; v20++ {
					if *(*uint32)(v47p) != 0 {
						v21 := xname(nox_xxx_spellNameByN_424870(v20))
						xlobyte(&v46, uint8(xstrlen(v21)))
						xrw(unsafe.Pointer(&v46), 1)
						xrw(v21, int(uint8(v46)))
						xrw(v47p, 4)
					}
					v47p = xat(v47p, 4)
				}
			}
		} else {
			xrw(xat(v1, 1488), 0x224)
		}
		// Pairs of 16-bit values; older versions store them as bytes.
		for _, off := range [4]int{1448, 1456, 1464, 1472} {
			if int16(v45) < 46 {
				xrw(unsafe.Pointer(&v44), 1)
				*xu16(v1, off) = uint16(uint8(v44))
				xrw(unsafe.Pointer(&v44), 1)
				*xu16(v1, off+2) = uint16(uint8(v44))
			} else {
				xrw(xat(v1, off), 2)
				xrw(xat(v1, off+2), 2)
			}
			if int16(v45) <= 32 {
				xrw(unsafe.Pointer(&v49), 4)
			}
		}
		if int16(v45) < 46 {
			xrw(unsafe.Pointer(&v44), 1)
			*xu16(v1, 1480) = uint16(uint8(v44))
			xrw(unsafe.Pointer(&v44), 1)
			*xu16(v1, 1482) = uint16(uint8(v44))
		} else {
			xrw(xat(v1, 1480), 2)
			xrw(xat(v1, 1482), 2)
		}
		c := int16(v45) > 32
		if !c {
			xrw(unsafe.Pointer(&v49), 4)
			c = int16(v45) >= 32
		}
		if c {
			xrw(xat(v1, 1316), 4)
		}
		if int16(v45) >= 33 {
			xrw(xat(v1, 2040), 4)
			xrw(xat(v1, 1320), 4)
			if int16(v45) < 42 {
				xrw(unsafe.Pointer(&v44), 2)
				if uint16(v44) == 0 {
					*xu8(v1, 1445) = 1
				}
			}
			if int16(v45) < 53 {
				xrw(xat(v1, 2044), 4)
				xrw(xat(v1, 2048), 4)
				xrw(xat(v1, 2052), 4)
			} else {
				v24 := xat(v1, 2044)
				for v46 = 3; v46 != 0; v46-- {
					if xro() {
						xrw(unsafe.Pointer(&v44), 1)
						xrw(unsafe.Pointer(&v53[0]), int(uint8(v44)))
						v53[uint8(v44)] = 0
						*xi32(v24, 0) = int32(nox_xxx_spellNameToN_4243F0(xcs(&v53)))
					} else {
						v25 := xname(nox_xxx_spellNameByN_424870(int(*xi32(v24, 0))))
						xlobyte(&v44, uint8(xstrlen(v25)))
						xrw(unsafe.Pointer(&v44), 1)
						xrw(v25, int(uint8(v44)))
					}
					v24 = xat(v24, 4)
				}
			}
		}
		if int16(v45) >= 34 {
			if xro() {
				xrw(unsafe.Pointer(&v44), 1)
				xrw(unsafe.Pointer(&v53[0]), int(uint8(v44)))
				v53[uint8(v44)] = 0
				*xu32(v1, 1360) = uint32(C.nox_xxx_actionNByNameMB_5345F0(xcs(&v53)))
			} else {
				v26 := unsafe.Pointer(C.sub_5345B0(C.int(*xu32(v1, 1360))))
				xlobyte(&v44, uint8(xstrlen(v26)))
				xrw(unsafe.Pointer(&v44), 1)
				xrw(v26, int(uint8(v44)))
			}
		}
	}
	if int16(v45) >= 41 {
		if nox_xxx_XFer_ActionData_529CE0(a1) == 0 {
			return 0
		}
	}
	if int16(v45) >= 42 {
		xrw(xat(v1, 1445), 1)
	}
	if int16(v45) >= 43 && *xu8(a1, 12)&8 != 0 {
		xlobyte(&v44, 0)
		v27 := *xpp(a1, 692)
		if int16(v45) >= 50 {
			xrw(xat(v27, 1716), 4)
		}
		if int16(v45) >= 61 {
			xrw(xat(v27, 1720), 4)
		}
		if int16(v45) >= 48 {
			var b uint8 = uint8(xstrlen(xat(v27, 1684)))
			v47 = int32(b)
			xrw(unsafe.Pointer(&v47), 1)
			xrw(xat(v27, 1684), int(uint8(v47)))
			*xu8(v27, int(uint8(v47))+1684) = 0
		}
		skip := false
		if xro() {
			xrw(unsafe.Pointer(&v44), 1)
			if v27 == nil {
				skip = true
			} else {
				*xu8(v27, 0) = uint8(v44)
			}
		} else {
			if v27 != nil {
				xlobyte(&v44, *xu8(v27, 0))
			}
			xrw(unsafe.Pointer(&v44), 1)
		}
		if !skip && v27 != nil {
			xlobyte(&v44, 0)
			if *xu8(v27, 0) != 0 {
				for {
					if xro() {
						Nox_xxx_XFer_ReadShopItem_52A840(xat(v27, 28*int(uint8(v44))+4), int(int16(v45)))
					} else {
						Nox_xxx_XFer_WriteShopItem_52A5F0(xat(v27, 28*int(uint8(v44))+4))
					}
					xlobyte(&v44, uint8(v44)+1)
					if !(uint8(v44) < *xu8(v27, 0)) {
						break
					}
				}
			}
		}
	}
	// LABEL_137
	if int16(v45) >= 44 {
		xrw(v1, 4)
	}
	if int16(v45) >= 45 {
		v50 = int32(*xu32(a1, 12) & 0x180)
		xrw(unsafe.Pointer(&v50), 4)
		*xu32(a1, 12) |= uint32(v50)
	}
	if int16(v45) >= 49 {
		xrw(*xpp(a1, 556), 2)
	}
	if int16(v45) >= 51 {
		xrw(xat(v1, 1348), 1)
		xrw(xat(v1, 1340), 1)
		xrw(xat(v1, 1444), 1)
		xrw(xat(v1, 2036), 1)
	}
	if *xu8(a1, 12)&0x20 != 0 && int16(v45) >= 54 {
		v28 := uint8(0)
		v29 := xat(v1, 2076)
		xlobyte(&v46, 0)
		for {
			xrw(v29, 3)
			if xro() {
				nox_xxx_setNPCColor_4E4A90(a1p, uint8(v46), v29)
			}
			v28++
			v29 = xat(v29, 3)
			xlobyte(&v46, v28)
			if !(v28 < 6) {
				break
			}
		}
	}
	if int16(v45) >= 55 && *xu8(a1, 12)&0x20 != 0 {
		C.nox_xxx_readNPCVoiceSet_52AD10(xint(a1))
	}
	if int16(v45) >= 62 {
		if C.nox_xxx_XFer_ReadMonsterBuffs_52AAB0((*C.uint32_t)(a1)) == 0 {
			return 0
		}
	}
	if int16(v45) >= 63 && *xu32(a1, 12)&0x80000 != 0 {
		C.nox_xxx_readNPCVoiceSet_52AD10(xint(a1))
	}
	if int16(v45) >= 64 {
		xlobyte(&v47, *xu8(a1, 540))
		xrw(unsafe.Pointer(&v47), 1)
		if xro() && uint8(v47) != 0 {
			C.nox_xxx_setSomePoisonData_4EEA90(xint(a1), C.int(uint8(v47)))
		}
	}
	if xro() {
		// LABEL_164
		if *xu8(v1, 1445) == 0 {
			if *xu32(a1, 16)&0x8000 != 0 && C.nox_xxx_unitIsZombie_534A40(xint(a1)) == 0 {
				*xu32(a1, 16) |= 0x40
			}
		} else {
			if xgame(1) {
				*xu16(*xpp(a1, 556), 4) = 0
				*xu16(*xpp(a1, 556), 0) = 0
			}
			if xro() {
				if *xu32(a1, 16)&0x8000 != 0 && C.nox_xxx_unitIsZombie_534A40(xint(a1)) == 0 {
					*xu32(a1, 16) |= 0x40
				}
			}
		}
	}
	// LABEL_171
	if *xu32(a1, 136) != 0 {
		if !xro() {
			*xu32(a1, 136) = v51
			return 1
		}
		if Nox_xxx_xfer_4F3E30(int(v45), asObjectS(a1p), *xu32(a1, 136)) == 0 {
			return 0
		}
	}
	if xro() {
		if xgame(0x200000) || !Nox_xxx_gameIsSwitchToSolo_4DB240() {
			nox_xxx_monsterOnSpawnSpellcaster_529BC0(a1)
		}
		if xro() {
			if *xu8(a1, 8)&2 != 0 {
				if *xu32(a1, 12)&0x2000 != 0 {
					v32 := 0
					v33 := 0
					for j := asObjectS(a1p).FirstItem(); j != nil; j = j.NextItem() {
						if uint32(*xu16(unsafe.Pointer(j), 4)) == *xglyphID() {
							v32 = 1
						}
					}
					v35 := xat(v1, 2044)
					for k := 0; k < 3; k++ {
						if *xu32(v35, 4*k) != 0 {
							v33++
						}
					}
					if v32 == 0 && v33 != 0 {
						v38 := GetServer().S().NewObjectByTypeID("Glyph")
						if v38 != nil {
							v39 := *xpp(unsafe.Pointer(v38), 692)
							if v33 > 0 {
								for k := 0; k < v33; k++ {
									*xu32(v39, 4*k) = *xu32(v35, 4*k)
								}
							}
							*xu8(v39, 20) = uint8(v33)
							*xu32(v39, 24) = 0
							*xu32(v39, 28) = *xu32(a1, 56)
							*xu32(v39, 32) = *xu32(a1, 60)
						}
						Nox_xxx_inventoryPutImpl_4F3070(asObjectS(a1p), v38, 1)
					}
				}
			}
		}
	}
	*xu32(a1, 136) = v51
	return 1
}

func nox_xxx_monsterOnSpawnSpellcaster_529BC0(a1 unsafe.Pointer) {
	if a1 == nil {
		return
	}
	v1 := *xpp(a1, 748)
	v2 := C.nox_xxx_monsterDefByTT_517560(C.int(*xu16(a1, 4)))
	if v2 == nil {
		return
	}
	if *xu8(v1, 1445) == 0 {
		*xu16(*xpp(a1, 556), 0) = *xu16(v2, 68)
		*xu16(*xpp(a1, 556), 4) = *xu16(v2, 68)
		*xu16(*xpp(a1, 556), 2) = *xu16(v2, 68)
	}
	if *xu8(v1, 1444) == 1 {
		*xu32(v1, 1440) = *xu32(v2, 92)
	} else {
		for i := 0; i < 22; i++ {
			v4 := uint32(1) << i
			if v4&0x19C40 == 0 {
				v5 := *xu32(v1, 1440)
				if v5&v4 != 0 && v4&*xu32(v2, 92) == 0 {
					*xu32(v1, 1440) = v5 &^ v4
				}
				v6 := *xu32(v1, 1440)
				if v6&v4 == 0 && v4&*xu32(v2, 92) != 0 {
					*xu32(v1, 1440) = v4 | v6
				}
			}
		}
	}
	v7 := *xu32(v1, 1440)
	if v7&0x20 == 0 {
		v7 &^= 0x1800 // BYTE1(v7) &= 0xE7
		*xu32(v1, 1440) = v7
	}
	if *xu8(v1, 1340) == 1 {
		*xu32(v1, 1336) = *xu32(v2, 80)
	}
	if *xu8(v1, 1348) == 1 {
		*xu32(v1, 1344) = *xu32(v2, 84)
	}
	if *xu8(v1, 2036) == 1 {
		nox_xxx_monsterAutoSpells_54C0C0(a1)
	}
}

func nox_xxx_AssignIfGreater_52A420(a1 unsafe.Pointer, a2 int32) int32 {
	result := a2 + *xi32(a1, 0)
	if result >= 1 {
		*xi32(a1, 0) = result
	} else {
		*xi32(a1, 0) = 1
	}
	return result
}

// sub_534650 returns the name of an AI action by index, or NULL.
func sub_534650(a1 int32) unsafe.Pointer {
	v1 := int32(0)
	for v1 != a1 {
		v1++
		if v1 >= 72 {
			return nil
		}
	}
	return *memmap.PtrPtr(0x587000, uintptr(261768+4*v1))
}

// nox_xxx_actionByName_534670 returns the index of an AI action by name.
func nox_xxx_actionByName_534670(a1 string) int32 {
	v1 := int32(0)
	v2 := 261768
	for GoStringP(*memmap.PtrPtr(0x587000, uintptr(v2))) != a1 {
		v2 += 4
		v1++
		if v2 >= 262056 {
			return 0
		}
	}
	return v1
}

func sub_52A440(a1, a2 unsafe.Pointer, a3 int32) int {
	var (
		v8, v9, v10 int32
		v11         [256]byte
	)
	// GUARD: the original passes NULL to strcpy() (crash) for an unknown action index.
	if s := sub_534650(int32(*xu32(a2, 0))); s != nil {
		n := xstrlen(s)
		copy(v11[:], unsafe.Slice((*byte)(s), n+1))
	}
	xlobyte(&v10, uint8(xstrlen(unsafe.Pointer(&v11[0]))))
	xrw(unsafe.Pointer(&v10), 1)
	xrw(unsafe.Pointer(&v11[0]), int(uint8(v10)))
	v11[uint8(v10)] = 0
	v3 := nox_xxx_actionByName_534670(GoStringP(unsafe.Pointer(&v11[0])))
	*xi32(a2, 0) = v3
	xlobyte(&v8, *memmap.PtrUint8(0x587000, uintptr(255604+16*v3)))
	xrw(unsafe.Pointer(&v8), 1)
	v4 := int32(0)
	if uint8(v8) > 0 {
		v5 := xat(a2, 4)
		for {
			result := *memmap.PtrUint32(0x587000, uintptr(255608+4*(v4+4*int32(*xu32(a2, 0)))))
			switch result {
			case 0:
				xrw(v5, 8)
			case 1:
				if xro() {
					xrw(v5, 4)
					break
				}
				if p := *xpp(v5, 0); p != nil {
					xrw(xat(p, 44), 4)
				} else {
					v9 = 0
					xrw(unsafe.Pointer(&v9), 4)
				}
			case 2:
				if xro() {
					xrw(v5, 4)
					break
				}
				if p := *xpp(v5, 0); p != nil {
					xrw(p, 4)
				} else {
					v9 = 0
					xrw(unsafe.Pointer(&v9), 4)
				}
			case 3, 4, 6:
				xrw(v5, 4)
			case 5:
				xrw(v5, 4)
				if xro() {
					nox_xxx_AssignIfGreater_52A420(v5, a3)
				}
			case 7:
				xrw(v5, 1)
			default:
				return int(result)
			}
			v4++
			v5 = xat(v5, 8)
			if v4 >= int32(uint8(v8)) {
				xrw(xat(a2, 20), 4)
				return 1
			}
		}
	}
	xrw(xat(a2, 20), 4)
	return 1
}

func nox_xxx_XFer_ActionData_529CE0(a1 unsafe.Pointer) int {
	var (
		v18, v20, v21, v22, v23 int32
		v19                     uint8
		v24                     [256]byte
	)
	v1 := *xpp(a1, 748)
	v21 = 4
	xrw(unsafe.Pointer(&v21), 2)
	if int16(v21) > 4 {
		return 0
	}
	v19 = 1
	if int16(v21) >= 2 {
		v19 = 0
		if xgame(1) && !xgame(0x400000) {
			v19 = 1
		}
		xrw(unsafe.Pointer(&v19), 1)
		if !(v19 != 0 || uint16(v21) == 1) {
			return 1
		}
	}
	// LABEL_67
	v23 = int32(gameFrame())
	xrw(unsafe.Pointer(&v23), 4)
	v3 := int32(gameFrame()) - v23
	v18 = 0
	xrw(unsafe.Pointer(&v18), 4)
	xrw(xat(v1, 8), 4)
	xrw(xat(v1, 12), 8*int(*xi32(v1, 8)))
	xrw(xat(v1, 268), 4)
	xrw(xat(v1, 272), 8)
	xrw(xat(v1, 280), 4)
	xrw(xat(v1, 284), 1)
	if xro() {
		nox_xxx_AssignIfGreater_52A420(xat(v1, 280), v3)
	}
	xrw(xat(v1, 296), 4)
	v18 = 0
	if *xi32(v1, 296) > 0 {
		v5 := xat(v1, 300)
		for {
			if xro() {
				xrw(unsafe.Pointer(&v22), 4)
				*xpp(v5, 0) = unsafe.Pointer(sub_579C80(uint32(v22)))
			} else {
				xrw(*xpp(v5, 0), 4)
			}
			v5 = xat(v5, 4)
			v18++
			if !(v18 < *xi32(v1, 296)) {
				break
			}
		}
	}
	xrw(xat(v1, 364), 4)
	xrw(xat(v1, 368), 8)
	xrw(xat(v1, 376), 4)
	xrw(xat(v1, 380), 8)
	{
		s := unsafe.Pointer(nox_xxx_getSndName_40AF80(int(*xi32(v1, 388))))
		n := xstrlen(s)
		copy(v24[:], unsafe.Slice((*byte)(s), n+1))
	}
	xlobyte(&v20, uint8(xstrlen(unsafe.Pointer(&v24[0]))))
	xrw(unsafe.Pointer(&v20), 1)
	xrw(unsafe.Pointer(&v24[0]), int(uint8(v20)))
	v24[uint8(v20)] = 0
	*xi32(v1, 388) = int32(nox_xxx_utilFindSound_40AF50(xcs(&v24)))
	xrw(xat(v1, 396), 8)
	xrw(xat(v1, 404), 4)
	if xro() {
		nox_xxx_AssignIfGreater_52A420(xat(v1, 404), v3)
	}
	xrw(xat(v1, 481), 1)
	xrw(xat(v1, 482), 1)
	xrw(xat(v1, 483), 1)
	if int16(v21) < 3 {
		xrw(unsafe.Pointer(&v18), 4)
	}
	xrw(xat(v1, 496), 4)
	xrw(xat(v1, 500), 8)
	if xro() {
		nox_xxx_AssignIfGreater_52A420(xat(v1, 496), v3)
	}
	xrw(unsafe.Pointer(&v18), 4)
	xrw(xat(v1, 536), 4)
	xrw(xat(v1, 540), 4)
	nox_xxx_AssignIfGreater_52A420(xat(v1, 536), v3)
	nox_xxx_AssignIfGreater_52A420(xat(v1, 540), v3)
	xrw(xat(v1, 544), 1)
	v6 := int32(0)
	if int8(*xu8(v1, 544)) >= 0 {
		v18 = int32(uintptr(xat(v1, 552)))
		for {
			sub_52A440(a1, unsafe.Pointer(uintptr(v18)), v3)
			v6++
			v18 += 24
			if !(v6 <= int32(int8(*xu8(v1, 544)))) {
				break
			}
		}
	}
	xrw(xat(v1, 1129), 1)
	v18 = 0
	if *xu8(v1, 1129) != 0 {
		v8 := xat(v1, 1132)
		for {
			if xro() {
				xrw(v8, 4)
			} else {
				if *memmap.PtrUint32(0x5D4594, 2487688) == 0 {
					*memmap.PtrUint32(0x5D4594, 2487688) = uint32(nox_xxx_getNameId_4E3AA0(internCStr("NewPlayer")))
				}
				xrw(xat(*xpp(v8, 0), 44), 4)
			}
			v9 := int32(*xu8(v1, 1129))
			v8 = xat(v8, 4)
			v18++
			if !(v18 < v9) {
				break
			}
		}
	}
	if xro() {
		xrw(xat(v1, 1196), 4)
	} else {
		if v10 := *xpp(v1, 1196); v10 != nil {
			v22 = int32(*xu32(v10, 44))
		} else {
			v22 = 0
		}
		xrw(unsafe.Pointer(&v22), 4)
	}
	xrw(xat(v1, 1204), 4)
	nox_xxx_AssignIfGreater_52A420(xat(v1, 1204), v3)
	xrw(xat(v1, 2096), 4)
	xrw(xat(v1, 2100), 4)
	xrw(xat(v1, 2104), 1)
	xrw(xat(v1, 2105), 1)
	xlobyte(&v20, uint8(xstrlen(xat(v1, 2106))))
	xrw(unsafe.Pointer(&v20), 1)
	xrw(xat(v1, 2106), int(uint8(v20)))
	*xu8(v1, int(uint8(v20))+2106) = 0
	if int16(v21) < 4 {
		return 1
	}
	xrw(xat(v1, 4), 4)
	xrw(xat(v1, 288), 4)
	xrw(xat(v1, 292), 4)
	// writes the script ID of an object referenced by net code (0 if none)
	xnet := func(off int) {
		if v := Nox_server_getObjectFromNetCode_4ECCB0(int(*xi32(v1, off))); v != nil {
			v18 = int32(*xu32(unsafe.Pointer(v), 44))
		} else {
			v18 = 0
		}
		xrw(unsafe.Pointer(&v18), 4)
		if xro() {
			*xi32(v1, off) = v18
		}
	}
	xnet(392)
	xrw(xat(v1, 492), 4)
	for _, off := range [...]int{508, 512, 516, 520, 528, 532} {
		xrw(xat(v1, off), 4)
		nox_xxx_AssignIfGreater_52A420(xat(v1, off), v3)
	}
	xrw(xat(v1, 524), 4)
	xrw(xat(v1, 548), 4)
	nox_xxx_AssignIfGreater_52A420(xat(v1, 548), v3)
	xrw(xat(v1, 1128), 1)
	xnet(1200)
	xrw(xat(v1, 1208), 4)
	nox_xxx_AssignIfGreater_52A420(xat(v1, 1208), v3)
	xrw(xat(v1, 1212), 4)
	nox_xxx_AssignIfGreater_52A420(xat(v1, 1212), v3)
	v14 := int32(0)
	if v13 := *xpp(v1, 1216); v13 != nil {
		v18 = int32(*xu32(v13, 44))
	} else {
		v18 = 0
	}
	xrw(unsafe.Pointer(&v18), 4)
	if xro() {
		*xi32(v1, 1216) = v18
	}
	v15 := xat(v1, 2172)
	xrw(v15, 1)
	if *xu8(v1, 2172) != 0 {
		v16 := xat(v1, 2140)
		for {
			if v17 := Nox_server_getObjectFromNetCode_4ECCB0(int(*xi32(v16, 0))); v17 != nil {
				v18 = int32(*xu32(unsafe.Pointer(v17), 44))
			} else {
				v18 = 0
			}
			xrw(unsafe.Pointer(&v18), 4)
			if xro() {
				*xi32(v16, 0) = v18
			}
			v14++
			v16 = xat(v16, 4)
			if !(v14 < int32(*xu8(v15, 0))) {
				break
			}
		}
	}
	return 1
}

//export nox_xxx_XFerNPC_52ADE0
func nox_xxx_XFerNPC_52ADE0(a1p *nox_object_t) int {
	a1 := unsafe.Pointer(a1p)
	var (
		v43, v44, v45, v47, v48, v49, v50, v51, v52 int32
		v46                                         uint8
		v53                                         [3]byte
		v55                                         [2]uint32
		v56, v57                                    [256]byte
	)
	v1 := *xpp(a1, 748)
	v2 := *xpp(a1, 756)
	v54 := *xu32(a1, 136)
	v44 = 62
	xrw(unsafe.Pointer(&v44), 2)
	if int16(v44) > 62 {
		return 0
	}
	if xro() {
		v4 := C.nox_xxx_monsterDefByTT_517560(C.int(*xu16(a1, 4)))
		*xpp(v1, 484) = v4
		if v4 != nil {
			*xu32(v1, 1440) = *xu32(v4, 92)
		}
	}
	if Nox_xxx_mapReadWriteObjData_4F4530(asObjectS(a1p), int(int16(v44))) == 0 {
		return 0
	}
	if !xro() {
		xdirIndexed(a1, &v55)
	}
	xrw(unsafe.Pointer(&v55), 8)
	xscript(v1, 1232, v2, 640)
	xscript(v1, 1264, v2, 896)
	xrw(xat(v1, 1220), 2)
	xscript(v1, 1224, v2, 768)
	if int16(v44) >= 32 {
		xscript(v1, 1240, v2, 1024)
		xscript(v1, 1248, v2, 1152)
		xscript(v1, 1256, v2, 1280)
		xscript(v1, 1272, v2, 1408)
		xscript(v1, 1280, v2, 1536)
		xscript(v1, 1288, v2, 1664)
		if int16(v44) >= 50 {
			xscript(v1, 1296, v2, 1792)
		}
	}
	v47 = 0 // FIXME in the original: set to direction? was uninitialized
	if !xro() {
		xrw(unsafe.Pointer(&v47), 4)
	} else {
		a := int16(xdirAngle(&v55))
		*xi16(a1, 126) = a
		*xi16(a1, 124) = a
		xrw(unsafe.Pointer(&v47), 4)
	}
	if xro() {
		v17 := uint8(0)
		xlobyte(&v48, 0)
		for {
			xrw(unsafe.Pointer(&v53[0]), 3)
			nox_xxx_setNPCColor_4E4A90(a1p, uint8(v48), unsafe.Pointer(&v53[0]))
			v17++
			xlobyte(&v48, v17)
			if !(v17 < 6) {
				break
			}
		}
	} else {
		v18 := xat(v1, 2076)
		for v19 := 6; v19 != 0; v19-- {
			xrw(v18, 3)
			v18 = xat(v18, 3)
		}
	}
	if xro() && uint16(v44) == 31 {
		xrw(unsafe.Pointer(&v49), 2)
		xrw(unsafe.Pointer(&v46), 1)
		for i := uint8(0); i < v46; i++ {
			xrw(unsafe.Pointer(&v50), 1)
			xrw(unsafe.Pointer(&v56[0]), int(uint8(v50)))
			v56[uint8(v50)] = 0
			xrw(unsafe.Pointer(&v43), 1)
			xrw(unsafe.Pointer(&v43), 1)
		}
	}
	if int16(v44) >= 32 {
		xrw(xat(v1, 1332), 1)
		if int16(v44) < 49 {
			xrw(unsafe.Pointer(&v47), 2)
			*xu32(v1, 1440) = uint32(uint16(v47))
		} else {
			xrw(xat(v1, 1440), 4)
		}
		xrw(xat(v1, 1352), 4)
		xrw(xat(v1, 1336), 4)
		xrw(xat(v1, 1344), 4)
		xrw(xat(v1, 1312), 4)
		v43 = 0
		if v22 := *xpp(a1, 556); v22 != nil {
			v43 = int32(*xu16(v22, 0)) // upper bits are undefined in the original
		}
		xrw(unsafe.Pointer(&v43), 2)
		if v23 := *xpp(a1, 556); v23 != nil {
			*xu16(v23, 0) = uint16(v43)
		}
		xrw(xat(v1, 1304), 4)
		*xu32(v1, 1308) = *xu32(v1, 1304)
		if int16(v44) < 35 {
			xrw(xat(v1, 1360), 4)
		}
		xlobyte(&v49, uint8(xstrlen(xat(v1, 1364))))
		xrw(unsafe.Pointer(&v49), 1)
		xrw(xat(v1, 1364), int(uint8(v49)))
		*xu8(v1, int(uint8(v49))+1364) = 0
		if int16(v44) >= 34 {
			if xro() {
				C.memset(xat(v1, 1488), 0, 0x224)
				xrw(unsafe.Pointer(&v43), 4)
				// QUIRK: j is an unsigned char compared against an int count (never ends for counts > 255)
				for j := uint8(0); int32(j) < v43; {
					xrw(unsafe.Pointer(&v49), 1)
					xrw(unsafe.Pointer(&v56[0]), int(uint8(v49)))
					v56[uint8(v49)] = 0
					v31 := nox_xxx_spellNameToN_4243F0(xcs(&v56))
					xrw(xat(v1, 4*v31+1488), 4)
					j++
					xlobyte(&v48, j)
				}
			} else {
				v43 = 0
				for k := 0; k < 137; k++ {
					if *xu32(v1, 1488+4*k) != 0 {
						v43++
					}
				}
				xrw(unsafe.Pointer(&v43), 4)
				v25 := xat(v1, 1488)
				v48 = 0
				for v47 = 137; ; {
					if *(*uint32)(v25) != 0 {
						v28 := xname(nox_xxx_spellNameByN_424870(int(v48)))
						xlobyte(&v52, uint8(xstrlen(v28)))
						xrw(unsafe.Pointer(&v52), 1)
						xrw(v28, int(uint8(v52)))
						xrw(v25, 4)
					}
					v25 = xat(v25, 4)
					v29 := v47 == 1
					v48++
					v47--
					if v29 {
						break
					}
				}
			}
		} else {
			xrw(xat(v1, 1488), 0x224)
		}
		for _, off := range [...]int{1448, 1456, 1464, 1472, 1480} {
			if int16(v44) < 47 {
				xrw(unsafe.Pointer(&v43), 1)
				*xu16(v1, off) = uint16(uint8(v43))
				xrw(unsafe.Pointer(&v43), 1)
				*xu16(v1, off+2) = uint16(uint8(v43))
			} else {
				xrw(xat(v1, off), 2)
				xrw(xat(v1, off+2), 2)
			}
			if int16(v44) < 34 {
				xrw(unsafe.Pointer(&v45), 4)
			}
		}
		if int16(v44) >= 33 {
			xrw(xat(v1, 1316), 4)
		}
		if int16(v44) >= 34 {
			xrw(xat(v1, 2040), 4)
			xrw(xat(v1, 1324), 1)
			xrw(xat(v1, 1328), 4)
			xrw(xat(v1, 1320), 4)
			if int16(v44) < 42 {
				xrw(unsafe.Pointer(&v43), 2)
				if uint16(v43) == 0 {
					*xu8(v1, 1445) = 1
				}
			}
			v32 := xat(v1, 2044)
			for v47 = 3; v47 != 0; v47-- {
				if xro() {
					xrw(unsafe.Pointer(&v43), 1)
					xrw(unsafe.Pointer(&v56[0]), int(uint8(v43)))
					v56[uint8(v43)] = 0
					*xi32(v32, 0) = int32(nox_xxx_spellNameToN_4243F0(xcs(&v56)))
				} else {
					v33 := xname(nox_xxx_spellNameByN_424870(int(*xi32(v32, 0))))
					xlobyte(&v43, uint8(xstrlen(v33)))
					xrw(unsafe.Pointer(&v43), 1)
					xrw(v33, int(uint8(v43)))
				}
				v32 = xat(v32, 4)
			}
		}
		if int16(v44) >= 35 {
			if xro() {
				xrw(unsafe.Pointer(&v43), 1)
				xrw(unsafe.Pointer(&v57[0]), int(uint8(v43)))
				v57[uint8(v43)] = 0
				*xu32(v1, 1360) = uint32(C.nox_xxx_actionNByNameMB_5345F0(xcs(&v57)))
			} else {
				v34 := unsafe.Pointer(C.sub_5345B0(C.int(*xu32(v1, 1360))))
				xlobyte(&v43, uint8(xstrlen(v34)))
				xrw(unsafe.Pointer(&v43), 1)
				xrw(v34, int(uint8(v43)))
			}
		}
	}
	v35 := a1
	if int16(v44) >= 41 {
		if nox_xxx_XFer_ActionData_529CE0(a1) == 0 {
			return 0
		}
	}
	if int16(v44) >= 42 {
		xrw(xat(v1, 1445), 1)
	}
	if int16(v44) >= 44 {
		xrw(v1, 4)
	}
	if int16(v44) >= 45 {
		v36 := *xpp(v35, 556)
		v45 = 0
		if v36 != nil {
			v45 = int32(*xu16(v36, 4))
		}
		xrw(unsafe.Pointer(&v45), 4)
		if v37 := *xpp(v35, 556); v37 != nil {
			*xu16(v37, 4) = uint16(v45)
		}
	}
	if int16(v44) >= 46 {
		v51 = int32(*xu32(v35, 12) & 0x180)
		xrw(unsafe.Pointer(&v51), 4)
		*xu32(v35, 12) |= uint32(v51)
	}
	if int16(v44) >= 48 {
		xrw(*xpp(v35, 556), 2)
	}
	if int16(v44) >= 51 {
		xrw(xat(v35, 28), 4)
	}
	if int16(v44) >= 52 {
		C.nox_xxx_readNPCVoiceSet_52AD10(xint(v35))
	}
	if int16(v44) >= 61 {
		if C.nox_xxx_XFer_ReadMonsterBuffs_52AAB0((*C.uint32_t)(v35)) == 0 {
			return 0
		}
	}
	// Control flow of the original (LABEL_149, LABEL_170, LABEL_156):
	// 0 = goto 156, 1 = goto 149
	goTo := 0
	if int16(v44) < 62 {
		if xro() {
			goTo = 1
		}
	} else {
		xlobyte(&v45, *xu8(v35, 540))
		xrw(unsafe.Pointer(&v45), 1)
		if xro() {
			if uint8(v45) != 0 {
				C.nox_xxx_setSomePoisonData_4EEA90(xint(v35), C.int(uint8(v45)))
			}
			goTo = 1
		}
	}
	if goTo == 1 {
		// LABEL_149
		do170 := true
		if *xu8(v1, 1445) != 0 {
			if xgame(1) {
				if v38 := *xpp(v35, 556); v38 != nil {
					*xu16(v38, 4) = 0
					*xu16(*xpp(v35, 556), 0) = 0
				}
			}
			if !xro() {
				do170 = false
			}
		}
		if do170 {
			// LABEL_170
			if *xu32(v35, 16)&0x8000 != 0 {
				*xu32(v35, 16) |= 0x40
			}
		}
	}
	// LABEL_156
	if *xu32(v35, 136) == 0 || !xro() || Nox_xxx_xfer_4F3E30(int(v44), asObjectS(a1p), *xu32(v35, 136)) != 0 {
		sub_52BA70(v35)
		if xro() {
			if xgame(1) {
				for k := *xpp(v35, 504); k != nil; k = *xpp(k, 496) {
					if *xu32(k, 16)&0x100 != 0 {
						v42 := *xu32(k, 8)
						*xu32(k, 16) &= 0xFFFFFEFF
						if v42&0x1001000 != 0 {
							C.nox_xxx_NPCEquipWeapon_53A2C0(xint(v35), (*C.nox_object_t)(k))
						} else {
							C.nox_xxx_NPCEquipArmor_53E520(xint(v35), (*C.uint32_t)(k))
						}
					}
				}
			}
		}
		*xu32(v35, 136) = v54
		return 1
	}
	return 0
}

func sub_52BA70(a1 unsafe.Pointer) {
	v2 := 0
	v3 := 0
	if a1 == nil {
		return
	}
	for r := *xpp(a1, 504); r != nil; r = *xpp(r, 496) {
		v4 := *xu32(r, 16)
		if v4&0x100 != 0 {
			v5 := *xu32(r, 8)
			if v5&0x1001000 != 0 && *xu32(r, 12)&0x7FFE40C != 0 {
				if v3 == 1 {
					v4 &^= 0x100 // BYTE1(v4) &= 0xFE
					*xu32(r, 16) = v4
				} else {
					v2 = 1
				}
			} else if v5&0x2000000 != 0 && *xu8(r, 12)&2 != 0 {
				if v2 == 1 {
					v4 &^= 0x100
					*xu32(r, 16) = v4
				} else {
					v3 = 1
				}
			}
		}
	}
}

func nox_xxx_monsterAutoSpells_54C0C0(a1 unsafe.Pointer) {
	v1 := *xpp(a1, 748)
	mv := func(i uintptr) *uint32 { return memmap.PtrUint32(0x5D4594, 2491624+4*i) }
	if *mv(4) == 0 {
		for i, n := range [...]string{"UrchinShaman", "Wizard", "WizardWhite", "Beholder", "Lich", "LichLord", "Demon", "WizardGreen", "WillOWisp"} {
			*mv(uintptr(i)) = uint32(nox_xxx_getNameId_4E3AA0(internCStr(n)))
		}
	}
	v2 := uint32(*xu16(a1, 4))
	set := func(off int, v uint32) { *xu32(v1, off) = v }
	if v2 != *mv(0) {
		if v2 == *mv(1) || v2 == *mv(2) {
			set(2040, 3)
			set(1540, 0x8000000)
			set(1640, 0x8000000)
			set(1644, 0x10000000)
			set(1776, 0x20000000)
			set(1596, 0x40000000)
			set(1688, 0x40000000)
			set(1660, 0x40000000)
			set(1584, 0x40000000)
			set(1504, 2147483648)
		} else if v2 == *mv(3) {
			set(2040, 3)
			set(1584, 0x40000000)
			if !xgame(4096) {
				set(1504, 2147483648)
				*xu8(v1, 2036) = 1
				return
			}
		} else {
			if v2 == *mv(4) {
				set(2040, 3)
				set(1540, 0x8000000)
				set(1640, 0x8000000)
				set(1772, 0x10000000)
				set(1776, 0x20000000)
				set(1620, 0x20000000)
				set(1644, 2147483648)
				set(1596, 0x40000000)
				set(1660, 0x40000000)
				set(1584, 0x40000000)
				*xu8(v1, 2036) = 1
				return
			}
			if v2 == *mv(5) {
				set(2040, 3)
				set(1540, 0x8000000)
				set(1640, 0x8000000)
				set(1772, 0x10000000)
				*xu16(v1, 1472) = 3 * uint16(gameFPS())
				fps := uint16(gameFPS())
				set(1612, 0x40000000)
				set(1644, 2147483648)
				*xu8(v1, 2036) = 1
				*xu16(v1, 1474) = 5 * fps
				return
			}
			if v2 == *mv(6) {
				set(2040, 3)
				set(1540, 0x8000000)
				set(1784, 0x20000000)
				set(1596, 0x40000000)
				set(1528, 0x40000000)
				set(1644, 2147483648)
				*xu8(v1, 2036) = 1
				return
			}
			if v2 == *mv(7) {
				set(2040, 3)
				set(1640, 0x8000000)
				set(1536, 0x20000000)
				set(1720, 0x40000000)
				set(1728, 0x40000000)
				set(1688, 0x40000000)
				set(1504, 2147483648)
				*xu8(v1, 2036) = 1
				return
			}
			if v2 == *mv(8) {
				set(2040, 3)
				set(1660, 0x40000000)
				*xu8(v1, 2036) = 1
				return
			}
		}
		*xu8(v1, 2036) = 1
		return
	}
	set(2040, 2)
	set(1640, 0x8000000)
	*xu16(v1, 1448) = 0
	*xu16(v1, 1450) = uint16(gameFPS() >> 1)
	set(1536, 0x20000000)
	*xu16(v1, 1464) = 3 * uint16(gameFPS())
	v3 := uint16(gameFPS())
	*xu16(v1, 1472) = 0
	*xu16(v1, 1466) = 5 * v3
	set(1720, 0x40000000)
	set(1696, 0x40000000)
	v4 := uint16(gameFPS())
	set(1752, 2147483648)
	*xu16(v1, 1474) = 3 * v4
	*xu16(v1, 1480) = uint16(gameFPS())
	v2b := uint16(gameFPS())
	*xu8(v1, 2036) = 1
	*xu16(v1, 1482) = 3 * v2b
}
