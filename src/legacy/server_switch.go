package legacy

import (
	"unsafe"


	"github.com/opennox/opennox/v1/common/memmap"
)

// This file is a literal port of sub_4DBA30 (server__system__server.c) and the functions only it uses:
// sub_4ECF10 (GAME3_3.c), sub_52BAF0 (GAME4_2.c), sub_516FC0 (GAME4_1.c).
// The original works on raw object memory with offsets and 32-bit pointers stored as ints; this port does the same.

func ssu8(p unsafe.Pointer, off int) *uint8   { return (*uint8)(unsafe.Add(p, off)) }
func ssu16(p unsafe.Pointer, off int) *uint16 { return (*uint16)(unsafe.Add(p, off)) }
func ssu32(p unsafe.Pointer, off int) *uint32 { return (*uint32)(unsafe.Add(p, off)) }

// ssp converts a 32-bit pointer stored in an int slot into a pointer (decompiler artefact: pointers as ints).
func ssp(v uint32) unsafe.Pointer { return unsafe.Pointer(uintptr(v)) }

// ssv converts a pointer into a 32-bit int slot value.
func ssv(p unsafe.Pointer) uint32 { return uint32(uintptr(p)) }

func ssobj(p unsafe.Pointer) *nox_object_t { return (*nox_object_t)(p) }

// sub_4ECF10 finds an object by its extent id (+44).
func sub_4ECF10(a1 uint32) uint32 {
	result := unsafe.Pointer(nox_server_getFirstObject_4DA790())
	for result != nil {
		if !(*ssu8(result, 16)&0x20 != 0 || *ssu32(result, 44) != a1) {
			return ssv(result)
		}
		v2 := *ssu32(result, 504)
		for v2 != 0 {
			if !(*ssu8(ssp(v2), 16)&0x20 != 0 || *ssu32(ssp(v2), 44) != a1) {
				return v2
			}
			v2 = *ssu32(ssp(v2), 496)
		}
		result = unsafe.Pointer(nox_server_getNextObject_4DA7A0(ssobj(result)))
	}
	result = unsafe.Pointer(nox_server_getFirstObjectUninited_4DA870())
	for result != nil {
		if *ssu8(result, 16)&0x20 == 0 && *ssu32(result, 44) == a1 {
			return ssv(result)
		}
		result = unsafe.Pointer(nox_server_getNextObjectUninited_4DA880(ssobj(result)))
	}
	obj := unsafe.Pointer(nox_xxx_getFirstUpdatable2Object_4DA840())
	for obj != nil {
		if *ssu8(obj, 16)&0x20 == 0 && *ssu32(obj, 44) == a1 {
			return ssv(obj)
		}
		obj = unsafe.Pointer(nox_xxx_getNextUpdatable2Object_4DA850(ssobj(obj)))
	}
	return ssv(obj)
}

// sub_52BAF0 resolves ids stored in the monster action stack into pointers.
func sub_52BAF0(a1 unsafe.Pointer) uint32 {
	v8 := 0
	v7 := *ssu32(a1, 748)
	result := v7
	if *ssu8(ssp(result), 544)&0x80 != 0 {
		return result
	}
	v3 := unsafe.Add(ssp(result), 552)
	for {
		v4 := 0
		v5 := int32(memmap.Uint32(0x587000, uintptr(255604+16*int64(*ssu32(v3, 0)))))
		if v5 > 0 {
			v6 := unsafe.Add(v3, 4)
			for {
				t := memmap.Uint32(0x587000, uintptr(255608+4*(int64(v4)+4*int64(*ssu32(v3, 0)))))
				if t == 1 {
					if *ssu32(v6, 0) != 0 {
						*ssu32(v6, 0) = sub_4ECF10(*ssu32(v6, 0))
						goto LABEL_12
					}
				} else {
					if t != 2 {
						goto LABEL_12
					}
					if *ssu32(v6, 0) != 0 {
						*ssu32(v6, 0) = ssv(unsafe.Pointer(nox_server_getWaypointById_579C40(int(int32(*ssu32(v6, 0))))))
						goto LABEL_12
					}
				}
				*ssu32(v6, 0) = 0
			LABEL_12:
				v4++
				v6 = unsafe.Add(v6, 8)
				if !(v4 < int(v5)) {
					break
				}
			}
			result = v7
		}
		// LABEL_14
		v3 = unsafe.Add(v3, 24)
		v8++
		if !(v8 <= int(int8(*ssu8(ssp(result), 544)))) {
			break
		}
	}
	return result
}

// sub_516FC0 applies pending owner assignments (list in C global dword_5d4594_2386920).
func sub_516FC0() {
	v0 := ssp(Get_dword_5d4594_2386920())
	if v0 != nil {
		for {
			v1 := sub_4ECF10(*ssu32(v0, 0))
			v2 := sub_4ECF10(*ssu32(v0, 4))
			if v1 != 0 && v2 != 0 {
				nox_xxx_unitSetOwner_4EC290(ssobj(ssp(v1)), ssobj(ssp(v2)))
			}
			v0 = ssp(*ssu32(v0, 8))
			if v0 == nil {
				break
			}
		}
	}
	Sub_516F30()
}

// Sub_4DBA30 is sub_4DBA30: fixes up object references after a map switch/load.
func Sub_4DBA30(a1 bool) {
	var a1i int
	if a1 {
		a1i = 1
	}
	v3 := 0 // always 0 where it is read (null pointer / zero constant in the original)
	v2 := unsafe.Pointer(nox_common_playerInfoFromNum_417090(31))
	if v2 == nil {
		return
	}
	v4 := *ssu32(v2, 2056)
	if v4 == 0 {
		return
	}
	if *memmap.PtrUint32(0x5D4594, 1563128) == 0 {
		*memmap.PtrUint32(0x5D4594, 1563128) = uint32(nox_xxx_getNameId_4E3AA0(internCStr("SaveGameLocation")))
		*memmap.PtrUint32(0x5D4594, 1563132) = uint32(nox_xxx_getNameId_4E3AA0(internCStr("Glyph")))
	}
	if a1i == 1 {
		v5 := unsafe.Pointer(nox_server_getFirstObject_4DA790())
		for v5 != nil {
			v26 := unsafe.Pointer(nox_server_getNextObject_4DA7A0(ssobj(v5)))
			if uint32(*ssu16(v5, 4)) == *memmap.PtrUint32(0x5D4594, 1563128) {
				Nox_xxx_unitMove_4E7010(AsObjectP(ssp(*ssu32(v2, 2056))), AsPointf(unsafe.Add(v5, 56)))
				*ssu32(ssp(*ssu32(v2, 2056)), 44) = *ssu32(v5, 44)
				*ssu32(v5, 44) = 0
				nox_xxx_delayedDeleteObject_4E5CC0(ssobj(v5))
			} else {
				v6 := int32(*ssu32(v5, 8))
				if v6&2 != 0 {
					v7 := ssp(*ssu32(v5, 748))
					v8 := 0
					if *ssu8(v7, 1129) != 0 {
						v9 := unsafe.Add(v7, 1132)
						for {
							v10 := sub_4ECF10(*ssu32(v9, 0))
							*ssu32(v9, 0) = v10
							if v10 == 0 {
								*ssu8(v7, 1129) = 0
							}
							v8++
							v9 = unsafe.Add(v9, 4)
							if !(v8 < int(*ssu8(v7, 1129))) {
								break
							}
						}
						v3 = 0
					}
					*ssu32(v7, 1196) = sub_4ECF10(*ssu32(v7, 1196))
					if int32(*ssu32(v5, 16)) >= 0 {
						sub_52BAF0(v5)
					}
					v11 := sub_4ECF10(*ssu32(v7, 392))
					if v11 != 0 {
						*ssu32(v7, 392) = *ssu32(ssp(v11), 36)
					} else {
						*ssu32(v7, 392) = 0
					}
					v12 := sub_4ECF10(*ssu32(v7, 1200))
					if v12 != 0 {
						*ssu32(v7, 1200) = *ssu32(ssp(v12), 36)
					} else {
						*ssu32(v7, 1200) = 0
					}
					*ssu32(v7, 1216) = sub_4ECF10(*ssu32(v7, 1216))
					v13 := 0
					if *ssu8(v7, 2172) != 0 {
						v14 := unsafe.Add(v7, 2140)
						for {
							v15 := sub_4ECF10(*ssu32(v14, 0))
							if v15 != 0 {
								*ssu32(v14, 0) = *ssu32(ssp(v15), 36)
							} else {
								*ssu8(v7, 2172) = 0
							}
							v13++
							v14 = unsafe.Add(v14, 4)
							if !(v13 < int(*ssu8(v7, 2172))) {
								break
							}
						}
					}
				} else if v6&0x4000 != 0 {
					if *ssu32(ssp(*ssu32(v5, 748)), 16) != 0 {
						AsObjectP(v5).NeedSync() // nox_xxx_unitNeedSync_4E44F0
					}
				} else if v6&0x8000 == 0 {
					if uint8(v6)&0x80 != 0 && *ssu32(ssp(*ssu32(v5, 748)), 12) != *ssu32(ssp(*ssu32(v5, 748)), 4) {
						nox_xxx_unitAddToUpdatable_4DA8D0(ssobj(v5))
					}
				} else {
					v23 := *ssu32(ssp(*ssu32(v5, 748)), 4)
					if v23 != 0 && *ssu32(ssp(*ssu32(ssp(v23), 748)), 16) != 0 {
						AsObjectP(v5).NeedSync() // nox_xxx_unitNeedSync_4E44F0
					}
				}
			}
			v5 = v26
		}
		nox_script_activatorResolveObjs_51B0C0()
		sub_516FC0()
		if Get_dword_5d4594_1563096() != 0 {
			v16 := unsafe.Pointer(nox_server_getFirstObject_4DA790())
			for v16 != nil {
				v3p := unsafe.Pointer(nox_server_getNextObject_4DA7A0(ssobj(v16)))
				if int32(*ssu32(v16, 16)) >= 0 && Nox_xxx_isUnit_4E5B50(AsObjectP(v16)) != 0 {
					if *ssu8(v16, 8)&2 != 0 {
						v17 := int32(*ssu32(v16, 12))
						if v17&0x2000 != 0 {
							v18 := *ssu32(v16, 504) // nox_xxx_inventoryGetFirst_4E7980
							for v18 != 0 {
								v19 := *ssu32(ssp(v18), 496) // nox_xxx_inventoryGetNext_4E7990
								if uint32(*ssu16(ssp(v18), 4)) == *memmap.PtrUint32(0x5D4594, 1563132) {
									nox_xxx_delayedDeleteObject_4E5CC0(ssobj(ssp(v18)))
								}
								v18 = v19
							}
						}
					}
					nox_xxx_delayedDeleteObject_4E5CC0(ssobj(v16))
				}
				v16 = v3p
			}
			obj := unsafe.Pointer(nox_xxx_getFirstUpdatable2Object_4DA840())
			for obj != nil {
				v21 := unsafe.Pointer(nox_xxx_getNextUpdatable2Object_4DA850(ssobj(obj)))
				if int32(*ssu32(obj, 16)) >= 0 {
					if Sub_4E5B80(AsObjectP(obj)) != 0 {
						nox_xxx_delayedDeleteObject_4E5CC0(ssobj(obj))
					}
				}
				obj = v21
			}
		}
	}
	for i := *ssu32(ssp(v4), 516); i != uint32(v3); i = *ssu32(ssp(i), 512) {
		if *ssu8(ssp(i), 8)&2 != 0 {
			if *ssu8(ssp(*ssu32(ssp(i), 748)), 1440)&0x80 != 0 {
				Nox_xxx_netReportAcquireCreature_4D91A0(int(*ssu8(v2, 2064)), AsObjectP(ssp(i)))
				nox_xxx_netMarkMinimapObject_417190(int(*ssu8(v2, 2064)), ssobj(ssp(i)), 1)
				continue
			}
			if *ssu8(ssp(i), 12)&0x80 != 0 {
				Nox_xxx_netMonitorCreature_4D9250(int(*ssu8(v2, 2064)), AsObjectP(ssp(i)))
				nox_xxx_netMarkMinimapObject_417190(int(*ssu8(v2, 2064)), ssobj(ssp(i)), 1)
				continue
			}
		}
	}
	nox_xxx_gameSetSwitchSolo_4DB220(v3)
	Set_dword_5d4594_1563096(v3)
	nox_ticks_reset_416D40()
}


