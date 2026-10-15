package opennox

import (
	"github.com/opennox/libs/object"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func init() {
	legacy.Nox_xxx_comJournalEntryAdd_427500 = func(u *server.Object, name string, typ uint16) {
		noxServer.journalEntryAdd(u, name, typ)
	}
}

func journalPlayer(u *server.Object) *server.Player {
	if u == nil || !u.Class().Has(object.ClassPlayer) {
		return nil
	}
	return u.UpdateDataPlayer().Player
}

func (s *Server) journalEntryAdd(u *server.Object, name string, typ uint16) { // nox_xxx_comJournalEntryAdd_427500
	pl := journalPlayer(u)
	if pl == nil {
		return
	}
	e := pl.JournalAdd(name, typ)
	if pl.PlayerIndex() == server.HostPlayerIndex {
		legacy.Nox_xxx_cliBuildJournalString_469BC0()
	} else {
		s.NetSendJournalAdd(pl.PlayerIndex(), e)
	}
}

func (s *Server) journalEntryRemove(u *server.Object, name string) { // nox_xxx_comJournalEntryRemove_427630
	pl := journalPlayer(u)
	if pl == nil || !pl.JournalRemove(name) {
		return
	}
	if pl.PlayerIndex() == server.HostPlayerIndex {
		legacy.Nox_xxx_cliBuildJournalString_469BC0()
	} else {
		s.NetSendJournalRemove(pl.PlayerIndex(), name)
	}
}

func (s *Server) journalEntryRemoveAll(name string) { // nox_xxx_comRemoveEntryAll_427680
	for _, u := range s.Players.ListUnits() {
		s.journalEntryRemove(u, name)
	}
}

func (s *Server) journalEntryUpdate(u *server.Object, name string, typ uint16) { // nox_xxx_comJournalEntryUpdate_427720
	pl := journalPlayer(u)
	if pl == nil {
		return
	}
	e := pl.JournalUpdate(name, typ)
	// The original code doesn't rebuild the journal string for the host.
	if e != nil && pl.PlayerIndex() != server.HostPlayerIndex {
		s.NetSendJournalUpdate(pl.PlayerIndex(), e)
	}
}

func (s *Server) journalEntryUpdateAll(name string, typ uint16) { // nox_xxx_comUpdateEntryAll_427770
	for _, u := range s.Players.ListUnits() {
		s.journalEntryUpdate(u, name, typ)
	}
}
