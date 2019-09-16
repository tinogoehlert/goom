package goom

import (
	"encoding/binary"
	"strings"
)

//DoomStr DOOM string, 8 bytes, filled with 0's
type DoomStr [8]byte

// ToString to string
func (ds *DoomStr) ToString() string {
	return strings.TrimRight(string(ds[:]), "\x00")
}

// MagicU32 uint32 with special MSB
type MagicU32 uint32

// MagicBit gets the Magic bit as bool
func (child MagicU32) MagicBit() bool {
	return (child & (1 << 31)) > 0
}

// Uint32 gets the number without magic
func (child MagicU32) Uint32() uint32 {
	return uint32(child &^ (1 << 31))
}

// Fixed32 float16.16 packed into int32
type Fixed32 int32

// ToFloat32 convert to float32
func (num Fixed32) ToFloat32() float32 {
	return float32(num) / 65536.0
}

func i16Tof(buff []byte) float32 {
	return float32(int16(binary.LittleEndian.Uint16(buff)))
}

func i16(buff []byte) int16 {
	return int16(binary.LittleEndian.Uint16(buff))
}
