package legacy

/*
#include "defs.h"
#include "GAME2_3.h"
#include "GAME3_1.h"
#include "client__draw__canidraw.h"
#include "client__draw__drawrays.h"
#include "client__draw__fx.h"
#include "client__draw__maidendraw.h"
#include "client__draw__mgendraw.h"
#include "client__draw__partrain.h"
#include "client__draw__playerdraw.h"
#include "client__draw__udeddraw.h"
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/client"
)

var (
	Nox_thing_maiden_draw          DrawFunc
	Nox_thing_player_waypoint_draw DrawFunc
	Nox_thing_monster_gen_draw     DrawFunc
	Nox_thing_cond_animate_draw    DrawFunc
	Nox_thing_blue_rain_draw       DrawFunc
	Nox_thing_levelup_draw         DrawFunc
	Nox_thing_oblivion_up_draw     DrawFunc
	Nox_thing_undead_killer_draw   DrawFunc
	Nox_thing_falling_sparks_draw  func(typ int, vp unsafe.Pointer, dr *client.Drawable)
	Nox_xxx_netDrawRays_49BDD0     func(data []byte)
)

//export nox_thing_maiden_draw
func nox_thing_maiden_draw(vp *C.uint32_t, dr *nox_drawable) C.int {
	return C.int(Nox_thing_maiden_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_player_waypoint_draw
func nox_thing_player_waypoint_draw(vp C.int, dr *nox_drawable) C.int {
	return C.int(Nox_thing_player_waypoint_draw(asViewportP(unsafe.Pointer(uintptr(vp))), asDrawable(dr)))
}

//export nox_thing_monster_gen_draw
func nox_thing_monster_gen_draw(vp *C.int, dr *nox_drawable) C.int {
	return C.int(Nox_thing_monster_gen_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_cond_animate_draw
func nox_thing_cond_animate_draw(vp *C.uint32_t, dr *nox_drawable) C.int {
	return C.int(Nox_thing_cond_animate_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_blue_rain_draw
func nox_thing_blue_rain_draw(vp C.int, dr *nox_drawable) C.int {
	return C.int(Nox_thing_blue_rain_draw(asViewportP(unsafe.Pointer(uintptr(vp))), asDrawable(dr)))
}

//export nox_thing_levelup_draw
func nox_thing_levelup_draw(vp C.int, dr *nox_drawable) C.int {
	return C.int(Nox_thing_levelup_draw(asViewportP(unsafe.Pointer(uintptr(vp))), asDrawable(dr)))
}

//export nox_thing_oblivion_up_draw
func nox_thing_oblivion_up_draw(vp C.int, dr *nox_drawable) C.int {
	return C.int(Nox_thing_oblivion_up_draw(asViewportP(unsafe.Pointer(uintptr(vp))), asDrawable(dr)))
}

//export nox_thing_undead_killer_draw
func nox_thing_undead_killer_draw(vp *C.int, dr *nox_drawable) C.int {
	return C.int(Nox_thing_undead_killer_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_xxx_netDrawRays_49BDD0
func nox_xxx_netDrawRays_49BDD0(data *C.uchar) C.int {
	// The return value is not used by any callers.
	Nox_xxx_netDrawRays_49BDD0(unsafe.Slice((*byte)(unsafe.Pointer(data)), 9))
	return 0
}

func Get_nox_thing_player_waypoint_draw() unsafe.Pointer { return C.nox_thing_player_waypoint_draw }
func Get_nox_thing_monster_gen_draw() unsafe.Pointer     { return C.nox_thing_monster_gen_draw }
func Get_nox_thing_cond_animate_draw() unsafe.Pointer    { return C.nox_thing_cond_animate_draw }
func Get_nox_thing_blue_rain_draw() unsafe.Pointer       { return C.nox_thing_blue_rain_draw }
func Get_nox_thing_levelup_draw() unsafe.Pointer         { return C.nox_thing_levelup_draw }
func Get_nox_thing_oblivion_up_draw() unsafe.Pointer     { return C.nox_thing_oblivion_up_draw }
func Get_nox_thing_undead_killer_draw() unsafe.Pointer   { return C.nox_thing_undead_killer_draw }

func Nox_xxx_spriteDrawCircleMB_4C32A0(x, y, r int, cl uint32) {
	C.nox_xxx_spriteDrawCircleMB_4C32A0(C.int(x), C.int(y), C.int(r), C.int(cl))
}

func Sub_48C690(x1, y1, x2, y2 int) uint32 {
	return uint32(C.sub_48C690(C.int(x1), C.int(y1), C.int(x2), C.int(y2)))
}
