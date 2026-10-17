package gui

import "image"

// SetEnabled sets or clears the enabled flag for the window and all its children.
func (win *Window) SetEnabled(enabled bool) { // nox_xxx_wnd_46ABB0
	if enabled {
		win.Flags |= StatusEnabled
	} else {
		win.Flags &^= StatusEnabled
	}
	for c := win.Field100Ptr; c != nil; c = c.prev {
		c.SetEnabled(enabled)
	}
}

// SetSize changes the window size and notifies the window about it.
func (win *Window) SetSize(sz image.Point) { // sub_46AB20
	win.EndPos.X = sz.X + win.Off.X
	win.SizeVal = sz
	win.EndPos.Y = sz.Y + win.Off.Y
	win.Func94(AsWindowEvent(0x4004, uintptr(sz.X), uintptr(sz.Y)))
}

// PointInGlobal checks if the point in screen coordinates is inside the window.
func (win *Window) PointInGlobal(p image.Point) bool { // nox_xxx_wndPointInWnd_46AAB0
	pos := win.GlobalPos()
	sz := win.Size()
	return p.X >= pos.X && p.X <= pos.X+sz.X && p.Y >= pos.Y && p.Y <= pos.Y+sz.Y
}

// SetDrawDataWindow sets the window that is referenced by the draw data.
// Passing nil makes the window reference itself.
func (win *Window) SetDrawDataWindow(w *Window) { // nox_xxx_wnd_46B280
	if w == nil {
		w = win
	}
	win.drawData.Window = w
}
