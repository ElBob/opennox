package opennox

import (
	"fmt"
	"strconv"
	"strings"
	"unsafe"

	"github.com/opennox/libs/strman"

	"github.com/opennox/opennox/v1/client"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/binfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func init() {
	// TODO: pass strman
	client.RegisterThingParse("PRETTYNAME", func(typ *client.ObjectType, f *binfile.MemFile, str string, buf []byte) error {
		typ.PrettyName = alloc.InternCString16(noxClient.Strings().GetStringInFile(strman.ID(str), "drawdb.c"))
		return nil
	})
	client.RegisterThingParse("DESCRIPTION", func(typ *client.ObjectType, f *binfile.MemFile, str string, buf []byte) error {
		typ.Desc = alloc.InternCString16(noxClient.Strings().GetStringInFile(strman.ID(str), "drawdb.c"))
		return nil
	})
	client.RegisterThingParse("LIGHTDIRECTION", parseThingLightDir)
	client.RegisterThingParse("LIGHTPENUMBRA", parseThingLightPenumbra)
	client.RegisterThingParse("CLIENTUPDATE", parseThingClientUpdate)
	client.RegisterThingParse("PRETTYIMAGE", parseThingPrettyImage)
}

// scanInt parses a leading decimal integer the same way as sscanf("%d").
func scanInt(s string) (int, bool) {
	s = strings.TrimLeft(s, " \t\n\v\f\r")
	i := 0
	if i < len(s) && (s[i] == '+' || s[i] == '-') {
		i++
	}
	j := i
	for j < len(s) && s[j] >= '0' && s[j] <= '9' {
		j++
	}
	if j == i {
		return 0, false
	}
	v, err := strconv.Atoi(s[:j])
	if err != nil {
		return 0, false
	}
	return v, true
}

// thingLightAngle converts degrees to the light angle (0x10000 is a full circle).
func thingLightAngle(deg int) uint16 {
	return uint16(int64(float64(deg)*memmap.Float64(0x581450, 9560)*legacy.Get_qword_581450_9552() + legacy.Get_qword_581450_9544()))
}

func parseThingLightDir(typ *client.ObjectType, f *binfile.MemFile, str string, buf []byte) error { // nox_parse_thing_light_dir
	deg, ok := scanInt(str)
	if !ok || deg < 0 || deg >= 360 {
		return fmt.Errorf("invalid light direction: %q", str)
	}
	typ.LightDir = thingLightAngle(deg)
	typ.Field_10 = 0
	return nil
}

func parseThingLightPenumbra(typ *client.ObjectType, f *binfile.MemFile, str string, buf []byte) error { // nox_parse_thing_light_penumbra
	deg, ok := scanInt(str)
	if !ok || deg < 0 || deg >= 180 {
		return fmt.Errorf("invalid light penumbra: %q", str)
	}
	typ.LightPenumbra = thingLightAngle(deg)
	return nil
}

func parseThingClientUpdate(typ *client.ObjectType, f *binfile.MemFile, str string, buf []byte) error { // nox_parse_thing_client_update
	// first token, as strtok(str, " \t\n\r")
	fields := strings.FieldsFunc(str, func(r rune) bool { return r == ' ' || r == '\t' || r == '\n' || r == '\r' })
	if len(fields) == 0 {
		return fmt.Errorf("invalid client update: %q", str)
	}
	name := fields[0]
	// table of {name, update func} pairs
	for off := uintptr(175072); ; off += 8 {
		p := *memmap.PtrPtr(0x587000, off)
		if p == nil {
			return fmt.Errorf("unknown client update: %q", name)
		}
		if alloc.GoString((*byte)(p)) == name {
			typ.ClientUpdate = *memmap.PtrPtr(0x587000, off+4)
			return nil
		}
	}
}

func parseThingPrettyImage(typ *client.ObjectType, f *binfile.MemFile, str string, buf []byte) error { // nox_parse_thing_pretty_image
	typ.PrettyImage = uint32(uintptr(unsafe.Pointer(readThingImage(f))))
	return nil
}

func (c *Client) Nox_things_free_44C580() {
	c.Things.Each(func(it *client.ObjectType) {
		if it.DrawData != nil {
			freeDrawableData(it.DrawData, it.DrawFunc)
		}
	})
	c.Things.Nox_things_free_44C580_B()
	if !noxflags.HasGame(noxflags.GameHost) {
		nox_xxx_free_42BF80()
	}
}

func freeDrawableData(ptr unsafe.Pointer, fnc unsafe.Pointer) {
	legacy.Nox_xxx_draw_44C650_free_kind(ptr, client.DrawableDataKind(fnc))
}
