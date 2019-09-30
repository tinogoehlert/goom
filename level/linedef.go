package level

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/tinogoehlert/goom/wad"
)

const (
	linedefSize = 14
)

// LineDef is what make up the 'shape' (for lack of a better word) of a map.
type LineDef struct {
	Start       int16
	End         int16
	Flags       int16
	SpecialType int16
	SectorTag   int16
	Right       int16
	Left        int16
}

func newLinedefsFromLump(lump *wad.Lump) ([]LineDef, error) {
	if lump.Size%linedefSize != 0 {
		return nil, fmt.Errorf("size missmatch")
	}
	linesDefs := make([]LineDef, lump.Size/linedefSize)
	r := bytes.NewBuffer(lump.Data)
	if err := binary.Read(r, binary.LittleEndian, linesDefs); err != nil {
		return nil, err
	}
	return linesDefs, nil
}
