package opennox

import (
	"github.com/opennox/noxscript/ns/v4"

	"github.com/opennox/opennox/v1/common/sound"
	"github.com/opennox/opennox/v1/server"
)

func (s noxScriptNS) GetQuestStatus(name string) int {
	return int(s.s.QuestStatus.GetInt(name))
}

func (s noxScriptNS) GetQuestStatusFloat(name string) float32 {
	return s.s.QuestStatus.GetFloat(name)
}

func (s noxScriptNS) SetQuestStatus(status int, name string) {
	s.s.QuestStatus.SetInt(name, int32(status))
}

func (s noxScriptNS) SetQuestStatusFloat(status float32, name string) {
	s.s.QuestStatus.SetFloat(name, status)
}

func (s noxScriptNS) ResetQuestStatus(name string) {
	s.s.QuestStatus.Reset(name)
}

func (s noxScriptNS) JournalEntry(obj ns.Obj, msg ns.StringID, typ ns.EntryType) {
	if obj == nil {
		for _, it := range s.s.Players.ListUnits() {
			s.s.journalEntryAdd(it, msg, uint16(typ))
		}
	} else {
		s.s.journalEntryAdd(server.ToObject(obj.(server.Obj)), msg, uint16(typ))
		if (typ & 0xB) != 0 {
			s.s.Audio.EventObj(sound.SoundJournalEntryAdd, obj.(server.Obj), 0, 0)
		}
	}
}

func (s noxScriptNS) JournalEdit(obj ns.Obj, message ns.StringID, typ ns.EntryType) {
	if obj == nil {
		s.s.journalEntryUpdateAll(message, uint16(typ))
	} else {
		s.s.journalEntryUpdate(server.ToObject(obj.(server.Obj)), message, uint16(typ))
	}
}

func (s noxScriptNS) JournalDelete(obj ns.Obj, message ns.StringID) {
	if obj == nil {
		s.s.journalEntryRemoveAll(message)
	} else {
		s.s.journalEntryRemove(server.ToObject(obj.(server.Obj)), message)
	}
}

func (s noxScriptNS) JournalEntryStr(obj ns.Obj, msg string, typ ns.EntryType) {
	//TODO implement me
	panic("implement me")
}

func (s noxScriptNS) JournalEditStr(obj ns.Obj, message string, typ ns.EntryType) {
	//TODO implement me
	panic("implement me")
}

func (s noxScriptNS) JournalDeleteStr(obj ns.Obj, message string) {
	//TODO implement me
	panic("implement me")
}
