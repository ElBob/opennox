package opennox

import (
	"encoding/binary"
	"unicode/utf16"
	"unsafe"

	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/opennox/v1/internal/netlist"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func init() {
	legacy.ScriptCallbackReadWrite = scriptCallbackReadWrite
}

// netSendChat sends a chat bubble message for the object to all players.
func (s *Server) netSendChat(u *server.Object, msg string, dt uint16) { // nox_xxx_netSendChat_528AC0
	var buf [520]byte
	buf[0] = byte(netmsg.MSG_TEXT_MESSAGE)
	if legacy.Nox_xxx_cliCanTalkMB_4100F0(msg) {
		buf[3] |= 2
	} else {
		buf[3] |= 4
	}
	binary.LittleEndian.PutUint16(buf[1:], uint16(s.GetUnitNetCode(u)))
	binary.LittleEndian.PutUint16(buf[4:], uint16(int64(u.PosVec.X)))
	binary.LittleEndian.PutUint16(buf[6:], uint16(int64(u.PosVec.Y)))
	binary.LittleEndian.PutUint16(buf[9:], dt)
	wmsg := append(utf16.Encode([]rune(msg)), 0)
	buf[8] = byte(len(wmsg))
	sz := 1
	if buf[3]&4 != 0 {
		sz = 2
		for i, c := range wmsg {
			binary.LittleEndian.PutUint16(buf[11+2*i:], c)
		}
	} else {
		// The original code uses "%S", which only keeps the lowest byte.
		for i, c := range wmsg {
			buf[11+i] = byte(c)
		}
	}
	n := 11 + sz*int(buf[8])
	for _, it := range s.Players.ListUnits() {
		s.NetList.AddToMsgListCli(it.UpdateDataPlayer().Player.PlayerIndex(), netlist.Kind1, buf[:n])
	}
}

// scriptCallbackReadWrite reads or writes a script callback name and flags.
func scriptCallbackReadWrite(cf *cryptfile.CryptFile, cb *server.ScriptCallback, name *byte) bool { // nox_xxx_xferReadScriptHandler_4F5580
	const editor = noxflags.GameFlag(0x600000)
	if vers, _ := cf.ReadWriteU16(1); int16(vers) > 1 {
		return false
	}
	if cf.ReadOnly() {
		sz, _ := cf.ReadU32()
		if sz >= 0x400 {
			return false
		}
		str := make([]byte, sz)
		_, _ = cf.Read(str)
		if sz != 0 {
			str := alloc.GoStringS(append(str, 0))
			if !noxflags.HasGame(editor) {
				cb.Func = int32(noxServer.Server.NoxScriptVM.ScriptIndexByName(str))
			} else if name != nil {
				alloc.StrCopyZero(unsafe.Slice(name, len(str)+1), str)
			}
		}
	} else {
		var str string
		if noxflags.HasGame(editor) {
			if name != nil {
				str = alloc.GoString(name)
			}
		} else if cb.Func != -1 {
			str = noxServer.Server.NoxScriptVM.ScriptNameByIndex(int(cb.Func))
		}
		_ = cf.WriteU32(uint32(len(str)))
		if str != "" {
			_, _ = cf.Write([]byte(str))
		}
	}
	cb.Flags, _ = cf.ReadWriteU32(cb.Flags)
	return true
}
