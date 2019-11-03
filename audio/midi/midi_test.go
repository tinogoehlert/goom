package midi_test

import (
	"fmt"
	"testing"

	gmidi "github.com/tinogoehlert/goom/audio/midi"
	"github.com/tinogoehlert/goom/test"
)

func TestVarInt(t *testing.T) {
	for _, v := range []uint32{0, 127, 128, 255, 1000, 100000} {
		data := gmidi.EncodeVarInt(v)
		// fmt.Printf("%d -> %x\n", v, data)
		dec := gmidi.DecodeVarInt(data)
		test.Assert(v == dec, fmt.Sprintf("%d != %d", v, dec), t)
	}
}
