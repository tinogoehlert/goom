package utils

import (
	"encoding/binary"
	"strings"
)

// Int16Tof32 converts int16 byte buffer to float32
func Int16Tof32(buff []byte) float32 {
	return float32(int16(binary.LittleEndian.Uint16(buff)))
}

func I16(buff []byte) int16 {
	return int16(binary.LittleEndian.Uint16(buff))
}

func Int(buff []byte) int {
	return int(int16(binary.LittleEndian.Uint16(buff)))
}

func I16Tof(buff []byte) float32 {
	return float32(int16(binary.LittleEndian.Uint16(buff)))
}

//DoomStr DOOM string, 8 bytes, filled with 0's
type DoomStr [8]byte

// ToString to string
func (ds *DoomStr) String() string {
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

// WadString converts WAD string to go string
func WadString(buff []byte) string {
	return strings.TrimRight(string(buff), "\x00")
}
