package goom

import (
	"bytes"
	"encoding/binary"
)

const (
	segSize   = 12
	glSegSize = 16
)

// Segment  of linedefs, they describe the portion of a linedef that borders
// the subsector that the seg belongs to.
type Segment interface {
	GetStartVert() uint32
	GetEndVert() uint32
	GetLineDef() int16
	GetDirection() int16
	GetPartnerSeg() int32
}

// DoomSegment represents an original DOOM segment
type DoomSegment struct {
	Start     int16
	End       int16
	Angle     int16
	Linedef   int16
	Direction int16
	Offset    int16
}

func (ds *DoomSegment) GetStartVert() uint32 { return uint32(ds.Start) }
func (ds *DoomSegment) GetEndVert() uint32   { return uint32(ds.End) }
func (ds *DoomSegment) GetLineDef() int16    { return ds.Linedef }
func (ds *DoomSegment) GetDirection() int16  { return ds.Direction }
func (gs *DoomSegment) GetPartnerSeg() int32 { return -1 }

// GLSegment GL segment
type GLSegment struct {
	Start      uint32
	End        uint32
	Linedef    int16
	Direction  int16
	PartnerSeg int32
}

func (gs *GLSegment) GetStartVert() uint32 { return gs.Start }
func (gs *GLSegment) GetEndVert() uint32   { return gs.End }
func (gs *GLSegment) GetLineDef() int16    { return gs.Linedef }
func (gs *GLSegment) GetDirection() int16  { return gs.Direction }
func (gs *GLSegment) GetPartnerSeg() int32 { return gs.PartnerSeg }

func newSegmentsFromLump(lump *Lump) ([]Segment, error) {
	var segs = make([]Segment, (lump.Size)/segSize)
	r := bytes.NewBuffer(lump.Data)

	for i := 0; i < (lump.Size)/segSize; i++ {
		dseg := DoomSegment{}
		if err := binary.Read(r, binary.LittleEndian, &dseg); err != nil {
			return nil, err
		}
		segs[i] = &dseg
	}
	return segs, nil
}

func newGLSegmentsFromLump(lump *Lump) ([]Segment, error) {
	var segs = make([]Segment, (lump.Size)/glSegSize)
	r := bytes.NewBuffer(lump.Data)
	for i := 0; i < (lump.Size)/glSegSize; i++ {
		gseg := GLSegment{}
		if err := binary.Read(r, binary.LittleEndian, &gseg); err != nil {
			return nil, err
		}
		segs[i] = &gseg
	}
	return segs, nil
}
