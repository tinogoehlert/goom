package convert_test

import (
	"encoding/binary"
	"testing"
)

type I uint32

func TestLittleEndian(t *testing.T) {
	b := []byte{2, 3, 4, 0}
	li := binary.LittleEndian.Uint32(b)
	ii := I(2) | I(3)<<8 | I(4)<<16
	oi := make([]byte, 4)
	binary.LittleEndian.PutUint32(oi, uint32(ii))
	if string(b) != string(oi) {
		t.Fail()
	}
	if string(li) != string(ii) {
		t.Fail()
	}

	h := []byte{0, 0, 0, 0}
	for _, n := range []uint32{0, 1, 7, 16, 47, 1678, 23532} {
		// TODO: double check track header byte order
		h[3] = byte(n>>24) & 0xff
		h[2] = byte(n>>16) & 0xff
		h[1] = byte(n>>8) & 0xff
		h[0] = byte(n) & 0xff
		v := binary.LittleEndian.Uint32(h)
		if n != v {
			t.Errorf("%d != %v", n, v)
		}
	}
}