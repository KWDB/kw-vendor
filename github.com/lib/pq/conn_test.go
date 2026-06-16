package pq

import (
	"encoding/binary"
	"math"
	"testing"
	"time"

	"github.com/lib/pq/oid"
)

func TestDecodeCompressedNumericIntLayout(t *testing.T) {
	var payload [9]byte
	payload[0] = 0
	binary.LittleEndian.PutUint64(payload[1:], uint64(11))

	kdc := &KwDataChunk{}
	got := kdc.DepressGetData(nil, 0, 0, oid.T_numeric, formatBinary, nil, payload[:])
	if got != int64(11) {
		t.Fatalf("expected numeric int layout to decode to 11, got %#v", got)
	}
}

func TestDecodeCompressedNumericFloatLayout(t *testing.T) {
	var payload [9]byte
	payload[0] = 1
	binary.LittleEndian.PutUint64(payload[1:], math.Float64bits(3.5))

	kdc := &KwDataChunk{}
	got := kdc.DepressGetData(nil, 0, 0, oid.T_numeric, formatBinary, nil, payload[:])
	if got != float64(3.5) {
		t.Fatalf("expected numeric float layout to decode to 3.5, got %#v", got)
	}
}

func TestCompressedDecodeOIDNumericRawFloat64Layout(t *testing.T) {
	var payload [8]byte
	binary.LittleEndian.PutUint64(payload[:], math.Float64bits(5.5))

	gotOID := compressedDecodeOID(oid.T_numeric, uint32(len(payload)))
	if gotOID != oid.T_float8 {
		t.Fatalf("expected decode OID %v, got %v", oid.T_float8, gotOID)
	}

	kdc := &KwDataChunk{}
	got := kdc.DepressGetData(nil, 0, 0, gotOID, formatBinary, nil, payload[:])
	if got != float64(5.5) {
		t.Fatalf("expected raw float64 layout to decode to 5.5, got %#v", got)
	}
}

func TestCompressedDecodeOIDInt8WithInt4Storage(t *testing.T) {
	var payload [4]byte
	binary.LittleEndian.PutUint32(payload[:], uint32(10))

	gotOID := compressedDecodeOID(oid.T_int8, uint32(len(payload)))
	if gotOID != oid.T_int4 {
		t.Fatalf("expected decode OID %v, got %v", oid.T_int4, gotOID)
	}

	kdc := &KwDataChunk{}
	got := kdc.DepressGetData(nil, 0, 0, gotOID, formatBinary, nil, payload[:])
	if got != int64(10) {
		t.Fatalf("expected int4 storage to decode to 10, got %#v", got)
	}
}
func TestDecodeCompressedFloat4FormattingValues(t *testing.T) {
	var payload [4]byte
	binary.LittleEndian.PutUint32(payload[:], math.Float32bits(float32(1.1)))

	kdc := &KwDataChunk{}
	got := kdc.DepressGetData(nil, 0, 0, oid.T_float4, formatBinary, nil, payload[:])
	if _, ok := got.(float32); !ok {
		t.Fatalf("expected ordinary compressed float4 to decode as float32, got %T", got)
	}

	binary.LittleEndian.PutUint32(payload[:], math.Float32bits(math.MaxFloat32))

	got = kdc.DepressGetData(nil, 0, 0, oid.T_float4, formatBinary, nil, payload[:])
	if _, ok := got.(float64); !ok {
		t.Fatalf("expected max compressed float4 to decode as float64, got %T", got)
	}
	if got != float64(float32(math.MaxFloat32)) {
		t.Fatalf("expected max compressed float4 to decode to %v, got %#v", float64(float32(math.MaxFloat32)), got)
	}
}

func TestDecodeCompressedTimestampUsesUTC(t *testing.T) {
	var payload [8]byte
	ts := time.Date(2026, 6, 10, 8, 0, 0, 0, time.UTC).UnixMilli()
	binary.LittleEndian.PutUint64(payload[:], uint64(ts))

	kdc := &KwDataChunk{}
	got := kdc.DepressGetData(nil, 0, 0, oid.T_timestamptz, formatBinary, nil, payload[:])
	if string(got.([]byte)) != "2026-06-10 08:00:00+00:00" {
		t.Fatalf("expected UTC timestamp text, got %#v", got)
	}
}
