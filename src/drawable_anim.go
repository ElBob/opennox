package opennox

import (
	"encoding/binary"
	"fmt"
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
		n := int(int16(ani.Cnt40))
		if n < 0 {
			// calloc fails on a negative count in the original
			return fmt.Errorf("invalid animation frame count: %d", n)
		}
		arr := readThingImages(f, n)
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

func (c *Client) drawArrowWithTail(vp *noxrender.Viewport, dr *client.Drawable, tailType *uint32, tailName string) int {
	typ := int(*tailType)
	if *tailType == 0 {
		typ = c.Cli().Things.IndByID(tailName)
		*tailType = uint32(typ)
	}
	dx := uint32(dr.PosVec.X) - dr.Field_81
	dy := uint32(dr.PosVec.Y) - dr.Field_82
	if dx*dx+dy*dy > 200 {
		tail := c.Nox_xxx_spriteLoadAdd_45A360_drawable(typ, image.Pt(int(int32(dr.Field_81)), int(int32(dr.Field_82))))
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
	return c.drawArrowWithTail(vp, dr, memmap.PtrUint32(0x5D4594, 1313720), "ArrowTailLink")
}

func (c *Client) DrawWeakArrow(vp *noxrender.Viewport, dr *client.Drawable) int { // nox_thing_weak_arrow_draw
	return c.drawArrowWithTail(vp, dr, memmap.PtrUint32(0x5D4594, 1313724), "WeakArrowTailLink")
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
