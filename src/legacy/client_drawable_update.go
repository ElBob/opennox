package legacy

/*
#include "defs.h"
#include "GAME3_2.h"
#include "client__draw__fx.h"
#include "client__draw__glowdraw.h"
#include "client__drawable__update__charmup.h"
#include "client__drawable__update__cloud.h"
#include "client__drawable__update__dball.h"
#include "client__drawable__update__drainup.h"
#include "client__drawable__update__fireball.h"
#include "client__drawable__update__healup.h"
#include "client__drawable__update__manabomb.h"
#include "client__drawable__update__mmislup.h"
#include "client__drawable__update__mtailup.h"
#include "client__drawable__update__sparklup.h"
#include "client__drawable__update__telwake.h"
#include "client__drawable__update__vortexup.h"
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
)

type UpdateFunc = func(vp *noxrender.Viewport, dr *client.Drawable)

var (
	Nox_xxx_updDrawTeleportWake_4CD8D0   UpdateFunc
	Nox_xxx_updDrawVortexSource_4CC950   UpdateFunc
	Nox_xxx_updDrawManabombCharge_4CCAC0 UpdateFunc
	Nox_xxx_updDrawSparkleTrail_4CDBF0   UpdateFunc
	Nox_xxx_updDrawMagic_4CDD80          UpdateFunc
	Nox_xxx_updDrawMagicMissile_4CD9E0   UpdateFunc
	Nox_xxx_updDrawDBallCharge_4CE0C0    UpdateFunc
	Sub_4CD450                           UpdateFunc
	Sub_4CD690                           UpdateFunc
	Sub_4CE200                           func(dr *client.Drawable, n, rad int)
	Nox_xxx_updDrawAddRndSpark_4CDFA0    func(dr *client.Drawable, n int) *client.Drawable
	Sub_4CCEA0                           func(dr *client.Drawable, mult int)
	Sub_4CD150                           func(typ int, vp *noxrender.Viewport, dr *client.Drawable, dir bool)
)

func asDrawableInt(p C.int) *client.Drawable {
	return asDrawable((*nox_drawable)(unsafe.Pointer(uintptr(p))))
}

//export nox_xxx_updDrawTeleportWake_4CD8D0
func nox_xxx_updDrawTeleportWake_4CD8D0(vp, dr C.int) C.int {
	Nox_xxx_updDrawTeleportWake_4CD8D0(asViewportP(unsafe.Pointer(uintptr(vp))), asDrawableInt(dr))
	return 1
}

//export nox_xxx_updDrawVortexSource_4CC950
func nox_xxx_updDrawVortexSource_4CC950(vp, dr C.int) C.int {
	Nox_xxx_updDrawVortexSource_4CC950(asViewportP(unsafe.Pointer(uintptr(vp))), asDrawableInt(dr))
	return 1
}

//export nox_xxx_updDrawManabombCharge_4CCAC0
func nox_xxx_updDrawManabombCharge_4CCAC0(vp C.int, dr *C.uint32_t) C.int {
	Nox_xxx_updDrawManabombCharge_4CCAC0(asViewportP(unsafe.Pointer(uintptr(vp))), asDrawable((*nox_drawable)(unsafe.Pointer(dr))))
	return 1
}

//export nox_xxx_updDrawSparkleTrail_4CDBF0
func nox_xxx_updDrawSparkleTrail_4CDBF0(vp C.int, dr *C.uint32_t) C.int {
	Nox_xxx_updDrawSparkleTrail_4CDBF0(asViewportP(unsafe.Pointer(uintptr(vp))), asDrawable((*nox_drawable)(unsafe.Pointer(dr))))
	return 1
}

//export nox_xxx_updDrawMagic_4CDD80
func nox_xxx_updDrawMagic_4CDD80(vp C.int, dr *C.uint32_t) {
	Nox_xxx_updDrawMagic_4CDD80(asViewportP(unsafe.Pointer(uintptr(vp))), asDrawable((*nox_drawable)(unsafe.Pointer(dr))))
}

//export nox_xxx_updDrawMagicMissile_4CD9E0
func nox_xxx_updDrawMagicMissile_4CD9E0(vp C.int, dr *C.uint32_t) C.int {
	Nox_xxx_updDrawMagicMissile_4CD9E0(asViewportP(unsafe.Pointer(uintptr(vp))), asDrawable((*nox_drawable)(unsafe.Pointer(dr))))
	return 1
}

//export nox_xxx_updDrawDBallCharge_4CE0C0
func nox_xxx_updDrawDBallCharge_4CE0C0(vp, dr C.int) C.int {
	Nox_xxx_updDrawDBallCharge_4CE0C0(asViewportP(unsafe.Pointer(uintptr(vp))), asDrawableInt(dr))
	return 1
}

//export sub_4CD450
func sub_4CD450(vp *C.uint32_t, dr C.int) C.int {
	Sub_4CD450(asViewportP(unsafe.Pointer(vp)), asDrawableInt(dr))
	return 1
}

//export sub_4CD690
func sub_4CD690(vp *C.uint32_t, dr C.int) C.int {
	Sub_4CD690(asViewportP(unsafe.Pointer(vp)), asDrawableInt(dr))
	return 1
}

//export sub_4CE200
func sub_4CE200(a1, dr, n, rad C.int) C.int {
	Sub_4CE200(asDrawableInt(dr), int(n), int(rad))
	if n > 0 {
		return 0
	}
	return n
}

//export nox_xxx_updDrawAddRndSpark_4CDFA0
func nox_xxx_updDrawAddRndSpark_4CDFA0(dr C.int, n *C.uint32_t) *C.uint32_t {
	cnt := int(int32(uintptr(unsafe.Pointer(n))))
	last := Nox_xxx_updDrawAddRndSpark_4CDFA0(asDrawableInt(dr), cnt)
	if cnt <= 0 {
		// the original returns the count argument if nothing was spawned
		return n
	}
	return (*C.uint32_t)(last.C())
}

//export sub_4CCEA0
func sub_4CCEA0(dr *C.uint32_t, mult C.int) {
	Sub_4CCEA0(asDrawable((*nox_drawable)(unsafe.Pointer(dr))), int(mult))
}

//export sub_4CD150
func sub_4CD150(typ C.int, vp *C.uint32_t, dr, dir C.int) {
	Sub_4CD150(int(typ), asViewportP(unsafe.Pointer(vp)), asDrawableInt(dr), dir != 0)
}

func Get_nox_thing_pixie_dust_draw() unsafe.Pointer    { return C.nox_thing_pixie_dust_draw }
func Get_nox_thing_magic_sparkle_draw() unsafe.Pointer { return C.nox_thing_magic_sparkle_draw }
func Get_sub_4CE340() unsafe.Pointer                   { return C.sub_4CE340 }
