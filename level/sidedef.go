package level

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/tinogoehlert/goom/internal/utils"
	"github.com/tinogoehlert/goom/wad"
)

const (
	sidedefSize = 30
)

// SideDef contains the wall texture data for each linedef
type SideDef struct {
	X          int16
	Y          int16
	UpperName  utils.DoomStr
	LowerName  utils.DoomStr
	MiddleName utils.DoomStr
	Sector     int16
}

func (s *SideDef) Upper() string {
	return strings.ToUpper(s.UpperName.String())
}

func (s *SideDef) Middle() string {
	return strings.ToUpper(s.MiddleName.String())
}

func (s *SideDef) Lower() string {
	return strings.ToUpper(s.LowerName.String())
}

func newSidesDefFromLump(lump *wad.Lump) ([]SideDef, error) {
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
