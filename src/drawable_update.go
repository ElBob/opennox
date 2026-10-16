package opennox

import (
	"encoding/binary"
	"image"

	noxcolor "github.com/opennox/libs/color"
	"github.com/opennox/libs/types"

	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
)

func init() {
	legacy.Nox_xxx_updDrawTeleportWake_4CD8D0 = func(vp *noxrender.Viewport, dr *client.Drawable) { noxClient.updTeleportWake(dr) }
	legacy.Nox_xxx_updDrawVortexSource_4CC950 = func(vp *noxrender.Viewport, dr *client.Drawable) { noxClient.updVortexSource(dr) }
	legacy.Nox_xxx_updDrawManabombCharge_4CCAC0 = func(vp *noxrender.Viewport, dr *client.Drawable) { noxClient.updManabombCharge(dr) }
	legacy.Nox_xxx_updDrawSparkleTrail_4CDBF0 = func(vp *noxrender.Viewport, dr *client.Drawable) { noxClient.updSparkleTrail(dr) }
	legacy.Nox_xxx_updDrawMagic_4CDD80 = func(vp *noxrender.Viewport, dr *client.Drawable) { noxClient.updMagic(dr) }
	legacy.Nox_xxx_updDrawMagicMissile_4CD9E0 = func(vp *noxrender.Viewport, dr *client.Drawable) { noxClient.updMagicMissile(dr) }
	legacy.Nox_xxx_updDrawDBallCharge_4CE0C0 = func(vp *noxrender.Viewport, dr *client.Drawable) { noxClient.updDeathBallCharge(dr) }
	legacy.Sub_4CD450 = func(vp *noxrender.Viewport, dr *client.Drawable) { noxClient.updHealOrbs(vp, dr) }
	legacy.Sub_4CD690 = func(vp *noxrender.Viewport, dr *client.Drawable) { noxClient.updDrainManaOrbs(vp, dr) }
	legacy.Sub_4CE200 = func(dr *client.Drawable, n, rad int) { noxClient.spawnCloudPuffs(dr, n, rad) }
	legacy.Nox_xxx_updDrawAddRndSpark_4CDFA0 = func(dr *client.Drawable, n int) *client.Drawable { return noxClient.spawnDeathBallSparks(dr, n) }
	legacy.Sub_4CCEA0 = func(dr *client.Drawable, mult int) { noxClient.spawnFireballSparks(dr, mult) }
	legacy.Sub_4CD150 = func(typ int, vp *noxrender.Viewport, dr *client.Drawable, dir bool) {
		noxClient.spawnCharmOrb(typ, vp, dr, dir)
	}
}

// spriteTypeAt returns a client object type index cached in the original memory.
func spriteTypeAt(off uintptr) int {
	return int(memmap.Int32(0x5D4594, off))
}

// spriteTypeCached returns a client object type index cached in the original memory, looking it up by name if not set.
func (c *Client) spriteTypeCached(off uintptr, name string) int {
	p := memmap.PtrUint32(0x5D4594, off)
	if *p == 0 {
		*p = uint32(c.Cli().Things.IndByID(name))
	}
	return int(int32(*p))
}

func (c *Client) rnd(min, max int) int { // nox_common_randomIntMinMax_415FF0
	return c.srv.Rand.Other.Int(min, max)
}

// dirTable returns a 1/16 scaled direction vector from the angle table.
func dirTable(ang int) (int, int) {
	return int(memmap.Int32(0x587000, 192088+8*uintptr(ang))), int(memmap.Int32(0x587000, 192092+8*uintptr(ang)))
}

// setSparkMotion initializes a particle that moves away from its spawn point.
func (c *Client) setSparkMotion(s *client.Drawable, pos image.Point, dir byte, speed uint32, deadline uint32) {
	u := s.UnionEffect()
	u.Field_108 = uint32(pos.X) << 12
	u.Field_109 = uint32(pos.Y) << 12
	s.Field_74_4 = dir
	u.Field_110 = speed
	u.Field_112 = deadline
	u.Field_111 = c.srv.Frame()
}

func (c *Client) updTeleportWake(dr *client.Drawable) { // nox_xxx_updDrawTeleportWake_4CD8D0
	typ := c.spriteTypeCached(1522980, "BlueSpark")
	x := dr.PosVec.X + c.rnd(-5, 5)
	dy := c.rnd(-5, 5)
	s := c.Nox_xxx_spriteLoadAdd_45A360_drawable(typ, image.Pt(x, dr.PosVec.Y+dy))
	if s == nil {
		return
	}
	dir := byte(c.rnd(0, 255))
	speed := uint32(c.rnd(1, 100))
	c.setSparkMotion(s, s.PosVec, dir, speed, c.srv.Frame()+uint32(c.rnd(10, 32)))
	s.ZVal = 0
	s.VelZ = int8(c.rnd(3, 8))
	c.Cli().Objs.List34Add(s)
}

func (c *Client) updVortexSource(dr *client.Drawable) { // nox_xxx_updDrawVortexSource_4CC950
	if memmap.Uint32(0x5D4594, 1522952) == 0 {
		*memmap.PtrUint32(0x5D4594, 1522952) = uint32(c.Cli().Things.IndByID("WhiteVortexOrb"))
		*memmap.PtrUint32(0x5D4594, 1522944) = noxcolor.RGB5551Color(200, 200, 200).Color32()
		*memmap.PtrUint32(0x5D4594, 1522948) = noxcolor.RGB5551Color(255, 255, 255).Color32()
	}
	gray := memmap.Uint32(0x5D4594, 1522944)
	typ := spriteTypeAt(1522952)
	white := memmap.Uint32(0x5D4594, 1522948)
	ang := c.rnd(0, 255)
	dx, dy := dirTable(ang)
	s := c.Nox_xxx_spriteLoadAdd_45A360_drawable(typ, image.Pt(dr.PosVec.X+50*dx/16, dr.PosVec.Y+50*dy/16))
	if s == nil {
		return
	}
	u := effectBytes(s)
	u[16] = byte(ang)
	dr.ZVal = 0
	u[17] = byte(c.rnd(2, 3))
	if c.rnd(0, 100) > 50 {
		u[17] = -u[17]
	}
	u[19] = 1
	u[18] = 50
	binary.LittleEndian.PutUint32(u[8:], uint32(dr.PosVec.X))
	binary.LittleEndian.PutUint32(u[12:], uint32(dr.PosVec.Y))
	binary.LittleEndian.PutUint32(u[0:], gray)
	binary.LittleEndian.PutUint32(u[4:], white)
	c.Cli().Objs.List34Add(s)
	c.Cli().Objs.List6Add(s)
}

func (c *Client) updManabombCharge(dr *client.Drawable) { // nox_xxx_updDrawManabombCharge_4CCAC0
	rad := int(c.srv.Balance.Float("ManaBombOutRadius"))
	orbP := legacy.Get_dword_5d4594_1522956_ptr()
	if *orbP == 0 {
		*orbP = uint32(c.Cli().Things.IndByID("ManaBombOrb"))
		*memmap.PtrUint32(0x5D4594, 1522960) = uint32(c.Cli().Things.IndByID("VioletSpark"))
	}
	orb := int(int32(*orbP))
	spark := spriteTypeAt(1522960)
	for i := 0; i < 20; i++ {
		d := rad/4 + c.rnd(0, rad)
		if d > rad {
			d = rad
		}
		dx, dy := dirTable(c.rnd(0, 255))
		s := c.Nox_xxx_spriteLoadAdd_45A360_drawable(spark, image.Pt(dr.PosVec.X+d*dx/16, dr.PosVec.Y+d*dy/16))
		if s == nil {
			continue
		}
		c.setSparkMotion(s, s.PosVec, 0, 0, c.srv.Frame()+uint32(c.rnd(10, 30)))
		s.ZVal = 0
		s.VelZ = int8(c.rnd(2, 8))
		c.Cli().Objs.List34Add(s)
	}
	frame := c.srv.Frame()
	if frame&1 == 0 || frame-dr.Field_80 >= 10 {
		return
	}
	var pts [4]int16
	pts[0] = int16(dr.PosVec.X)
	pts[1] = int16(dr.PosVec.Y)
	for ang := int(frame % 51); ang < 256; ang += 51 {
		pts[2] = pts[0] + int16(rad/16)*int16(memmap.Uint16(0x587000, 192088+8*uintptr(ang)))
		pts[3] = pts[1] + int16(rad/16)*int16(memmap.Uint16(0x587000, 192088+8*uintptr(ang)+4))
		legacy.Sub_499520(orb, pts, int16(ang), 0, 0)
		legacy.Sub_499520(orb, pts, int16(ang), 1, 0)
	}
}

func (c *Client) updSparkleTrail(dr *client.Drawable) { // nox_xxx_updDrawSparkleTrail_4CDBF0
	dx := dr.PosVec.X - int(int32(dr.Field_8))
	dy := dr.PosVec.Y - int(int32(dr.Field_9))
	typ := c.spriteTypeCached(1522996, "BlueSpark")
	var ax, ay int
	for i := 0; i < 5; i++ {
		x := int(int32(dr.Field_8)) + ax/5 + c.rnd(-3, 3)
		y := int(int32(dr.Field_9)) + c.rnd(-3, 3) + ay/5
		if s := c.Nox_xxx_spriteLoadAdd_45A360_drawable(typ, image.Pt(x, y)); s != nil {
			s.DrawFuncPtr = legacy.Get_nox_thing_pixie_dust_draw()
			s.SetLightColor(255, 200, 75)
			c.setSparkMotion(s, image.Pt(x, y), 0, 0, c.srv.Frame()+uint32(c.rnd(2, 10)))
			s.ZVal = dr.ZVal
			s.ZVal2 = dr.ZVal2
			s.VelZ = 0
			c.Cli().Objs.List34Add(s)
		}
		ax += dx
		ay += dy
	}
}

// spawnTailLink adds a tail link from the previous position of the drawable, if it moved far enough.
func (c *Client) spawnTailLink(dr *client.Drawable, typ func() int, decay int) {
	u := dr.UnionEffect()
	dx := uint32(dr.PosVec.X) - u.Field_108
	dy := uint32(dr.PosVec.Y) - u.Field_109
	if dx*dx+dy*dy <= 200 {
		return
	}
	tail := c.Nox_xxx_spriteLoadAdd_45A360_drawable(typ(), image.Pt(int(int32(u.Field_108)), int(int32(u.Field_109))))
	tu := tail.UnionEffect()
	tu.Field_108 = uint32(dr.PosVec.X)
	tu.Field_109 = uint32(dr.PosVec.Y)
	c.Cli().Objs.List34Add(tail)
	u.Field_108 = uint32(dr.PosVec.X)
	u.Field_109 = uint32(dr.PosVec.Y)
	c.Cli().Objs.TransparentDecay(tail, decay)
}

func (c *Client) updMagic(dr *client.Drawable) { // nox_xxx_updDrawMagic_4CDD80
	c.spawnTailLink(dr, func() int { return c.spriteTypeCached(1523000, "MagicTailLink") }, int(c.srv.TickRate()))
	dx := dr.PosVec.X - int(int32(dr.Field_8))
	dy := dr.PosVec.Y - int(int32(dr.Field_9))
	typ := c.spriteTypeCached(1523004, "BlueSpark")
	var ax, ay int
	for i := 0; i < 4; i++ {
		x := int(int32(dr.Field_8)) + ax/4 + c.rnd(-8, 8)
		ry := c.rnd(-8, 8)
		if s := c.Nox_xxx_spriteLoadAdd_45A360_drawable(typ, image.Pt(x, int(int32(dr.Field_9))+ay/4+ry)); s != nil {
			s.DrawFuncPtr = legacy.Get_nox_thing_magic_sparkle_draw()
			s.SetLightColor(128, 128, 255)
			c.setSparkMotion(s, dr.PosVec, 0, 0, c.srv.Frame()+uint32(c.rnd(10, 20)))
			s.ZVal = dr.ZVal
			s.ZVal2 = dr.ZVal2
			s.VelZ = 0
			c.Cli().Objs.List34Add(s)
		}
		ax += dx
		ay += dy
	}
}

func (c *Client) updMagicMissile(dr *client.Drawable) { // nox_xxx_updDrawMagicMissile_4CD9E0
	if memmap.Uint32(0x5D4594, 1522984) == 0 {
		*memmap.PtrUint32(0x5D4594, 1522984) = uint32(c.Cli().Things.IndByID("Spark"))
		*memmap.PtrUint32(0x5D4594, 1522988) = uint32(c.Cli().Things.IndByID("Puff"))
		*memmap.PtrUint32(0x5D4594, 1522992) = uint32(c.Cli().Things.IndByID("MagicMissileTailLink"))
	}
	dx := dr.PosVec.X - int(int32(dr.Field_8))
	dy := dr.PosVec.Y - int(int32(dr.Field_9))
	// The original tests the count as unsigned first, and then loops with a signed test at the end.
	if memmap.Uint32(0x587000, 190108) > 0 {
		for i := 0; ; {
			x := dx*c.rnd(0, 100)/100 + int(int32(dr.Field_8))
			y := dy*c.rnd(0, 100)/100 + int(int32(dr.Field_9))
			if s := c.Nox_xxx_spriteLoadAdd_45A360_drawable(spriteTypeAt(1522984), image.Pt(x, y)); s != nil {
				dir := byte(c.rnd(0, 255))
				c.setSparkMotion(s, image.Pt(x, y), dir, 0, c.srv.Frame()+uint32(c.rnd(3, 10)))
				s.ZVal = 20
				s.VelZ = int8(c.rnd(0, 6))
				c.Cli().Objs.List34Add(s)
			}
			i++
			if i >= int(memmap.Int32(0x587000, 190108)) {
				break
			}
		}
	}
	c.spawnTailLink(dr, func() int { return spriteTypeAt(1522992) }, int(c.srv.TickRate()/3))
}

// spawnDeathBallSparks returns the last spawned sprite (nil if it failed), like the original.
func (c *Client) spawnDeathBallSparks(dr *client.Drawable, n int) *client.Drawable { // nox_xxx_updDrawAddRndSpark_4CDFA0
	typ := c.spriteTypeCached(1523008, "DeathBallSpark")
	var last *client.Drawable
	for ; n > 0; n-- {
		s := c.Nox_xxx_spriteLoadAdd_45A360_drawable(typ, dr.PosVec)
		last = s
		if s == nil {
			continue
		}
		dir := byte(c.rnd(0, 255))
		speed := uint32(c.rnd(1000, 3000))
		c.setSparkMotion(s, dr.PosVec, dir, speed, c.srv.Frame()+uint32(c.rnd(10, 40)))
		s.ZVal = 22
		s.VelZ = int8(c.rnd(0, 4))
		c.Cli().Objs.List34Add(s)
	}
	return last
}

func (c *Client) updDeathBallCharge(dr *client.Drawable) { // nox_xxx_updDrawDBallCharge_4CE0C0
	typ := c.spriteTypeCached(1523012, "CharmOrb")
	var pts [4]uint16
	pts[0] = uint16(dr.PosVec.X)
	pts[1] = uint16(dr.PosVec.Y)
	for i := 0; i < 10; i++ {
		ang := c.rnd(0, 255)
		d := c.rnd(2, 8)
		pts[3] = uint16(dr.PosVec.Y) + uint16(d*int(int16(memmap.Uint16(0x587000, 192092+8*uintptr(ang)))))
		pts[2] = uint16(dr.PosVec.X) + uint16(d*int(int16(memmap.Uint16(0x587000, 192088+8*uintptr(ang)))))
		if c.rnd(0, 100) < 50 {
			a5 := c.rnd(6, 10)
			a4 := c.rnd(-20, 20)
			a3 := c.rnd(-20, 20)
			legacy.Sub_499490(typ, pts, a3, a4, byte(a5), 0)
		}
	}
}

func (c *Client) spawnCloudPuffs(dr *client.Drawable, n, rad int) { // sub_4CE200
	if memmap.Uint32(0x5D4594, 1523016) == 0 {
		*memmap.PtrUint32(0x5D4594, 1523016) = uint32(c.Cli().Things.IndByID("GreenPuff"))
		*memmap.PtrUint32(0x5D4594, 1523020) = uint32(c.Cli().Things.IndByID("GreenSmoke"))
	}
	puff := spriteTypeAt(1523016)
	smoke := spriteTypeAt(1523020)
	for ; n > 0; n-- {
		dx, dy := dirTable(c.rnd(0, 255))
		d := c.rnd(0, rad)
		pos := image.Pt(dr.PosVec.X+d*dx/16, dr.PosVec.Y+d*dy/16)
		typ := puff
		if c.rnd(0, 10) < 3 {
			typ = smoke
		}
		s := c.Nox_xxx_spriteLoadAdd_45A360_drawable(typ, pos)
		if s == nil {
			continue
		}
		s.ZVal = 0
		c.Cli().Objs.List34Add(s)
		effectBytes(s)[0] = byte(c.rnd(1, 3))
		s.Field_115 = legacy.Get_sub_4CE340()
		c.Cli().Objs.TransparentDecay(s, c.rnd(10, 32))
		c.Cli().Objs.List5Add(s)
		c.Cli().Objs.List6Add(s)
	}
}

func (c *Client) spawnFireballSparks(dr *client.Drawable, mult int) { // sub_4CCEA0
	typ := c.spriteTypeCached(1522964, "Spark")
	dx := dr.PosVec.X - int(int32(dr.Field_8))
	dy := dr.PosVec.Y - int(int32(dr.Field_9))
	n := (absInt(dx) + absInt(dy)) / 7
	for ; n > 0; n-- {
		t := c.rnd(0, 100)
		x := int(int32(dr.Field_8)) + dx*t/100
		y := int(int32(dr.Field_9)) + dy*t/100
		s := c.Nox_xxx_spriteLoadAdd_45A360_drawable(typ, image.Pt(x, y))
		if s == nil {
			continue
		}
		u := s.UnionEffect()
		u.Field_108 = uint32(x) << 12
		u.Field_109 = uint32(y) << 12
		dir := legacy.Nox_xxx_math_509ED0(types.Ptf(float32(-dx), float32(-dy))) + c.rnd(-25, 25)
		if dir < 0 {
			dir += int(uint32(255-dir) >> 8 << 8)
		}
		s.Field_74_4 = byte(dir)
		u.Field_110 = uint32(mult * c.rnd(100, 300))
		u.Field_112 = c.srv.Frame() + uint32(c.rnd(30, 45))
		u.Field_111 = c.srv.Frame()
		s.ZVal = 28
		s.ZVal2 = 0
		s.VelZ = int8(c.rnd(-2, 4))
		c.Cli().Objs.List34Add(s)
	}
}

func absInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

// byNetCodeHigh finds a drawable by a net code with the static flag in the high bit.
func (c *Client) byNetCodeHigh(code uint32) *client.Drawable {
	id := int(nox_xxx_netClearHighBit_578B30(uint16(code)))
	if nox_xxx_netTestHighBit_578B70(uint16(code)) {
		return c.Cli().Objs.ByNetCodeStatic(id)
	}
	return c.Cli().Objs.ByNetCodeDynamic(id)
}

// clampOrbToViewport keeps the spawn position of an orb on the screen.
func clampOrbToViewport(vp *noxrender.Viewport, dr *client.Drawable, x, y int) (int, int) {
	z := int(int16(dr.ZVal))
	sx := x + vp.Screen.Min.X - vp.World.Min.X
	sy := y + vp.Screen.Min.Y - vp.World.Min.Y - z
	if sx < 0 {
		x = vp.World.Min.X + vp.Screen.Min.X + 1
	}
	if sy < 0 {
		y = vp.World.Min.Y - z + vp.Screen.Min.Y + 1
	}
	if sx >= vp.Size.X {
		x = vp.Screen.Max.X + vp.World.Min.X - 1
	}
	if sy >= vp.Size.Y {
		y = vp.Screen.Max.Y - z + vp.World.Min.Y - 1
	}
	return x, y
}

func (c *Client) spawnOrb(typ int, vp *noxrender.Viewport, dr *client.Drawable, x, y int, tx, ty uint16) {
	x, y = clampOrbToViewport(vp, dr, x, y)
	s := c.Nox_xxx_spriteLoadAdd_45A360_drawable(typ, image.Pt(x, y))
	if s == nil {
		return
	}
	v := byte(c.rnd(6, 12))
	u := effectBytes(s)
	binary.LittleEndian.PutUint16(u[0:], tx)
	binary.LittleEndian.PutUint16(u[2:], ty)
	u[11] = v
	u[12] = byte(c.rnd(3, 10))
}

// updTargetOrbs spawns orbs flying from the target to the source.
func (c *Client) updTargetOrbs(typ int, vp *noxrender.Viewport, dr *client.Drawable) {
	if c.rnd(0, 100) >= 50 {
		return
	}
	u := effectBytes(dr)
	var (
		x, y   int
		tx, ty uint16
	)
	if u[0] != 0 {
		from := c.byNetCodeHigh(binary.LittleEndian.Uint32(u[5:]))
		to := c.byNetCodeHigh(binary.LittleEndian.Uint32(u[9:]))
		if from == nil || to == nil {
			return
		}
		x = to.PosVec.X + c.rnd(-20, 20)
		y = to.PosVec.Y + c.rnd(-20, 20)
		tx, ty = uint16(from.PosVec.X), uint16(from.PosVec.Y)
	} else {
		x = int(binary.LittleEndian.Uint16(u[9:])) + c.rnd(-20, 20)
		y = int(binary.LittleEndian.Uint16(u[11:])) + c.rnd(-20, 20)
		tx, ty = binary.LittleEndian.Uint16(u[5:]), binary.LittleEndian.Uint16(u[7:])
	}
	c.spawnOrb(typ, vp, dr, x, y, tx, ty)
}

func (c *Client) updHealOrbs(vp *noxrender.Viewport, dr *client.Drawable) { // sub_4CD450
	c.updTargetOrbs(c.spriteTypeCached(1522972, "HealOrb"), vp, dr)
}

func (c *Client) updDrainManaOrbs(vp *noxrender.Viewport, dr *client.Drawable) { // sub_4CD690
	c.updTargetOrbs(c.spriteTypeCached(1522976, "DrainManaOrb"), vp, dr)
}

func (c *Client) spawnCharmOrb(typ int, vp *noxrender.Viewport, dr *client.Drawable, rev bool) { // sub_4CD150
	if c.rnd(0, 100) >= 50 {
		return
	}
	u := effectBytes(dr)
	var (
		x, y   int
		tx, ty uint16
	)
	if u[0] == 0 {
		if rev {
			x = int(binary.LittleEndian.Uint16(u[9:])) + c.rnd(-20, 20)
			y = int(binary.LittleEndian.Uint16(u[11:])) + c.rnd(-20, 20)
			tx, ty = binary.LittleEndian.Uint16(u[5:]), binary.LittleEndian.Uint16(u[7:])
		} else {
			x = int(binary.LittleEndian.Uint16(u[5:])) + c.rnd(-20, 20)
			y = int(binary.LittleEndian.Uint16(u[7:])) + c.rnd(-20, 20)
			tx, ty = binary.LittleEndian.Uint16(u[9:]), binary.LittleEndian.Uint16(u[11:])
		}
	} else {
		from := c.byNetCodeHigh(binary.LittleEndian.Uint32(u[5:]))
		to := c.byNetCodeHigh(binary.LittleEndian.Uint32(u[9:]))
		if from == nil || to == nil {
			return
		}
		if rev {
			from, to = to, from
		}
		x = from.PosVec.X + c.rnd(-20, 20)
		y = from.PosVec.Y + c.rnd(-20, 20)
		tx, ty = uint16(to.PosVec.X), uint16(to.PosVec.Y)
	}
	c.spawnOrb(typ, vp, dr, x, y, tx, ty)
}
