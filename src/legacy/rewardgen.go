package legacy

/*
#include <string.h>
#include "defs.h"
#include "GAME1_1.h"
#include "GAME3_3.h"
#include "GAME5_2.h"
*/
import "C"
import (
	"math"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
)

// Reward marker loot generator (ported from GAME3_3.c). All tables are read from the 0x587000 blob exactly as the
// original C does, since mods may patch them.

func rwD(off uintptr) uint32 { return memmap.Uint32(0x587000, off) }
func rwB(off uintptr) byte   { return memmap.Uint8(0x587000, off) }

// rwPD/rwPB read from a pointer stored in a table (e.g. a modifier effect definition).
func rwPD(p uint32, off uintptr) uint32 { return *(*uint32)(unsafe.Pointer(uintptr(p) + off)) }
func rwPB(p uint32, off uintptr) byte   { return *(*byte)(unsafe.Pointer(uintptr(p) + off)) }

func rwInitD(obj *server.Object, off uintptr) uint32 {
	return *(*uint32)(unsafe.Add(obj.InitData, off))
}
func rwInitB(obj *server.Object, off uintptr) byte {
	return *(*byte)(unsafe.Add(obj.InitData, off))
}

func rwRand(min, max int) int { return GetServer().S().Rand.Logic.IntClamp(min, max) }

// rwModCount mirrors the do/while counting loops over 24-byte modifier tables (terminated by a zero at +28).
func rwModCount(e uintptr, cond func(e uintptr) bool) int {
	n := 0
	for {
		if cond(e) {
			n++
		}
		next := rwD(e + 28)
		e += 24
		if next == 0 {
			break
		}
	}
	return n
}

// rwModPick mirrors the loops that pick the n-th matching modifier table entry. Returns its index.
func rwModPick(e uintptr, n int, cond func(e uintptr) bool) (int, bool) {
	k, idx := 0, 0
	for {
		if cond(e) {
			if k == n {
				return idx, true
			}
			k++
		}
		next := rwD(e + 28)
		e += 24
		idx++
		if next == 0 {
			return 0, false
		}
	}
}

//export nox_server_rewardgen_activateMarker_4F0720
func nox_server_rewardgen_activateMarker_4F0720(a1 C.int, a2 C.uint) *C.uint32_t {
	obj := rewardgenActivateMarker((*server.Object)(unsafe.Pointer(uintptr(a1))), uint32(a2))
	return (*C.uint32_t)(unsafe.Pointer(obj))
}

func rewardgenActivateMarker(a1 *server.Object, a2 uint32) *server.Object { // nox_server_rewardgen_activateMarker_4F0720
	s := GetServer().S()
	p := memmap.PtrUint32(0x5D4594, 1568276)
	v2 := int32(*p)
	v3 := a1
	if *p == 0 {
		v2 = int32(s.Types.IndByID("RewardMarkerPlus"))
		*p = uint32(v2)
	}
	if int32(a1.TypeInd) == v2 {
		a2 += 2
	}
	switch int32(rwInitD(a1, 212)) {
	case 1:
		if rwRand(0, 100) > 75 {
			return nil
		}
	case 2:
		if rwRand(0, 100) > 50 {
			return nil
		}
	case 3:
		if rwRand(0, 100) > 25 {
			return nil
		}
	case 4:
		if rwRand(0, 100) > 5 {
			return nil
		}
	}
	flags := int32(rwInitD(a1, 0))
	v6 := 0
	v7 := 0
	for v8 := uintptr(207044); ; {
		if (1<<v7)&flags != 0 {
			v6 += int(rwB(v8))
		}
		v8 += 8
		v7++
		if !(v8 < 207108) {
			break
		}
	}
	if v6 == 0 {
		return nil
	}
	v9 := uint32(rwRand(1, v6))
	v10 := int32(rwInitD(a1, 0))
	v11 := uint32(0)
	v12 := 0
	v13 := uintptr(207044)
	var v14 uint32
	for {
		v14 = uint32(1) << v12
		if v10&(1<<v12) != 0 {
			v11 += uint32(rwB(v13))
			v3 = a1
			if v11 >= v9 {
				break
			}
		}
		v13 += 8
		v12++
		if !(v13 < 207108) {
			v14 = a2
			break
		}
	}
	switch v14 {
	case 1:
		return rewardSpellBook(v3, a2)
	case 2:
		return rewardAbilityBook(v3)
	case 4:
		return rewardFieldGuide(v3, a2)
	case 8:
		return rewardMakeWeapon(v3, a2)
	case 0x10:
		return rewardMakeArmor(v3, a2)
	case 0x20:
		return rewardCreateGem(v3, a2)
	case 0x40:
		return rewardMakePotion(v3, a2)
	case 0x80:
		return rewardCreateGem2(v3, a2)
	default:
		return rewardCreateGem(v3, a2)
	}
}

func rewardSpellBook(a1 *server.Object, a2 uint32) *server.Object { // nox_xxx_rewardSpellBook_4F09F0
	s := GetServer().S()
	var v9 int
	if rwInitB(a1, 4)&1 != 0 {
		v3 := 0
		for i := 0; i < 137; i++ {
			if rwInitB(a1, uintptr(i+8)) == 1 {
				v3++
			}
		}
		if v3 == 0 {
			return nil
		}
		v5 := rwRand(0, v3-1)
		v6 := 0
		v7 := 0
		for {
			if rwInitB(a1, uintptr(v7+8)) == 1 {
				if v6 == v5 {
					v9 = v7
					break
				}
				v6++
			}
			v7++
			if v7 >= 137 {
				return nil
			}
		}
	} else {
		v10 := uint32(rewardPickRandomSlots(a2))
		v11 := 0
		if rwD(207108) == 0 {
			return nil
		}
		for v12 := uintptr(207104); ; {
			if v10&rwD(v12+8) != 0 {
				v11 += int(rwB(v12))
			}
			v13 := rwD(v12 + 16)
			v12 += 12
			if v13 == 0 {
				break
			}
		}
		if v11 == 0 {
			return nil
		}
		v14 := rwRand(0, v11-1)
		v15 := 0
		v16 := 0
		if rwD(207108) == 0 {
			return nil
		}
		for j := uintptr(207104); ; j += 12 {
			if v10&rwD(j+8) != 0 {
				v15 += int(rwB(j))
				if v14 < v15 {
					break
				}
			}
			v18 := rwD(j + 16)
			v16++
			if v18 == 0 {
				return nil
			}
		}
		v9 = int(int32(rwD(uintptr(207108 + 12*v16))))
	}
	if v9 == 0 {
		return nil
	}
	var result *server.Object
	if C.nox_xxx_playerCheckSpellClass_57AEA0(1, C.int(v9)) != 0 || C.nox_xxx_playerCheckSpellClass_57AEA0(2, C.int(v9)) != 0 {
		if C.nox_xxx_playerCheckSpellClass_57AEA0(1, C.int(v9)) != 0 {
			if C.nox_xxx_playerCheckSpellClass_57AEA0(2, C.int(v9)) != 0 {
				return nil
			}
			result = s.NewObjectByTypeID("ConjurerSpellBook")
		} else {
			result = s.NewObjectByTypeID("WizardSpellBook")
		}
	} else {
		result = s.NewObjectByTypeID("CommonSpellBook")
	}
	if result == nil {
		return nil
	}
	result.UseData.AsSpellReward().Spell = byte(v9)
	return result
}

func rewardPickRandomSlots(a1 uint32) int { // nox_server_rewardGen_pickRandomSlots_4F0B60
	var v9 [5]float32
	if a1 > 0xA {
		return 16
	}
	switch a1 {
	case 0:
		return 1
	case 1:
		v9[0] = 87.5
		v9[1] = 12.5
	case 9:
		v9[3] = 12.5
		v9[4] = 87.5
	case 0xA:
		return 16
	default:
		if a1&1 != 0 {
			v2 := a1 >> 1
			v9[v2] = 75.0
			// Quirk: the original overwrites the 75.0 with the float bits of this int constant (12.5).
			v9[v2] = math.Float32frombits(1095237632)
			v9[v2+1] = 12.5
		} else {
			v3 := a1 >> 1
			v9[v3-1] = 50.0
			v9[v3] = 50.0
		}
	}
	v4 := rwRand(0, 200)
	v5 := 0
	v7 := float64(0.0)
	for {
		v7 = v7 + float64(v9[v5])
		v10 := float32(float64(v4) * 0.5)
		if float64(v10) <= v7 {
			break
		}
		v5++
		if v5 >= 5 {
			return 1
		}
	}
	return 1 << v5
}

func rewardAbilityBook(a1 *server.Object) *server.Object { // nox_xxx_rewardAbilityBook_4F0C70
	s := GetServer().S()
	var v8 int
	if rwInitB(a1, 4)&2 == 0 {
		v8 = rwRand(1, 5)
	} else {
		v2 := 0
		for i := 0; i < 6; i++ {
			if rwInitB(a1, uintptr(i+145)) == 1 {
				v2++
			}
		}
		if v2 == 0 {
			return nil
		}
		v5 := rwRand(0, v2-1)
		v6 := 0
		v7 := 0
		for {
			if rwInitB(a1, uintptr(v7+145)) == 1 {
				if v6 != v5 {
					v6++
				} else {
					v8 = v7
					break
				}
			}
			v7++
			if v7 >= 6 {
				return nil
			}
		}
	}
	if v8 == 0 {
		return nil
	}
	result := s.NewObjectByTypeID("AbilityBook")
	if result != nil {
		result.UseData.AsAbilityReward().Ability = byte(v8)
	}
	return result
}

func rewardFieldGuide(a1 *server.Object, a2 uint32) *server.Object { // nox_xxx_rewardFieldGuide_4F0D20
	s := GetServer().S()
	var v9 int
	if rwInitB(a1, 4)&4 != 0 {
		v3 := 0
		for i := 0; i < 41; i++ {
			if rwInitB(a1, uintptr(i+151)) == 1 {
				v3++
			}
		}
		if v3 == 0 {
			return nil
		}
		v6 := rwRand(0, v3-1)
		v7 := 0
		v8 := 0
		for {
			if rwInitB(a1, uintptr(v8+151)) == 1 {
				if v7 == v6 {
					v9 = v8
					break
				}
				v7++
			}
			v8++
			if v8 >= 41 {
				return nil
			}
		}
	} else {
		v10 := uint32(rewardPickRandomSlots(a2))
		v11 := 0
		if rwD(207796) == 0 {
			return nil
		}
		for v12 := uintptr(207792); ; {
			if v10&rwD(v12+8) != 0 {
				v11 += int(rwB(v12))
			}
			v13 := rwD(v12 + 16)
			v12 += 12
			if v13 == 0 {
				break
			}
		}
		if v11 == 0 {
			return nil
		}
		v14 := rwRand(0, v11-1)
		v15 := 0
		v16 := 0
		if rwD(207796) == 0 {
			return nil
		}
		for j := uintptr(207792); ; j += 12 {
			if v10&rwD(j+8) != 0 {
				v15 += int(rwB(j))
				if v14 < v15 {
					break
				}
			}
			v18 := rwD(j + 16)
			v16++
			if v18 == 0 {
				return nil
			}
		}
		v9 = int(int32(rwD(uintptr(207796 + 12*v16))))
	}
	if v9 == 0 {
		return nil
	}
	result := s.NewObjectByTypeID("FieldGuide")
	if result != nil {
		C.strcpy((*C.char)(result.UseData.Ptr), C.nox_xxx_guideNameByN_427230(C.int(v9)))
	}
	return result
}

func rewardMakeArmor(a1 *server.Object, a2 uint32) *server.Object { // nox_xxx_rewardMakeArmor_4F0E80
	s := GetServer().S()
	v65 := 0
	v64 := uint32(rewardPickRandomSlots(a2))
	v3 := 0
	if rwD(208180) == 0 {
		return nil
	}
	for v4 := uintptr(208192); ; {
		if rwB(v4-4)&2 != 0 && v64&rwD(v4) != 0 && s.Types.ByInd(int(rwD(v4-8))).Allowed() {
			v3 += int(rwB(v4 - 16))
		}
		v5 := rwD(v4 + 8)
		v4 += 20
		if v5 == 0 {
			break
		}
	}
	if v3 == 0 {
		return nil
	}
	v7 := rwRand(0, v3-1)
	v8 := 0
	v9 := 0
	if rwD(208180) == 0 {
		return nil
	}
	for i := uintptr(208192); ; i += 20 {
		if rwB(i-4)&2 != 0 {
			if v64&rwD(i) != 0 {
				if s.Types.ByInd(int(rwD(i - 8))).Allowed() {
					v8 += int(rwB(i - 16))
					if v7 < v8 {
						break
					}
				}
			}
		}
		v11 := rwD(i + 8)
		v9++
		if v11 == 0 {
			return nil
		}
	}
	v12 := int(rwD(uintptr(208184 + 20*v9)))
	if v12 == 0 {
		return nil
	}
	v13 := s.Armor.Sub_415D10(int(rwD(uintptr(208184 + 20*v9))))
	v14 := s.NewObjectByTypeInd(v12)
	if v14 == nil {
		return nil
	}
	var v15 int
	switch v64 {
	case 2:
		v15 = rwRand(0, 1)
	case 4:
		v15 = rwRand(0, 2)
	case 8:
		v15 = rwRand(1, 3)
	case 16:
		v15 = rwRand(2, 4)
	default:
		return v14
	}
	if v15 == 0 {
		return v14
	}
	var v67 [5]uint32 // 20 bytes: 4 dwords + 2 words
	var v17 int
	switch v15 {
	case 1:
		v16 := rwRand(1, 100)
		if v16 > 20 {
			v17 = 1
			if v16 > 50 {
				v17 = 2
			}
		} else {
			v17 = 4
		}
	case 2:
		v18 := rwRand(1, 100)
		if v18 > 12 {
			if v18 > 25 {
				v17 = 3
			} else {
				v17 = 6
			}
		} else {
			v17 = 5
		}
	case 3:
		v17 = 7
	default:
		v17 = 15
		if v15 != 4 {
			v17 = int(int16(a2))
		}
	}
	// modifier table conditions (armor: effect def +32 mask, table +12 exclusion)
	armorCond := func(e uintptr) bool {
		return v64&rwD(e+8) != 0 && v13&rwPD(rwD(e), 32) != 0 && v13&rwD(e+12) == 0
	}
	armorCond1 := func(e uintptr) bool {
		return armorCond(e) && rwPB(rwD(e), 36)&1 != 0
	}
	if v17&1 != 0 {
		v19 := 0
		if rwD(210856) != 0 {
			v19 = rwModCount(210852, armorCond)
		}
		if v19 == 0 {
			if v17&2 != 0 {
				if v17&4 != 0 {
					if v17&8 == 0 {
						v17 |= 8
					}
				} else {
					v17 |= 4
				}
			} else {
				v17 |= 2
			}
			v17 &= 0xFFFE
		}
	}
	if v17&2 != 0 {
		v22 := 0
		if rwD(211000) != 0 {
			v22 = rwModCount(210996, armorCond)
		}
		if v22 == 0 {
			if v17&4 != 0 {
				if v17&8 == 0 {
					v17 |= 8
				}
			} else {
				v17 |= 4
			}
			v17 &= 0xFFFD
		}
	}
	if v17&4 != 0 {
		v25 := 0
		if rwD(209344) != 0 {
			v25 = rwModCount(209340, armorCond1)
		}
		if v25 == 0 {
			v17 &= 0xFFF3
		}
	}
	if v17 == 0 {
		return v14
	}
	if v17&1 != 0 && rwD(210856) != 0 {
		if v28 := rwModCount(210852, armorCond); v28 != 0 {
			v31 := rwRand(0, v28-1)
			if rwD(210856) != 0 {
				if idx, ok := rwModPick(210852, v31, armorCond); ok {
					v67[0] = rwD(uintptr(210852 + 24*idx))
				}
			}
		}
	}
	if v17&2 != 0 {
		if rwD(211000) != 0 {
			if v36 := rwModCount(210996, armorCond); v36 != 0 {
				v39 := rwRand(0, v36-1)
				if rwD(211000) != 0 {
					if idx, ok := rwModPick(210996, v39, armorCond); ok {
						v67[1] = rwD(uintptr(210996 + 24*idx))
					}
				}
			}
		}
	}
	if v17&4 != 0 {
		if rwD(209344) != 0 {
			if v44 := rwModCount(209340, armorCond1); v44 != 0 {
				v47 := rwRand(0, v44-1)
				if rwD(209344) != 0 {
					if idx, ok := rwModPick(209340, v47, armorCond1); ok {
						v65 = idx
						v67[2] = rwD(uintptr(209340 + 24*idx))
					}
				}
			}
		}
	}
	if v17&8 != 0 {
		v52 := int(int32(a2 >> 1))
		if v52 >= 1 {
			if v52 >= 5 {
				v52 = 4
			}
		} else {
			v52 = 1
		}
		v53 := v52 - 1
		v54 := v52 + 1
		if v53 < 1 {
			v53 = 1
		}
		if v54 >= 5 {
			v54 = 4
		}
		v55 := uint32(rwRand(v53, v54))
		cond2 := func(e uintptr) bool {
			return v55&rwD(e+8) != 0 && v13&rwPD(rwD(e), 32) != 0 && v13&rwD(e+12) == 0 && rwPB(rwD(e), 36)&2 != 0
		}
		if rwD(209344) != 0 {
			if v56 := rwModCount(209340, cond2); v56 != 0 {
				v59 := rwRand(0, v56-1)
				if rwD(209344) != 0 {
					if idx, ok := rwModPick(209340, v59, cond2); ok {
						if rwB(uintptr(209336+24*v65)) != rwB(uintptr(209336+24*idx)) {
							v67[3] = rwD(uintptr(209340 + 24*idx))
						}
					}
				}
			}
		}
	}
	Nox_xxx_modifSetItemAttrs_4E4990(v14, unsafe.Pointer(&v67[0]))
	return v14
}

func rewardMakeWeapon(a1 *server.Object, a2 uint32) *server.Object { // nox_xxx_rewardMakeWeapon_4F14E0
	s := GetServer().S()
	v74 := 0
	v72 := uint32(rewardPickRandomSlots(a2))
	v3 := 0
	if rwD(208180) == 0 {
		return nil
	}
	for v4 := uintptr(208192); ; {
		if rwB(v4-4)&1 != 0 && v72&rwD(v4) != 0 && s.Types.ByInd(int(rwD(v4-8))).Allowed() {
			v3 += int(rwB(v4 - 16))
		}
		v5 := rwD(v4 + 8)
		v4 += 20
		if v5 == 0 {
			break
		}
	}
	if v3 == 0 {
		return nil
	}
	v7 := rwRand(0, v3-1)
	v8 := 0
	v9 := 0
	var v12 uint32
	if rwD(208180) != 0 {
		for v10 := uintptr(208192); ; {
			if rwB(v10-4)&1 != 0 {
				if v72&rwD(v10) != 0 {
					if s.Types.ByInd(int(rwD(v10 - 8))).Allowed() {
						v8 += int(rwB(v10 - 16))
						if v7 < v8 {
							v12 = rwD(uintptr(208184 + 20*v9))
							break
						}
					}
				}
			}
			v11 := rwD(v10 + 8)
			v10 += 20
			v9++
			if v11 == 0 {
				v12 = a2 // quirk: falls back to a2 as the type index
				break
			}
		}
	} else {
		v12 = a2
	}
	if v12 == 0 {
		return nil
	}
	v73 := s.Weapons.Nox_xxx_ammoCheck_415880(int(v12))
	v14 := s.NewObjectByTypeInd(int(v12))
	if v14 == nil {
		return nil
	}
	if uint32(v14.ObjClass)&0x1000 != 0 && uint32(v14.ObjSubClass)&0x47F0000 != 0 {
		// GUARD: when the fallback path (v12 = a2) is taken, v9 indexes the terminator entry whose name pointer is NULL;
		// the original would dereference NULL here.
		if namep := rwD(uintptr(208180 + 20*v9)); namep != 0 {
			if *(*byte)(unsafe.Pointer(uintptr(namep))) == 35 {
				var v76 [5]uint32
				v16 := s.Modif.Nox_xxx_modifGetIdByName413290("Replenishment1")
				v76[2] = uint32(uintptr(s.Modif.Nox_xxx_modifGetDescById413330(v16).C()))
				Nox_xxx_modifSetItemAttrs_4E4990(v14, unsafe.Pointer(&v76[0]))
			}
		}
		return v14
	}
	var v17 int
	switch v72 {
	case 2:
		v17 = rwRand(0, 1)
	case 4:
		v17 = rwRand(0, 2)
	case 8:
		v17 = rwRand(1, 3)
	case 16:
		v17 = rwRand(2, 4)
	default:
		return v14
	}
	if v17 == 0 {
		return v14
	}
	var v76 [5]uint32
	// Note: the original mirrors v19 into a separate variable v71 (low byte); they always hold the same value,
	// except in the unreachable default case, so only v19 is kept here.
	var v19 int
	switch v17 {
	case 1:
		v18 := rwRand(1, 100)
		if v18 > 20 {
			v19 = 1
			if v18 > 50 {
				v19 = 2
			}
		} else {
			v19 = 4
		}
	case 2:
		v20 := rwRand(1, 100)
		if v20 > 12 {
			if v20 > 25 {
				v19 = 3
			} else {
				v19 = 6
			}
		} else {
			v19 = 5
		}
	case 3:
		v19 = 7
	case 4:
		v19 = 15
	}
	v22 := v73
	wpnCond := func(e uintptr) bool {
		return v72&rwD(e+8) != 0 && v22&rwPD(rwD(e), 28) != 0 && v22&rwD(e+16) == 0
	}
	wpnCond1 := func(e uintptr) bool {
		return wpnCond(e) && rwPB(rwD(e), 36)&1 != 0
	}
	if v19&1 != 0 {
		v21 := 0
		if rwD(210712) != 0 {
			v21 = rwModCount(210708, wpnCond)
		}
		if v21 == 0 {
			if v19&2 != 0 {
				if v19&4 != 0 {
					if v19&8 == 0 {
						v19 |= 8
					}
					v19 &= 0xFFFE
				} else {
					v19 = (v19 | 4) & 0xFFFE
				}
			} else {
				v19 = (v19 | 2) & 0xFFFE
			}
		}
	}
	if v19&2 != 0 {
		v25 := 0
		if rwD(211000) != 0 {
			v25 = rwModCount(210996, wpnCond)
		}
		if v25 == 0 {
			if v19&4 != 0 {
				if v19&8 == 0 {
					v19 |= 8
				}
			} else {
				v19 |= 4
			}
			v19 &= 0xFFFD
		}
	}
	if v19&4 != 0 {
		v28 := 0
		if rwD(209344) != 0 {
			v28 = rwModCount(209340, wpnCond1)
		}
		if v28 == 0 {
			v19 &= 0xFFF3
		}
	}
	if v19 == 0 {
		return v14
	}
	if v19&1 != 0 {
		if rwD(210712) != 0 {
			if v31 := rwModCount(210708, wpnCond); v31 != 0 {
				v34 := rwRand(0, v31-1)
				if rwD(210712) != 0 {
					if idx, ok := rwModPick(210708, v34, wpnCond); ok {
						v76[0] = rwD(uintptr(210708 + 24*idx))
					}
				}
			}
		}
	}
	if v19&2 != 0 {
		if rwD(211000) != 0 {
			if v39 := rwModCount(210996, wpnCond); v39 != 0 {
				v42 := rwRand(0, v39-1)
				if rwD(211000) != 0 {
					if idx, ok := rwModPick(210996, v42, wpnCond); ok {
						v76[1] = rwD(uintptr(210996 + 24*idx))
					}
				}
			}
		}
	}
	rangeLevel := func() uint32 {
		v47 := int(int32(a2 >> 1))
		if v47 >= 1 {
			if v47 >= 5 {
				v47 = 4
			}
		} else {
			v47 = 1
		}
		v48 := v47 - 1
		v49 := v47 + 1
		if v48 < 1 {
			v48 = 1
		}
		if v49 >= 5 {
			v49 = 4
		}
		return uint32(rwRand(v48, v49))
	}
	if v19&4 != 0 {
		v50 := rangeLevel()
		cond := func(e uintptr) bool {
			return v50&rwD(e+8) != 0 && v73&rwPD(rwD(e), 28) != 0 && v73&rwD(e+16) == 0 && rwPB(rwD(e), 36)&1 != 0
		}
		if rwD(209344) != 0 {
			if v51 := rwModCount(209340, cond); v51 != 0 {
				v54 := rwRand(0, v51-1)
				if rwD(209344) != 0 {
					if idx, ok := rwModPick(209340, v54, cond); ok {
						v74 = idx
						v76[2] = rwD(uintptr(209340 + 24*idx))
					}
				}
			}
		}
	}
	if v19&8 != 0 {
		v62 := rangeLevel()
		cond := func(e uintptr) bool {
			return v62&rwD(e+8) != 0 && v73&rwPD(rwD(e), 28) != 0 && v73&rwD(e+16) == 0 && rwPB(rwD(e), 36)&2 != 0
		}
		if rwD(209344) != 0 {
			if v63 := rwModCount(209340, cond); v63 != 0 {
				v66 := rwRand(0, v63-1)
				if rwD(209344) != 0 {
					if idx, ok := rwModPick(209340, v66, cond); ok {
						if rwB(uintptr(209336+24*v74)) != rwB(uintptr(209336+24*idx)) {
							v76[3] = rwD(uintptr(209340 + 24*idx))
						}
					}
				}
			}
		}
	}
	Nox_xxx_modifSetItemAttrs_4E4990(v14, unsafe.Pointer(&v76[0]))
	return v14
}

func rewardMakePotion(a1 *server.Object, a2 uint32) *server.Object { // nox_xxx_rewardMakePotion_4F1C40
	s := GetServer().S()
	v12 := uint32(rewardPickRandomSlots(a2))
	v3 := 0
	if rwD(208180) == 0 {
		return nil
	}
	for v4 := uintptr(208192); ; {
		if rwB(v4-4)&4 != 0 && v12&rwD(v4) != 0 && s.Types.ByInd(int(rwD(v4-8))).Allowed() {
			v3 += int(rwB(v4 - 16))
		}
		v5 := rwD(v4 + 8)
		v4 += 20
		if v5 == 0 {
			break
		}
	}
	if v3 == 0 {
		return nil
	}
	v7 := rwRand(0, v3-1)
	v8 := 0
	v9 := 0
	if rwD(208180) == 0 {
		return nil
	}
	for i := uintptr(208192); ; i += 20 {
		if rwB(i-4)&4 != 0 {
			if rwD(i)&v12 != 0 {
				if s.Types.ByInd(int(rwD(i - 8))).Allowed() {
					v8 += int(rwB(i - 16))
					if v7 < v8 {
						break
					}
				}
			}
		}
		v11 := rwD(i + 8)
		v9++
		if v11 == 0 {
			return nil
		}
	}
	if ind := rwD(uintptr(208184 + 20*v9)); ind != 0 {
		return s.NewObjectByTypeInd(int(ind))
	}
	return nil
}

func rewardCreateGem(a1 *server.Object, a2 uint32) *server.Object { // nox_xxx_createGem_4F1D30
	s := GetServer().S()
	v2 := uint32(rewardPickRandomSlots(a2))
	var result *server.Object
	if v2 < 4 || rwRand(1, 100) <= 90 {
		if rwRand(1, 2) == 1 {
			result = s.NewObjectByTypeID("QuestGoldChest")
		} else {
			result = s.NewObjectByTypeID("QuestGoldPile")
		}
		if result != nil {
			v6 := (*int32)(result.InitData)
			gi := func(off uintptr) int { return int(memmap.Int32(0x587000, off)) }
			switch v2 {
			case 2:
				*v6 = int32(rwRand(gi(211144), gi(211148)))
			case 4:
				*v6 = int32(rwRand(gi(211152), gi(211156)))
			case 8:
				*v6 = int32(rwRand(gi(211160), gi(211164)))
			case 0x10:
				*v6 = int32(rwRand(gi(211168), gi(211172)))
			default:
				*v6 = int32(rwRand(gi(211136), gi(211140)))
			}
		}
	} else {
		v3 := rwRand(1, 100)
		if v3 >= 50 {
			if v3 >= 90 {
				result = s.NewObjectByTypeID("DiamondGem")
			} else {
				result = s.NewObjectByTypeID("EmeraldGem")
			}
		} else {
			result = s.NewObjectByTypeID("RubyGem")
		}
	}
	return result
}

func rewardCreateGem2(a1 *server.Object, a2 uint32) *server.Object { // nox_xxx_createGem2_4F1F00
	return rewardCreateGem(a1, a2)
}
