package legacy

/*
#include "defs.h"
#include "client__draw__animdraw.h"
#include "client__draw__arrowdraw.h"
#include "client__draw__harpoondraw.h"
#include "client__draw__armordraw.h"
#include "client__draw__basedraw.h"
#include "client__draw__boulderdraw.h"
#include "client__draw__doordraw.h"
#include "client__draw__flagdraw.h"
#include "client__draw__glyphdraw.h"
#include "client__draw__powderdraw.h"
#include "client__draw__pressureplatedraw.h"
#include "client__draw__slavedraw.h"
#include "client__draw__staticdraw.h"
#include "client__draw__triggerdraw.h"
#include "client__draw__weapondraw.h"
#include "GAME1_1.h"
#include "GAME3_1.h"
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
)

type DrawFunc = func(vp *noxrender.Viewport, dr *client.Drawable) int

var (
	Nox_thing_static_draw           DrawFunc
	Nox_thing_static_random_draw    DrawFunc
	Nox_thing_slave_draw            DrawFunc
	Nox_thing_trigger_draw          DrawFunc
	Nox_thing_pressure_plate_draw   DrawFunc
	Nox_thing_black_powder_draw     DrawFunc
	Nox_thing_weapon_draw           DrawFunc
	Nox_thing_weapon_animate_draw   DrawFunc
	Nox_thing_spherical_shield_draw DrawFunc
	Nox_thing_armor_draw            DrawFunc
	Nox_thing_armor_animate_draw    DrawFunc
	Nox_thing_base_draw             DrawFunc
	Nox_thing_glyph_draw            DrawFunc
	Nox_thing_flag_draw             DrawFunc
	Nox_thing_boulder_draw          DrawFunc
	Nox_thing_door_draw             DrawFunc
)

//export nox_thing_static_draw
func nox_thing_static_draw(vp *C.uint32_t, dr *nox_drawable) C.int {
	return C.int(Nox_thing_static_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_static_random_draw
func nox_thing_static_random_draw(vp *C.uint32_t, dr *nox_drawable) C.int {
	return C.int(Nox_thing_static_random_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_slave_draw
func nox_thing_slave_draw(vp *C.int, dr *nox_drawable) C.int {
	return C.int(Nox_thing_slave_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_trigger_draw
func nox_thing_trigger_draw(vp *C.uint32_t, dr *nox_drawable) C.int {
	return C.int(Nox_thing_trigger_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_pressure_plate_draw
func nox_thing_pressure_plate_draw(vp *C.uint32_t, dr *nox_drawable) C.int {
	return C.int(Nox_thing_pressure_plate_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_black_powder_draw
func nox_thing_black_powder_draw(vp *C.uint32_t, dr *nox_drawable) C.int {
	return C.int(Nox_thing_black_powder_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_weapon_draw
func nox_thing_weapon_draw(vp *C.int, dr *nox_drawable) C.int {
	return C.int(Nox_thing_weapon_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_weapon_animate_draw
func nox_thing_weapon_animate_draw(vp *C.int, dr *nox_drawable) C.int {
	return C.int(Nox_thing_weapon_animate_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_spherical_shield_draw
func nox_thing_spherical_shield_draw(vp *C.int, dr *nox_drawable) C.int {
	return C.int(Nox_thing_spherical_shield_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_armor_draw
func nox_thing_armor_draw(vp *C.int, dr *nox_drawable) C.int {
	return C.int(Nox_thing_armor_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_armor_animate_draw
func nox_thing_armor_animate_draw(vp *C.int, dr *nox_drawable) C.int {
	return C.int(Nox_thing_armor_animate_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_base_draw
func nox_thing_base_draw(vp *C.int, dr *nox_drawable) C.int {
	return C.int(Nox_thing_base_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_glyph_draw
func nox_thing_glyph_draw(vp *C.int, dr *nox_drawable) C.int {
	return C.int(Nox_thing_glyph_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_flag_draw
func nox_thing_flag_draw(vp *C.int, dr *nox_drawable) C.int {
	return C.int(Nox_thing_flag_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_boulder_draw
func nox_thing_boulder_draw(vp *C.int, dr *nox_drawable) C.int {
	return C.int(Nox_thing_boulder_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_door_draw
func nox_thing_door_draw(vp *C.uint32_t, dr *nox_drawable) C.int {
	return C.int(Nox_thing_door_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

func Get_nox_thing_static_draw() unsafe.Pointer           { return C.nox_thing_static_draw }
func Get_nox_thing_slave_draw() unsafe.Pointer            { return C.nox_thing_slave_draw }
func Get_nox_thing_trigger_draw() unsafe.Pointer          { return C.nox_thing_trigger_draw }
func Get_nox_thing_pressure_plate_draw() unsafe.Pointer   { return C.nox_thing_pressure_plate_draw }
func Get_nox_thing_black_powder_draw() unsafe.Pointer     { return C.nox_thing_black_powder_draw }
func Get_nox_thing_weapon_draw() unsafe.Pointer           { return C.nox_thing_weapon_draw }
func Get_nox_thing_weapon_animate_draw() unsafe.Pointer   { return C.nox_thing_weapon_animate_draw }
func Get_nox_thing_spherical_shield_draw() unsafe.Pointer { return C.nox_thing_spherical_shield_draw }
func Get_nox_thing_armor_draw() unsafe.Pointer            { return C.nox_thing_armor_draw }
func Get_nox_thing_armor_animate_draw() unsafe.Pointer    { return C.nox_thing_armor_animate_draw }
func Get_nox_thing_base_draw() unsafe.Pointer             { return C.nox_thing_base_draw }
func Get_nox_thing_glyph_draw() unsafe.Pointer            { return C.nox_thing_glyph_draw }
func Get_nox_thing_flag_draw() unsafe.Pointer             { return C.nox_thing_flag_draw }
func Get_nox_thing_boulder_draw() unsafe.Pointer          { return C.nox_thing_boulder_draw }
func Get_nox_thing_door_draw() unsafe.Pointer             { return C.nox_thing_door_draw }
func Get_nox_thing_arrow_draw() unsafe.Pointer            { return C.nox_thing_arrow_draw }
func Get_nox_thing_weak_arrow_draw() unsafe.Pointer       { return C.nox_thing_weak_arrow_draw }
func Get_nox_thing_harpoon_draw() unsafe.Pointer          { return C.nox_thing_harpoon_draw }

func Sub_4B95D0(dr *client.Drawable) {
	C.sub_4B95D0((*nox_drawable)(dr.C()))
}

func Sub_4B96F0(dr *client.Drawable) {
	C.sub_4B96F0((*nox_drawable)(dr.C()))
}

func Sub_4B94E0(dr *client.Drawable) int {
	return int(C.sub_4B94E0((*nox_drawable)(dr.C())))
}

func Sub_418A80(a1 int) (string, bool) {
	p := C.sub_418A80(C.int(a1))
	if p == nil {
		return "", false
	}
	return GoWString((*wchar2_t)(unsafe.Pointer(p))), true
}
