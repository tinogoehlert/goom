package goom

import (
	"encoding/binary"
)

func i16Tof(buff []byte) float32 {
	return float32(int16(binary.LittleEndian.Uint16(buff)))
}

func i16(buff []byte) int16 {
	return int16(binary.LittleEndian.Uint16(buff))
}
