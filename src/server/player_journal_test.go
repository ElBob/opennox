package server

import (
	"fmt"
	"strings"
	"testing"

	"github.com/shoenig/test/must"
)

func journalDump(p *Player) string {
	var out []string
	for e := p.Journal; e != nil; e = e.Next {
		out = append(out, fmt.Sprintf("%s %04x", e.Name(), e.Field3))
	}
	return strings.Join(out, "; ")
}

// Expected results are taken from the original implementation.
func TestPlayerJournal(t *testing.T) {
	long, trunc := strings.Repeat("A", 70), strings.Repeat("A", 63)
	p := new(Player)
	p.JournalAdd("Hint1", 1)
	p.JournalAdd("FindKey", 2)
	p.JournalAdd(long, 4)
	p.JournalAdd("Hint2", 8)
	p.JournalAdd("Hint1", 2)
	must.EqOp(t, "Hint1 0002; Hint2 0008; "+trunc+" 0004; FindKey 0002; Hint1 0001", journalDump(p))

	must.True(t, p.JournalRemove("Hint1"))
	must.False(t, p.JournalRemove(long))
	must.True(t, p.JournalRemove(trunc))
	must.NotNil(t, p.JournalUpdate("FindKey", 0xE))
	must.Nil(t, p.JournalUpdate("hint2", 1))
	must.EqOp(t, "Hint2 0008; FindKey 000e; Hint1 0001", journalDump(p))

	p.JournalRemoveByType(0xC)
	must.EqOp(t, "Hint1 0001", journalDump(p))
}
