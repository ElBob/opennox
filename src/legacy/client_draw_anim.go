package legacy

/*
#include "defs.h"
#include "GAME3.h"
#include "client__draw__animdraw.h"
#include "client__draw__arrowdraw.h"
#include "client__draw__harpoondraw.h"
#include "client__draw__spiderspitdraw.h"
#include "client__draw__summondraw.h"
*/
import "C"
import (
	"image"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

var (
	Nox_thing_animate_draw              DrawFunc
	Nox_thing_arrow_draw                DrawFunc
	Nox_thing_weak_arrow_draw           DrawFunc
	Nox_thing_arrow_tail_link_draw      DrawFunc
	Nox_thing_weak_arrow_tail_link_draw DrawFunc
	Nox_thing_harpoon_draw              DrawFunc
	Nox_thing_harpoon_rope_draw         DrawFunc
	Nox_thing_spider_spit_draw          DrawFunc
)

//export nox_thing_animate_draw
func nox_thing_animate_draw(vp *C.uint32_t, dr *nox_drawable) C.int {
	return C.int(Nox_thing_animate_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_arrow_draw
func nox_thing_arrow_draw(vp *C.int, dr *nox_drawable) C.int {
	return C.int(Nox_thing_arrow_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_weak_arrow_draw
func nox_thing_weak_arrow_draw(vp *C.int, dr *nox_drawable) C.int {
	return C.int(Nox_thing_weak_arrow_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_arrow_tail_link_draw
func nox_thing_arrow_tail_link_draw(vp *C.uint32_t, dr *nox_drawable) C.int {
	return C.int(Nox_thing_arrow_tail_link_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_weak_arrow_tail_link_draw
func nox_thing_weak_arrow_tail_link_draw(vp *C.uint32_t, dr *nox_drawable) C.int {
	return C.int(Nox_thing_weak_arrow_tail_link_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_harpoon_draw
func nox_thing_harpoon_draw(vp *C.int, dr *nox_drawable) C.int {
	return C.int(Nox_thing_harpoon_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_harpoon_rope_draw
func nox_thing_harpoon_rope_draw(vp *C.int, dr *nox_drawable) C.int {
	return C.int(Nox_thing_harpoon_rope_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_spider_spit_draw
func nox_thing_spider_spit_draw(vp *C.uint32_t, dr *nox_drawable) C.int {
	return C.int(Nox_thing_spider_spit_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

func Get_nox_thing_animate_state_draw() unsafe.Pointer   { return C.nox_thing_animate_state_draw }
func Get_nox_thing_summon_effect_draw() unsafe.Pointer   { return C.nox_thing_summon_effect_draw }
func Get_nox_thing_arrow_tail_link_draw() unsafe.Pointer { return C.nox_thing_arrow_tail_link_draw }
func Get_nox_thing_weak_arrow_tail_link_draw() unsafe.Pointer {
	return C.nox_thing_weak_arrow_tail_link_draw
}
func Get_nox_thing_harpoon_rope_draw() unsafe.Pointer { return C.nox_thing_harpoon_rope_draw }
func Get_nox_thing_spider_spit_draw() unsafe.Pointer  { return C.nox_thing_spider_spit_draw }

func Sub_4B63B0(p1, p2 image.Point) {
	pts, free := alloc.Make([]C.int2{}, 2)
	defer free()
	pts[0] = C.int2{field_0: C.int(p1.X), field_4: C.int(p1.Y)}
	pts[1] = C.int2{field_0: C.int(p2.X), field_4: C.int(p2.Y)}
	C.sub_4B63B0(&pts[0], &pts[1])
}
