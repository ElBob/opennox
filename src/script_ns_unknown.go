package opennox

import (
	"unsafe"

	"github.com/opennox/libs/strman"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func (s noxScriptNS) Unused5e(id string) int {
	str := s.s.Strings().GetStringInFile(strman.ID(id), "CScrFunc.c")
	return scriptAddString(str)
}

// scriptAddString stores a string in a table that the game never reads and returns its index.
func scriptAddString(str string) int { // sub_512E80
	n := legacy.Get_dword_5d4594_1599628()
	if n < 1024 {
		p, _ := alloc.CString16(str)
		*memmap.PtrUint32(0x973F18, 26664+4*uintptr(n)) = uint32(uintptr(unsafe.Pointer(p)))
		n++
		legacy.Set_dword_5d4594_1599628(n)
	}
	return n - 1
}

func (s noxScriptNS) Unused74(arg1 int, arg2 int) {}

// Unknownb8 checks the subclass bit set by MakeFriendly.
func (s noxScriptNS) Unknownb8(id int) bool { // nox_script_builtin_516790
	obj := s.s.noxScript.ScriptToObject(id)
	// For an invalid object the original code pushes an uninitialized value
	// (the unused "this" argument); we push 0 instead.
	return obj != nil && obj.ObjSubClass&0x100 != 0
}

// Unknownb9 checks the subclass bit set by BecomePet.
func (s noxScriptNS) Unknownb9(id int) bool { // nox_script_builtin_516850
	obj := s.s.noxScript.ScriptToObject(id)
	// For an invalid object the original code pushes an uninitialized value
	// (the unused "this" argument); we push 0 instead.
	return obj != nil && obj.ObjSubClass&0x80 != 0
}

func (s noxScriptNS) Unknownc4() {
	noxClient.Nox_video_stopAllFades44E040()
}
