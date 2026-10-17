package opennox

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf16"
	"unsafe"

	"github.com/opennox/libs/console"
	"github.com/opennox/libs/spell"
	"github.com/opennox/libs/strman"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/sound"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

// Console commands ported from client__system__parsecmd.c.
// They keep the legacy calling convention: tokens include the command words, tokInd points to the first argument.

// cmdStr loads a string for parsecmd.c and converts it to a Go format.
func cmdStr(c *console.Console, id string) string {
	return strings.ReplaceAll(c.Strings().GetStringInFile(strman.ID(id), "parsecmd.c"), "%S", "%s")
}

// cSprintf formats like C printf: arguments without a matching verb are ignored.
func cSprintf(format string, args ...any) string {
	n := 0
	for i := 0; i < len(format); i++ {
		if format[i] != '%' {
			continue
		}
		if i+1 < len(format) && format[i+1] == '%' {
			i++
			continue
		}
		n++
	}
	if n < len(args) {
		args = args[:n]
	}
	return fmt.Sprintf(format, args...)
}

func cmdPrintf(c *console.Console, format string, args ...any) {
	c.Print(console.ColorRed, cSprintf(format, args...))
}

// blobWStrSet writes a NUL-terminated wide string to the static data, like wcscpy/vswprintf into the original buffers.
func blobWStrSet(blob, off uintptr, s string) {
	u := utf16.Encode([]rune(s))
	dst := unsafe.Slice((*uint16)(memmap.PtrOff(blob, off)), len(u)+1)
	copy(dst, u)
	dst[len(u)] = 0
}

// blobWStrCat appends a wide string to a NUL-terminated wide string in the static data, like wcscat.
func blobWStrCat(blob, off uintptr, s string) {
	cur := blobWStr(blob, off)
	blobWStrSet(blob, off+2*uintptr(len(utf16.Encode([]rune(cur)))), s)
}

func blobWStr(blob, off uintptr) string {
	return alloc.GoString16((*uint16)(memmap.PtrOff(blob, off)))
}

// cmdPrintCentered prints a message in the center of the screen (sub_440A20).
func cmdPrintCentered(format string, args ...any) {
	blobWStrSet(0x5D4594, 822660, cSprintf(format, args...))
	legacy.Nox_xxx_printCentered_445490(blobWStr(0x5D4594, 822660))
}

// wcsieq compares strings case-insensitively, as _nox_wcsicmp(a, b) == 0.
func wcsieq(a, b string) bool {
	return strings.ToLower(a) == strings.ToLower(b)
}

// wcstol parses a leading base 10 integer, as wcstol(s, 0, 10).
func wcstol(s string) int {
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
		return 0
	}
	v, err := strconv.ParseInt(s[:j], 10, 32)
	if err != nil {
		if s[0] == '-' {
			return -1 << 31
		}
		return 1<<31 - 1
	}
	return int(v)
}

// cmdTokenName returns a command token from the static table at 0x587000+94468 (6 = "respawn", 7 = "all").
func cmdTokenName(i int) string {
	return alloc.GoString16((*uint16)(*memmap.PtrPtr(0x587000, 94468+4*uintptr(i))))
}

// playerByName finds a player by name (nox_xxx_playerByName_4170D0).
func playerByName(name string) *server.Player {
	for pl := noxServer.Players.First(); pl != nil; pl = noxServer.Players.Next(pl) {
		if wcsieq(pl.Name(), name) {
			return pl
		}
	}
	return nil
}

func noxCmdSetSysop(ctx context.Context, c *console.Console, tokInd int, tokens []string) bool { // nox_cmd_set_sysop
	if len(tokens) != 3 {
		return false
	}
	legacy.Nox_xxx_sysopSetPass_40A610(tokens[tokInd])
	cmdPrintf(c, cmdStr(c, "sysoppasswordset"))
	return true
}

func noxCmdSetCycle(ctx context.Context, c *console.Console, tokInd int, tokens []string) bool { // nox_cmd_set_cycle
	if len(tokens) != 3 {
		return false
	}
	switch {
	case wcsieq(tokens[tokInd], "on"):
		legacy.Sub_4D0D90(1)
		cmdPrintf(c, cmdStr(c, "MapCycleOn"))
	case wcsieq(tokens[tokInd], "off"):
		legacy.Sub_4D0D90(0)
		cmdPrintf(c, cmdStr(c, "MapCycleOff"))
	default:
		return false
	}
	legacy.Sub_4AD840()
	return true
}

// cmdSetRule switches a game rule flag on or off and prints the result.
func cmdSetRule(c *console.Console, arg string, flag int, name string) bool {
	switch {
	case wcsieq(arg, "on"):
		legacy.Sub_409E70(flag)
		legacy.Nox_server_gameSettingsUpdated_40A670()
		v := cmdStr(c, "cmd_token:on")
		cmdPrintf(c, cmdStr(c, name), v)
		return true
	case wcsieq(arg, "off"):
		legacy.Sub_409EC0(flag)
		legacy.Nox_server_gameSettingsUpdated_40A670()
		v := cmdStr(c, "cmd_token:off")
		cmdPrintf(c, cmdStr(c, name), v)
		return true
	}
	return false
}

func noxCmdSetWeapons(ctx context.Context, c *console.Console, tokInd int, tokens []string) bool { // nox_cmd_set_weapons
	if len(tokens) != 3 {
		return false
	}
	return cmdSetRule(c, tokens[tokInd], 1, "weapons")
}

func noxCmdSetStaffs(ctx context.Context, c *console.Console, tokInd int, tokens []string) bool { // nox_cmd_set_staffs
	if len(tokens) != 3 {
		return false
	}
	return cmdSetRule(c, tokens[tokInd], 16, "staffs")
}

func noxCmdSetName(ctx context.Context, c *console.Console, tokInd int, tokens []string) bool { // nox_cmd_set_name
	if len(tokens) < 3 {
		return false
	}
	// separator is a short C string in the static data
	sep := memmap.Uint16(0x587000, 104484)
	sepStr := ""
	if b := byte(sep); b != 0 {
		sepStr = string(b)
		if b2 := byte(sep >> 8); b2 != 0 {
			sepStr += string(b2)
		}
	}
	var name strings.Builder
	for i := tokInd; i < len(tokens); i++ {
		name.WriteString(tokens[i])
		if i+1 < len(tokens) {
			name.WriteString(sepStr)
		}
	}
	if s := name.String(); s != "" {
		legacy.Nox_xxx_gameSetServername_40A440(s)
		cmdPrintf(c, cmdStr(c, "setgamename"), s)
	}
	return true
}

func noxCmdSetMnstrs(ctx context.Context, c *console.Console, tokInd int, tokens []string) bool { // nox_cmd_set_mnstrs
	if len(tokens) != 3 && len(tokens) != 4 {
		return false
	}
	if !wcsieq(tokens[2], cmdTokenName(6)) {
		return cmdSetRule(c, tokens[tokInd], 4, "monsters")
	}
	if len(tokens) == 3 {
		return false
	}
	switch {
	case wcsieq(tokens[3], "on"):
		legacy.Sub_409E70(8)
		legacy.Nox_server_gameSettingsUpdated_40A670()
		v := cmdStr(c, "cmd_token:on")
		cmdPrintf(c, cmdStr(c, "monsterrespawn"), v)
		return true
	case wcsieq(tokens[3], "off"):
		// the original sets the flag here as well
		legacy.Sub_409E70(8)
		legacy.Nox_server_gameSettingsUpdated_40A670()
		v := cmdStr(c, "cmd_token:off")
		cmdPrintf(c, cmdStr(c, "monsterrespawn"), v)
		return true
	}
	return false
}

func noxCmdSetSpell(ctx context.Context, c *console.Console, tokInd int, tokens []string) bool { // nox_cmd_set_spell
	if len(tokens) != 4 {
		return false
	}
	s := noxServer
	if noxflags.HasGame(noxflags.GameModeChat) {
		cmdPrintf(c, cmdStr(c, "NotInChat"), tokens[2])
		return true
	}
	id := s.Spells.ByTitle(tokens[2])
	if id == 0 {
		if v := spell.ParseID(tokens[2]); v > 0 {
			id = v
		}
	}
	if id == 0 {
		cmdPrintf(c, cmdStr(c, "invalidspell"), tokens[2])
		return false
	}
	switch {
	case wcsieq(tokens[3], "on"):
		sett := getCurrentSettings2()
		if (noxflags.HasGame(noxflags.GameModeFlagBall) || byte(sett.Field52)&0x40 != 0) && id == 132 {
			return true
		}
		if s.Spells.DefByInd(id).IsEnabled() {
			return true
		}
		s.Spells.Enable(id, true)
		legacy.Nox_server_gameSettingsUpdated_40A670()
		cmdPrintf(c, cmdStr(c, "spellenabled"), tokens[2])
		return true
	case wcsieq(tokens[3], "off"):
		if !s.Spells.DefByInd(id).IsEnabled() {
			return true
		}
		s.Spells.Enable(id, false)
		legacy.Nox_server_gameSettingsUpdated_40A670()
		cmdPrintf(c, cmdStr(c, "spelldisabled"), tokens[2])
		return true
	}
	return false
}

func noxCmdBan(ctx context.Context, c *console.Console, tokInd int, tokens []string) bool { // nox_cmd_ban
	if len(tokens) != 2 {
		return false
	}
	name := tokens[tokInd]
	pl := playerByName(name)
	if pl == nil {
		legacy.Sub_416770(0, name, nil)
		cmdPrintCentered(cmdStr(c, "banDisallow"), name)
		return true
	}
	if pl.PlayerIndex() == server.HostPlayerIndex {
		cmdPrintCentered(cmdStr(c, "cantbanyourself"))
		return true
	}
	if noxflags.HasGame(noxflags.GameModeQuest) {
		sub_4DCFB0(pl.PlayerUnit)
	} else {
		noxServer.PlayerDisconnectByIndCode4(pl.PlayerIndex())
	}
	legacy.Sub_416770(0, name, &pl.SerialBuf[0])
	cmdPrintCentered(cmdStr(c, "banned"), pl.Name())
	return true
}

func noxCmdNotImplemented(ctx context.Context, c *console.Console, tokInd int, tokens []string) bool { // nox_cmd_allow_user, nox_cmd_allow_ip, nox_cmd_set_spellpts
	cmdPrintf(c, cmdStr(c, "notyetimplemented"))
	return true
}

func noxCmdKick(ctx context.Context, c *console.Console, tokInd int, tokens []string) bool { // nox_cmd_kick
	if len(tokens) != 2 {
		return false
	}
	pl := playerByName(tokens[tokInd])
	if pl == nil {
		return true
	}
	if pl.PlayerIndex() == server.HostPlayerIndex {
		cmdPrintCentered(cmdStr(c, "cantkickyourself"))
		return true
	}
	if noxflags.HasGame(noxflags.GameModeQuest) {
		sub_4DCFB0(pl.PlayerUnit)
	} else {
		noxServer.PlayerDisconnectByInd(pl.PlayerIndex(), 4)
		cmdPrintCentered(cmdStr(c, "kicked"), pl.Name())
	}
	return true
}

func noxCmdSetPlayers(ctx context.Context, c *console.Console, tokInd int, tokens []string) bool { // nox_cmd_set_players
	if len(tokens) != 3 {
		return false
	}
	n := min(max(wcstol(tokens[2]), 0), 999)
	if legacy.Nox_xxx_servGetPlrLimit_409FA0() == n {
		return true
	}
	legacy.Nox_xxx_servSetPlrLimit_409F80(n)
	legacy.Sub_455800()
	cmdPrintCentered(cmdStr(c, "playersset"), n)
	return true
}

func noxCmdListUsers(ctx context.Context, c *console.Console, tokInd int, tokens []string) bool { // nox_cmd_list_users
	isServer := !console.IsClient(ctx)
	cmdPrintf(c, cmdStr(c, "userslist"))
	format := strings.ReplaceAll(alloc.GoString16((*uint16)(memmap.PtrOff(0x587000, 106604))), "%S", "%s")
	for pl := noxServer.Players.First(); pl != nil; pl = noxServer.Players.Next(pl) {
		s := pl.Name()
		if isServer && pl.Field3680&4 != 0 {
			s += ", " + cmdStr(c, "SysMuted")
		}
		if pl.Field3680&8 != 0 {
			s += ", " + cmdStr(c, "ClientMuted")
		}
		cmdPrintf(c, format, s)
	}
	return true
}

// cmdMute mutes or unmutes a player: client mute (8) or system mute (4).
func cmdMute(ctx context.Context, c *console.Console, tokInd int, tokens []string, mute bool) bool {
	if len(tokens) < 2 || len(tokens) > 3 {
		return false
	}
	idx := tokInd
	system := false
	if !console.IsClient(ctx) && wcsieq(tokens[tokInd], cmdTokenName(7)) {
		if tokInd+1 != len(tokens)-1 {
			return false
		}
		idx = tokInd + 1
		system = true
	}
	name := tokens[idx]
	ok := false
	if pl := playerByName(name); pl != nil {
		switch {
		case system && mute: // sub_57A0C0
			legacy.Nox_xxx_netNeedTimestampStatus_4174F0(pl, 4)
			ok = true
		case system: // sub_57A130
			legacy.Nox_xxx_playerUnsetStatus_417530(pl, 4)
			ok = true
		case !noxflags.HasGame(noxflags.GameClient):
		case mute: // sub_57A080
			if pl.PlayerIndex() != server.HostPlayerIndex {
				legacy.Nox_xxx_netNeedTimestampStatus_4174F0(pl, 8)
				ok = true
			}
		default: // sub_57A0F0
			legacy.Nox_xxx_playerUnsetStatus_417530(pl, 8)
			ok = true
		}
	}
	var str string
	switch {
	case !ok:
		str = cmdStr(c, "UserNotFound")
	case mute:
		str = cmdStr(c, "Muted")
	default:
		str = cmdStr(c, "UnMuted")
	}
	cmdPrintf(c, str, name)
	return true
}

func noxCmdMute(ctx context.Context, c *console.Console, tokInd int, tokens []string) bool { // nox_cmd_mute
	return cmdMute(ctx, c, tokInd, tokens, true)
}

func noxCmdUnmute(ctx context.Context, c *console.Console, tokInd int, tokens []string) bool { // nox_cmd_unmute
	return cmdMute(ctx, c, tokInd, tokens, false)
}

func noxCmdOffOnly1(ctx context.Context, c *console.Console, tokInd int, tokens []string) bool { // nox_cmd_offonly1
	if len(tokens) != 2 {
		return false
	}
	legacy.Nox_xxx_wndGuiTeamCreate_4185B0()
	return true
}

func noxCmdOffOnly2(ctx context.Context, c *console.Console, tokInd int, tokens []string) bool { // nox_cmd_offonly2
	if len(tokens) != 3 {
		return false
	}
	sett := getSettings2ByInd(1)
	// table of {name, mode flags} pairs
	for off := uintptr(94400); ; off += 8 {
		p := *memmap.PtrPtr(0x587000, off)
		if p == nil {
			break
		}
		if wcsieq(alloc.GoString16((*uint16)(p)), tokens[tokInd]) {
			sett.Field52 &= 0xE80F
			sett.Field52 |= memmap.Uint16(0x587000, off+4)
			break
		}
	}
	return true
}

func noxCmdSetNetDebug(ctx context.Context, c *console.Console, tokInd int, tokens []string) bool { // nox_cmd_set_net_debug
	if len(tokens) != 2 {
		return false
	}
	noxflags.SetEngine(noxflags.EngineNetDebug)
	return true
}

func noxCmdUnsetNetDebug(ctx context.Context, c *console.Console, tokInd int, tokens []string) bool { // nox_cmd_unset_net_debug
	if len(tokens) != 2 {
		return false
	}
	noxflags.UnsetEngine(noxflags.EngineNetDebug)
	return true
}

func noxCmdShowInfo(ctx context.Context, c *console.Console, tokInd int, tokens []string) bool { // nox_cmd_show_info
	legacy.Nox_xxx_clientPlaySoundSpecial_452D80(sound.SoundShellClick, 100)
	legacy.Sub_435F60()
	return true
}

func noxCmdShowMem(ctx context.Context, c *console.Console, tokInd int, tokens []string) bool { // nox_cmd_show_mem
	legacy.Nox_server_currentMapGetFilename_409B30()
	legacy.Nox_xxx_gameLoopMemDump_413E30()
	return true
}

func noxCmdShowRank(ctx context.Context, c *console.Console, tokInd int, tokens []string) bool { // nox_cmd_show_rank
	if noxflags.HasGame(noxflags.GameOnline) {
		sub_4703F0()
	}
	return true
}

func noxCmdShowMotd(ctx context.Context, c *console.Console, tokInd int, tokens []string) bool { // nox_cmd_show_motd
	if len(tokens) != 2 {
		return false
	}
	legacy.CmdShowMotd()
	return true
}

func noxCmdShowSeq(ctx context.Context, c *console.Console, tokInd int, tokens []string) bool { // nox_cmd_show_seq
	if len(tokens) != 2 {
		return false
	}
	legacy.Sub_48D7B0()
	return true
}

func noxCmdListMaps(ctx context.Context, c *console.Console, tokInd int, tokens []string) bool { // nox_cmd_list_maps
	// the line is built in the original buffer at 0x5D4594+822404
	clear(memmap.Slice(0x5D4594, 822404)[:0x100])
	// map file extension is a short C string in the static data
	ext := alloc.GoStringS(memmap.Slice(0x587000, 103292)[:5])
	fmtLine := strings.ReplaceAll(blobWStr(0x587000, 103276), "%S", "%s")
	fmtLast := strings.ReplaceAll(blobWStr(0x587000, 103284), "%S", "%s")
	i := 1
	for it := legacy.Nox_common_maplist_first_4D09B0(); it != nil; it = legacy.Nox_common_maplist_next_4D09C0(it) {
		// sub_4417E0: the original formats it with "%-20.20S\t\t", but our C printf ignored width and precision for strings
		blobWStrCat(0x5D4594, 822404, alloc.GoStringS(it.Name[:])+ext+"\t\t")
		if i%4 == 0 {
			cmdPrintf(c, fmtLine, blobWStr(0x5D4594, 822404))
			*memmap.PtrUint16(0x5D4594, 822404) = 0
		}
		i++
	}
	if (i-1)%4 != 0 {
		cmdPrintf(c, fmtLast, blobWStr(0x5D4594, 822404))
	}
	return true
}

func noxCmdCheatAbility(ctx context.Context, c *console.Console, tokInd int, tokens []string) bool { // nox_cmd_cheat_ability
	if !noxflags.HasGame(noxflags.GameOnline) {
		for pl := noxServer.Players.First(); pl != nil; pl = noxServer.Players.Next(pl) {
			if pl.PlayerUnit != nil {
				nox_xxx_playerCancelAbils_4FC180(pl.PlayerUnit)
			}
		}
	}
	return true
}

func noxCmdCheatLevel(ctx context.Context, c *console.Console, tokInd int, tokens []string) bool { // nox_cmd_cheat_level
	if noxflags.HasGame(noxflags.GameOnline) {
		return true
	}
	if len(tokens) < 3 {
		return false
	}
	for pl := noxServer.Players.First(); pl != nil; pl = noxServer.Players.Next(pl) {
		if pl.PlayerUnit != nil {
			legacy.Sub_4EF410(pl.PlayerUnit, byte(wcstol(tokens[2])))
		}
	}
	return true
}

func noxCmdWindow(ctx context.Context, c *console.Console, tokInd int, tokens []string) bool { // nox_cmd_window
	if len(tokens) > 1 {
		v := tokens[tokInd]
		if v == "" || (v[0] != '+' && v[0] != '-') {
			nox_draw_setCutSize_476700(wcstol(v), 0)
		} else {
			nox_draw_setCutSize_476700(0, wcstol(v))
		}
	}
	return true
}

// cmdSetQuality sets the network connection quality (nox_cmd_set_qual_*).
func cmdSetQuality(rate int, lan bool) console.CommandLegacyFunc {
	return func(ctx context.Context, c *console.Console, tokInd int, tokens []string) bool {
		if lan {
			legacy.Set_dword_5d4594_2650652(0)
		} else {
			legacy.Set_dword_5d4594_2650652(1)
		}
		legacy.Nox_xxx_rateUpdate_40A6D0(legacy.CmdRateForQuality(rate))
		legacy.Set_nox_server_connectionType_3596(rate)
		legacy.Sub_4AD840()
		return true
	}
}

func noxCmdSetTime(ctx context.Context, c *console.Console, tokInd int, tokens []string) bool { // nox_cmd_set_time
	if len(tokens) != 3 {
		return false
	}
	sett := getCurrentSettings2()
	noxServer.Sub40A040settings(int(int16(sett.Field52)), int(byte(wcstol(tokens[2]))))
	return true
}

func noxCmdSetLessons(ctx context.Context, c *console.Console, tokInd int, tokens []string) bool { // nox_cmd_set_lessons
	if len(tokens) != 3 {
		return false
	}
	sett := getCurrentSettings2()
	legacy.Sub_409FB0_settings(sett.Field52, uint16(wcstol(tokens[2])))
	return true
}

func noxCmdMenuOptions(ctx context.Context, c *console.Console, tokInd int, tokens []string) bool { // nox_cmd_menu_options
	if !noxflags.HasGame(noxflags.GameFlag4) && noxflags.HasGame(noxflags.GameOnline) {
		legacy.Nox_xxx_guiServerOptsLoad_457500()
	}
	return true
}

func noxCmdMenuVidopt(ctx context.Context, c *console.Console, tokInd int, tokens []string) bool { // nox_cmd_menu_vidopt
	legacy.Sub_4ADA40()
	return true
}

func noxCmdReenter(ctx context.Context, c *console.Console, tokInd int, tokens []string) bool { // nox_cmd_reenter
	if !noxflags.HasGame(noxflags.GameOnline) {
		legacy.Sub_40AA60(1)
	}
	return true
}
