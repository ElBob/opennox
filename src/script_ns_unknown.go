package opennox

import (
	"unsafe"

	"github.com/opennox/libs/strman"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func (s noxScriptNS) Unused5e(id string) int {
	str := s.s.Strings().GetStringInFile(strman.ID(id), "CScrFunc.c")
	return scriptAddString(str)
}

var scriptStringsCnt int // dword_5d4594_1599628

// scriptAddString stores a string in a table that the game never reads and returns its index.
func scriptAddString(str string) int { // sub_512E80
	if scriptStringsCnt >= 1024 {
		return 1023
	}
	p, _ := alloc.CString16(str)
	*memmap.PtrUint32(0x973F18, 26664+4*uintptr(scriptStringsCnt)) = uint32(uintptr(unsafe.Pointer(p)))
	scriptStringsCnt++
	return scriptStringsCnt - 1
}

func (s noxScriptNS) Unused74(arg1 int, arg2 int) {}

func (s noxScriptNS) Unknownb8(id int) bool {
	//TODO implement me
	panic("implement me")
}

func (s noxScriptNS) Unknownb9(id int) bool {
	//TODO implement me
	panic("implement me")
}

func (s noxScriptNS) Unknownc4() {
	noxClient.Nox_video_stopAllFades44E040()
}
