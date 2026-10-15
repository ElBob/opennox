package server

import (
	"encoding/binary"
	"math"
	"strings"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

const (
	questStatusInt   = 0
	questStatusFloat = 1
)

// questStatusVar is a node of the quest status list, allocated with calloc(1, 0x94).
type questStatusVar struct {
	name [132]byte       // 0, 0
	typ  uint32          // 33, 132
	val  uint32          // 34, 136
	next *questStatusVar // 35, 140
	prev *questStatusVar // 36, 144
}

// questStatusHead points to the list head (dword_5d4594_1570272).
// It is replaced with a pointer to the C global by the legacy package.
var questStatusHead = new(*questStatusVar)

// SetQuestStatusHead sets the location of the quest status list head.
func SetQuestStatusHead(p unsafe.Pointer) {
	questStatusHead = (**questStatusVar)(p)
}

func questStatusMapName() []byte { // 0x5D4594+1570008
	return memmap.Slice(0x5D4594, 1570008)
}

func questStatusScratch() []byte { // 0x5D4594+1570140
	return memmap.Slice(0x5D4594, 1570140)
}

// QuestStatusFile is the save file used by QuestStatus.Save and QuestStatus.Load.
type QuestStatusFile interface {
	ReadWrite(p []byte) (int, error)
}

// QuestStatus stores named quest variables set by map scripts.
//
// Variable names without a colon are prefixed with the current map name ("map:name").
// Names are compared case-insensitively.
//
// As in the original code, variables are kept in a linked list starting from
// dword_5d4594_1570272 (newest first), and names are built in the buffer at 0x5D4594+1570140.
type QuestStatus struct{}

func cstrCopy(dst []byte, s string) int {
	n := copy(dst, s)
	dst[n] = 0
	return n
}

func cstrName(s string) string {
	if i := strings.IndexByte(s, 0); i >= 0 {
		s = s[:i]
	}
	return s
}

// SetMapName sets the map name used as a prefix for quest variables.
func (q *QuestStatus) SetMapName(name string) { // sub_500510
	cstrCopy(questStatusMapName(), cstrName(name))
}

// find builds the full name of the variable in the scratch buffer and finds the variable.
func (q *QuestStatus) find(name string) *questStatusVar { // nox_xxx_scriptGetJournal_5005E0
	name = cstrName(name)
	buf := questStatusScratch()
	if strings.IndexByte(name, ':') >= 0 {
		cstrCopy(buf, name)
	} else {
		n := cstrCopy(buf, alloc.GoStringS(questStatusMapName()))
		*(*uint16)(unsafe.Pointer(&buf[n])) = memmap.Uint16(0x587000, 217952) // ":"
		n = len(alloc.GoStringS(buf))
		cstrCopy(buf[n:], name)
	}
	full := alloc.GoStringS(buf)
	for it := *questStatusHead; it != nil; it = it.next {
		if equalFoldASCII(alloc.GoStringS(it.name[:]), full) { // nox_strcmpi
			return it
		}
	}
	return nil
}

// fullName builds the full name of the variable in the scratch buffer.
func (q *QuestStatus) fullName(name string) { // sub_5009B0
	name = cstrName(name)
	buf := questStatusScratch()
	if strings.IndexByte(name, ':') >= 0 {
		cstrCopy(buf, name)
		return
	}
	n := cstrCopy(buf, alloc.GoStringS(questStatusMapName()))
	*(*uint16)(unsafe.Pointer(&buf[n])) = memmap.Uint16(0x587000, 217960) // ":"
	n = len(alloc.GoStringS(buf))
	cstrCopy(buf[n:], name)
}

func (q *QuestStatus) set(name string, typ uint32, val uint32) *questStatusVar {
	if v := q.find(name); v != nil {
		// The original code only updates the value, not the type.
		v.val = val
		return v
	}
	p, _ := alloc.Calloc(1, 0x94)
	v := (*questStatusVar)(p)
	// The original code uses strcpy here and may overflow the node, so we truncate the name.
	full := alloc.GoStringS(questStatusScratch())
	if len(full) >= len(v.name) {
		full = full[:len(v.name)-1]
	}
	cstrCopy(v.name[:], full)
	v.typ = typ
	v.val = val
	v.prev = nil
	v.next = *questStatusHead
	if *questStatusHead != nil {
		(*questStatusHead).prev = v
	}
	*questStatusHead = v
	return v
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
	if v := q.find(name); v != nil {
		return int32(v.val)
	}
	return 0
}

// GetFloat returns a float quest variable, or zero if it is not set.
func (q *QuestStatus) GetFloat(name string) float32 { // sub_500770
	if v := q.find(name); v != nil {
		return math.Float32frombits(v.val)
	}
	return 0
}

func (q *QuestStatus) delete(v *questStatusVar) { // sub_500790
	if v.prev != nil {
		v.prev.next = v.next
	}
	if v.next != nil {
		v.next.prev = v.prev
	}
	if v == *questStatusHead {
		*questStatusHead = v.next
	}
	alloc.Free(v)
}

// Reset removes quest variables matching the pattern.
//
// Supported patterns are "*:*" (all variables), "prefix*", "*suffix" and "map:*suffix".
func (q *QuestStatus) Reset(pattern string) { // sub_5007E0
	q.fullName(pattern)
	buf := questStatusScratch()
	full := alloc.GoStringS(buf)
	star := strings.IndexByte(full, '*')
	if star < 0 {
		if v := q.find(pattern); v != nil {
			q.delete(v)
		}
		return
	}
	v3 := len(full) + 1
	if full == "*:*" {
		for v := *questStatusHead; v != nil; {
			next := v.next
			q.delete(v)
			v = next
		}
	} else if star == v3-2 {
		prefix := full[:v3-2]
		for v := *questStatusHead; v != nil; {
			next := v.next
			if hasPrefixFoldASCII(alloc.GoStringS(v.name[:]), prefix) { // nox_strnicmp
				q.delete(v)
			}
			v = next
		}
	} else if star == 0 {
		suffix := alloc.GoStringS(buf[1:])
		for v := *questStatusHead; v != nil; {
			next := v.next
			if firstIndexIsSuffix(alloc.GoStringS(v.name[:]), suffix) {
				q.delete(v)
			}
			v = next
		}
	} else {
		// The original code assumes the pattern is "map:*suffix",
		// and always skips one character after the colon.
		// It may read past the end of the strings, so we do the same.
		colon := strings.IndexByte(full, ':')
		prefix := full[:colon+1]
		suffix := alloc.GoStringS(buf[colon+2:])
		for v := *questStatusHead; v != nil; {
			next := v.next
			if hasPrefixFoldASCII(alloc.GoStringS(v.name[:]), prefix) { // nox_strnicmp
				rest := unsafe.Slice((*byte)(unsafe.Add(unsafe.Pointer(v), colon+2)), 0x94-(colon+2))
				if firstIndexIsSuffix(alloc.GoStringS(rest), suffix) {
					q.delete(v)
				}
			}
			v = next
		}
	}
}

// Save writes quest variables. Variables are only saved in coop mode.
func (q *QuestStatus) Save(cf QuestStatusFile, coop bool) int { // sub_500A60
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], 1)
	cf.ReadWrite(buf[:2])
	if int16(binary.LittleEndian.Uint16(buf[:2])) > 1 {
		return 0
	}
	cnt := uint32(0)
	for v := *questStatusHead; v != nil; v = v.next {
		cnt++
	}
	if !coop {
		binary.LittleEndian.PutUint32(buf[:], 0)
		cf.ReadWrite(buf[:4])
		return 1
	}
	binary.LittleEndian.PutUint32(buf[:], cnt)
	cf.ReadWrite(buf[:4])
	for v := *questStatusHead; v != nil; v = v.next {
		buf[0] = byte(len(alloc.GoStringS(v.name[:])))
		cf.ReadWrite(buf[:1])
		cf.ReadWrite(v.name[:buf[0]])
		cf.ReadWrite(unsafe.Slice((*byte)(unsafe.Pointer(&v.typ)), 4))
		if typ := v.typ; typ == questStatusInt || typ == questStatusFloat {
			cf.ReadWrite(unsafe.Slice((*byte)(unsafe.Pointer(&v.val)), 4))
		}
	}
	return 1
}

// Load replaces all quest variables with ones read from cf.
func (q *QuestStatus) Load(cf QuestStatusFile) int { // sub_500B70
	q.Reset("*:*")
	var (
		v2   [1]byte
		v3   [2]byte
		v4   [4]byte
		v5   [4]byte
		v6   [4]byte
		v7   [4]byte
		name [256]byte
	)
	binary.LittleEndian.PutUint16(v3[:], 1)
	cf.ReadWrite(v3[:])
	if int16(binary.LittleEndian.Uint16(v3[:])) > 1 {
		return 0
	}
	cf.ReadWrite(v4[:])
	cnt := binary.LittleEndian.Uint32(v4[:])
	for i := uint32(0); i < cnt; i++ {
		cf.ReadWrite(v2[:])
		cf.ReadWrite(name[:v2[0]])
		name[v2[0]] = 0
		cf.ReadWrite(v5[:])
		switch binary.LittleEndian.Uint32(v5[:]) {
		case questStatusInt:
			cf.ReadWrite(v6[:])
			q.set(alloc.GoStringS(name[:]), questStatusInt, binary.LittleEndian.Uint32(v6[:])) // nox_xxx_journalQuestSet_500540
		case questStatusFloat:
			cf.ReadWrite(v7[:])
			q.set(alloc.GoStringS(name[:]), questStatusFloat, binary.LittleEndian.Uint32(v7[:])) // nox_xxx_journalQuestSetBool_5006B0
		}
	}
	return 1
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
	i := strings.Index(s, sub) // strstr
	return i >= 0 && len(s)-i == len(sub)
}
