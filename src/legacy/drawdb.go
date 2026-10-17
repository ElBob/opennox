package legacy

/*
#include <stdlib.h>
#include "defs.h"
extern uint32_t nox_tile_def_cnt;
extern nox_tileDef_t nox_tile_defs_arr[176];
extern uint32_t dword_5d4594_251572;
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
)

// Sub_485CF0 frees images of tile definitions (nox_free_tile_defs).
func Sub_485CF0() {
	for i := 0; i < int(C.nox_tile_def_cnt); i++ {
		it := &C.nox_tile_defs_arr[i]
		if it.data_32 != nil {
			C.free(unsafe.Pointer(it.data_32))
			it.data_32 = nil
		}
	}
}

// Sub_485F30 frees data of edge definitions.
func Sub_485F30() {
	n := int(int32(C.dword_5d4594_251572))
	for i := 0; i < n; i++ {
		p := memmap.PtrPtr(0x85B3FC, 28676+60*uintptr(i))
		if *p != nil {
			C.free(*p)
			*p = nil
		}
	}
}

func freePtrAt(p unsafe.Pointer, off uintptr) unsafe.Pointer {
	v := *(*unsafe.Pointer)(unsafe.Add(p, off))
	if v != nil {
		C.free(v)
	}
	return v
}

// freeAnimFrames frees frames of an animation (nox_xxx_draw_44C780). Offset 16 is not a pointer.
func freeAnimFrames(p unsafe.Pointer) {
	for i := uintptr(0); i < 32; i += 4 {
		off := i
		if i >= 16 {
			off = i + 4
		}
		freePtrAt(p, off)
	}
}

func freeAnimAt(p unsafe.Pointer, off uintptr) {
	pp := (*unsafe.Pointer)(unsafe.Add(p, off))
	if *pp != nil {
		freeAnimFrames(unsafe.Add(*pp, 4))
		C.free(*pp)
	}
}

// freeMonsterAnims frees animations of monster draw data (sub_44C7B0).
func freeMonsterAnims(p unsafe.Pointer) {
	base := unsafe.Add(p, 52)
	for i := uintptr(0); i < 55; i++ {
		blk := unsafe.Add(base, i*66*4)
		freeAnimAt(blk, 0)
		for k := uintptr(0); k < 26; k++ {
			freeAnimAt(blk, (1+k)*4)
		}
		for k := uintptr(0); k < 27; k++ {
			freeAnimAt(blk, (27+k)*4)
		}
	}
}

// Nox_xxx_draw_44C650_free_kind frees draw data of a given kind.
func Nox_xxx_draw_44C650_free_kind(ptr unsafe.Pointer, kind int) {
	switch kind {
	case 2, 3:
		freePtrAt(ptr, 4)
	case 4:
		for i := uintptr(0); i < 5; i++ {
			freePtrAt(ptr, 4+4*i)
		}
	case 5:
		freeAnimFrames(unsafe.Add(ptr, 4))
	case 6:
		freeMonsterAnims(ptr)
	case 7:
		for i := uintptr(0); i < 16; i++ {
			freeAnimFrames(unsafe.Add(ptr, 8+48*i))
		}
	case 8:
		for i := uintptr(0); i < 3; i++ {
			freeAnimFrames(unsafe.Add(ptr, 8+48*i))
		}
	}
	C.free(ptr)
}
