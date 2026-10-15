package server

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/shoenig/test/must"
)

func questStatusDump(q *QuestStatus) string {
	var out []string
	for i := len(q.vars) - 1; i >= 0; i-- {
		v := q.vars[i]
		out = append(out, fmt.Sprintf("%s %d %08x", v.name, v.typ, v.val))
	}
	return strings.Join(out, "; ")
}

func newTestQuestStatus() *QuestStatus {
	q := &QuestStatus{}
	q.SetMapName("Con01")
	q.SetInt("foo", 1)
	q.SetInt("foofoo", 2)
	q.SetInt("barFoo", 3)
	q.SetInt("Xfoo", 4)
	q.SetInt("con02A:foo", 5)
	q.SetInt("con02A:xfoo", 6)
	q.SetInt("Con01:a1", 7)
	return q
}

func TestQuestStatusNames(t *testing.T) {
	q := &QuestStatus{}
	q.SetMapName("Con01")
	q.SetInt("foo", 5)
	must.EqOp(t, 5, q.GetInt("CON01:FOO"))
	q.SetInt("con01:FOO", 7)
	must.EqOp(t, 7, q.GetInt("foo"))
	must.EqOp(t, 0, q.GetInt("con02:foo"))
	q.SetFloat("Bar", 1.5)
	must.EqOp(t, 1.5, q.GetFloat("bar"))
	must.EqOp(t, "Con01:Bar 1 3fc00000; Con01:foo 0 00000007", questStatusDump(q))
}

// Expected results are taken from the original implementation.
func TestQuestStatusReset(t *testing.T) {
	const (
		a1      = "Con01:a1 0 00000007"
		xfoo2   = "con02A:xfoo 0 00000006"
		foo2    = "con02A:foo 0 00000005"
		xfoo    = "Con01:Xfoo 0 00000004"
		barFoo  = "Con01:barFoo 0 00000003"
		foofoo  = "Con01:foofoo 0 00000002"
		foo     = "Con01:foo 0 00000001"
		all     = a1 + "; " + xfoo2 + "; " + foo2 + "; " + xfoo + "; " + barFoo + "; " + foofoo + "; " + foo
		noFoo   = a1 + "; " + xfoo2 + "; " + foo2 + "; " + xfoo + "; " + barFoo + "; " + foofoo
		noFooFF = a1 + "; " + xfoo2 + "; " + foo2 + "; " + xfoo + "; " + barFoo
	)
	cases := []struct {
		pattern string
		exp     string
	}{
		{"*:*", ""},
		{"*", xfoo2 + "; " + foo2},
		{"foo", noFoo},
		{"FOO", noFoo},
		{"foo*", noFooFF},
		{"Con01:*", xfoo2 + "; " + foo2},
		{"CON01:f*", noFooFF},
		{"*foo", a1 + "; " + xfoo2 + "; " + foo2 + "; " + barFoo + "; " + foo},
		{"*:foo", a1 + "; " + xfoo2 + "; " + xfoo + "; " + barFoo + "; " + foofoo},
		{"*Foo", a1 + "; " + xfoo2 + "; " + foo2 + "; " + xfoo + "; " + foofoo + "; " + foo},
		{"Con01:*foo", a1 + "; " + xfoo2 + "; " + foo2 + "; " + barFoo + "; " + foo},
		{"con02A:*oo", a1 + "; " + xfoo + "; " + barFoo + "; " + foofoo + "; " + foo},
		{"con02A:*foo", a1 + "; " + foo2 + "; " + xfoo + "; " + barFoo + "; " + foofoo + "; " + foo},
		{"Co*", all},
		{"*:*foo", all},
	}
	for _, c := range cases {
		t.Run(c.pattern, func(t *testing.T) {
			q := newTestQuestStatus()
			q.Reset(c.pattern)
			must.EqOp(t, c.exp, questStatusDump(q))
		})
	}
}

func TestQuestStatusSave(t *testing.T) {
	q := &QuestStatus{}
	q.SetMapName("Con01")
	q.SetFloat("x", 1)
	q.SetInt("x", 3) // keeps the float type
	q.SetInt("y", -1)
	must.EqOp(t, "Con01:y 0 ffffffff; Con01:x 1 00000003", questStatusDump(q))

	var buf bytes.Buffer
	must.NoError(t, q.Save(&buf, true))
	must.EqOp(t, "01000200000007436f6e30313a7900000000ffffffff07436f6e30313a780100000003000000", hex.EncodeToString(buf.Bytes()))

	must.NoError(t, q.Load(&buf))
	// The original code reverses the order on load.
	must.EqOp(t, "Con01:x 1 00000003; Con01:y 0 ffffffff", questStatusDump(q))

	buf.Reset()
	must.NoError(t, q.Save(&buf, false))
	must.EqOp(t, "010000000000", hex.EncodeToString(buf.Bytes()))
	must.NoError(t, q.Load(&buf))
	must.EqOp(t, "", questStatusDump(q))

	must.Error(t, q.Load(bytes.NewReader([]byte{2, 0, 0, 0, 0, 0})))
}
