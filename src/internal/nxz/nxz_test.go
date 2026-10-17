package nxz

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/opennox/libs/noxtest"
	"github.com/stretchr/testify/require"
)

func TestDecompressMaps(t *testing.T) {
	maps := noxtest.DataPath(t, "maps")
	files, err := os.ReadDir(maps)
	require.NoError(t, err)
	cnt := 0
	for _, fi := range files {
		mname := filepath.Join(maps, fi.Name(), fi.Name()+".map")
		zname := filepath.Join(maps, fi.Name(), fi.Name()+".nxz")
		zdata, err := readFile(zname)
		if err != nil {
			continue
		}
		cnt++
		t.Run(fi.Name(), func(t *testing.T) {
			exp, err := readFile(mname)
			require.NoError(t, err)
			got, err := Decompress(zdata)
			require.NoError(t, err)
			require.True(t, bytes.Equal(exp, got))
		})
	}
	require.NotZero(t, cnt)
}

// TestCompressMaps checks that the output is byte-exact with the original map files.
func TestCompressMaps(t *testing.T) {
	maps := noxtest.DataPath(t, "maps")
	files, err := os.ReadDir(maps)
	require.NoError(t, err)
	cnt := 0
	for _, fi := range files {
		mname := filepath.Join(maps, fi.Name(), fi.Name()+".map")
		zname := filepath.Join(maps, fi.Name(), fi.Name()+".nxz")
		exp, err := readFile(zname)
		if err != nil {
			continue
		}
		cnt++
		t.Run(fi.Name(), func(t *testing.T) {
			data, err := readFile(mname)
			require.NoError(t, err)
			got := Compress(data)
			require.True(t, bytes.Equal(exp, got))
		})
	}
	require.NotZero(t, cnt)
}

// TestCompressGolden checks the output against hashes recorded from the original C implementation.
func TestCompressGolden(t *testing.T) {
	golden := map[string]string{
		"abcd":   "cce2fd96b3f46238a6a8c22bc4a576d13bb3a885",
		"mixed":  "4aba7274cd55060af4f0c39272ae27430a38f0d9",
		"rand":   "42ee1fea5c2fa3651da29e472e6a7c59b18e5b8a",
		"repeat": "8fcbfa6ad4d5ac738cb3cad3dfc2cc090e4bbf1e",
		"runs":   "4ecbfcf7a2482a434b156a78ae6b59940e9b41af",
		"zero":   "a67a3088b5bcead884f1507f8348b69a798f85d7",
	}
	in := testInputs()
	var keys []string
	for k := range in {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	require.Len(t, keys, len(golden))
	for _, k := range keys {
		t.Run(k, func(t *testing.T) {
			h := sha1.New()
			for _, b := range in[k] {
				z := Compress(b)
				h.Write(z)
				got, err := Decompress(z)
				require.NoError(t, err)
				require.True(t, bytes.Equal(b, got), "len=%d", len(b))
			}
			require.Equal(t, golden[k], hex.EncodeToString(h.Sum(nil)))
		})
	}
}

func TestDecompressInvalid(t *testing.T) {
	data := Compress(bytes.Repeat([]byte("nox map data "), 1000))
	for n := 0; n < len(data); n++ {
		_, err := Decompress(data[:n])
		require.Error(t, err, "n=%d", n)
	}
}

// testInputs generates deterministic inputs that cover all encoder paths:
// literals, short and long matches, far distances, long runs and multiple chunks.
func testInputs() map[string][][]byte {
	r := rand.New(rand.NewSource(1))
	gen := map[string]func(n int) []byte{
		"zero": func(n int) []byte { return make([]byte, n) },
		"rand": func(n int) []byte {
			b := make([]byte, n)
			r.Read(b)
			return b
		},
		"abcd": func(n int) []byte {
			b := make([]byte, n)
			for i := range b {
				b[i] = "abcd"[r.Intn(4)]
			}
			return b
		},
		"runs": func(n int) []byte {
			b := make([]byte, 0, n)
			for len(b) < n {
				c := byte(r.Intn(256))
				k := r.Intn(600)
				for j := 0; j < k && len(b) < n; j++ {
					b = append(b, c)
				}
			}
			return b
		},
		"repeat": func(n int) []byte {
			b := make([]byte, 0, n)
			for len(b) < n {
				if len(b) > 0 && r.Intn(3) != 0 {
					// copy from the past, possibly far away
					d := 1 + r.Intn(min(len(b), 70000))
					k := 1 + r.Intn(700)
					for j := 0; j < k && len(b) < n; j++ {
						b = append(b, b[len(b)-d])
					}
				} else {
					b = append(b, byte(r.Intn(256)))
				}
			}
			return b
		},
	}
	names := []string{"zero", "rand", "abcd", "runs", "repeat"}
	sizes := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 15, 63, 64, 65, 68, 69, 70, 100, 521, 1000, 4096, 10000, 70000, 140000, 499999, 500000, 500001, 500005, 1100000}
	out := make(map[string][][]byte)
	for _, name := range names {
		for _, n := range sizes {
			out[name] = append(out[name], gen[name](n))
		}
	}
	for i := 0; i < 300; i++ {
		n := r.Intn(20000)
		name := names[2+i%3]
		out["mixed"] = append(out["mixed"], gen[name](n))
	}
	return out
}
