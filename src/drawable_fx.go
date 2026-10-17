package opennox

import (
	"encoding/binary"
	"image"
	"math"
	"unsafe"

	noxcolor "github.com/opennox/libs/color"
	"github.com/opennox/libs/types"

	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/ccall"
)

func init() {
	legacy.Sub_499490 = func(typ int, pts [4]uint16, a3, a4 int, a5, a6 byte) {
		noxClient.spawnFx499490(typ, pts, a3, a4, a5, a6)
	}
	legacy.Sub_499520 = func(typ int, pts [4]int16, a3 int16, a4, a5 byte) { noxClient.spawnFx499520(typ, pts, a3, a4, a5) }
	legacy.Nox_xxx_makePointFxCli_499610 = func(typ, n, a3, a4 int, pos image.Point) { noxClient.makePointFx(typ, n, a3, a4, pos) }
	legacy.Nox_xxx_drawEnergyBolt_499710 = func(pos image.Point, z int16, typ int) { noxClient.drawEnergyBolt(pos, z, typ) }
	legacy.Sub_499950 = func(typ int, p1, p2 image.Point, a4 uint16, a5 int8) *client.Drawable {
		return noxClient.spawnFx499950(typ, p1, p2, a4, a5)
	}
	legacy.Nox_xxx_makeLightningParticles_4999D0 = func(typ int, p1, p2 image.Point) int { return noxClient.makeLightningParticles(typ, p1, p2) }
	legacy.Nox_xxx_draw_499E70 = func(a1, a2, a3, a4, a5, a6, a7 int) { noxClient.spawnScreenParticles(a1, a2, a3, a4, a5, a6, a7) }
	legacy.Sub_49A150 = func(pos image.Point, typ int, a3 byte) { noxClient.spawnBurstFx(pos, typ, a3) }

	draws := []struct {
		hook *legacy.DrawFunc
		fnc  func(c *Client, vp *noxrender.Viewport, dr *client.Drawable) int
	}{
		{&legacy.Nox_thing_magic_sparkle_draw, (*Client).DrawMagicSparkle},
		{&legacy.Nox_thing_pixie_draw, (*Client).DrawPixie},
		{&legacy.Nox_thing_pixie_dust_draw, (*Client).DrawPixieDust},
		{&legacy.Nox_thing_blue_rain_spark_draw, (*Client).DrawBlueRainSpark},
		{&legacy.Nox_thing_rain_orb_draw, (*Client).DrawRainOrb},
		{&legacy.Nox_thing_red_spark_draw, sparkDraw(1313532, 1313528)},
		{&legacy.Nox_thing_blue_spark_draw, sparkDraw(1313540, 1313536)},
		{&legacy.Nox_thing_cyan_spark_draw, sparkDraw(1313548, 1313544)},
		{&legacy.Nox_thing_green_spark_draw, sparkDraw(1313584, 1313580)},
		{&legacy.Nox_thing_yellow_spark_draw, sparkDraw(1313532, 1313576)},
		{&legacy.Nox_thing_violet_spark_draw, sparkDraw(1313556, 1313552)},
		{&legacy.Nox_thing_death_ball_spark_draw, sparkDraw(1313572, 1313568)},
		{&legacy.Nox_thing_white_spark_draw, (*Client).DrawWhiteSpark},
		{&legacy.Nox_thing_particle_draw, (*Client).DrawParticle},
		{&legacy.Nox_thing_glow_orb_draw, func(c *Client, vp *noxrender.Viewport, dr *client.Drawable) int {
			return legacy.Sub_4B6B80(vp.C(), dr, 0)
		}},
		{&legacy.Nox_thing_glow_orb_move_draw, func(c *Client, vp *noxrender.Viewport, dr *client.Drawable) int {
			return legacy.Sub_4B6B80(vp.C(), dr, 1)
		}},
	}
	for _, d := range draws {
		fnc := d.fnc
		*d.hook = func(vp *noxrender.Viewport, dr *client.Drawable) int { return fnc(noxClient, vp, dr) }
	}
	client.RegisterDraw("YellowSparkDraw", legacy.Get_nox_thing_yellow_spark_draw(), 0, nil)
	client.RegisterDraw("RedSparkDraw", legacy.Get_nox_thing_red_spark_draw(), 0, nil)
	client.RegisterDraw("BlueSparkDraw", legacy.Get_nox_thing_blue_spark_draw(), 0, nil)
	client.RegisterDraw("CyanSparkDraw", legacy.Get_nox_thing_cyan_spark_draw(), 0, nil)
	client.RegisterDraw("GreenSparkDraw", legacy.Get_nox_thing_green_spark_draw(), 0, nil)
	client.RegisterDraw("VioletSparkDraw", legacy.Get_nox_thing_violet_spark_draw(), 0, nil)
	client.RegisterDraw("WhiteSparkDraw", legacy.Get_nox_thing_white_spark_draw(), 0, nil)
	client.RegisterDraw("BlueRainSparkDraw", legacy.Get_nox_thing_blue_rain_spark_draw(), 0, nil)
	client.RegisterDraw("DeathBallSparkDraw", legacy.Get_nox_thing_death_ball_spark_draw(), 0, nil)
	client.RegisterDraw("PixieDraw", legacy.Get_nox_thing_pixie_draw(), 0, nil)
	client.RegisterDraw("PixieDustDraw", legacy.Get_nox_thing_pixie_dust_draw(), 0, nil)
	client.RegisterDraw("ParticleDraw", legacy.Get_nox_thing_particle_draw(), 0, nil)
	client.RegisterDraw("GlowOrbDraw", legacy.Get_nox_thing_glow_orb_draw(), 0, nil)
	client.RegisterDraw("GlowOrbMoveDraw", legacy.Get_nox_thing_glow_orb_move_draw(), 0, nil)
	client.RegisterDraw("MagicSparkleDraw", legacy.Get_nox_thing_magic_sparkle_draw(), 0, nil)
	client.RegisterDraw("RainOrbDraw", legacy.Get_nox_thing_rain_orb_draw(), 0, nil)
}

// fxColor returns one of the effect colors initialized by the legacy code.
func fxColor(off uintptr) uint32 {
	// Some of these are separate C variables instead of being stored in the memory blob.
	switch off {
	case 1313532:
		return legacy.Get_dword_5d4594_1313532()
	case 1313536:
		return legacy.Get_dword_5d4594_1313536()
	case 1313540:
		return legacy.Get_dword_5d4594_1313540()
	case 1313564:
		return legacy.Get_dword_5d4594_1313564()
	}
	return memmap.Uint32(0x5D4594, off)
}

func sparkDraw(cl1, cl2 uintptr) func(c *Client, vp *noxrender.Viewport, dr *client.Drawable) int {
	return func(c *Client, vp *noxrender.Viewport, dr *client.Drawable) int {
		return legacy.Sub_4B6970(vp.C(), dr, fxColor(cl1), fxColor(cl2))
	}
}

func (c *Client) spawnFx499490(typ int, pts [4]uint16, a3, a4 int, a5, a6 byte) { // sub_499490
	s := c.Nox_xxx_spriteLoadAdd_45A360_drawable(typ, image.Pt(a3+int(pts[2]), a4+int(pts[3])))
	if s == nil {
		return
	}
	u := effectBytes(s)
	binary.LittleEndian.PutUint16(u[0:], pts[0])
	binary.LittleEndian.PutUint16(u[2:], pts[1])
	u[11] = a5
	u[12] = byte(c.rnd(3, 10))
	u[14] = a6
	u[13] = a6
	c.Cli().Objs.List34Add(s)
}

func (c *Client) spawnFx499520(typ int, pts [4]int16, a3 int16, a4, a5 byte) { // sub_499520
	x, y := int(pts[2]), int(pts[3])
	s := c.Nox_xxx_spriteLoadAdd_45A360_drawable(typ, image.Pt(x, y))
	if s == nil {
		return
	}
	u := effectBytes(s)
	binary.LittleEndian.PutUint16(u[0:], uint16(pts[0]))
	binary.LittleEndian.PutUint16(u[2:], uint16(pts[1]))
	binary.LittleEndian.PutUint16(u[4:], uint16(x))
	binary.LittleEndian.PutUint16(u[6:], uint16(y))
	u[10] = byte(a3)
	dx, dy := x-int(pts[0]), y-int(pts[1])
	binary.LittleEndian.PutUint16(u[8:], uint16(int64(math.Sqrt(float64(dx*dx+dy*dy)))))
	u[11] = a4
	u[12] = byte(c.rnd(3, 10))
	u[14] = a5
	u[13] = a5
	s.ClientUpdateFuncPtr = legacy.Get_sub_4CA720()
	*(*uint16)(unsafe.Pointer(&s.Field_127)) = uint16(a3)
	c.Cli().Objs.List34Add(s)
}

func (c *Client) makePointFx(typ, n, a3, a4 int, pos image.Point) { // nox_xxx_makePointFxCli_499610
	for ; n > 0; n-- {
		s := c.Nox_xxx_spriteLoadAdd_45A360_drawable(typ, pos)
		if s == nil {
			continue
		}
		dir := byte(c.rnd(0, 255))
		speed := uint32(c.rnd(1, a3))
		c.setSparkMotion(s, s.PosVec, dir, speed, c.srv.Frame()+uint32(c.rnd(a4, 64)))
		s.ZVal = 0
		s.VelZ = int8(c.rnd(2, 10))
		c.Cli().Objs.List34Add(s)
	}
}

func (c *Client) drawEnergyBolt(pos image.Point, z int16, typ int) { // nox_xxx_drawEnergyBolt_499710
	for i := 0; i < 2; i++ {
		s := c.Nox_xxx_spriteLoadAdd_45A360_drawable(typ, pos)
		if s == nil {
			continue
		}
		dir := byte(c.rnd(0, 255))
		speed := uint32(c.rnd(1, 3000))
		c.setSparkMotion(s, pos, dir, speed, c.srv.Frame()+uint32(c.rnd(5, 20)))
		s.ZVal = uint16(z + int16(c.rnd(0, 20)))
		s.VelZ = int8(c.rnd(0, 4))
		c.Cli().Objs.List34Add(s)
	}
}

func (c *Client) spawnFx499950(typ int, p1, p2 image.Point, a4 uint16, a5 int8) *client.Drawable { // sub_499950
	s := c.Nox_xxx_spriteLoadAdd_45A360_drawable(typ, p1)
	if s == nil {
		return nil
	}
	s.ZVal = a4
	s.ZVal2 = 0
	s.VelZ = a5
	u := effectBytes(s)
	binary.LittleEndian.PutUint16(u[8:], a4)
	u[10] = byte(c.rnd(3, 10))
	eu := s.UnionEffect()
	eu.Field_108 = uint32(p2.X)
	eu.Field_109 = uint32(p2.Y)
	c.Cli().Objs.List34Add(s)
	return s
}

func (c *Client) makeLightningParticles(typ int, p1, p2 image.Point) int { // nox_xxx_makeLightningParticles_4999D0
	d := p2.Sub(p1)
	l := int(int64(math.Sqrt(float64(d.X*d.X + d.Y*d.Y))))
	ret := l
	if l <= 0 {
		return ret
	}
	t := c.rnd(0, l)
	ret = t
	for pos := t; pos <= l; {
		x := p1.X + d.X*pos/l
		y := p1.Y + d.Y*pos/l
		if s := c.Nox_xxx_spriteLoadAdd_45A360_drawable(typ, image.Pt(x, y)); s != nil {
			dir := byte(c.rnd(0, 255))
			speed := uint32(c.rnd(1, 3000))
			c.setSparkMotion(s, image.Pt(x, y), dir, speed, c.srv.Frame()+uint32(c.rnd(5, 20)))
			s.ZVal = uint16(c.rnd(15, 30))
			s.VelZ = int8(c.rnd(-4, 4))
			c.Cli().Objs.List34Add(s)
		}
		ret = c.rnd(8, 100)
		pos += ret
	}
	return ret
}

func (c *Client) spawnScreenParticles(a1, x0, y0, w, h, mode, side int) { // nox_xxx_draw_499E70
	for i := 0; i < 100; i++ {
		x, y := x0, y0
		if mode == 2 {
			x = x0 + c.rnd(0, w)
		} else {
			y = y0 + c.rnd(0, h)
		}
		size := c.rnd(2, 5)
		vy := c.rnd(-40, -20)
		var vx int
		if side == 1 {
			vx = c.rnd(-20, 0)
		} else {
			vx = c.rnd(0, 20)
		}
		legacy.Nox_client_newScreenParticle_431540(a1, x, y, vx, vy, 1, size, 0, 0, 1)
	}
}

func (c *Client) spawnBurstFx(pos image.Point, typ int, a3 byte) { // sub_49A150
	maxSpeed := 2400*int(a3)/255 + 200
	minLife := 10*int(a3)/255 + 5
	for n := 180*int(a3)/255 + 10; n > 0; n-- {
		s := c.Nox_xxx_spriteLoadAdd_45A360_drawable(typ, pos)
		if s == nil {
			continue
		}
		dir := byte(c.rnd(0, 255))
		speed := uint32(c.rnd(1, maxSpeed))
		c.setSparkMotion(s, s.PosVec, dir, speed, c.srv.Frame()+uint32(c.rnd(minLife, 96)))
		s.ZVal = uint16(c.rnd(5, 15))
		s.ZVal2 = 0
		s.VelZ = int8(c.rnd(0, 8))
		c.Cli().Objs.List34Add(s)
	}
}

// drawFadingSpark draws a spark that fades out by the deadline.
func (c *Client) drawFadingSpark(vp *noxrender.Viewport, dr *client.Drawable, cl, glow uint32) int { // sub_4B6770
	r := c.r
	u := dr.UnionEffect()
	total := int(int32(u.Field_112 - u.Field_111))
	left := int(int32(u.Field_112 - c.srv.Frame()))
	if left == total {
		left--
	}
	if left <= 0 {
		c.Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr)
		return 0
	}
	p := vp.ToScreenPos(dr.Pos())
	p.Y -= int(int16(dr.ZVal2)) + int(int16(dr.ZVal))
	if p.X-10 >= vp.Screen.Min.X && p.Y-10 >= vp.Screen.Min.Y && p.X+10 < vp.Screen.Max.X && p.Y+10 < vp.Screen.Max.Y {
		if sz := left * c.rnd(0, 4) / total; sz != 0 {
			r.DrawGlow(p, noxcolor.RGBA5551(glow), 2*sz+1, sz+1)
			r.Data().SetColor2(noxcolor.RGBA5551(cl))
			r.DrawPoint(p, sz, r.Data().Color2())
		}
	}
	return 1
}

func (c *Client) DrawMagicSparkle(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_magic_sparkle_draw
	if c.rnd(0, 10) >= 5 {
		return c.drawFadingSpark(vp, dr, nox_color_white_2523948.Color32(), fxColor(1313540))
	}
	return c.drawFadingSpark(vp, dr, fxColor(1313540), fxColor(1313536))
}

func (c *Client) DrawPixieDust(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_pixie_dust_draw
	if c.rnd(0, 10) >= 5 {
		return c.drawFadingSpark(vp, dr, nox_color_white_2523948.Color32(), fxColor(1313564))
	}
	return c.drawFadingSpark(vp, dr, fxColor(1313564), fxColor(1313560))
}

func (c *Client) DrawPixie(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_pixie_draw
	r := c.r
	z := int16(dr.ZVal)
	if c.rnd(0, 100) < 50 {
		if z > 0 {
			dr.ZVal--
		}
	} else if z < 35 {
		dr.ZVal++
	}
	p := vp.ToScreenPos(dr.Pos())
	p.Y -= int(int16(dr.ZVal2)) + int(int16(dr.ZVal))
	if p.X-10 < vp.Screen.Min.X || p.Y-10 < vp.Screen.Min.Y || p.X+10 >= vp.Screen.Max.X || p.Y+10 >= vp.Screen.Max.Y {
		return 1
	}
	r.DrawGlow(p, noxcolor.RGBA5551(fxColor(1313560)), 10, 4)
	tail := vp.ToScreenPos(image.Pt(int(int32(dr.Field_8)), int(int32(dr.Field_9))))
	tail.Y -= int(int16(dr.ZVal))
	d := p.Sub(tail)
	if dist := d.X*d.X + d.Y*d.Y; dist > 400 {
		l := int(int64(math.Sqrt(float64(dist))))
		tail = image.Pt(p.X-20*d.X/l, p.Y-20*d.Y/l)
	}
	r.Data().SetColor2(noxcolor.RGBA5551(fxColor(1313564)))
	r.AddPoint(p)
	r.AddPoint(tail)
	r.DrawLineFromPoints(r.Data().Color2())
	return 1
}

var (
	typeWhiteSpark    = spriteTypeCache{id: "WhiteSpark"}
	typeRainOrbWhite2 = spriteTypeCache{id: "RainOrbWhite"}
	typeWhiteMoveOrb  = spriteTypeCache{id: "WhiteMoveOrb"}
	typeBlueMoveOrb   = spriteTypeCache{id: "BlueMoveOrb"}
)

func (c *Client) DrawWhiteSpark(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_white_spark_draw
	return legacy.Sub_4B6970(vp.C(), dr, nox_color_white_2523948.Color32(), fxColor(1313540))
}

func (c *Client) DrawBlueRainSpark(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_blue_rain_spark_draw
	res := legacy.Sub_4B6970(vp.C(), dr, nox_color_white_2523948.Color32(), fxColor(1313540))
	if res != 1 || byte(dr.VelZ) < 5 {
		return res
	}
	if s := c.Nox_xxx_spriteLoadAdd_45A360_drawable(c.spriteType(&typeWhiteSpark), dr.PosVec); s != nil {
		dir := byte(c.rnd(0, 255))
		speed := uint32(c.rnd(1, 1611))
		c.setSparkMotion(s, dr.PosVec, dir, speed, c.srv.Frame()+uint32(c.rnd(10, 96)))
		s.ZVal = uint16(c.rnd(5, 15))
		s.ZVal2 = 0
		s.VelZ = int8(c.rnd(0, 8))
		c.Cli().Objs.List34Add(s)
	}
	c.Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr)
	return 0
}

func (c *Client) DrawRainOrb(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_rain_orb_draw
	r := c.r
	white := c.spriteType(&typeRainOrbWhite2)
	u := effectBytes(dr)
	if z := int16(dr.ZVal); z > 0 {
		p := vp.ToScreenPos(dr.Pos())
		p.Y -= int(z)
		cl := fxColor(1313588)
		if int(dr.TypeIDVal) != white {
			cl = fxColor(1313536)
		}
		r.DrawGlow(p, noxcolor.RGBA5551(cl), int(u[10]), 5)
		dz := int(int16(dr.ZVal)) - int(int16(binary.LittleEndian.Uint16(u[8:])))
		r.Data().SetColor2(noxcolor.RGBA5551(cl))
		r.AddPoint(p)
		r.AddPointRel(image.Pt(0, dz))
		r.DrawLineFromPoints(r.Data().Color2())
		binary.LittleEndian.PutUint16(u[8:], dr.ZVal)
		dr.ZVal += uint16(dr.VelZ)
		r.Data().SetColor2(noxcolor.RGBA5551(fxColor(1313592)))
		r.DrawPoint(p, int(u[10])/3, r.Data().Color2())
		return 1
	}
	eu := dr.UnionEffect()
	x, y := dr.PosVec.X, dr.PosVec.Y
	dy := int(int32(uint32(y) - eu.Field_109))
	ang := legacy.Nox_xxx_math_509ED0(types.Ptf(float32(x-int(int32(eu.Field_108))), float32(dy)))
	var pts [4]uint16
	pts[2] = uint16(x)
	pts[3] = uint16(y + 20)
	pts[0] = uint16(int(float32(float64(memmap.Float32(0x587000, 194136+8*uintptr(ang)))*150.0 + float64(int32(eu.Field_108)))))
	pts[1] = uint16(int(float32(float64(memmap.Float32(0x587000, 194140+8*uintptr(ang)))*150.0 + float64(int32(eu.Field_109)))))
	typ := c.spriteType(&typeWhiteMoveOrb)
	if int(dr.TypeIDVal) != white {
		typ = c.spriteType(&typeBlueMoveOrb)
	}
	c.spawnFx499490(typ, pts, 0, 0, byte(c.rnd(6, 8)), 1)
	c.Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr)
	return 0
}

func (c *Client) DrawParticle(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_particle_draw
	// The original code expects a legacy particle structure here; Go particles never set it.
	p := unsafe.Pointer(uintptr(dr.UnionEffect().Field_108))
	dr.PosVec.X = int(*(*uint32)(unsafe.Add(p, 80)) >> 16)
	dr.PosVec.Y = int(*(*uint32)(unsafe.Add(p, 84)) >> 16)
	if fnc := *(*unsafe.Pointer)(unsafe.Add(p, 124)); fnc != nil {
		ccall.CallVoidPtr(fnc, p)
	}
	return 1
}
