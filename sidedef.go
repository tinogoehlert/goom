package goom

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	sidedefSize = 30
)

// SideDef contains the wall texture data for each linedef
type SideDef struct {
	X          int16
	Y          int16
	UpperName  DoomStr
	Lowername  DoomStr
	MiddleName DoomStr
	Sector     int16
}

func newSidesDefFromLump(lump *Lump) ([]SideDef, error) {
	if lump.Size%sidedefSize != 0 {
		return nil, fmt.Errorf("size missmatch")
	}

	var sideCount = lump.Size / sidedefSize
	sides := make([]SideDef, sideCount)
	for i := 0; i < sideCount; i++ {
		r := bytes.NewBuffer(lump.Data[(i * sidedefSize) : (i*sidedefSize)+sidedefSize])
		s := SideDef{}
		if err := binary.Read(r, binary.LittleEndian, &s); err != nil {
			return nil, err
		}
		sides[i] = s
	}
	return sides, nil
}
