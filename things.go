package goom

import (
	"encoding/binary"
	"fmt"
)

const (
	thingSize = 10
)

// Thing - A thing, presented in the Map
type Thing struct {
	X          float32
	Y          float32
	Angle      float32
	Type       int16
	Flags      int16
	SpriteName string
}

func loadThingsFromLump(lump *Lump) ([]Thing, error) {
	if lump.Size%thingSize != 0 {
		return nil, fmt.Errorf("size missmatch")
	}
	var thingCount = lump.Size / thingSize

	things := make([]Thing, thingCount)
	for i := 0; i < thingCount; i++ {
		buff := lump.Data[(i * thingSize) : (i*thingSize)+thingSize]
		things[i].X = float32(int16(binary.LittleEndian.Uint16(buff[0:2])))
		things[i].Y = float32(int16(binary.LittleEndian.Uint16(buff[2:4])))
		things[i].Angle = float32(int16(binary.LittleEndian.Uint16(buff[4:6])))
		things[i].Type = int16(binary.LittleEndian.Uint16(buff[6:8]))
		things[i].Flags = int16(binary.LittleEndian.Uint16(buff[8:10]))
	}
	return things, nil
}
