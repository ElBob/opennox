package opennox

import (
	"context"
	"strings"
	"unicode/utf16"
	"unsafe"

	"github.com/opennox/libs/console"
	"github.com/opennox/libs/strman"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/sound"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func init() {
	legacy.ServerHandleClientConsole = nox_xxx_serverHandleClientConsole_443E90
}

func clientConsoleStr(id string) string {
	return strings.ReplaceAll(strMan.GetStringInFile(strman.ID(id), "parsecmd.c"), "%S", "%s")
}

func clientConsolePrintf(format string, args ...any) {
	noxConsole.Print(console.ColorRed, cSprintf(format, args...))
}

// clientConsoleAllowed checks if the player can run scripts and server commands remotely.
func clientConsoleAllowed(pl *server.Player) bool {
	return legacy.Sub_4D12A0(pl.Index()) != 0 || pl.Index() == legacy.ClientPlayerNetCode() || noxflags.HasGame(noxflags.GameModeCoop)
}

// nox_xxx_serverHandleClientConsole_443E90 handles a console request from a player.
func nox_xxx_serverHandleClientConsole_443E90(pl *server.Player, typ byte, cmd string) int {
	if pl == nil || pl.PlayerUnit == nil {
		return 0
	}
	// the request is copied to the original buffer
	blobWStrSet(0x5D4594, 818228, cmd)
	if typ != 4 && typ != 5 && typ != 0 {
		if noxflags.HasGame(noxflags.GameFlag15 | noxflags.GameFlag16) {
			return 1
		}
	}
	switch typ {
	case 0: // observer mode
		if noxflags.HasGame(noxflags.GameFlag4) || noxflags.HasGame(noxflags.GameModeQuest) || pl.Field3680&1 != 0 {
			return 1
		}
		// a special code that skips the observer timestamp
		forced := memmap.Uint16(0x5D4594, 818228) == 0xF00D && memmap.Uint16(0x5D4594, 818230) == 0xF0AD && memmap.Uint16(0x5D4594, 818232) == 0
		if nox_xxx_playerGoObserver_4E6860(pl, bool2int(forced), 0) != 1 {
			return 1
		}
		on := clientConsoleStr("set")
		if legacy.Nox_xxx_gamePlayIsAnyPlayers_40A8A0() != 0 {
			if !forced {
				legacy.Nox_xxx_netNeedTimestampStatus_4174F0(pl, 256)
			}
			if noxflags.HasGame(noxflags.GameModeElimination) && legacy.Sub_40A770() == 1 {
				legacy.Sub_5095E0()
			}
		}
		if u := pl.PlayerUnit; u != nil {
			legacy.Nox_xxx_netChangeTeamMb_419570(u.TeamPtr(), uint32(pl.NetCode()))
		}
		clientConsolePrintf(clientConsoleStr("observermode"), on)
	case 1: // execute a map script function
		if !clientConsoleAllowed(pl) {
			return 1
		}
		// second space-separated token, split in the original buffer like wcstok; a NULL token prints as "(null)"
		buf := unsafe.Slice((*uint16)(memmap.PtrOff(0x5D4594, 818228)), len(utf16.Encode([]rune(cmd)))+1)
		_, next, _ := wcstokSpace(buf, 0)
		name, _, ok := wcstokSpace(buf, next)
		if !ok {
			name = "(null)"
		}
		if ind := noxServer.S().NoxScriptVM.ScriptIndexByName(name); ind != -1 {
			clientConsolePrintf(clientConsoleStr("ExecutingFunction"), name)
			if err := noxServer.S().NoxScriptVM.CallByIndex(ind, pl.PlayerUnit, pl.PlayerUnit); err != nil {
				gameLog.Println(err)
			}
		} else {
			clientConsolePrintf(clientConsoleStr("InvalidFunction"), name)
		}
	case 2: // remote server command
		if !clientConsoleAllowed(pl) {
			return 1
		}
		clientConsolePrintf(clientConsoleStr("RemoteSysop"), pl.Name(), cmd)
		execConsoleCmd(context.Background(), cmd)
	case 3:
		legacy.PrintToAllFlags(0, cmd)
	case 4: // watch another player
		if pl.Field3680&1 == 0 && !noxflags.HasEngine(noxflags.EngineReplayRead) {
			if noxflags.HasGame(noxflags.GameHost) {
				clientConsolePrintf(clientConsoleStr("notinobserver"))
			}
			return 1
		}
		if cmd == "" {
			nox_xxx_playerCameraUnlock_4E6040(pl.PlayerUnit)
			return 1
		}
		for it := noxServer.Players.First(); it != nil; it = noxServer.Players.Next(it) {
			if wcsieq(cmd, it.Name()) {
				nox_xxx_playerCameraFollow_4E6060(pl.PlayerUnit, it.PlayerUnit)
			}
		}
	case 5:
		legacy.PrintToAllFlags(16, cmd)
		for it := noxServer.Players.First(); it != nil; it = noxServer.Players.Next(it) {
			if it.PlayerUnit != nil {
				noxServer.Audio.EventObj(sound.ID(902), it.PlayerUnit, 0, 0)
			}
		}
	default:
		clientConsolePrintf(clientConsoleStr("invalidattempt"), pl.Name(), cmd)
	}
	return 1
}

// wcstokSpace splits a NUL-terminated wide string in place like wcstok(buf, L" "), starting at pos.
// It returns the token, the position to continue from and false if there are no more tokens.
func wcstokSpace(buf []uint16, pos int) (string, int, bool) {
	for pos < len(buf) && buf[pos] == ' ' {
		pos++
	}
	if pos >= len(buf) || buf[pos] == 0 {
		return "", pos, false
	}
	start := pos
	for pos < len(buf) && buf[pos] != ' ' && buf[pos] != 0 {
		pos++
	}
	tok := string(utf16.Decode(buf[start:pos]))
	if pos < len(buf) && buf[pos] == ' ' {
		buf[pos] = 0
		pos++
	}
	return tok, pos, true
}
