package opennox

import (
	"context"

	"github.com/opennox/libs/env"
	"github.com/opennox/libs/log"
	"github.com/opennox/libs/strman"

	"github.com/opennox/libs/console"

	noxflags "github.com/opennox/opennox/v1/common/flags"
)

var (
	consoleMux = console.NewMultiPrinter(
		console.LogPrinter(log.New("console")),
	)
	noxConsole = console.New(consoleMux)
)

func initConsole(sm *strman.StringManager) {
	if env.IsDevMode() {
		noxConsole.SetCheats(true)
	}
	noxConsole.Localize(sm, "on", "off", "force", "ctf", "coop", "team", "respawn", "all")
}

func init() {
	noxConsole.SetExec(execConsoleCmd)
	for _, c := range noxCommands {
		noxConsole.Register(c)
	}
}

var (
	noxCmdList = &console.Command{Token: "list", HelpID: "listhelp", Flags: console.ClientServer, Sub: []*console.Command{
		{Token: "maps", HelpID: "listmapshelp", Flags: console.ClientServer, LegacyFunc: noxCmdListMaps},
		{Token: "users", HelpID: "listusershelp", Flags: console.ClientServer, LegacyFunc: noxCmdListUsers},
	}}
	noxCmdSet = &console.Command{Token: "set", HelpID: "sethelp", Flags: console.ClientServer, Sub: []*console.Command{
		{Token: "cycle", HelpID: "setcyclehelp", Flags: console.Server, LegacyFunc: noxCmdSetCycle},
		{Token: "frameratelimiter", HelpID: "setfrhelp", Flags: console.Server | console.Cheat, Func: func(ctx context.Context, c *console.Console, tokens []string) bool {
			if len(tokens) != 0 {
				return false
			}
			setEnableFrameLimit(true)
			return true
		}},
		{Token: "lessons", HelpID: "setlessonshelp", Flags: console.Server, LegacyFunc: noxCmdSetLessons},
		{Token: "monsters", HelpID: "setmnstrshelp", Flags: console.Server | console.Cheat, LegacyFunc: noxCmdSetMnstrs},
		{Token: "name", HelpID: "setnamehelp", Flags: console.Server, LegacyFunc: noxCmdSetName},
		{Token: "netdebug", HelpID: "setnetdebughelp", Flags: console.ClientServer | console.Cheat, LegacyFunc: noxCmdSetNetDebug},
		{Token: "ob", HelpID: "setobshelp", Flags: console.ClientServer | console.Cheat, Func: func(ctx context.Context, c *console.Console, tokens []string) bool {
			cur := console.CurCommand(ctx)
			s := c.Strings().GetStringInFile("processingobs", "parsecmd.c")
			c.Print(console.ColorRed, s)
			if noxflags.HasGame(noxflags.GameHost) {
				nox_xxx_serverHandleClientConsole_443E90(clientPlayer(), 0, cur)
			} else {
				nox_xxx_netServerCmd_440950(0, cur)
			}
			return true
		}},
		{Token: "players", HelpID: "setplayershelp", Flags: console.Server, LegacyFunc: noxCmdSetPlayers},
		{Token: "quality", HelpID: "setqualityhelp", Flags: console.ClientServer, Sub: []*console.Command{
			{Token: "modem", HelpID: "setqualmodemhelp", Flags: console.ClientServer, LegacyFunc: cmdSetQuality(4, false)},
			{Token: "isdn", HelpID: "setqualisdnhelp", Flags: console.ClientServer, LegacyFunc: cmdSetQuality(3, false)},
			{Token: "cable", HelpID: "setqualcablehelp", Flags: console.ClientServer, LegacyFunc: cmdSetQuality(2, false)},
			{Token: "T1", HelpID: "setqualT1help", Flags: console.ClientServer, LegacyFunc: cmdSetQuality(1, false)},
			{Token: "LAN", HelpID: "setqualLANhelp", Flags: console.ClientServer, LegacyFunc: cmdSetQuality(1, true)},
		}},
		{Token: "savedebugcmd", HelpID: "setsavedebughelp", Flags: console.Server, Func: func(ctx context.Context, c *console.Console, tokens []string) bool {
			noxflags.SetEngine(noxflags.EngineSaveDebug)
			s := c.Strings().GetStringInFile("savedebugset", "parsecmd.c")
			c.Print(console.ColorRed, s)
			return true
		}},
		{Token: "spell", HelpID: "setspellhelp", Flags: console.Server, LegacyFunc: noxCmdSetSpell},
		{Token: "spellpoints", HelpID: "setspellptshelp", Flags: console.Server | console.Cheat, LegacyFunc: noxCmdNotImplemented},
		{Token: "staffs", HelpID: "setstaffshelp", Flags: console.Server | console.Cheat, LegacyFunc: noxCmdSetStaffs},
		{Token: "sysop", HelpID: "setsysophelp", Flags: console.Server, LegacyFunc: noxCmdSetSysop},
		{Token: "time", HelpID: "settimehelp", Flags: console.Server, LegacyFunc: noxCmdSetTime},
		{Token: "weapons", HelpID: "setweaponshelp", Flags: console.Server | console.Cheat, LegacyFunc: noxCmdSetWeapons},
		{Token: "team", HelpID: "officialonly", Flags: console.Server | console.Cheat | console.FlagDedicated, LegacyFunc: noxCmdOffOnly1},
		{Token: "mode", HelpID: "officialonly", Flags: console.Server | console.Cheat | console.FlagDedicated, LegacyFunc: noxCmdOffOnly2},
	}}
	noxCmdUnSet = &console.Command{Token: "unset", HelpID: "unsethelp", Flags: console.ClientServer, Sub: []*console.Command{
		{Token: "frameratelimiter", HelpID: "unsetfrhelp", Flags: console.Server, Func: func(ctx context.Context, c *console.Console, tokens []string) bool {
			if len(tokens) != 0 {
				return false
			}
			setEnableFrameLimit(false)
			return true
		}},
		{Token: "netdebug", HelpID: "unsetnetdebug", Flags: console.ClientServer, LegacyFunc: noxCmdUnsetNetDebug},
	}}
	noxCmdShow = &console.Command{Token: "show", HelpID: "showhelp", Flags: console.ClientServer, Sub: []*console.Command{
		{
			Token:  "bindings",
			HelpID: "showbindingshelp",
			Flags:  console.ClientServer,
			Func: func(ctx context.Context, c *console.Console, tokens []string) bool {
				if len(tokens) != 0 {
					return false
				}
				return c.Exec(ctx, "bindings")
			},
		},
		{Token: "motd", HelpID: "showmotdhelp", Flags: console.ClientServer, LegacyFunc: noxCmdShowMotd},
		{Token: "rank", HelpID: "showrankhelp", Flags: console.ClientServer, LegacyFunc: noxCmdShowRank},
		{Token: "ai", HelpID: "showaihelp", Flags: console.Server | console.Cheat, Func: func(ctx context.Context, c *console.Console, tokens []string) bool {
			noxflags.ToggleEngine(noxflags.EngineShowAI)
			return true
		}},
		{Token: "info", HelpID: "showinfohelp", Flags: console.ClientServer | console.NoHelp, LegacyFunc: noxCmdShowInfo},
		{Token: "mem", HelpID: "showmemhelp", Flags: console.ClientServer | console.Cheat, LegacyFunc: noxCmdShowMem},
		{Token: "netstat", HelpID: "shownetstathelp", Flags: console.ClientServer, Func: func(ctx context.Context, c *console.Console, tokens []string) bool {
			noxClient.guiFPS.sub_470A60()
			return true
		}},
		{Token: "seq", HelpID: "showseqhelp", Flags: console.ClientServer | console.Cheat, LegacyFunc: noxCmdShowSeq},
	}}
	noxCmdCheat = &console.Command{Token: "cheat", HelpID: "cheathelp", Flags: console.Server | console.Cheat, Sub: []*console.Command{
		{Token: "ability", HelpID: "cheatabilityhelp", Flags: console.Server, LegacyFunc: noxCmdCheatAbility},
		{Token: "level", HelpID: "cheatlevelhelp", Flags: console.Server, LegacyFunc: noxCmdCheatLevel},
		{Token: "re-enter", HelpID: "", Flags: console.Server, LegacyFunc: noxCmdReenter},
	}}
	noxCommands = []*console.Command{
		noxCmdList,
		noxCmdSet,
		noxCmdUnSet,
		noxCmdShow,
		noxCmdCheat,
		{Token: "allow", HelpID: "allowhelp", Flags: console.Server, Sub: []*console.Command{
			{Token: "user", HelpID: "allowuserhelp", Flags: console.Server, LegacyFunc: noxCmdNotImplemented},
			{Token: "IP", HelpID: "allowiphelp", Flags: console.Server, LegacyFunc: noxCmdNotImplemented},
		}},
		{Token: "audtest", HelpID: "sethelp", Flags: console.ClientServer, Func: func(ctx context.Context, c *console.Console, tokens []string) bool {
			return true
		}},
		{Token: "ban", HelpID: "banhelp", Flags: console.Server, LegacyFunc: noxCmdBan},
		{Token: "exit", HelpID: "exithelp", Flags: console.ClientServer, Func: func(ctx context.Context, c *console.Console, tokens []string) bool {
			if noxflags.HasGame(noxflags.GameModeQuest) && noxflags.HasGame(noxflags.GameHost) {
				noxServer.sub_4D6B10(false)
			}
			if noxflags.HasGame(noxflags.GameFlag26) {
				nox_client_quit_4460C0()
			}
			nox_game_exit_xxx2()
			return true
		}},
		{Token: "kick", HelpID: "kickhelp", Flags: console.Server, LegacyFunc: noxCmdKick},
		{Token: "log", HelpID: "loghelp", Flags: console.ClientServer | console.Cheat, Sub: []*console.Command{
			{Token: "console", HelpID: "logconsolehelp", Flags: console.ClientServer, Func: func(ctx context.Context, c *console.Console, tokens []string) bool {
				if len(tokens) != 0 {
					return false
				}
				noxflags.SetEngine(noxflags.EngineLogToConsole)
				return true
			}},
			{Token: "file", HelpID: "logfilehelp", Flags: console.ClientServer, Func: func(ctx context.Context, c *console.Console, tokens []string) bool {
				return true
			}},
			{Token: "stop", HelpID: "logstophelp", Flags: console.ClientServer, Func: func(ctx context.Context, c *console.Console, tokens []string) bool {
				return true
			}},
		}},
		{Token: "menu", HelpID: "menuhelp", Flags: console.ClientServer, Sub: []*console.Command{
			{Token: "vidopt", HelpID: "menuvidopthelp", Flags: console.ClientServer, LegacyFunc: noxCmdMenuVidopt},
			{Token: "options", HelpID: "menuoptionshelp", Flags: console.ClientServer, LegacyFunc: noxCmdMenuOptions},
		}},
		{Token: "mute", HelpID: "mutehelp", Flags: console.ClientServer, LegacyFunc: noxCmdMute},
		{Token: "quit", HelpID: "quithelp", Flags: console.ClientServer, Func: func(ctx context.Context, c *console.Console, tokens []string) bool {
			nox_client_quit_4460C0()
			return true
		}},
		{Token: "unmute", HelpID: "unmutehelp", Flags: console.ClientServer, LegacyFunc: noxCmdUnmute},
		{Token: "window", HelpID: "windowhelp", Flags: console.ClientServer, LegacyFunc: noxCmdWindow},
	}
)
