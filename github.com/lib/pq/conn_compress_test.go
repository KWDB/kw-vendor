package pq

import (
	"bytes"
	"testing"

	"github.com/lib/pq/oid"
)

func TestCompressedFixedStringData(t *testing.T) {
	raw := []byte{4, 0, 0xde, 0xad, 0xbe, 0xef, 0, 0}
	if got := compressedFixedStringData(raw); !bytes.Equal(got, raw[2:6]) {
		t.Fatalf("unexpected fixed string data: %x", got)
	}
}

func TestCompressedVarValuePreservesBytes(t *testing.T) {
	raw := []byte{0xde, 0xad, 0xbe, 0xef}
	got, ok := compressedVarValue(oid.T_varbytea, raw).(string)
	if !ok || got != `\xdeadbeef` {
		t.Fatalf("unexpected varbytea value: %#v", got)
	}

	raw[0] = 0
	if got != `\xdeadbeef` {
		t.Fatal("compressed varbytea value aliases the protocol buffer")
	}
}
