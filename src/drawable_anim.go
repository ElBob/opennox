package opennox

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
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
	client.RegisterDraw("AnimateDraw", legacy.Get_nox_thing_animate_draw(), 3, parseAnimateDraw)
	client.RegisterDraw("AnimateStateDraw", legacy.Get_nox_thing_animate_state_draw(), 8, parseAnimateStateDraw)
	client.RegisterDraw("SummonEffectDraw", legacy.Get_nox_thing_summon_effect_draw(), 3, parseAnimateDraw)
	client.RegisterDraw("ArrowTailLinkDraw", legacy.Get_nox_thing_arrow_tail_link_draw(), 0, nil)
	client.RegisterDraw("WeakArrowTailLinkDraw", legacy.Get_nox_thing_weak_arrow_tail_link_draw(), 0, nil)
	client.RegisterDraw("HarpoonRopeDraw", legacy.Get_nox_thing_harpoon_rope_draw(), 0, nil)
	client.RegisterDraw("SpiderSpitDraw", legacy.Get_nox_thing_spider_spit_draw(), 0, nil)
}

// animateDrawData is a draw data for AnimateDraw and similar.
type animateDrawData struct {
	Size   uint32
	Images unsafe.Pointer // *[Cnt]noxrender.ImageHandle
	Cnt    uint8
	Delay  uint8
	_      [2]byte
	Kind   client.AnimKind
}

func (d *animateDrawData) Frame(i int) noxrender.ImageHandle {
	return *(*noxrender.ImageHandle)(unsafe.Add(d.Images, uintptr(i)*unsafe.Sizeof(noxrender.ImageHandle(nil))))
}

func readThingImages(f *binfile.MemFile, n int) unsafe.Pointer {
	if n == 0 {
		return nil
	}
	arr, _ := alloc.Make([]noxrender.ImageHandle{}, n)
	for i := range arr {
		arr[i] = readThingImage(f)
	}
	return unsafe.Pointer(&arr[0])
}

func parseAnimateDraw(obj *client.ObjectType, f *binfile.MemFile, _ string, _ []byte) error { // nox_things_animate_draw_parse
	dd, _ := alloc.New(animateDrawData{})
	dd.Size = uint32(unsafe.Sizeof(animateDrawData{}))
	dd.Cnt = f.ReadU8()
	dd.Delay = f.ReadU8()
	kind, _ := f.ReadString8()
	dd.Kind = client.ParseAnimKind(kind)
	dd.Images = readThingImages(f, int(dd.Cnt))
	obj.DrawData = unsafe.Pointer(dd)
	obj.DrawFunc = legacy.Get_nox_thing_animate_draw()
	return nil
}

func parseAnimateStateDraw(obj *client.ObjectType, f *binfile.MemFile, _ string, _ []byte) error { // nox_things_animate_state_draw_parse
	const size = 0x94
	dd, _ := alloc.Make([]byte{}, size)
	binary.LittleEndian.PutUint32(dd, size)
	for {
		cmd := f.ReadU32()
		if cmd == 0x454E4420 { // "END "
			break
		}
		params := f.ReadU32()
		if params&0xE == 0 {
			return fmt.Errorf("invalid animation state params: 0x%x", params)
		}
		// TODO: the first string is always "NULL" and the second one is always empty
		f.SkipString8()
		f.SkipString8()
		ind := 0
		if params&2 != 0 {
			ind = 0
		} else if params&4 != 0 {
			ind = 1
		} else if params&8 != 0 {
			ind = 2
		}
		ani := (*client.AnimationVector)(unsafe.Pointer(&dd[4+48*ind]))
		if err := nox_xxx_loadVectorAnimated_44B8B0(ani, f); err != nil {
			return err
		}
		arr := readThingImages(f, int(int16(ani.Cnt40)))
		ani.Frames[0] = (*noxrender.ImageHandle)(arr)
	}
	obj.Field_54 = 2
	obj.DrawFunc = legacy.Get_nox_thing_animate_state_draw()
	obj.DrawData = unsafe.Pointer(&dd[0])
	return nil
}

func (c *Client) DrawAnimate(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_animate_draw
	r := c.r
	dd := (*animateDrawData)(dr.DrawData)
	frame := c.srv.Frame()
	delay := uint32(dd.Delay) + 1
	var ind int
	switch dd.Kind {
	case client.AnimOneShot:
		ind = int(int32((frame - dr.AnimStart) / delay))
		if ind >= int(dd.Cnt) {
			ind = int(dd.Cnt) - 1
		}
	case client.AnimOneShotRemove:
		ind = int(int32((frame - dr.AnimStart) / delay))
		if ind >= int(dd.Cnt) {
			c.Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr)
			return 0
		}
	case client.AnimLoop:
		loop := dr.ObjFlags.Has(object.FlagEnabled)
		if !loop && dr.ObjClass.Has(object.ClassFlag) {
			if !noxflags.HasGame(noxflags.GameFlag(0x20)) {
				return 1
			}
			loop = true
		}
		if loop {
			ind = int(int32((frame + dr.NetCode32) / delay))
			if ind >= int(dd.Cnt) {
				ind %= int(dd.Cnt)
			}
		}
	case client.AnimLoopAndFade:
		cnt := 2 * int(dd.Cnt)
		r.Data().SetAlphaEnabled(true)
		ind = int(int32((frame - dr.AnimStart) / delay))
		if ind >= cnt {
			c.Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr)
			return 0
		}
		r.Data().SetAlpha(byte(-56 - 200*ind/cnt))
		if ind >= int(dd.Cnt) {
			ind %= int(dd.Cnt)
		}
	case client.AnimRandom:
		ind = c.srv.Rand.Other.Int(0, int(dd.Cnt)-1)
	case client.AnimSlave:
		ind = int(dr.AnimFrameSlave)
	default:
		return 1
	}
	legacy.Nox_xxx_drawObject_4C4770_draw(vp, dr, dd.Frame(ind))
	if dd.Kind == client.AnimLoopAndFade {
		r.Data().SetAlphaEnabled(false)
	}
	return 1
}

var (
	arrowTailLinkType     int
	weakArrowTailLinkType int
)

func (c *Client) drawArrowWithTail(vp *noxrender.Viewport, dr *client.Drawable, tailType *int, tailName string) int {
	if *tailType == 0 {
		*tailType = c.Cli().Things.IndByID(tailName)
	}
	dx := uint32(dr.PosVec.X) - dr.Field_81
	dy := uint32(dr.PosVec.Y) - dr.Field_82
	if dx*dx+dy*dy > 200 {
		tail := c.Nox_xxx_spriteLoadAdd_45A360_drawable(*tailType, image.Pt(int(int32(dr.Field_81)), int(int32(dr.Field_82))))
		u := tail.UnionEffect()
		u.Field_108 = uint32(dr.PosVec.X)
		u.Field_109 = uint32(dr.PosVec.Y)
		c.Cli().Objs.List34Add(tail)
		dr.Field_81 = uint32(dr.PosVec.X)
		dr.Field_82 = uint32(dr.PosVec.Y)
		c.Cli().Objs.TransparentDecay(tail, int(c.srv.TickRate()/3))
	}
	return c.DrawSlave(vp, dr)
}

func (c *Client) DrawArrow(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_arrow_draw
	return c.drawArrowWithTail(vp, dr, &arrowTailLinkType, "ArrowTailLink")
}

func (c *Client) DrawWeakArrow(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_weak_arrow_draw
	return c.drawArrowWithTail(vp, dr, &weakArrowTailLinkType, "WeakArrowTailLink")
}

func (c *Client) drawArrowTailLink(vp *noxrender.Viewport, dr *client.Drawable, colors uintptr) int {
	r := c.r
	u := dr.UnionEffect()
	p := vp.ToScreenPos(dr.Pos())
	p.Y -= int(int16(dr.ZVal2)) + int(int16(dr.ZVal))
	dy := int(int32(u.Field_109 - uint32(dr.PosVec.Y)))
	x2 := int(int32(u.Field_108)) + vp.Screen.Min.X - vp.World.Min.X
	left := dr.Deadline - c.srv.Frame()
	if left == 0 {
		return 1
	}
	ci := int(int32((left << 6) / (c.srv.TickRate() / 3)))
	if ci >= 64 {
		ci = 63
	}
	r.Data().SetColor2(noxcolor.RGBA5551(memmap.Uint32(0x5D4594, colors+4*uintptr(ci))))
	r.Data().SetAlphaEnabled(true)
	r.Data().SetAlpha(0x80)
	r.AddPoint(image.Pt(p.X, p.Y-4))
	r.AddPoint(image.Pt(x2, dy+p.Y-4))
	r.DrawLineFromPoints(r.Data().Color2())
	r.Data().SetAlphaEnabled(false)
	return 1
}

func (c *Client) DrawArrowTailLink(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_arrow_tail_link_draw
	return c.drawArrowTailLink(vp, dr, 1313012)
}

func (c *Client) DrawWeakArrowTailLink(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_weak_arrow_tail_link_draw
	return c.drawArrowTailLink(vp, dr, 1313268)
}

func (c *Client) DrawHarpoon(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_harpoon_draw
	return c.DrawSlave(vp, dr)
}

func (c *Client) DrawHarpoonRope(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_harpoon_rope_draw
	u := unsafe.Slice((*byte)(unsafe.Pointer(&dr.Union)), unsafe.Sizeof(dr.Union))
	off := vp.Screen.Min.Sub(vp.World.Min)
	var p1, p2 image.Point
	if u[0] == 0 {
		p1 = image.Pt(int(binary.LittleEndian.Uint16(u[5:])), int(binary.LittleEndian.Uint16(u[7:]))-20).Add(off)
		p2 = image.Pt(int(binary.LittleEndian.Uint16(u[9:])), int(binary.LittleEndian.Uint16(u[11:]))-20).Add(off)
	} else {
		byCode := func(code uint32) *client.Drawable {
			if nox_xxx_netTestHighBit_578B70(uint16(code)) {
				return c.Cli().Objs.ByNetCodeStatic(int(code))
			}
			return c.Cli().Objs.ByNetCodeDynamic(int(code))
		}
		from := byCode(binary.LittleEndian.Uint32(u[5:]))
		to := byCode(binary.LittleEndian.Uint32(u[9:]))
		if from == nil || to == nil {
			return 1
		}
		p1 = from.Pos().Add(off)
		p2 = to.Pos().Add(off)
		p2.Y -= int(int16(to.ZVal2)) + int(int16(to.ZVal))
		p1.X += int(memmap.Int32(0x587000, 175864+8*uintptr(from.AnimDir)))
		p1.Y += int(memmap.Int32(0x587000, 175868+8*uintptr(from.AnimDir)))
		p2.Y -= 8
	}
	*memmap.PtrUint32(0x5D4594, 1312492) = noxcolor.RGB5551Color(144, 104, 64).Color32()
	*memmap.PtrUint32(0x5D4594, 1312496) = noxcolor.RGB5551Color(24, 16, 0).Color32()
	legacy.Sub_4B63B0(p1, p2)
	return 1
}

func (c *Client) DrawSpiderSpit(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_spider_spit_draw
	r := c.r
	z := int(int16(dr.ZVal))
	p1 := vp.ToScreenPos(dr.Pos())
	p1.Y -= z
	p2 := vp.ToScreenPos(image.Pt(int(int32(dr.Field_8)), int(int32(dr.Field_9))))
	p2.Y -= z
	dx := int(int32(uint32(dr.PosVec.X) - dr.Field_8))
	if dx < 0 {
		dx = -dx
	}
	dy := int(int32(uint32(dr.PosVec.Y) - dr.Field_9))
	if dy < 0 {
		dy = -dy
	}
	line := func(a, b image.Point) {
		r.AddPoint(a)
		r.AddPoint(b)
		r.DrawLineFromPoints(r.Data().Color2())
	}
	r.Data().SetColor2(noxcolor.RGBA5551(memmap.Uint32(0x85B3FC, 956)))
	if dx <= dy {
		line(p1.Add(image.Pt(1, 0)), p2.Add(image.Pt(-1, 0)))
		line(p1.Add(image.Pt(-1, 0)), p2.Add(image.Pt(-1, 0)))
	} else {
		line(p1.Add(image.Pt(0, 1)), p2.Add(image.Pt(0, 1)))
		line(p1.Add(image.Pt(0, 1)), p2.Add(image.Pt(0, -1)))
	}
	r.DrawRectFilledOpaque(p1.X-1, p1.Y-1, 4, 4, r.Data().Color2())
	r.Data().SetColor2(nox_color_white_2523948)
	line(p1, p2)
	r.DrawRectFilledOpaque(p1.X, p1.Y, 2, 2, r.Data().Color2())
	return 1
}

func init() {
	client.RegisterDraw("BubbleDraw", legacy.Get_nox_thing_bubble_draw(), 0, nil)
	client.RegisterDraw("VortexDraw", legacy.Get_nox_thing_vortex_draw(), 0, nil)
}

// effectBytes returns the drawable union as raw bytes.
func effectBytes(dr *client.Drawable) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(&dr.Union)), unsafe.Sizeof(dr.Union))
}

func (c *Client) DrawVortex(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_vortex_draw
	r := c.r
	u := effectBytes(dr)
	glowCl := binary.LittleEndian.Uint32(u[0:])
	pointCl := binary.LittleEndian.Uint32(u[4:])
	cx := binary.LittleEndian.Uint32(u[8:])
	cy := binary.LittleEndian.Uint32(u[12:])
	gray := noxcolor.RGB5551Color(170, 170, 170).Color32()

	ang := int(u[16])
	rad := int(u[18])
	p := image.Pt(
		int(int32(cx+uint32(int32(rad*int(memmap.Int32(0x587000, 192088+8*uintptr(ang)))/16)))),
		int(int32(cy+uint32(int32(rad*int(memmap.Int32(0x587000, 192092+8*uintptr(ang)))/16)))),
	)
	sp := vp.ToScreenPos(p)
	sp.Y -= int(int16(dr.ZVal))
	if sp.X <= vp.Screen.Min.X || sp.X >= vp.Screen.Max.X || sp.Y <= vp.Screen.Min.Y || sp.Y >= vp.Screen.Max.Y {
		c.Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr)
		return 0
	}
	if p.Y >= int(int32(cy)) {
		r.DrawGlow(sp, noxcolor.RGBA5551(glowCl), 3, 5)
		r.Data().SetColor2(noxcolor.RGBA5551(pointCl))
		r.DrawPoint(sp, 3, r.Data().Color2())
	} else {
		r.DrawGlow(sp, noxcolor.RGBA5551(gray), 2, 4)
		r.Data().SetColor2(noxcolor.RGBA5551(gray))
		r.DrawPoint(sp, 2, r.Data().Color2())
	}
	r.AddPoint(sp)

	speed := int(int8(u[17]))
	ang2 := int(u[16]) - 2*speed
	if speed > 0 {
		if ang2 < 0 {
			ang2 += 0x100
		}
	} else if ang2 >= 256 {
		ang2 -= 256
	}
	// The original code uses unsigned arithmetic for the second point.
	p = image.Pt(
		int(int32(cx+uint32(rad)*memmap.Uint32(0x587000, 192088+8*uintptr(ang2))/16)),
		int(int32(cy+uint32(rad)*memmap.Uint32(0x587000, 192092+8*uintptr(ang2))/16)),
	)
	sp = vp.ToScreenPos(p)
	sp.Y -= int(int16(dr.ZVal))
	r.AddPoint(sp)
	if p.Y >= int(int32(cy)) {
		r.Data().SetColor2(noxcolor.RGBA5551(pointCl))
	} else {
		r.Data().SetColor2(noxcolor.RGBA5551(gray))
	}
	r.DrawLineFromPoints(r.Data().Color2())

	u[16] += u[17]
	dr.ZVal += uint16(u[19])
	v := int(float32(float64(int16(dr.ZVal)) * 0.0024999999 * 50.0))
	if 50-v <= 0 {
		c.Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr)
		return 0
	}
	u[18] = byte(50 - v)
	return 1
}

func (c *Client) DrawBubble(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_bubble_draw
	r := c.r
	u := effectBytes(dr)
	const (
		size     = 8
		state    = 9
		cnt      = 10
		cntReset = 11
		dzReset  = 12
		dzCnt    = 13
		dz       = 14
	)
	if u[state] != 3 {
		if d := dr.Deadline; d != 0 && d <= c.srv.Frame() {
			u[state] = 3
			u[cnt] = 4
			u[cntReset] = 4
			c.Cli().Objs.TransparentDecay(dr, int(c.srv.TickRate()))
		}
	}
	if u[state] == 3 && u[size] == 0 {
		c.Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr)
		return 0
	}
	p := vp.ToScreenPos(dr.Pos())
	p.Y -= int(int16(dr.ZVal))
	r.DrawGlow(p, noxcolor.RGBA5551(binary.LittleEndian.Uint32(u[0:])), int(u[size]), int(u[size])+3)
	r.Data().SetColor2(noxcolor.RGBA5551(binary.LittleEndian.Uint32(u[4:])))
	r.DrawPoint(p, int(u[size]>>1), r.Data().Color2())
	if byte(c.srv.Frame())&3 != 0 {
		dr.ZVal += uint16(int8(u[dz]))
	}
	if u[cnt] != 0 {
		u[cnt]--
		if u[cnt] == 0 {
			switch u[state] {
			case 1:
				u[size]++
				if u[size] >= 5 {
					u[state] = 2
				}
			case 2:
				u[size]--
				if u[size] == 0 {
					u[state] = 1
				}
			default:
				u[size]--
				if u[size] == 0 {
					c.Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr)
					return 0
				}
			}
			u[cnt] = u[cntReset]
		}
	}
	if u[dzCnt] != 0 {
		u[dzCnt]--
		if u[dzCnt] == 0 {
			u[dz] = -u[dz]
			u[dzCnt] = u[dzReset]
		}
	}
	if int16(dr.ZVal) >= 0 {
		return 1
	}
	c.Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr)
	return 0
}

var blueSparkType int

func (c *Client) DrawSummonEffect(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_summon_effect_draw
	r := c.r
	dd := (*animateDrawData)(dr.DrawData)
	orig := dr.PosVec
	origSlave := dr.AnimFrameSlave
	if blueSparkType == 0 {
		blueSparkType = c.Cli().Things.IndByID("BlueSpark")
	}
	u := effectBytes(dr)
	summoned := (*client.Drawable)(unsafe.Pointer(uintptr(binary.LittleEndian.Uint32(u[0:]))))
	dur := int(binary.LittleEndian.Uint16(u[4:]))
	dt := int(int32(c.srv.Frame() - dr.AnimStart))
	if uint32(dt) >= uint32(dur) {
		c.Nox_xxx_spriteDelete_45A4B0(summoned)
		c.Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr)
		return 0
	}
	if dt >= dur-1 {
		c.makePointFx(blueSparkType, 50, 1000, 30, dr.PosVec)
	}
	c.DrawAnimate(vp, dr)
	j := 0
	for off := uintptr(192092); off < 194140; off += 80 {
		if j >= int(dd.Cnt) {
			j = 0
		}
		fi := (uint32(j) + c.srv.Frame() + dr.NetCode32) / (uint32(dd.Delay) + 1)
		dr.AnimFrameSlave = fi
		if fi >= uint32(dd.Cnt) {
			dr.AnimFrameSlave = fi % uint32(dd.Cnt)
		}
		dd.Kind = client.AnimSlave
		dr.PosVec.X = orig.X + 2*int(memmap.Int32(0x587000, off-4))
		dr.PosVec.Y = orig.Y + 2*int(memmap.Int32(0x587000, off))
		if dr.PosVec.X >= 0 && dr.PosVec.X < 5888 && dr.PosVec.Y >= 0 && dr.PosVec.Y < 5888 {
			c.DrawAnimate(vp, dr)
		}
		j++
	}
	dr.PosVec = orig
	dr.AnimFrameSlave = origSlave
	dd.Kind = client.AnimLoop
	r.Data().SetAlphaEnabled(true)
	r.Data().SetAlpha(byte(int64((float64(c.srv.Frame()) - float64(dr.AnimStart)) / float64(dur) * 255.0)))
	legacy.CallDrawFunc(summoned, vp)
	r.Data().SetAlphaEnabled(false)
	return 1
}

func init() {
	client.RegisterDraw("DrainManaDraw", legacy.Get_nox_thing_drain_mana_draw(), 0, nil)
	client.RegisterDraw("MagicDraw", legacy.Get_nox_thing_magic_draw(), 0, nil)
	client.RegisterDraw("MagicMissileDraw", legacy.Get_nox_thing_magic_missle_draw(), 0, nil)
	client.RegisterDraw("MagicTailLinkDraw", legacy.Get_nox_thing_magic_tail_link_draw(), 0, nil)
	client.RegisterDraw("MagicMissileTailLinkDraw", legacy.Get_nox_thing_magic_missle_tail_link_draw(), 0, nil)
}

func (c *Client) drawMagicSpark(vp *noxrender.Viewport, dr *client.Drawable, glow color.Color, lr, lg, lb byte) int {
	r := c.r
	p := vp.ToScreenPos(dr.Pos())
	p.Y -= int(int16(dr.ZVal2)) + int(int16(dr.ZVal))
	if p.X-10 < vp.Screen.Min.X || p.Y-10 < vp.Screen.Min.Y || p.X+10 >= vp.Screen.Max.X || p.Y+10 >= vp.Screen.Max.Y {
		return 1
	}
	sz := c.srv.Rand.Other.Int(1, 4)
	r.DrawGlow(p, glow, 2*sz+1, (sz>>1)+3)
	r.Data().SetColor2(nox_color_white_2523948)
	r.DrawRectFilledOpaque(p.X-(sz>>1), p.Y-(sz>>1), sz, sz, r.Data().Color2())
	dr.SetLightColor(lr, lg, lb)
	dr.SetLightIntensity(float32(c.srv.Rand.Other.Float(0.0, 100.0)))
	return 1
}

func (c *Client) DrawMagic(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_magic_draw
	return c.drawMagicSpark(vp, dr, noxcolor.RGB5551Color(0, 200, 255), 200, 200, 255)
}

func (c *Client) DrawMagicMissile(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_magic_missle_draw
	return c.drawMagicSpark(vp, dr, noxcolor.RGB5551Color(255, 255, 50), 255, 180, 50)
}

func (c *Client) drawMagicTailLine(vp *noxrender.Viewport, dr *client.Drawable, cl uint32) {
	r := c.r
	u := dr.UnionEffect()
	off := vp.Screen.Min.Sub(vp.World.Min)
	z := int(int16(dr.ZVal2)) + int(int16(dr.ZVal))
	p1 := dr.Pos().Add(off)
	p1.Y -= z
	p2 := image.Pt(int(int32(u.Field_108)), int(int32(u.Field_109))).Add(off)
	p2.Y -= z
	r.Data().SetColor2(noxcolor.RGBA5551(cl))
	r.Data().SetAlphaEnabled(true)
	r.Data().SetAlpha(0x80)
	r.AddPoint(p1)
	r.AddPoint(p2)
	r.DrawLineFromPoints(r.Data().Color2())
	r.Data().SetAlphaEnabled(false)
}

func (c *Client) DrawMagicMissileTailLink(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_magic_missle_tail_link_draw
	left := int(int32(dr.Deadline - c.srv.Frame()))
	if left <= 0 {
		return 1
	}
	fps := int(c.srv.TickRate() / 3)
	ci := (left << 6) / fps
	if ci >= 64 {
		ci = 63
	}
	cl := memmap.Uint32(0x5D4594, 1312756+4*uintptr(ci))
	dr.SetLightColor(255, 128, 50)
	dr.SetLightIntensity(float32(float64(left) * 20.0 / float64(fps)))
	c.drawMagicTailLine(vp, dr, cl)
	return 1
}

func (c *Client) DrawMagicTailLink(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_magic_tail_link_draw
	left := int(int32(dr.Deadline - c.srv.Frame()))
	if left < 0 || dr.Deadline == c.srv.Frame() {
		return 1
	}
	fps := c.srv.TickRate()
	ci := uint32(left<<6) / fps
	if int(int32(uint32(left<<6)/fps)) >= 64 {
		ci = 63
	}
	cl := memmap.Uint32(0x5D4594, 1312500+4*uintptr(ci))
	dr.SetLightColor(128, 128, 255)
	dr.SetLightIntensity(float32(float64(left) * 20.0 / float64(int(fps))))
	c.drawMagicTailLine(vp, dr, cl)
	return 1
}
