package legacy

/*
#include "GAME2_3.h"
#include "client__draw__debugdraw.h"
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/client"
)

func init() {
	client.ThingDrawDefault = C.nox_thing_debug_draw
}

type nox_thing = C.nox_thing

//export nox_xxx_getTTByNameSpriteMB_44CFC0
func nox_xxx_getTTByNameSpriteMB_44CFC0(cstr *C.char) int {
	id := GoString(cstr)
	return GetClient().Cli().Things.IndByID(id)
}

//export sub_44D330
func sub_44D330(cstr *C.char) *nox_thing {
	id := GoString(cstr)
	return (*C.nox_thing)(GetClient().Cli().Things.TypeByID(id).C())
}

//export nox_get_thing_name
func nox_get_thing_name(i int) *C.char {
	t := GetClient().Cli().Things.TypeByInd(i)
	if t == nil {
		return nil
	}
	return (*C.char)(unsafe.Pointer(t.Name))
}

//export nox_get_thing
func nox_get_thing(i int) *nox_thing {
	t := GetClient().Cli().Things.TypeByInd(i)
	if t == nil {
		return nil
	}
	return (*nox_thing)(t.C())
}

//export nox_get_thing_pretty_name
func nox_get_thing_pretty_name(i int) *wchar2_t {
	t := GetClient().Cli().Things.TypeByInd(i)
	if t == nil {
		return nil
	}
	return (*wchar2_t)(unsafe.Pointer(t.PrettyName))
}

//export nox_get_thing_desc
func nox_get_thing_desc(i int) *wchar2_t {
	t := GetClient().Cli().Things.TypeByInd(i)
	if t == nil {
		return nil
	}
	return (*wchar2_t)(unsafe.Pointer(t.Desc))
}

//export nox_get_thing_pretty_image
func nox_get_thing_pretty_image(i int) int {
	t := GetClient().Cli().Things.TypeByInd(i)
	if t == nil {
		return 0
	}
	return int(t.PrettyImage)
}

//export nox_drawable_link_thing
func nox_drawable_link_thing(a1c *nox_drawable, i int) int {
	return GetClient().Cli().DrawableLinkThing(asDrawable(a1c), i)
}
