package opennox

import (
	"image"
	"unsafe"

	noxcolor "github.com/opennox/libs/color"
	"github.com/opennox/libs/object"

	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/binfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func init() {
	client.RegisterDraw("StaticDraw", legacy.Get_nox_thing_static_draw(), 1, parseStaticDraw)
	client.RegisterDraw("StaticRandomDraw", legacy.Get_nox_thing_static_random_draw(), 2, parseStaticRandomDraw)
	client.RegisterDraw("DoorDraw", legacy.Get_nox_thing_door_draw(), 2, parseDoorDraw)
	client.RegisterDraw("SlaveDraw", legacy.Get_nox_thing_slave_draw(), 2, parseSlaveDraw)
	client.RegisterDraw("TriggerDraw", legacy.Get_nox_thing_trigger_draw(), 0, nil)
	client.RegisterDraw("PressurePlateDraw", legacy.Get_nox_thing_pressure_plate_draw(), 0, nil)
	client.RegisterDraw("BlackPowderDraw", legacy.Get_nox_thing_black_powder_draw(), 0, nil)
	client.RegisterDraw("GlyphDraw", legacy.Get_nox_thing_glyph_draw(), 3, parseAnimateDraw)
	client.RegisterDraw("BoulderDraw", legacy.Get_nox_thing_boulder_draw(), 2, parseSlaveDraw)
	client.RegisterDraw("WeaponDraw", legacy.Get_nox_thing_weapon_draw(), 1, parseStaticDraw)
	client.RegisterDraw("ArmorDraw", legacy.Get_nox_thing_armor_draw(), 1, parseStaticDraw)
	client.RegisterDraw("WeaponAnimateDraw", legacy.Get_nox_thing_weapon_animate_draw(), 3, parseAnimateDraw)
	client.RegisterDraw("ArmorAnimateDraw", legacy.Get_nox_thing_armor_animate_draw(), 3, parseAnimateDraw)
	client.RegisterDraw("FlagDraw", legacy.Get_nox_thing_flag_draw(), 3, parseAnimateDraw)
	client.RegisterDraw("BaseDraw", legacy.Get_nox_thing_base_draw(), 1, parseStaticDraw)
	client.RegisterDraw("SphericalShieldDraw", legacy.Get_nox_thing_spherical_shield_draw(), 3, parseAnimateDraw)
	client.RegisterDraw("ArrowDraw", legacy.Get_nox_thing_arrow_draw(), 2, parseSlaveDraw)
	client.RegisterDraw("WeakArrowDraw", legacy.Get_nox_thing_weak_arrow_draw(), 2, parseSlaveDraw)
	client.RegisterDraw("HarpoonDraw", legacy.Get_nox_thing_harpoon_draw(), 2, parseSlaveDraw)
}

// staticDrawData is a draw data for StaticDraw.
type staticDrawData struct {
	Size uint32
	Img  noxrender.ImageHandle
}

// staticRandomDrawData is a draw data for StaticRandomDraw, SlaveDraw and similar.
type staticRandomDrawData struct {
	Size   uint32
	Images unsafe.Pointer // *[Cnt]noxrender.ImageHandle
	Cnt    uint8
	_      [3]byte
}

func (d *staticRandomDrawData) Frame(i int) noxrender.ImageHandle {
	return *(*noxrender.ImageHandle)(unsafe.Add(d.Images, uintptr(i)*unsafe.Sizeof(noxrender.ImageHandle(nil))))
}

func readThingImage(f *binfile.MemFile) noxrender.ImageHandle {
	id := f.ReadI32()
	var (
		typ   byte
		name2 string
	)
	if id == -1 {
		typ = f.ReadU8()
		name2, _ = f.ReadString8()
	}
	return noxClient.r.Bag.ImageRef(int(id), typ, name2).C()
}

func parseStaticDraw(obj *client.ObjectType, f *binfile.MemFile, _ string, _ []byte) error { // nox_things_static_draw_parse
	dd, _ := alloc.New(staticDrawData{})
	dd.Size = uint32(unsafe.Sizeof(staticDrawData{}))
	dd.Img = readThingImage(f)
	obj.DrawFunc = legacy.Get_nox_thing_static_draw()
	obj.DrawData = unsafe.Pointer(dd)
	return nil
}

func loadStaticRandomData(f *binfile.MemFile) *staticRandomDrawData { // nox_xxx_spriteLoadStaticRandomData_44C000
	dd, _ := alloc.New(staticRandomDrawData{})
	dd.Size = uint32(unsafe.Sizeof(staticRandomDrawData{}))
	dd.Cnt = f.ReadU8()
	if dd.Cnt == 0 {
		return dd
	}
	arr, _ := alloc.Make([]noxrender.ImageHandle{}, int(dd.Cnt))
	dd.Images = unsafe.Pointer(&arr[0])
	for i := range arr {
		arr[i] = readThingImage(f)
	}
	return dd
}

func parseStaticRandomDraw(obj *client.ObjectType, f *binfile.MemFile, _ string, _ []byte) error { // nox_things_static_random_draw_parse
	obj.DrawFunc = legacy.Get_nox_thing_static_random_draw()
	obj.DrawData = unsafe.Pointer(loadStaticRandomData(f))
	return nil
}

func parseSlaveDraw(obj *client.ObjectType, f *binfile.MemFile, _ string, _ []byte) error { // nox_things_slave_draw_parse
	obj.DrawFunc = legacy.Get_nox_thing_slave_draw()
	obj.DrawData = unsafe.Pointer(loadStaticRandomData(f))
	obj.Field_60 = 0
	return nil
}

func parseDoorDraw(obj *client.ObjectType, f *binfile.MemFile, _ string, _ []byte) error { // nox_things_door_draw_parse
	obj.DrawFunc = legacy.Get_nox_thing_door_draw()
	obj.DrawData = unsafe.Pointer(loadStaticRandomData(f))
	return nil
}

func (c *Client) DrawStatic(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_static_draw
	if !dr.ObjClass.Has(object.ClassReadable) || dr.ObjFlags.Has(object.FlagEnabled) {
		dd := (*staticDrawData)(dr.DrawData)
		legacy.Nox_xxx_drawObject_4C4770_draw(vp, dr, dd.Img)
	}
	return 1
}

func (c *Client) DrawStaticRandom(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_static_random_draw
	dd := (*staticRandomDrawData)(dr.DrawData)
	legacy.Nox_xxx_drawObject_4C4770_draw(vp, dr, dd.Frame(int(dr.AnimFrameSlave)))
	return 1
}

func (c *Client) DrawSlave(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_slave_draw
	return c.DrawStaticRandom(vp, dr)
}

func (c *Client) DrawTrigger(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_trigger_draw
	r := c.r
	p := vp.ToScreenPos(dr.Pos())
	box := &dr.Shape.Box
	p1 := image.Pt(p.X+int(box.LeftTop), p.Y+int(box.LeftBottom))
	p2 := image.Pt(p.X+int(box.RightBottom2), p.Y+int(box.RightTop2))
	p3 := image.Pt(p.X+int(box.LeftBottom2), p.Y+int(box.LeftTop2))
	p4 := image.Pt(p.X+int(box.RightTop), p.Y+int(box.RightBottom))
	r.Data().SetColor2(nox_color_black_2650656)
	r.Data().SetAlphaEnabled(true)
	for _, l := range [][2]image.Point{{p1, p4}, {p2, p4}, {p1, p3}, {p2, p3}} {
		r.AddPoint(l[0])
		r.AddPoint(l[1])
		r.DrawLineFromPoints(r.Data().Color2())
	}
	r.Data().SetAlphaEnabled(false)
	return 1
}

func (c *Client) DrawPressurePlate(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_pressure_plate_draw
	r := c.r
	cl := (*[6]byte)(unsafe.Pointer(&dr.Union))
	if *cl == [6]byte{} {
		return 1
	}
	p := vp.ToScreenPos(dr.Pos())
	box := &dr.Shape.Box
	p1 := image.Pt(p.X+int(box.LeftTop), p.Y+int(box.LeftBottom))
	p2 := image.Pt(p.X+int(box.RightBottom2), p.Y+int(box.RightTop2))
	p3 := image.Pt(p.X+int(box.LeftBottom2), p.Y+int(box.LeftTop2))
	p4 := image.Pt(p.X+int(box.RightTop), p.Y+int(box.RightBottom))
	for i := 0; i < 2; i++ {
		r.Data().SetColor2(noxcolor.RGB5551Color(cl[3*i+0], cl[3*i+1], cl[3*i+2]))
		d := image.Pt(i, 0)
		for _, l := range [][2]image.Point{{p1, p4}, {p2, p4}, {p1, p3}, {p2, p3}} {
			r.AddPoint(l[0].Add(d))
			r.AddPoint(l[1].Add(d))
			r.DrawLineFromPoints(r.Data().Color2())
		}
	}
	return 1
}

func (c *Client) DrawBlackPowder(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_black_powder_draw
	r := c.r
	p := vp.ToScreenPos(dr.Pos())
	r.Data().SetColor2(noxcolor.RGBA5551(memmap.Uint32(0x85B3FC, 956)))
	r.DrawRectFilledOpaque(p.X-1, p.Y-1, 3, 3, r.Data().Color2())
	r.DrawRectFilledOpaque(p.X-5, p.Y, 1, 1, r.Data().Color2())
	r.DrawRectFilledOpaque(p.X, p.Y+7, 1, 1, r.Data().Color2())
	r.DrawRectFilledOpaque(p.X+8, p.Y-2, 1, 1, r.Data().Color2())
	return 1
}

func (c *Client) DrawWeaponItem(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_weapon_draw
	legacy.Sub_4B95D0(dr)
	return c.DrawStatic(vp, dr)
}

func (c *Client) DrawWeaponAnimate(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_weapon_animate_draw
	legacy.Sub_4B95D0(dr)
	return c.DrawAnimate(vp, dr)
}

func (c *Client) DrawArmorItem(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_armor_draw
	legacy.Sub_4B96F0(dr)
	return c.DrawStatic(vp, dr)
}

func (c *Client) DrawArmorAnimate(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_armor_animate_draw
	legacy.Sub_4B96F0(dr)
	return c.DrawAnimate(vp, dr)
}

func (c *Client) DrawBase(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_base_draw
	c.DrawWeaponItem(vp, dr)
	return 1
}

func (c *Client) DrawSphericalShield(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_spherical_shield_draw
	code := uint16(dr.UnionEffect().Field_108)
	var owner *client.Drawable
	if nox_xxx_netTestHighBit_578B70(code) {
		owner = c.Cli().Objs.ByNetCodeStatic(int(nox_xxx_netClearHighBit_578B30(code)))
	} else {
		owner = c.Cli().Objs.ByNetCodeDynamic(int(nox_xxx_netClearHighBit_578B30(code)))
	}
	if owner == nil {
		c.Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr)
		return 0
	}
	c.Cli().Nox_xxx_updateSpritePosition_49AA90(dr, owner.PosVec.X, owner.PosVec.Y+3)
	return c.DrawAnimate(vp, dr)
}

func (c *Client) DrawGlyph(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_glyph_draw
	r := c.r
	// The original code uses the low byte of the drawable pointer as alpha here.
	alpha := byte(uintptr(dr.C()))
	if pu := c.ClientPlayerUnit(); noxflags.HasGame(noxflags.GameClient) && pu != nil {
		if dr.ObjFlags.Has(object.FlagSelected) {
			alpha = 0xff
		} else if pu.HasEnchant(server.EnchantID(21)) {
			r.Data().SetColorize17(1)
			cl := noxcolor.RGBA5551(legacy.Get_dword_8531A0_2572()).ColorNRGBA()
			r.SetColorMultAndIntensityRGB(cl.R, cl.G, cl.B)
			alpha = 0xff
		} else {
			d := dr.PosVec.Sub(pu.PosVec)
			dist := d.X*d.X + d.Y*d.Y
			if dist >= 22500 {
				return 1
			}
			alpha = byte(-56 - 200*dist/22500)
		}
	}
	r.Data().SetAlphaEnabled(true)
	r.Data().SetAlpha(alpha)
	res := c.DrawAnimate(vp, dr)
	r.Data().SetAlphaEnabled(false)
	r.Data().SetColorize17(0)
	return res
}

func (c *Client) DrawFlag(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_flag_draw
	r := c.r
	c.DrawWeaponAnimate(vp, dr)
	if !noxflags.HasGame(noxflags.GameModeChat) || !dr.ObjFlags.Has(object.FlagEnabled) {
		return 1
	}
	p := vp.ToScreenPos(dr.Pos())
	p.Y -= int(int16(dr.ZVal)) + int(int64(dr.ZSizeMax))
	str, ok := legacy.Sub_418A80(legacy.Sub_4B94E0(dr))
	if !ok {
		return 1
	}
	fnt := r.GetFonts().AsFont(nil)
	r.Data().SetTextColor(noxcolor.RGBA5551(legacy.Get_dword_8531A0_2572()))
	w := r.GetStringSizeWrapped(fnt, str, 0).X
	r.DrawString(fnt, str, image.Pt(w/-2+p.X, p.Y))
	return 1
}

func (c *Client) DrawBoulder(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_boulder_draw
	dd := (*staticRandomDrawData)(dr.DrawData)
	u := dr.UnionEffect()
	if u.Field_108 == 0 && u.Field_109 == 0 {
		u.Field_108 = uint32(dr.PosVec.X)
		u.Field_109 = uint32(dr.PosVec.Y)
	}
	dx := int(int32(uint32(dr.PosVec.X) - u.Field_108))
	dy := int(int32(uint32(dr.PosVec.Y) - u.Field_109))
	if dx*dx+dy*dy >= 100 {
		// Roll the boulder in the direction of movement.
		forward := false
		if dx <= 0 {
			if dy > 0 {
				u.Field_111 = 0
				forward = true
			} else {
				u.Field_111 = 16
			}
		} else {
			if dy > 0 {
				u.Field_111 = 16
				forward = true
			} else {
				u.Field_111 = 0
			}
		}
		if forward {
			u.Field_110++
			if u.Field_110 >= 16 {
				u.Field_110 = 0
			}
		} else if u.Field_110 != 0 {
			u.Field_110--
		} else {
			u.Field_110 = 15
		}
		u.Field_108 = uint32(dr.PosVec.X)
		u.Field_109 = uint32(dr.PosVec.Y)
	}
	legacy.Nox_xxx_drawObject_4C4770_draw(vp, dr, dd.Frame(int(u.Field_110+u.Field_111)))
	return 1
}

func (c *Client) DrawDoor(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_door_draw
	r := c.r
	dd := (*staticRandomDrawData)(dr.DrawData)
	legacy.Nox_xxx_drawObject_4C4770_draw(vp, dr, dd.Frame(int(dr.Field_74_4)))
	u := dr.UnionDoor()
	if !noxflags.HasGame(noxflags.GameModeQuest) || u.Field_108_1 != 1 {
		return 1
	}
	if u.Field_109 == nil {
		u.Field_110 = nox_xxx_gLoadImg("DoorLockSilverSW").C()
		u.Field_109 = nox_xxx_gLoadImg("DoorLockSilverSE").C()
		u.Field_111 = nox_xxx_gLoadImg("DoorLockGoldSW").C()
		u.Field_112 = nox_xxx_gLoadImg("DoorLockGoldSE").C()
	}
	p := vp.ToScreenPos(dr.Pos())
	p.X -= 64
	p.Y -= 79
	light := c.Sub469920(dr.Pos())
	r.Data().SetMultiply14(1)
	r.SetColorMultAndIntensityRGB(byte(light[0]), byte(light[1]), byte(light[2]))
	silver := byte(u.Field_108_2) == 1
	var img noxrender.ImageHandle
	switch dr.Field_74_4 {
	case 0:
		p = p.Add(image.Pt(-15, -20))
		img = u.Field_111
		if silver {
			img = u.Field_110
		}
	case 8:
		p = p.Add(image.Pt(15, -20))
		img = u.Field_112
		if silver {
			img = u.Field_109
		}
	case 16:
		p = p.Add(image.Pt(8, 2))
		img = u.Field_111
		if silver {
			img = u.Field_110
		}
	default:
		p = p.Add(image.Pt(-8, 2))
		img = u.Field_112
		if silver {
			img = u.Field_109
		}
	}
	r.DrawImageAt(r.Bag.AsImage(img), p)
	return 1
}
