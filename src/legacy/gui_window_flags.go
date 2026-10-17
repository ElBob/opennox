package legacy

/*
#include <stdbool.h>
#include "client__gui__window.h"
*/
import "C"
import (
	"image"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
)

func asWindowInt(p C.int) *gui.Window {
	return AsWindowP(unsafe.Pointer(uintptr(p)))
}

//export sub_46ACE0
func sub_46ACE0(a1 *C.uint, from, to, hidden C.int) {
	win := AsWindowP(unsafe.Pointer(a1))
	for i := from; i <= to; i++ {
		nox_window_set_hidden((*nox_window)(win.ChildByID(uint(i)).C()), int(hidden))
	}
}

//export sub_46AD20
func sub_46AD20(a1 *C.uint, from, to, enabled C.int) {
	win := AsWindowP(unsafe.Pointer(a1))
	for i := from; i <= to; i++ {
		nox_xxx_wnd_46ABB0((*nox_window)(win.ChildByID(uint(i)).C()), enabled)
	}
}

//export nox_gui_getWindowOffs_46AA20
func nox_gui_getWindowOffs_46AA20(p *nox_window, px, py *C.uint) C.int {
	if p == nil {
		*px, *py = 0, 0
		return -2
	}
	off := asWindow(p).Off
	*px, *py = C.uint(off.X), C.uint(off.Y)
	return 0
}

//export nox_client_wndGetPosition_46AA60
func nox_client_wndGetPosition_46AA60(p *nox_window, px, py *C.uint) C.int {
	if p == nil {
		return -2
	}
	pos := asWindow(p).GlobalPos()
	*px, *py = C.uint(pos.X), C.uint(pos.Y)
	return 0
}

//export nox_xxx_wndPointInWnd_46AAB0
func nox_xxx_wndPointInWnd_46AAB0(p *C.uint, x, y C.int) C.bool {
	return C.bool(AsWindowP(unsafe.Pointer(p)).PointInGlobal(image.Pt(int(x), int(y))))
}

//export sub_46AB20
func sub_46AB20(p *C.uint, w, h C.int) C.int {
	if p == nil {
		return -2
	}
	AsWindowP(unsafe.Pointer(p)).SetSize(image.Pt(int(w), int(h)))
	return 0
}

//export nox_window_get_size
func nox_window_get_size(p *nox_window, outW, outH *C.int) C.int {
	if p == nil {
		*outW, *outH = 0, 0
		return -2
	}
	sz := asWindow(p).SizeVal
	*outW, *outH = C.int(sz.X), C.int(sz.Y)
	return 0
}

//export nox_xxx_wnd_46ABB0
func nox_xxx_wnd_46ABB0(p *nox_window, enabled C.int) C.int {
	if p == nil {
		return -2
	}
	asWindow(p).SetEnabled(enabled != 0)
	return 0
}

//export nox_xxx_wnd_46AD60
func nox_xxx_wnd_46AD60(p, flags C.int) C.int {
	if p == 0 {
		return -2
	}
	win := asWindowInt(p)
	prev := win.Flags
	win.Flags |= gui.StatusFlags(flags)
	return C.int(prev)
}

//export nox_xxx_wndClearFlag_46AD80
func nox_xxx_wndClearFlag_46AD80(p, flags C.int) C.int {
	if p == 0 {
		return -2
	}
	win := asWindowInt(p)
	prev := win.Flags
	win.Flags &^= gui.StatusFlags(flags)
	return C.int(prev)
}

//export nox_xxx_wndGetFlags_46ADA0
func nox_xxx_wndGetFlags_46ADA0(p C.int) C.int {
	if p == 0 {
		return -2
	}
	return C.int(asWindowInt(p).Flags)
}

//export nox_window_is_child
func nox_window_is_child(a1, a2 *nox_window) C.int {
	if a1 == nil || a2 == nil {
		return 0
	}
	for cur := asWindow(a2).Parent(); cur != nil; cur = cur.Parent() {
		if cur == asWindow(a1) {
			return 1
		}
	}
	return 0
}

//export nox_xxx_wnd_46B280
func nox_xxx_wnd_46B280(p, w C.int) C.int {
	if p == 0 {
		return -2
	}
	asWindowInt(p).SetDrawDataWindow(asWindowInt(w))
	return 0
}
