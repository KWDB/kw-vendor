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

func TestDataChunkIsNull(t *testing.T) {
	chunk := &DataChunk{nullBitmap: []byte{0b00000100}}
	if chunk.isNull(1) {
		t.Fatal("row 1 should not be NULL")
	}
	if !chunk.isNull(2) {
		t.Fatal("row 2 should be NULL")
	}
	if chunk.isNull(8) {
		t.Fatal("row outside the bitmap should not be NULL")
	}
}

func TestCompressedVarData(t *testing.T) {
	data := []byte{3, 0, 'a', 'b', 'c', 0, 0}

	got, err := compressedVarData(data, 0)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, []byte("abc")) {
		t.Fatalf("unexpected variable data: %q", got)
	}

	got, err = compressedVarData(data, 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty variable data, got %q", got)
	}

	if _, err := compressedVarData(data, uint32(len(data))); err == nil {
		t.Fatal("expected invalid offset error")
	}
	if _, err := compressedVarData([]byte{4, 0, 'a'}, 0); err == nil {
		t.Fatal("expected invalid length error")
	}
}
