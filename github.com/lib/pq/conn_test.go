package pq

import (
	"encoding/binary"
	"math"
	"testing"

	"github.com/lib/pq/oid"
)

func TestKwDataChunkDepressGetDataNumericIntLayout(t *testing.T) {
	var payload [9]byte
	payload[0] = 0
	binary.LittleEndian.PutUint64(payload[1:], uint64(11))

	kdc := &KwDataChunk{}
	got := kdc.DepressGetData(nil, 0, 0, oid.T_numeric, formatBinary, nil, payload[:])
	if got != "11" {
		t.Fatalf("expected numeric int layout to decode to 11, got %#v", got)
	}
}

func TestKwDataChunkDepressGetDataNumericFloatLayout(t *testing.T) {
	var payload [9]byte
	payload[0] = 1
	binary.LittleEndian.PutUint64(payload[1:], math.Float64bits(3.5))

	kdc := &KwDataChunk{}
	got := kdc.DepressGetData(nil, 0, 0, oid.T_numeric, formatBinary, nil, payload[:])
	if got != "3.5" {
		t.Fatalf("expected numeric float layout to decode to 3.5, got %#v", got)
	}
}

func TestKwDataChunkDepressGetDataNumericRawFloat64Layout(t *testing.T) {
	var payload [8]byte
	binary.LittleEndian.PutUint64(payload[:], math.Float64bits(5.5))

	kdc := &KwDataChunk{}
	got := kdc.DepressGetData(nil, 0, 0, oid.T_numeric, formatBinary, nil, payload[:])
	if got != "5.5" {
		t.Fatalf("expected raw float64 layout to decode to 5.5, got %#v", got)
	}
}
