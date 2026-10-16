package legacy

/*
#include <stdlib.h>
#include "defs.h"
#include "GAME3_1.h"
#include "client__draw__lightning.h"
#include "client__draw__plasma.h"
extern uint32_t dword_5d4594_1316408;
extern uint32_t dword_5d4594_1316412;
*/
import "C"
import (
	"image"
	"unsafe"
)

var (
	Nox_thing_lightning_draw            DrawFunc
	Nox_thing_chain_lightning_bolt_draw DrawFunc
	Nox_thing_energy_bolt_draw          DrawFunc
	Nox_thing_green_bolt_draw           DrawFunc
	Nox_thing_plasma_draw               DrawFunc
	Sub_4BA8B0                          func(p1, p2 image.Point, ribbon int) byte
)

//export nox_thing_lightning_draw
func nox_thing_lightning_draw(vp *C.int, dr *nox_drawable) C.int {
	return C.int(Nox_thing_lightning_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_chain_lightning_bolt_draw
func nox_thing_chain_lightning_bolt_draw(vp *C.int, dr *nox_drawable) C.int {
	return C.int(Nox_thing_chain_lightning_bolt_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_energy_bolt_draw
func nox_thing_energy_bolt_draw(vp *C.int, dr *nox_drawable) C.int {
	return C.int(Nox_thing_energy_bolt_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_green_bolt_draw
func nox_thing_green_bolt_draw(vp *C.int, dr *nox_drawable) C.int {
	return C.int(Nox_thing_green_bolt_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_plasma_draw
func nox_thing_plasma_draw(vp *C.int, dr *nox_drawable) C.int {
	return C.int(Nox_thing_plasma_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export sub_4BA8B0
func sub_4BA8B0(a1, a2, a3 *C.int) C.char {
	p1 := *(*[2]int32)(unsafe.Pointer(a1))
	p2 := *(*[2]int32)(unsafe.Pointer(a2))
	return C.char(Sub_4BA8B0(image.Pt(int(p1[0]), int(p1[1])), image.Pt(int(p2[0]), int(p2[1])), int(*a3)))
}

func Get_nox_thing_lightning_draw() unsafe.Pointer { return C.nox_thing_lightning_draw }
func Get_nox_thing_chain_lightning_bolt_draw() unsafe.Pointer {
	return C.nox_thing_chain_lightning_bolt_draw
}
func Get_nox_thing_energy_bolt_draw() unsafe.Pointer { return C.nox_thing_energy_bolt_draw }
func Get_nox_thing_green_bolt_draw() unsafe.Pointer  { return C.nox_thing_green_bolt_draw }
func Get_nox_thing_plasma_draw() unsafe.Pointer      { return C.nox_thing_plasma_draw }

func Get_dword_5d4594_1316408() int32  { return int32(C.dword_5d4594_1316408) }
func Get_dword_5d4594_1316412() int32  { return int32(C.dword_5d4594_1316412) }
func Set_dword_5d4594_1316412(v int32) { C.dword_5d4594_1316412 = C.uint32_t(v) }

func Sub_4BA670(a1, a2, a3, a4, a5 int) {
	C.sub_4BA670(C.int(a1), C.int(a2), C.int(a3), C.int(a4), C.int(a5))
}

func Sub_4BE800(cl uint32) {
	C.sub_4BE800(C.int(cl))
}

func Sub_4BE810(a1 int, cl uint32, a3 int, a4 byte) {
	C.sub_4BE810(C.int(a1), C.int(cl), C.int(a3), C.char(a4))
}

func Sub_4BEAD0(p1, p2, p3, p4 image.Point, a5, a6 int) {
	pts := (*[4]C.int2)(C.calloc(4, C.size_t(unsafe.Sizeof(C.int2{}))))
	defer C.free(unsafe.Pointer(pts))
	for i, p := range [4]image.Point{p1, p2, p3, p4} {
		pts[i] = C.int2{field_0: C.int(p.X), field_4: C.int(p.Y)}
	}
	C.sub_4BEAD0(&pts[0], &pts[1], &pts[2], &pts[3], C.int(a5), C.int(a6))
}
