package noxscript

import "github.com/opennox/noxscript/ns/asm"

func init() {
	Register(asm.BuiltinSetQuestStatus, nsSetQuestStatus)
	Register(asm.BuiltinSetQuestStatusFloat, nsSetQuestStatusFloat)
	Register(asm.BuiltinGetQuestStatus, nsGetQuestStatus)
	Register(asm.BuiltinGetQuestStatusFloat, nsGetQuestStatusFloat)
	Register(asm.BuiltinResetQuestStatus, nsResetQuestStatus)
}

func nsSetQuestStatus(s VM) int {
	name := s.PopString()
	v := int(s.PopI32())
	s.NoxScript().SetQuestStatus(v, name)
	return 0
}

func nsSetQuestStatusFloat(s VM) int {
	name := s.PopString()
	v := s.PopF32()
	s.NoxScript().SetQuestStatusFloat(v, name)
	return 0
}

func nsGetQuestStatus(s VM) int {
	name := s.PopString()
	s.PushI32(int32(s.NoxScript().GetQuestStatus(name)))
	return 0
}

func nsGetQuestStatusFloat(s VM) int {
	name := s.PopString()
	s.PushF32(s.NoxScript().GetQuestStatusFloat(name))
	return 0
}

func nsResetQuestStatus(s VM) int {
	name := s.PopString()
	s.NoxScript().ResetQuestStatus(name)
	return 0
}
