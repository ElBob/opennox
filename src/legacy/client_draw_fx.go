package legacy

/*
#include "defs.h"
#include "GAME1_2.h"
#include "GAME3.h"
#include "GAME3_1.h"
#include "client__draw__fx.h"
#include "client__draw__glowdraw.h"
extern uint32_t dword_5d4594_1313532;
extern uint32_t dword_5d4594_1313536;
extern uint32_t dword_5d4594_1313540;
extern uint32_t dword_5d4594_1313564;
*/
import "C"
import (
	"image"
	"unsafe"

	"github.com/opennox/opennox/v1/client"
)

var (
	Sub_499490                            func(typ int, pts [4]uint16, a3, a4 int, a5, a6 byte)
	Sub_499520                            func(typ int, pts [4]int16, a3 int16, a4, a5 byte)
	Nox_xxx_makePointFxCli_499610         func(typ, n, a3, a4 int, pos image.Point)
	Nox_xxx_drawEnergyBolt_499710         func(pos image.Point, z int16, typ int)
	Sub_499950                            func(typ int, p1, p2 image.Point, a4 uint16, a5 int8) *client.Drawable
	Nox_xxx_makeLightningParticles_4999D0 func(typ int, p1, p2 image.Point) int
	Nox_xxx_draw_499E70                   func(a1, a2, a3, a4, a5, a6, a7 int)
	Sub_49A150                            func(pos image.Point, typ int, a3 byte)

	Nox_thing_magic_sparkle_draw    DrawFunc
	Nox_thing_pixie_draw            DrawFunc
	Nox_thing_pixie_dust_draw       DrawFunc
	Nox_thing_blue_rain_spark_draw  DrawFunc
	Nox_thing_rain_orb_draw         DrawFunc
	Nox_thing_red_spark_draw        DrawFunc
	Nox_thing_blue_spark_draw       DrawFunc
	Nox_thing_cyan_spark_draw       DrawFunc
	Nox_thing_green_spark_draw      DrawFunc
	Nox_thing_yellow_spark_draw     DrawFunc
	Nox_thing_violet_spark_draw     DrawFunc
	Nox_thing_death_ball_spark_draw DrawFunc
	Nox_thing_white_spark_draw      DrawFunc
	Nox_thing_particle_draw         DrawFunc
	Nox_thing_glow_orb_draw         DrawFunc
	Nox_thing_glow_orb_move_draw    DrawFunc
)

func asPoint2(p *C.int2) image.Point {
	return image.Pt(int(p.field_0), int(p.field_4))
}

//export sub_499490
func sub_499490(typ C.int, pts *C.uint16_t, a3, a4 C.int, a5, a6 C.char) {
	Sub_499490(int(typ), *(*[4]uint16)(unsafe.Pointer(pts)), int(a3), int(a4), byte(a5), byte(a6))
}

//export sub_499520
func sub_499520(typ C.int, pts *C.short, a3 C.short, a4, a5 C.char) {
	Sub_499520(int(typ), *(*[4]int16)(unsafe.Pointer(pts)), int16(a3), byte(a4), byte(a5))
}

//export nox_xxx_makePointFxCli_499610
func nox_xxx_makePointFxCli_499610(typ, n, a3, a4, x, y C.int) C.int {
	Nox_xxx_makePointFxCli_499610(int(typ), int(n), int(a3), int(a4), image.Pt(int(x), int(y)))
	if n > 0 {
		return 0
	}
	return n
}

//export nox_xxx_drawEnergyBolt_499710
func nox_xxx_drawEnergyBolt_499710(x, y C.int, z C.short, typ C.int) C.int {
	Nox_xxx_drawEnergyBolt_499710(image.Pt(int(x), int(y)), int16(z), int(typ))
	return 0
}

//export sub_499950
func sub_499950(typ C.int, p1, p2 *C.int2, a4 C.ushort, a5 C.char) C.int {
	dr := Sub_499950(int(typ), asPoint2(p1), asPoint2(p2), uint16(a4), int8(a5))
	return C.int(uintptr(dr.C()))
}

//export nox_xxx_makeLightningParticles_4999D0
func nox_xxx_makeLightningParticles_4999D0(typ C.int, p1, p2 *C.int2) C.int {
	return C.int(Nox_xxx_makeLightningParticles_4999D0(int(typ), asPoint2(p1), asPoint2(p2)))
}

//export nox_xxx_draw_499E70
func nox_xxx_draw_499E70(a1, a2, a3, a4, a5, a6, a7 C.int) C.int {
	Nox_xxx_draw_499E70(int(a1), int(a2), int(a3), int(a4), int(a5), int(a6), int(a7))
	return 0
}

//export sub_49A150
func sub_49A150(pos *C.int2, typ C.int, a3 C.uchar) C.int {
	Sub_49A150(asPoint2(pos), int(typ), byte(a3))
	return 0
}

//export nox_thing_magic_sparkle_draw
func nox_thing_magic_sparkle_draw(vp *C.int, dr *nox_drawable) C.int {
	return C.int(Nox_thing_magic_sparkle_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_pixie_draw
func nox_thing_pixie_draw(vp *C.int, dr *nox_drawable) C.int {
	return C.int(Nox_thing_pixie_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_pixie_dust_draw
func nox_thing_pixie_dust_draw(vp *C.int, dr *nox_drawable) C.int {
	return C.int(Nox_thing_pixie_dust_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_blue_rain_spark_draw
func nox_thing_blue_rain_spark_draw(vp *C.uint32_t, dr *nox_drawable) C.int {
	return C.int(Nox_thing_blue_rain_spark_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_rain_orb_draw
func nox_thing_rain_orb_draw(vp *C.uint32_t, dr *nox_drawable) C.int {
	return C.int(Nox_thing_rain_orb_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_red_spark_draw
func nox_thing_red_spark_draw(vp *C.uint32_t, dr *nox_drawable) C.int {
	return C.int(Nox_thing_red_spark_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_blue_spark_draw
func nox_thing_blue_spark_draw(vp *C.uint32_t, dr *nox_drawable) C.int {
	return C.int(Nox_thing_blue_spark_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_cyan_spark_draw
func nox_thing_cyan_spark_draw(vp *C.uint32_t, dr *nox_drawable) C.int {
	return C.int(Nox_thing_cyan_spark_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_green_spark_draw
func nox_thing_green_spark_draw(vp *C.uint32_t, dr *nox_drawable) C.int {
	return C.int(Nox_thing_green_spark_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_yellow_spark_draw
func nox_thing_yellow_spark_draw(vp *C.uint32_t, dr *nox_drawable) C.int {
	return C.int(Nox_thing_yellow_spark_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_violet_spark_draw
func nox_thing_violet_spark_draw(vp *C.uint32_t, dr *nox_drawable) C.int {
	return C.int(Nox_thing_violet_spark_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_death_ball_spark_draw
func nox_thing_death_ball_spark_draw(vp *C.uint32_t, dr *nox_drawable) C.int {
	return C.int(Nox_thing_death_ball_spark_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_white_spark_draw
func nox_thing_white_spark_draw(vp *C.uint32_t, dr *nox_drawable) C.int {
	return C.int(Nox_thing_white_spark_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_particle_draw
func nox_thing_particle_draw(vp C.int, dr *nox_drawable) C.int {
	return C.int(Nox_thing_particle_draw(asViewportP(unsafe.Pointer(uintptr(vp))), asDrawable(dr)))
}

//export nox_thing_glow_orb_draw
func nox_thing_glow_orb_draw(vp *C.int, dr *nox_drawable) C.int {
	return C.int(Nox_thing_glow_orb_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

//export nox_thing_glow_orb_move_draw
func nox_thing_glow_orb_move_draw(vp *C.int, dr *nox_drawable) C.int {
	return C.int(Nox_thing_glow_orb_move_draw(asViewportP(unsafe.Pointer(vp)), asDrawable(dr)))
}

func Get_nox_thing_pixie_draw() unsafe.Pointer           { return C.nox_thing_pixie_draw }
func Get_nox_thing_blue_rain_spark_draw() unsafe.Pointer { return C.nox_thing_blue_rain_spark_draw }
func Get_nox_thing_rain_orb_draw() unsafe.Pointer        { return C.nox_thing_rain_orb_draw }
func Get_nox_thing_violet_spark_draw() unsafe.Pointer    { return C.nox_thing_violet_spark_draw }
func Get_nox_thing_death_ball_spark_draw() unsafe.Pointer {
	return C.nox_thing_death_ball_spark_draw
}
func Get_nox_thing_particle_draw() unsafe.Pointer      { return C.nox_thing_particle_draw }
func Get_nox_thing_white_spark_draw() unsafe.Pointer   { return C.nox_thing_white_spark_draw }
func Get_nox_thing_glow_orb_draw() unsafe.Pointer      { return C.nox_thing_glow_orb_draw }
func Get_nox_thing_glow_orb_move_draw() unsafe.Pointer { return C.nox_thing_glow_orb_move_draw }
func Get_sub_4CA720() unsafe.Pointer                   { return C.sub_4CA720 }

func Sub_4B6970(vp unsafe.Pointer, dr *client.Drawable, a3, a4 uint32) int {
	return int(C.sub_4B6970((*C.uint32_t)(vp), (*nox_drawable)(dr.C()), C.int(a3), C.int(a4)))
}

func Sub_4B6B80(vp unsafe.Pointer, dr *client.Drawable, a3 int) int {
	return int(C.sub_4B6B80((*C.int)(vp), (*nox_drawable)(dr.C()), C.int(a3)))
}

func Nox_client_newScreenParticle_431540(a1, a2, a3, a4, a5, a6, a7, a8, a9, a10 int) {
	C.nox_client_newScreenParticle_431540(C.int(a1), C.int(a2), C.int(a3), C.int(a4), C.int(a5), C.int(a6), C.char(a7), C.char(a8), C.char(a9), C.char(a10))
}

func Get_dword_5d4594_1313532() uint32 { return uint32(C.dword_5d4594_1313532) }
func Get_dword_5d4594_1313536() uint32 { return uint32(C.dword_5d4594_1313536) }
func Get_dword_5d4594_1313540() uint32 { return uint32(C.dword_5d4594_1313540) }
func Get_dword_5d4594_1313564() uint32 { return uint32(C.dword_5d4594_1313564) }
