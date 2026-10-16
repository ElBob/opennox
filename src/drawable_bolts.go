package opennox

import (
	"encoding/binary"
	"image"
	"math"

	noxcolor "github.com/opennox/libs/color"
	"github.com/opennox/libs/types"

	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
)

func init() {
	client.RegisterDraw("LightningDraw", legacy.Get_nox_thing_lightning_draw(), 0, nil)
	client.RegisterDraw("ChainLightningBoltDraw", legacy.Get_nox_thing_chain_lightning_bolt_draw(), 0, nil)
	client.RegisterDraw("EnergyBoltDraw", legacy.Get_nox_thing_energy_bolt_draw(), 0, nil)
	client.RegisterDraw("GreenBoltDraw", legacy.Get_nox_thing_green_bolt_draw(), 0, nil)
	client.RegisterDraw("PlasmaDraw", legacy.Get_nox_thing_plasma_draw(), 0, nil)
	legacy.Nox_thing_lightning_draw = func(vp *noxrender.Viewport, dr *client.Drawable) int { return noxClient.DrawLightning(vp, dr) }
	legacy.Nox_thing_chain_lightning_bolt_draw = func(vp *noxrender.Viewport, dr *client.Drawable) int { return noxClient.DrawLightning(vp, dr) }
	legacy.Nox_thing_energy_bolt_draw = func(vp *noxrender.Viewport, dr *client.Drawable) int { return noxClient.DrawEnergyBolt(vp, dr) }
	legacy.Nox_thing_green_bolt_draw = func(vp *noxrender.Viewport, dr *client.Drawable) int { return noxClient.DrawGreenBolt(vp, dr) }
	legacy.Nox_thing_plasma_draw = func(vp *noxrender.Viewport, dr *client.Drawable) int { return noxClient.DrawPlasma(vp, dr) }
	legacy.Sub_4BA8B0 = func(p1, p2 image.Point, ribbon int) byte { return noxClient.plasmaAddSegment(p1, p2, ribbon) }
}

// lightningState holds the state of the recursive lightning bolt renderer.
var lightningState struct {
	steps    int    // nox_xxx_lightningSteps_587000_178216
	maxDepth int    // dword_5d4594_1316448
	depth    int    // dword_5d4594_1316492
	jitter   int    // dword_5d4594_1316476
	color    uint32 // dword_5d4594_1316472
	color1   uint32 // dword_5d4594_1316452
	color2   uint32 // dword_5d4594_1316456
	color3   uint32 // dword_5d4594_1316436
	color3b  uint32 // dword_5d4594_1316484
}

func init() {
	lightningState.steps = 8
}

// packPoint packs a point as two 16 bit values like the original code.
func packPoint(p image.Point) uint32 {
	return uint32(uint16(p.X)) | uint32(uint16(p.Y))<<16
}

func (c *Client) drawLightningStep(a, b uint32) { // nox_xxx_drawLightningStep_4BB070
	r := c.r
	st := &lightningState
	st.depth++
	defer func() { st.depth-- }()
	ax, ay := int(int16(a)), int(int16(a>>16))
	bx, by := int(int16(b)), int(int16(b>>16))
	if st.depth != st.maxDepth {
		div, mul := 1, 1
		for i := 1; i < st.depth; i++ {
			div *= int(memmap.Int32(0x587000, 178212)) + 9
			mul *= 10
		}
		mx := uint16(mul*c.rnd(-st.jitter, st.jitter)/div + (ax+bx)>>1)
		my := uint16(mul*c.rnd(-st.jitter, st.jitter)/div + (ay+by)>>1)
		mid := uint32(mx) | uint32(my)<<16
		c.drawLightningStep(a, mid)
		c.drawLightningStep(mid, b)
		return
	}
	glow := memmap.Uint32(0x5D4594, 1316508) != 0
	if glow {
		setColorInt54(r, memmap.Uint32(0x5D4594, 1316440))
		r.Data().SetField262(12)
		r.AddPoint(image.Pt(ax, ay))
		r.AddPoint(image.Pt(bx, by))
		r.DrawParticles49ED80(int(memmap.Uint8(0x5D4594, 1316420)) + 48)
	} else {
		setColorInt54(r, st.color)
		r.Data().SetField262(3)
		r.AddPoint(image.Pt(ax, ay))
		r.AddPoint(image.Pt(bx, by))
		r.DrawParticles49ED80(32)
	}
	r.Data().SetColor2(noxcolor.RGBA5551(st.color))
	r.AddPoint(image.Pt(ax, ay))
	r.AddPoint(image.Pt(bx, by))
	r.DrawLineFromPoints(r.Data().Color2())
	if !glow {
		return
	}
	width := memmap.Uint8(0x5D4594, 1316420)
	if width&0xFE <= 2 {
		return
	}
	dx, dy := absInt(ax-bx), absInt(ay-by)
	for _, sign := range []int{+1, -1} {
		for k := 1; ; {
			r.ClearPoints()
			d := sign * k
			if dx <= dy {
				r.AddPoint(image.Pt(bx+d, by))
				r.AddPoint(image.Pt(ax+d, ay))
			} else {
				r.AddPoint(image.Pt(bx, by+d))
				r.AddPoint(image.Pt(ax, ay+d))
			}
			r.DrawLineFromPoints(r.Data().Color2())
			k++
			if k >= int(memmap.Uint8(0x5D4594, 1316420)>>1) {
				break
			}
		}
	}
}

func setColorInt54(r *NoxRender, cl uint32) { // nox_draw_set54RGB32_434040
	c := noxrender.SplitColor(noxcolor.RGBA5551(cl))
	r.Data().SetColorInt54(noxrender.RGB{R: int(c.R), G: int(c.G), B: int(c.B)})
}

func (c *Client) lightningBolt(p1, p2 image.Point, mode int, outer, inner bool) { // nox_xxx_lightningProc2_4BAE60
	st := &lightningState
	d := p1.Sub(p2)
	dx, dy := absInt(d.X), absInt(d.Y)
	l := int(int64(math.Sqrt(float64(dx*dx + dy*dy))))
	maxJitter := int(memmap.Int32(0x587000, 178204))
	minJitter := int(memmap.Int32(0x587000, 178208))
	if l >= 512 {
		st.jitter = maxJitter
		st.maxDepth = st.steps
	} else {
		st.jitter = minJitter + l*(maxJitter-minJitter)/512
		steps := st.steps
		switch {
		case l < 64:
			steps -= 3
		case l < 128:
			steps -= 2
		case l < 256:
			steps--
		}
		if steps < 1 && l < 256 {
			steps = 1
		}
		st.maxDepth = steps
	}
	*memmap.PtrUint32(0x5D4594, 1316532) = uint32(mode)
	a, b := packPoint(p1), packPoint(p2)
	if outer {
		st.depth = 1
		st.color = st.color2
		*memmap.PtrUint32(0x5D4594, 1316508) = 0
		c.drawLightningStep(a, b)
		st.depth = 1
		st.color = st.color1
		c.drawLightningStep(a, b)
	}
	if inner {
		st.depth = 1
		st.color = st.color3
		*memmap.PtrUint32(0x5D4594, 1316440) = st.color3b
		*memmap.PtrUint32(0x5D4594, 1316508) = 1
		c.drawLightningStep(a, b)
	}
}

// boltEnds returns the screen and world end points of a bolt drawable.
func (c *Client) boltEnds(vp *noxrender.Viewport, dr *client.Drawable, clearHigh bool) (s1, s2, w1, w2 image.Point, ok bool) {
	u := effectBytes(dr)
	off := vp.Screen.Min.Sub(vp.World.Min)
	if u[0] == 0 {
		w1 = image.Pt(int(binary.LittleEndian.Uint16(u[5:])), int(binary.LittleEndian.Uint16(u[7:])))
		w2 = image.Pt(int(binary.LittleEndian.Uint16(u[9:])), int(binary.LittleEndian.Uint16(u[11:])))
	} else {
		byCode := func(code uint32) *client.Drawable {
			if clearHigh {
				return c.byNetCodeHigh(code)
			}
			if nox_xxx_netTestHighBit_578B70(uint16(code)) {
				return c.Cli().Objs.ByNetCodeStatic(int(code))
			}
			return c.Cli().Objs.ByNetCodeDynamic(int(code))
		}
		from := byCode(binary.LittleEndian.Uint32(u[5:]))
		to := byCode(binary.LittleEndian.Uint32(u[9:]))
		if from == nil || to == nil {
			return
		}
		w1, w2 = from.PosVec, to.PosVec
	}
	s1 = w1.Add(off).Sub(image.Pt(0, 20))
	s2 = w2.Add(off).Sub(image.Pt(0, 20))
	return s1, s2, w1, w2, true
}

func (c *Client) DrawLightning(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_lightning_draw, nox_thing_chain_lightning_bolt_draw
	s1, s2, w1, w2, ok := c.boltEnds(vp, dr, false)
	if !ok {
		return 1
	}
	st := &lightningState
	st.color1 = memmap.Uint32(0x5D4594, 1316428)
	st.color3 = memmap.Uint32(0x5D4594, 1316464)
	st.color2 = memmap.Uint32(0x5D4594, 1316424)
	st.color3b = memmap.Uint32(0x5D4594, 1316488)
	*memmap.PtrUint8(0x5D4594, 1316420) = 1
	c.lightningBolt(s1, s2, 2, true, true)
	if !nox_xxx_checkGameFlagPause_413A50() {
		c.makeLightningParticles(int(memmap.Int32(0x5D4594, 1316520)), w1, w2)
	}
	return 1
}

func (c *Client) DrawEnergyBolt(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_energy_bolt_draw
	width := int8(effectBytes(dr)[1])
	s1, s2, w1, w2, ok := c.boltEnds(vp, dr, false)
	if !ok {
		return 1
	}
	st := &lightningState
	*memmap.PtrUint8(0x5D4594, 1316420) = byte(2 * (int(width) + 127))
	st.color3 = memmap.Uint32(0x5D4594, 1316496)
	st.color3b = memmap.Uint32(0x5D4594, 1316468)
	c.lightningBolt(s1, s2, 2, false, true)
	if !nox_xxx_checkGameFlagPause_413A50() {
		c.makeLightningParticles(int(memmap.Int32(0x5D4594, 1316524)), w1, w2)
	}
	return 1
}

func (c *Client) DrawGreenBolt(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_green_bolt_draw
	u := effectBytes(dr)
	if u[0] == 0 {
		if left := binary.LittleEndian.Uint32(u[1:]); left != 0 {
			left--
			binary.LittleEndian.PutUint32(u[1:], left)
			if left == 0 {
				c.Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr)
				return 0
			}
		}
	}
	s1, s2, w1, w2, ok := c.boltEnds(vp, dr, false)
	if !ok {
		return 1
	}
	st := &lightningState
	st.color1 = memmap.Uint32(0x5D4594, 1316444)
	st.color3 = memmap.Uint32(0x5D4594, 1316504)
	st.color2 = memmap.Uint32(0x5D4594, 1316460)
	st.color3b = memmap.Uint32(0x5D4594, 1316480)
	*memmap.PtrUint8(0x5D4594, 1316420) = 1
	c.lightningBolt(s1, s2, 2, true, true)
	if !nox_xxx_checkGameFlagPause_413A50() {
		c.makeLightningParticles(int(memmap.Int32(0x5D4594, 1316528)), w1, w2)
	}
	return 1
}

// plasmaSeg returns a field of a plasma ribbon segment.
// Each of the 3 ribbons has 30 segments of 7 values: x, y, radius, angle, angle speed, angle accel, timer.
func plasmaSeg(ribbon, seg, field int) *int32 {
	return memmap.PtrInt32(0x5D4594, 1313884+uintptr(28*(30*ribbon+seg)+4*field))
}

const (
	plasmaX = iota
	plasmaY
	plasmaRad
	plasmaAng
	plasmaAngSpeed
	plasmaAngAccel
	plasmaTimer
)

var typePlasmaSpark = spriteTypeCache{id: "CyanSpark"}

func (c *Client) plasmaAddSegment(p1, p2 image.Point, ribbon int) byte { // sub_4BA8B0
	vp := c.Cli().Viewport()
	typ := c.spriteType(&typePlasmaSpark)
	seg := int(legacy.Get_dword_5d4594_1316412())
	*plasmaSeg(ribbon, seg, plasmaX) = int32(p1.X)
	*plasmaSeg(ribbon, seg, plasmaY) = int32(p1.Y)
	*plasmaSeg(ribbon, seg+1, plasmaX) = int32(p2.X)
	*plasmaSeg(ribbon, seg+1, plasmaY) = int32(p2.Y)
	legacy.Set_dword_5d4594_1316412(int32(seg + 1))
	ret := byte(c.srv.Frame())
	if c.srv.Frame()&4 != 0 {
		v := c.rnd(0, 10)
		ret = byte(v)
		if v > 5 {
			c.drawEnergyBolt(p2.Add(vp.World.Min).Sub(vp.Screen.Min), 8, typ)
			ret = 0
		}
	}
	return ret
}

func (c *Client) drawPlasma(dir byte, x1, y1, x2, y2 int) { // sub_4BA230
	if memmap.Uint32(0x5D4594, 1316404) == 0 {
		*memmap.PtrUint32(0x5D4594, 1313828) = noxcolor.RGB5551Color(40, 180, 255).Color32()
		*memmap.PtrUint32(0x5D4594, 1313832) = noxcolor.RGB5551Color(255, 255, 255).Color32()
		*memmap.PtrUint32(0x5D4594, 1313836) = noxcolor.RGB5551Color(40, 180, 255).Color32()
		*memmap.PtrUint32(0x5D4594, 1313844) = 8
		*memmap.PtrUint32(0x5D4594, 1313852) = 8
		*memmap.PtrUint32(0x5D4594, 1313856) = 0
		*memmap.PtrUint32(0x5D4594, 1313860) = 0
		*memmap.PtrUint8(0x5D4594, 1313840) = 16
		*memmap.PtrUint8(0x5D4594, 1313842) = 16
		*memmap.PtrUint32(0x5D4594, 1313848) = 12
		*memmap.PtrUint8(0x5D4594, 1313841) = 24
		*memmap.PtrUint32(0x5D4594, 1313864) = 0
		for i := 0; i < 3; i++ {
			for j := 0; j < 30; j++ {
				for f := plasmaRad; f <= plasmaTimer; f++ {
					*plasmaSeg(i, j, f) = 0
				}
			}
		}
		*memmap.PtrUint32(0x5D4594, 1316404) = 1
	}
	legacy.Sub_4BA670(int(dir), x1, y1, x2, y2)
	for i := 0; i < 3; i++ {
		for j := 0; j < int(legacy.Get_dword_5d4594_1316408())+1; j++ {
			ang := plasmaSeg(i, j, plasmaAng)
			*ang += *plasmaSeg(i, j, plasmaAngSpeed)
			if *ang >= 256 {
				*ang -= 256
			} else if *ang < 0 {
				// The original code wraps by 255 here.
				*ang += 255
			}
			if j == 0 {
				*plasmaSeg(i, 0, plasmaAng) = int32(dir)
			}
			speed := plasmaSeg(i, j, plasmaAngSpeed)
			*speed += *plasmaSeg(i, j, plasmaAngAccel)
			lim := int32(25)
			if i == 1 {
				lim = 10
			}
			if *speed > lim {
				*speed = lim
			} else if *speed < -lim {
				*speed = -lim
			}
			timer := plasmaSeg(i, j, plasmaTimer)
			*timer--
			if *timer > 0 {
				continue
			}
			*timer = int32(c.rnd(10, 90))
			accel := plasmaSeg(i, j, plasmaAngAccel)
			*accel = int32(c.rnd(4, 8))
			if c.rnd(0, 1) != 0 {
				*accel = -*accel
			}
			rad := plasmaSeg(i, j, plasmaRad)
			if i == 1 {
				*rad = int32(c.rnd(40, 50))
			} else {
				*rad = int32(c.rnd(80, 110))
				if j < 4 && c.rnd(0, 100) > 90 {
					*rad = int32(c.rnd(150, 200))
					*timer *= 2
				}
			}
		}
	}
	n := int(legacy.Get_dword_5d4594_1316408())
	for i := 0; i < 3; i++ {
		cl := memmap.Uint32(0x5D4594, 1313828+4*uintptr(i))
		for j := 0; j < n; j++ {
			rad := float64(float32(*plasmaSeg(i, j, plasmaRad)))
			ang := *plasmaSeg(i, j, plasmaAng)
			ang2 := *plasmaSeg(i, j+1, plasmaAng)
			p1 := image.Pt(int(*plasmaSeg(i, j, plasmaX)), int(*plasmaSeg(i, j, plasmaY)))
			p2 := image.Pt(int(*plasmaSeg(i, j+1, plasmaX)), int(*plasmaSeg(i, j+1, plasmaY)))
			d1 := image.Pt(
				int(float32(rad*float64(memmap.Float32(0x587000, 194136+8*uintptr(ang))))),
				int(float32(rad*float64(memmap.Float32(0x587000, 194140+8*uintptr(ang))))),
			)
			d2 := image.Pt(
				int(float32(rad*float64(memmap.Float32(0x587000, 194136+8*uintptr(ang2))))),
				int(float32(rad*float64(memmap.Float32(0x587000, 194140+8*uintptr(ang2))))),
			)
			legacy.Sub_4BE800(cl)
			legacy.Sub_4BE810(1, cl, int(memmap.Int32(0x5D4594, 1313844+4*uintptr(i))), memmap.Uint8(0x5D4594, 1313840+uintptr(i)))
			thin := 0
			if i == 1 {
				thin = 1
			}
			legacy.Sub_4BEAD0(p1, p2, d1, d2, 8, thin)
		}
	}
}

func (c *Client) DrawPlasma(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_plasma_draw
	mpos := c.Cli().GetMousePos()
	u := effectBytes(dr)
	off := vp.Screen.Min.Sub(vp.World.Min)
	var (
		dir            byte
		x1, y1, x2, y2 int
	)
	if u[0] != 0 {
		from := c.byNetCodeHigh(binary.LittleEndian.Uint32(u[5:]))
		to := c.byNetCodeHigh(binary.LittleEndian.Uint32(u[9:]))
		if from != nil && to != nil {
			x1, y1 = from.PosVec.X+off.X, from.PosVec.Y+off.Y-20
			x2, y2 = to.PosVec.X+off.X, to.PosVec.Y+off.Y-20
		} else {
			// The original code uses the drawable and viewport pointers as coordinates here.
			x1 = int(int32(uintptr(dr.C())))
			y1 = int(int32(uintptr(vp.C())))
			x2, y2 = y1, y1
		}
		dir = byte(legacy.Nox_xxx_math_509ED0(types.Ptf(float32(float64(mpos.X)-float64(x1)), float32(float64(mpos.Y)-float64(y1)))))
	} else {
		dir = u[1]
		x1 = int(binary.LittleEndian.Uint16(u[5:])) + off.X
		y1 = int(binary.LittleEndian.Uint16(u[7:])) + off.Y - 20
		x2 = int(binary.LittleEndian.Uint16(u[9:])) + off.X
		y2 = int(binary.LittleEndian.Uint16(u[11:])) + off.Y - 20
	}
	c.drawPlasma(dir, x1, y1, x2, y2)
	return 1
}
