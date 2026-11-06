package legacy

/*
#include <stdint.h>
#include "defs.h"
#include "GAME1_1.h"
#include <stdlib.h>
#include <string.h>
#include "server__script__script.h"

extern uint32_t nox_xxx_polygonNextIdx_587000_60352;
extern uint32_t nox_xxx_polygonNextAngle_587000_60356;
extern uint32_t dword_5d4594_588068;
extern uint32_t dword_5d4594_1556128;
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/cryptfile"
)

// All polygon data stays in the C-visible blob (0x5D4594) and in C globals:
// map points are 16-byte records at 535844, polygons are 140-byte records at 552228.

func polyMem(off uintptr) unsafe.Pointer { return memmap.PtrOff(0x5D4594, off) }

func polyRW(p unsafe.Pointer, n int) {
	// Result is ignored, as in the original.
	cryptfile.Global().ReadWrite(unsafe.Slice((*byte)(p), n))
}

func polyXfer(a1 unsafe.Pointer, a2 unsafe.Pointer) {
	C.nox_xxx_xferReadScriptHandler_4F5580(C.int(uintptr(a1)), (*C.char)(a2))
}

func polyF2I(v float32) uint32 { return uint32(C.nox_float2int(C.float(v))) }

func nox_xxx_polygonGetAngle_421030(a1 uint32) unsafe.Pointer {
	return polyMem(535844 + 16*uintptr(a1))
}

func sub_421040(a1 unsafe.Pointer) {
	v1 := 0
	if *(*uint16)(unsafe.Add(a1, 128)) != 0 {
		v2 := C.dword_5d4594_588068
		for {
			v3 := v2
			if v2 != 0 {
				v4 := *(*uintptr)(unsafe.Add(a1, 108))
				v5 := *(*uint32)(unsafe.Pointer(v4 + 4*uintptr(v1)))
				v6 := (*uint32)(unsafe.Pointer(v4 + 4*uintptr(v1)))
				ok := true
				for v5 != *(*uint32)(unsafe.Pointer(uintptr(v3) + 4)) {
					v3 = C.uint(*(*uint32)(unsafe.Pointer(uintptr(v3) + 8)))
					if v3 == 0 {
						ok = false
						break
					}
				}
				if ok {
					*v6 = *(*uint32)(unsafe.Pointer(uintptr(v3)))
					v2 = C.dword_5d4594_588068
				}
			}
			v1++
			if !(v1 < int(*(*uint16)(unsafe.Add(a1, 128)))) {
				break
			}
		}
	}
}

func nox_xxx_polygon_420CA0() unsafe.Pointer {
	v0 := uint32(1)
	if uint32(C.nox_xxx_polygonNextAngle_587000_60356) <= 1 {
		return nil
	}
	i := uintptr(535872)
	for *(*uint32)(polyMem(i)) == 0 {
		v0++
		if v0 >= uint32(C.nox_xxx_polygonNextAngle_587000_60356) {
			return nil
		}
		i += 16
	}
	return polyMem(535844 + 16*uintptr(v0))
}

func nox_xxx_polygon_420CD0(a1 unsafe.Pointer) unsafe.Pointer {
	v1 := *(*uint32)(a1) + 1
	if v1 >= uint32(C.nox_xxx_polygonNextAngle_587000_60356) {
		return nil
	}
	i := 535856 + 16*uintptr(v1)
	for *(*uint32)(polyMem(i)) == 0 {
		v1++
		if v1 >= uint32(C.nox_xxx_polygonNextAngle_587000_60356) {
			return nil
		}
		i += 16
	}
	return polyMem(535844 + 16*uintptr(v1))
}

func nox_xxx_polygonGetNext_4210A0() unsafe.Pointer {
	v0 := uint32(1)
	if uint32(C.nox_xxx_polygonNextIdx_587000_60352) <= 1 {
		return nil
	}
	i := uintptr(552452)
	for *(*uint32)(polyMem(i)) == 0 {
		v0++
		if v0 >= uint32(C.nox_xxx_polygonNextIdx_587000_60352) {
			return nil
		}
		i += 140
	}
	return polyMem(552228 + 140*uintptr(v0))
}

func sub_4210E0(a1 unsafe.Pointer) unsafe.Pointer {
	v1 := *(*uint32)(unsafe.Add(a1, 80)) + 1
	if v1 >= uint32(C.nox_xxx_polygonNextIdx_587000_60352) {
		return nil
	}
	i := 552312 + 140*uintptr(v1)
	for *(*uint32)(polyMem(i)) == 0 {
		v1++
		if v1 >= uint32(C.nox_xxx_polygonNextIdx_587000_60352) {
			return nil
		}
		i += 140
	}
	return polyMem(552228 + 140*uintptr(v1))
}

func sub_4D72D0(a1 int) int {
	result := int(C.dword_5d4594_1556128)
	*memmap.PtrUint32(0x5D4594, 1556132) = uint32(C.dword_5d4594_1556128)
	C.dword_5d4594_1556128 = C.uint(a1)
	return result
}

// nox_server_mapRWPolygons_428CD0 reads or writes the map polygons section.
func nox_server_mapRWPolygons_428CD0(a1 uint32) int {
	var (
		v25 uint32
		v26 uint32 = 4
		v27 uint32
		v28 uint32
		v29 [2]float32
		v30 int
		v31 uint32
		v32 [12]byte
	)
	if a1 != 0 {
		return 1
	}
	polyRW(unsafe.Pointer(&v26), 2)
	if int16(v26) > 4 {
		return 0
	}
	if cryptfile.Global().ReadOnly() {
		polyRW(unsafe.Pointer(&v25), 4)
		v9 := uint32(1)
		if v25 >= 1 {
			for {
				polyRW(unsafe.Pointer(&a1), 4)
				polyRW(unsafe.Pointer(&v29), 8)
				if C.nox_xxx_polygonSetAngle_420D40(C.int(*(*int32)(unsafe.Pointer(&v29[0]))), C.int(*(*int32)(unsafe.Pointer(&v29[1]))), C.uint(v9), C.int(a1)) == nil {
					return 0
				}
				if int16(v26) < 3 {
					polyRW(unsafe.Pointer(&v31), 4)
				}
				v9++
				if !(v9 <= v25) {
					break
				}
			}
		}
		polyRW(unsafe.Pointer(&v25), 4)
		v28 = 1
		if v25 >= 1 {
			for {
				v10 := unsafe.Pointer(C.sub_421230())
				if v10 == nil {
					return 0
				}
				polyRW(unsafe.Pointer(&v27), 1)
				polyRW(unsafe.Add(v10, 4), int(byte(v27)))
				*(*byte)(unsafe.Add(v10, int(byte(v27))+4)) = 0
				// Quirk: for versions < 3 the colors are read from the 12 byte block and the
				// single byte reads below are skipped.
				if int16(v26) >= 3 {
					polyRW(unsafe.Pointer(&a1), 1)
					*(*byte)(unsafe.Add(v10, 104)) = byte(a1)
					polyRW(unsafe.Pointer(&a1), 1)
					*(*byte)(unsafe.Add(v10, 105)) = byte(a1)
					polyRW(unsafe.Pointer(&a1), 1)
					*(*byte)(unsafe.Add(v10, 106)) = byte(a1)
				} else {
					polyRW(unsafe.Pointer(&v31), 4)
					polyRW(unsafe.Pointer(&v32[0]), 12)
					*(*byte)(unsafe.Add(v10, 104)) = v32[0]
					*(*byte)(unsafe.Add(v10, 105)) = v32[4]
					*(*byte)(unsafe.Add(v10, 106)) = v32[8]
				}
				polyRW(unsafe.Add(v10, 130), 1)
				v11 := (*uint16)(unsafe.Add(v10, 128))
				polyRW(unsafe.Add(v10, 128), 2)
				v12 := C.calloc(C.size_t(*(*uint16)(unsafe.Add(v10, 128))), 4)
				*(*unsafe.Pointer)(unsafe.Add(v10, 108)) = v12
				if v12 == nil {
					return 0
				}
				polyRW(v12, 4*int(*v11))
				sub_421040(v10)
				v13 := nox_xxx_polygonGetAngle_421030(**(**uint32)(unsafe.Add(v10, 108)))
				*(*uint32)(unsafe.Add(v10, 88)) = polyF2I(*(*float32)(unsafe.Add(v13, 4)))
				*(*uint32)(unsafe.Add(v10, 92)) = polyF2I(*(*float32)(unsafe.Add(v13, 8)))
				*(*uint32)(unsafe.Add(v10, 96)) = polyF2I(*(*float32)(unsafe.Add(v13, 4)))
				v14 := polyF2I(*(*float32)(unsafe.Add(v13, 8)))
				v15 := uint32(1)
				v16 := *v11 <= 1
				*(*uint32)(unsafe.Add(v10, 100)) = v14
				if !v16 {
					for {
						pts := *(*unsafe.Pointer)(unsafe.Add(v10, 108))
						v17 := nox_xxx_polygonGetAngle_421030(*(*uint32)(unsafe.Add(pts, 4*v15)))
						v18 := float64(*(*float32)(unsafe.Add(v17, 4)))
						v29[0] = *(*float32)(unsafe.Add(v17, 4))
						if v18 >= float64(*(*int32)(unsafe.Add(v10, 88))) {
							v19 := float64(*(*float32)(unsafe.Add(v17, 4)))
							v29[0] = *(*float32)(unsafe.Add(v17, 4))
							if v19 > float64(*(*int32)(unsafe.Add(v10, 96))) {
								*(*uint32)(unsafe.Add(v10, 96)) = polyF2I(v29[0])
							}
						} else {
							*(*uint32)(unsafe.Add(v10, 88)) = polyF2I(v29[0])
						}
						v20 := float64(*(*float32)(unsafe.Add(v17, 8)))
						v29[0] = *(*float32)(unsafe.Add(v17, 8))
						if v20 >= float64(*(*int32)(unsafe.Add(v10, 92))) {
							v21 := float64(*(*float32)(unsafe.Add(v17, 8)))
							v29[0] = *(*float32)(unsafe.Add(v17, 8))
							if v21 > float64(*(*int32)(unsafe.Add(v10, 100))) {
								*(*uint32)(unsafe.Add(v10, 100)) = polyF2I(v29[0])
							}
						} else {
							*(*uint32)(unsafe.Add(v10, 92)) = polyF2I(v29[0])
						}
						v15++
						if !(v15 < uint32(*v11)) {
							break
						}
					}
				}
				if int16(v26) >= 2 {
					v22 := *(*unsafe.Pointer)(v10)
					polyXfer(unsafe.Add(v10, 112), v22)
					var v23 unsafe.Pointer
					if v22 != nil {
						v23 = unsafe.Add(v22, 128)
					}
					polyXfer(unsafe.Add(v10, 120), v23)
				}
				if int16(v26) >= 4 {
					v24 := unsafe.Add(v10, 132)
					polyRW(v24, 4)
					if *(*byte)(v24)&1 != 0 {
						v30++
					}
				}
				v28++
				if !(v28 <= v25) {
					break
				}
			}
		}
		sub_4D72D0(v30)
		return 1
	}
	v25 = 0
	for i := nox_xxx_polygon_420CA0(); i != nil; v25++ {
		i = nox_xxx_polygon_420CD0(i)
	}
	polyRW(unsafe.Pointer(&v25), 4)
	v3 := nox_xxx_polygon_420CA0()
	for j := uint32(0); j < v25; j++ {
		polyRW(v3, 4)
		polyRW(unsafe.Add(v3, 4), 8)
		v3 = nox_xxx_polygon_420CD0(v3)
	}
	v25 = 0
	for k := nox_xxx_polygonGetNext_4210A0(); k != nil; v25++ {
		k = sub_4210E0(k)
	}
	polyRW(unsafe.Pointer(&v25), 4)
	v6 := nox_xxx_polygonGetNext_4210A0()
	v28 = 0
	if v25 <= 0 {
		return 1
	}
	for {
		v7 := *(*unsafe.Pointer)(v6)
		*(*byte)(unsafe.Pointer(&v27)) = byte(C.strlen((*C.char)(unsafe.Add(v6, 4))))
		polyRW(unsafe.Pointer(&v27), 1)
		polyRW(unsafe.Add(v6, 4), int(byte(v27)))
		*(*byte)(unsafe.Pointer(&a1)) = *(*byte)(unsafe.Add(v6, 104))
		polyRW(unsafe.Pointer(&a1), 1)
		*(*byte)(unsafe.Pointer(&a1)) = *(*byte)(unsafe.Add(v6, 105))
		polyRW(unsafe.Pointer(&a1), 1)
		*(*byte)(unsafe.Pointer(&a1)) = *(*byte)(unsafe.Add(v6, 106))
		polyRW(unsafe.Pointer(&a1), 1)
		polyRW(unsafe.Add(v6, 130), 1)
		polyRW(unsafe.Add(v6, 128), 2)
		polyRW(*(*unsafe.Pointer)(unsafe.Add(v6, 108)), 4*int(*(*uint16)(unsafe.Add(v6, 128))))
		polyXfer(unsafe.Add(v6, 112), v7)
		var v8 unsafe.Pointer
		if v7 != nil {
			v8 = unsafe.Add(v7, 128)
		}
		polyXfer(unsafe.Add(v6, 120), v8)
		polyRW(unsafe.Add(v6, 132), 4)
		v6 = sub_4210E0(v6)
		v28++
		if !(v28 < v25) {
			break
		}
	}
	return 1
}
