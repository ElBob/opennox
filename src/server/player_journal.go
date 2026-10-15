package server

import (
	"strings"

	"github.com/opennox/libs/noxnet/netmsg"

	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func (e *PlayerJournal) Name() string {
	return alloc.GoStringS(e.EntryBuf[:])
}

// JournalAdd adds a new entry to the beginning of the player's journal.
func (p *Player) JournalAdd(name string, typ uint16) *PlayerJournal { // nox_xxx_journalEntryAdd_427490
	if i := strings.IndexByte(name, 0); i >= 0 {
		name = name[:i]
	}
	e, _ := alloc.New(PlayerJournal{})
	alloc.StrCopyZero(e.EntryBuf[:], name)
	e.Field3 = typ
	e.Next = p.Journal
	if p.Journal != nil {
		p.Journal.Prev = e
	}
	p.Journal = e
	return e
}

// JournalFind finds a journal entry by name.
func (p *Player) JournalFind(name string) *PlayerJournal {
	for e := p.Journal; e != nil; e = e.Next {
		if e.Name() == name {
			return e
		}
	}
	return nil
}

func (p *Player) journalUnlink(e *PlayerJournal) {
	if e.Prev != nil {
		e.Prev.Next = e.Next
	}
	if e.Next != nil {
		e.Next.Prev = e.Prev
	}
	if e == p.Journal {
		p.Journal = e.Next
	}
	alloc.Free(e)
}

// JournalRemove removes a journal entry by name.
func (p *Player) JournalRemove(name string) bool { // nox_xxx_journalEntryRemove_427590
	e := p.JournalFind(name)
	if e == nil {
		return false
	}
	p.journalUnlink(e)
	return true
}

// JournalUpdate changes the type of a journal entry.
func (p *Player) JournalUpdate(name string, typ uint16) *PlayerJournal { // nox_xxx_journalUpdateEntry_4276B0
	e := p.JournalFind(name)
	if e == nil {
		return nil
	}
	e.Field3 = typ
	return e
}

// JournalRemoveByType removes all journal entries that have any of the type bits set.
func (p *Player) JournalRemoveByType(typ uint16) { // sub_4277B0
	for e := p.Journal; e != nil; {
		next := e.Next
		if e.Field3&typ != 0 {
			p.journalUnlink(e)
		}
		e = next
	}
}

func (s *Server) netSendJournal(pind ntype.PlayerInd, op byte, name string, typ uint16) {
	var buf [68]byte
	buf[0] = byte(netmsg.MSG_JOURNAL_MSG)
	buf[1] = op
	alloc.StrCopyZero(buf[2:66], name)
	buf[66] = byte(typ)
	buf[67] = byte(typ >> 8)
	s.NetSendPacketXxx0(int(pind), buf[:], 0, 1)
}

func (s *Server) NetSendJournalAdd(pind ntype.PlayerInd, e *PlayerJournal) { // nox_xxx_netSendJournalAdd_4D9440
	s.netSendJournal(pind, 1, e.Name(), e.Field3)
}

func (s *Server) NetSendJournalRemove(pind ntype.PlayerInd, name string) { // nox_xxx_netSendJournalRemove_4D94A0
	s.netSendJournal(pind, 2, name, 0)
}

func (s *Server) NetSendJournalUpdate(pind ntype.PlayerInd, e *PlayerJournal) { // nox_xxx_netSendJournalUpdate_4D9500
	s.netSendJournal(pind, 3, e.Name(), e.Field3)
}
