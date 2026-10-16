package opennox

import (
	"encoding/binary"
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
)

func init() {
	client.RegisterDraw("ConditionalAnimateDraw", legacy.Get_nox_thing_cond_animate_draw(), 4, parseCondAnimateDraw)
	client.RegisterDraw("MonsterGeneratorDraw", legacy.Get_nox_thing_monster_gen_draw(), 4, parseCondAnimateDraw)
	client.RegisterDraw("PlayerWaypointDraw", legacy.Get_nox_thing_player_waypoint_draw(), 0, nil)
	client.RegisterDraw("UndeadKillerDraw", legacy.Get_nox_thing_undead_killer_draw(), 0, nil)
	client.RegisterDraw("BlueRainDraw", legacy.Get_nox_thing_blue_rain_draw(), 0, nil)
	client.RegisterDraw("LevelUpDraw", legacy.Get_nox_thing_levelup_draw(), 0, nil)
	client.RegisterDraw("OblivionUpDraw", legacy.Get_nox_thing_oblivion_up_draw(), 0, nil)

	legacy.Nox_thing_maiden_draw = func(vp *noxrender.Viewport, dr *client.Drawable) int { return noxClient.DrawMaiden(vp, dr) }
	legacy.Nox_thing_player_waypoint_draw = func(vp *noxrender.Viewport, dr *client.Drawable) int { return noxClient.DrawPlayerWaypoint(vp, dr) }
	legacy.Nox_thing_monster_gen_draw = func(vp *noxrender.Viewport, dr *client.Drawable) int { return noxClient.DrawMonsterGen(vp, dr) }
	legacy.Nox_thing_cond_animate_draw = func(vp *noxrender.Viewport, dr *client.Drawable) int { return noxClient.DrawCondAnimate(vp, dr) }
	legacy.Nox_thing_blue_rain_draw = func(vp *noxrender.Viewport, dr *client.Drawable) int { return noxClient.DrawBlueRain(vp, dr) }
	legacy.Nox_thing_levelup_draw = func(vp *noxrender.Viewport, dr *client.Drawable) int { return noxClient.DrawLevelUp(vp, dr) }
	legacy.Nox_thing_oblivion_up_draw = func(vp *noxrender.Viewport, dr *client.Drawable) int { return noxClient.DrawOblivionUp(vp, dr) }
	legacy.Nox_thing_undead_killer_draw = func(vp *noxrender.Viewport, dr *client.Drawable) int { return noxClient.DrawUndeadKiller(vp, dr) }
	legacy.Nox_xxx_netDrawRays_49BDD0 = func(data []byte) { noxClient.netDrawRays(data) }
}

// condAnimateDrawData is a draw data for ConditionalAnimateDraw and MonsterGeneratorDraw.
type condAnimateDrawData struct {
	Size   uint32
	Images [5]unsafe.Pointer // *[Cnt]noxrender.ImageHandle
	Cnt    [5]uint8
	Delay  [5]uint8
	_      [2]byte
	Kind   [5]client.AnimKind
}

var _ = [1]struct{}{}[0x38-unsafe.Sizeof(condAnimateDrawData{})]

func (d *condAnimateDrawData) Frame(set, i int) noxrender.ImageHandle {
	return *(*noxrender.ImageHandle)(unsafe.Add(d.Images[set], uintptr(i)*unsafe.Sizeof(noxrender.ImageHandle(nil))))
}

func parseCondAnimateDraw(obj *client.ObjectType, f *binfile.MemFile, _ string, _ []byte) error { // nox_things_cond_animate_draw_parse
	dd, _ := alloc.New(condAnimateDrawData{})
	dd.Size = 16
	n := int(f.ReadU8())
	for i := 0; i < n; i++ {
		dd.Cnt[i] = f.ReadU8()
		dd.Delay[i] = f.ReadU8()
		kind, _ := f.ReadString8()
		dd.Kind[i] = client.ParseAnimKind(kind)
		dd.Images[i] = readThingImages(f, int(dd.Cnt[i]))
	}
	obj.DrawData = unsafe.Pointer(dd)
	obj.DrawFunc = legacy.Get_nox_thing_cond_animate_draw()
	obj.Field_60 = 0
	return nil
}

func (c *Client) DrawCondAnimate(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_cond_animate_draw
	dd := (*condAnimateDrawData)(dr.DrawData)
	set := 1
	if dr.ObjFlags.Has(object.FlagEnabled) {
		set = 0
	}
	cnt := int(dd.Cnt[set])
	var ind int
	switch dd.Kind[set] {
	case client.AnimLoop:
		ind = int(int32((c.srv.Frame() + dr.NetCode32) / (uint32(dd.Delay[set]) + 1)))
		if ind >= cnt {
			ind %= cnt
		}
	case client.AnimRandom:
		ind = c.rnd(0, cnt)
	case client.AnimSlave:
		ind = int(dr.AnimFrameSlave)
	default:
		return 0
	}
	legacy.Nox_xxx_drawObject_4C4770_draw(vp, dr, dd.Frame(set, ind))
	return 1
}

func (c *Client) DrawMonsterGen(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_monster_gen_draw
	dd := (*condAnimateDrawData)(dr.DrawData)
	set := 0
	if fl := dr.Flags70Val; fl&0x100 != 0 {
		set = 1
	} else if fl&0x200 != 0 {
		set = 2
	} else if fl&0xC00 != 0 {
		set = 3
	}
	cnt := int(dd.Cnt[set])
	var ind int
	switch dd.Kind[set] {
	case client.AnimOneShot:
		// The original code uses the address of the delay value as the frame index.
		ind = int(uintptr(unsafe.Pointer(&dd.Delay[set])))
	case client.AnimLoop:
		ind = int(int32((c.srv.Frame() + dr.NetCode32) / (uint32(dd.Delay[set]) + 1)))
		if ind >= cnt {
			ind %= cnt
		}
	case client.AnimRandom:
		ind = c.rnd(0, cnt)
	case client.AnimSlave:
		ind = int(dr.AnimFrameSlave)
	default:
		return 0
	}
	if dr.Flags70Val&0x800 != 0 {
		ind = cnt - 1
		dr.ObjClass &^= 0x80000
		dr.ObjFlags &^= 0x20000000
	}
	u := dr.UnionEffect()
	if left := u.Field_108; left != 0 {
		delay := uint32(dd.Delay[set]) + 1
		ind = int(int32((delay*uint32(dd.Cnt[set]) - left) / delay))
		if ind >= cnt {
			ind = cnt - 1
		}
		if ind < 0 {
			ind = 0
		}
		u.Field_108 = left - 1
		if u.Field_108 == 0 {
			dr.Flags70Val = dr.Flags70Val&^0x400 | 0x800
		}
	}
	legacy.Nox_xxx_drawObject_4C4770_draw(vp, dr, dd.Frame(set, ind))
	if dr.Flags70Val&0xC00 == 0 {
		cnt4 := int(dd.Cnt[4])
		ind4 := int(int32((c.srv.Frame() + dr.NetCode32) / (uint32(dd.Delay[4]) + 1)))
		if ind4 >= cnt4 {
			ind4 %= cnt4
		}
		legacy.Nox_xxx_drawObject_4C4770_draw(vp, dr, dd.Frame(4, ind4))
	}
	if dr.Flags70Val&0x800 != 0 {
		dr.ObjFlags |= 1
	}
	return 1
}

func (c *Client) DrawMaiden(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_maiden_draw
	r := c.r
	if !noxflags.HasGame(0x200000) {
		npc := c.srv.NPCs.ByID(int(dr.NetCode32))
		if npc == nil {
			return 1
		}
		for i, cl := range npc.Color8 {
			r.Data().SetMaterial(i+1, noxcolor.RGBA5551(cl))
		}
		return c.DrawMonster(vp, dr)
	}
	for obj := c.srv.Objs.First(); obj != nil; obj = obj.Next() {
		if obj.Extent != dr.NetCode32 {
			continue
		}
		cl := unsafe.Slice((*byte)(unsafe.Add(obj.UpdateData, 2076)), 18)
		for i := 0; i < 6; i++ {
			r.Data().SetMaterialRGB(i+1, int(cl[3*i+0]), int(cl[3*i+1]), int(cl[3*i+2]))
		}
		break
	}
	return c.DrawMonster(vp, dr)
}

func (c *Client) DrawPlayerWaypoint(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_player_waypoint_draw
	r := c.r
	cl := memmap.Uint32(0x85B3FC, 940)
	// The original code ignores the screen offset of the viewport here.
	p := dr.Pos().Sub(vp.World.Min)
	legacy.Nox_xxx_spriteDrawCircleMB_4C32A0(p.X, p.Y, 10, cl)
	ang := int(byte(2 * byte(c.srv.Frame())))
	r.Data().SetAlphaEnabled(true)
	r.Data().SetColor2(noxcolor.RGBA5551(cl))
	for i := 0; i < 5; i++ {
		dx1, dy1 := dirTable(ang)
		next := (ang + 102) % 256
		dx2, dy2 := dirTable(next)
		r.AddPoint(image.Pt(p.X+10*dx1/16, p.Y+10*dy1/16))
		r.AddPoint(image.Pt(p.X+10*dx2/16, p.Y+10*dy2/16))
		r.DrawLineFromPoints(r.Data().Color2())
		ang = next
	}
	r.Data().SetAlphaEnabled(false)
	return 1
}

var (
	typeBlueRainSpark = spriteTypeCache{id: "BlueRainSpark"}
	typeRainOrbWhite  = spriteTypeCache{id: "RainOrbWhite"}
	typeRainOrbBlue   = spriteTypeCache{id: "RainOrbBlue"}
	typeUndeadOrb     = spriteTypeCache{id: "WhiteOrb"}
)

func (c *Client) DrawBlueRain(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_blue_rain_draw
	if noxflags.HasGame(0x200000) {
		return 1
	}
	typ := c.spriteType(&typeBlueRainSpark)
	for i := 0; i < 2; i++ {
		x := dr.PosVec.X + c.rnd(-10, 10)
		y := dr.PosVec.Y + c.rnd(-10, 10)
		s := c.Nox_xxx_spriteLoadAdd_45A360_drawable(typ, image.Pt(x, y))
		if s == nil {
			continue
		}
		c.setSparkMotion(s, image.Pt(x, y), 0, 0, c.srv.Frame()+uint32(c.rnd(90, 120)))
		s.ZVal2 = 0
		s.VelZ = -5
		s.ZVal = uint16(y - vp.World.Min.Y)
		c.Cli().Objs.List34Add(s)
	}
	return 1
}

func (c *Client) drawFallingSparks(typ int, vp *noxrender.Viewport, dr *client.Drawable) { // nox_thing_falling_sparks_draw_4B7740
	for i := 0; i < 2; i++ {
		p := image.Pt(dr.PosVec.X+c.rnd(-15, 15), dr.PosVec.Y+c.rnd(-15, 15))
		z := p.Y - vp.World.Min.Y
		vz := c.rnd(8, 12)
		c.spawnFx499950(typ, p, dr.PosVec, uint16(z), -int8(vz))
	}
}

func (c *Client) DrawLevelUp(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_levelup_draw
	c.drawFallingSparks(c.spriteType(&typeRainOrbWhite), vp, dr)
	return 1
}

func (c *Client) DrawOblivionUp(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_oblivion_up_draw
	c.drawFallingSparks(c.spriteType(&typeRainOrbBlue), vp, dr)
	return 1
}

func (c *Client) DrawUndeadKiller(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_undead_killer_draw
	r := c.r
	typ := c.spriteType(&typeUndeadOrb)
	cl := noxcolor.RGB5551Color(100, 100, 255)
	if c.srv.Frame()-dr.AnimStart > 0x46 {
		c.Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr)
		return 0
	}
	if c.rnd(0, 100) > 85 {
		var pts [4]uint16
		pts[0] = uint16(dr.Field_81)
		pts[1] = uint16(dr.Field_82)
		pts[2] = uint16(dr.PosVec.X) + uint16(c.rnd(-5, 5))
		pts[3] = uint16(dr.PosVec.Y) + uint16(c.rnd(-5, 5))
		a5 := c.rnd(6, 10)
		a4 := c.rnd(-5, 5)
		c.spawnFx499490(typ, pts, 0, a4, byte(a5), 0)
	}
	p := vp.ToScreenPos(dr.Pos())
	p.Y -= int(int16(dr.ZVal2)) + int(int16(dr.ZVal))
	d := int(int32(legacy.Sub_48C690(dr.PosVec.X, dr.PosVec.Y, int(int32(dr.Field_81)), int(int32(dr.Field_82)))))
	rad := 8 - d/40
	if rad < 0 {
		rad = 1
	}
	r.DrawGlow(p, cl, rad, 12)
	return 1
}

var (
	typeRayLightning      = spriteTypeCache{id: "DynamicLightning"}
	typeRayChainLightning = spriteTypeCache{id: "DynamicChainLightning"}
	typeRayEnergyBolt     = spriteTypeCache{id: "DynamicEnergyBolt"}
	typeRayOrb            = spriteTypeCache{id: "OrbRay"}
	typeRayPlasma         = spriteTypeCache{id: "PlasmaRay"}
	typeRayDrainManaOrb   = spriteTypeCache{id: "DrainManaOrb"}
	typeRayHealOrb        = spriteTypeCache{id: "HealOrb"}
	typeRayCharmOrb       = spriteTypeCache{id: "CharmOrb"}
)

func (c *Client) netDrawRays(data []byte) { // nox_xxx_netDrawRays_49BDD0
	cnt := memmap.PtrInt32(0x5D4594, 1304308)
	if *cnt >= 96 {
		return
	}
	var pts [4]uint16
	for i := range pts {
		pts[i] = binary.LittleEndian.Uint16(data[1+2*i:])
	}
	x := int(pts[0]) + (int(pts[2])-int(pts[0]))/2
	y := int(pts[1]) + (int(pts[3])-int(pts[1]))/2
	orbs := func(orb *spriteTypeCache, pts [4]uint16, a5 int) {
		if c.rnd(0, 100) < 50 {
			a4 := c.rnd(-20, 20)
			a3 := c.rnd(-20, 20)
			c.spawnFx499490(c.spriteType(orb), pts, a3, a4, byte(a5), 0)
		}
	}
	var typ int
	switch data[0] {
	case 0x7D:
		typ = c.spriteType(&typeRayPlasma)
	case 0x8C:
		typ = c.spriteType(&typeRayLightning)
	case 0x8D:
		typ = c.spriteType(&typeRayEnergyBolt)
	case 0x8E:
		typ = c.spriteType(&typeRayChainLightning)
	case 0x8F:
		a5 := c.rnd(6, 12)
		typ = c.spriteType(&typeRayOrb)
		orbs(&typeRayDrainManaOrb, pts, a5)
	case 0x90:
		a5 := c.rnd(6, 12)
		typ = c.spriteType(&typeRayOrb)
		orbs(&typeRayCharmOrb, pts, a5)
		orbs(&typeRayCharmOrb, [4]uint16{pts[2], pts[3], pts[0], pts[1]}, a5)
	case 0x91:
		a5 := c.rnd(6, 12)
		typ = c.spriteType(&typeRayOrb)
		orbs(&typeRayHealOrb, pts, a5)
	default:
		return
	}
	s := c.Nox_xxx_spriteLoadAdd_45A360_drawable(typ, image.Pt(x, y))
	if s == nil {
		return
	}
	u := effectBytes(s)
	u[0] = 0
	copy(u[5:13], data[1:9])
	*memmap.PtrPtr(0x5D4594, 1303540+4*uintptr(*cnt)) = s.C()
	*cnt++
}
