package server

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"strings"
)

const (
	questStatusInt   = 0
	questStatusFloat = 1
)

type questStatusVar struct {
	name string
	typ  uint32
	val  uint32
}

// QuestStatus stores named quest variables set by map scripts.
//
// Variable names without a colon are prefixed with the current map name ("map:name").
// Names are compared case-insensitively.
type QuestStatus struct {
	mapName string
	// vars are ordered from the oldest to the newest.
	// The original code kept them in a linked list starting from the newest one.
	vars []questStatusVar
}

// SetMapName sets the map name used as a prefix for quest variables.
func (q *QuestStatus) SetMapName(name string) { // sub_500510
	q.mapName = name
}

func (q *QuestStatus) fullName(name string) string { // sub_5009B0
	if strings.IndexByte(name, ':') >= 0 {
		return name
	}
	return q.mapName + ":" + name
}

func (q *QuestStatus) find(name string) int { // nox_xxx_scriptGetJournal_5005E0
	full := q.fullName(name)
	for i := len(q.vars) - 1; i >= 0; i-- {
		if equalFoldASCII(q.vars[i].name, full) {
			return i
		}
	}
	return -1
}

func (q *QuestStatus) set(name string, typ uint32, val uint32) {
	if i := q.find(name); i >= 0 {
		// The original code only updates the value, not the type.
		q.vars[i].val = val
		return
	}
	q.vars = append(q.vars, questStatusVar{name: q.fullName(name), typ: typ, val: val})
}

// SetInt sets an integer quest variable.
func (q *QuestStatus) SetInt(name string, v int32) { // nox_xxx_journalQuestSet_500540
	q.set(name, questStatusInt, uint32(v))
}

// SetFloat sets a float quest variable.
func (q *QuestStatus) SetFloat(name string, v float32) { // nox_xxx_journalQuestSetBool_5006B0
	q.set(name, questStatusFloat, math.Float32bits(v))
}

// GetInt returns an integer quest variable, or zero if it is not set.
func (q *QuestStatus) GetInt(name string) int32 { // sub_500750
	if i := q.find(name); i >= 0 {
		return int32(q.vars[i].val)
	}
	return 0
}

// GetFloat returns a float quest variable, or zero if it is not set.
func (q *QuestStatus) GetFloat(name string) float32 { // sub_500770
	if i := q.find(name); i >= 0 {
		return math.Float32frombits(q.vars[i].val)
	}
	return 0
}

func (q *QuestStatus) deleteIf(fnc func(name string) bool) {
	out := q.vars[:0]
	for _, v := range q.vars {
		if !fnc(v.name) {
			out = append(out, v)
		}
	}
	clear(q.vars[len(out):])
	q.vars = out
}

// Reset removes quest variables matching the pattern.
//
// Supported patterns are "*:*" (all variables), "prefix*", "*suffix" and "map:*suffix".
func (q *QuestStatus) Reset(pattern string) { // sub_5007E0
	full := q.fullName(pattern)
	star := strings.IndexByte(full, '*')
	if star < 0 {
		if i := q.find(pattern); i >= 0 {
			q.vars = append(q.vars[:i], q.vars[i+1:]...)
		}
		return
	}
	switch {
	case full == "*:*":
		q.vars = nil
	case star == len(full)-1:
		prefix := full[:star]
		q.deleteIf(func(name string) bool {
			return hasPrefixFoldASCII(name, prefix)
		})
	case star == 0:
		suffix := full[1:]
		q.deleteIf(func(name string) bool {
			return firstIndexIsSuffix(name, suffix)
		})
	default:
		// The original code assumes the pattern is "map:*suffix",
		// and always skips one character after the colon.
		colon := strings.IndexByte(full, ':')
		prefix := full[:colon+1]
		suffix := skipBytes(full, colon+2)
		q.deleteIf(func(name string) bool {
			if !hasPrefixFoldASCII(name, prefix) {
				return false
			}
			return firstIndexIsSuffix(skipBytes(name, colon+2), suffix)
		})
	}
}

// Save writes quest variables. Variables are only saved in coop mode.
func (q *QuestStatus) Save(w io.Writer, coop bool) error { // sub_500A60
	buf := binary.LittleEndian.AppendUint16(nil, 1)
	if !coop {
		buf = binary.LittleEndian.AppendUint32(buf, 0)
		_, err := w.Write(buf)
		return err
	}
	buf = binary.LittleEndian.AppendUint32(buf, uint32(len(q.vars)))
	for i := len(q.vars) - 1; i >= 0; i-- {
		v := q.vars[i]
		n := byte(len(v.name))
		buf = append(buf, n)
		buf = append(buf, v.name[:n]...)
		buf = binary.LittleEndian.AppendUint32(buf, v.typ)
		buf = binary.LittleEndian.AppendUint32(buf, v.val)
	}
	_, err := w.Write(buf)
	return err
}

// Load replaces all quest variables with ones read from r.
func (q *QuestStatus) Load(r io.Reader) error { // sub_500B70
	q.vars = nil
	var b [4]byte
	if _, err := io.ReadFull(r, b[:2]); err != nil {
		return err
	}
	if vers := int16(binary.LittleEndian.Uint16(b[:2])); vers > 1 {
		return fmt.Errorf("unsupported quest status version: %d", vers)
	}
	if _, err := io.ReadFull(r, b[:4]); err != nil {
		return err
	}
	cnt := binary.LittleEndian.Uint32(b[:4])
	var name [256]byte
	for range cnt {
		if _, err := io.ReadFull(r, b[:1]); err != nil {
			return err
		}
		n := name[:b[0]]
		if _, err := io.ReadFull(r, n); err != nil {
			return err
		}
		if i := strings.IndexByte(string(n), 0); i >= 0 {
			n = n[:i]
		}
		if _, err := io.ReadFull(r, b[:4]); err != nil {
			return err
		}
		typ := binary.LittleEndian.Uint32(b[:4])
		if typ != questStatusInt && typ != questStatusFloat {
			continue
		}
		if _, err := io.ReadFull(r, b[:4]); err != nil {
			return err
		}
		q.set(string(n), typ, binary.LittleEndian.Uint32(b[:4]))
	}
	return nil
}

func lowerASCII(c byte) byte {
	if 'A' <= c && c <= 'Z' {
		return c + 'a' - 'A'
	}
	return c
}

func equalFoldASCII(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	return hasPrefixFoldASCII(a, b)
}

func hasPrefixFoldASCII(s, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}
	for i := 0; i < len(prefix); i++ {
		if lowerASCII(s[i]) != lowerASCII(prefix[i]) {
			return false
		}
	}
	return true
}

// firstIndexIsSuffix checks that the first occurrence of sub in s is at the end of s.
func firstIndexIsSuffix(s, sub string) bool {
	i := strings.Index(s, sub)
	return i >= 0 && len(s)-i == len(sub)
}

func skipBytes(s string, n int) string {
	if n >= len(s) {
		return ""
	}
	return s[n:]
}
