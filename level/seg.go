package level

import (
	"bytes"
	"encoding/binary"

	"github.com/tinogoehlert/goom/wad"
)

const (
	segSize   = 12
	glSegSize = 16
)

// Segment  of linedefs, they describe the portion of a linedef that borders
// the subsector that the seg belongs to.
type Segment interface {
	StartVert() uint32
	EndVert() uint32
	LineDef() int16
	Direction() int16
	PartnerSeg() int32
}

// ClassicSegment represents an original DOOM segment
type ClassicSegment struct {
	Start   int16
	End     int16
	Angle   int16
	Linedef int16
	Dir     int16
	Offset  int16
}

//StartVert Start Vertex of the Segment
func (ds *ClassicSegment) StartVert() uint32 { return uint32(ds.Start) }

//EndVert End Vertex of the Segment
func (ds *ClassicSegment) EndVert() uint32 { return uint32(ds.End) }

// LineDef linedef the Segment belongs to
func (ds *ClassicSegment) LineDef() int16 { return ds.Linedef }

// Direction direction of the Segment
func (ds *ClassicSegment) Direction() int16 { return ds.Dir }

// PartnerSeg the partnerseg, not set for classic segs
func (ds *ClassicSegment) PartnerSeg() int32 { return -1 }

// GLSegment GL segment
type GLSegment struct {
	Start      uint32
	End        uint32
	Linedef    int16
	Dir        int16
	Partnerseg int32
}

//StartVert Start Vertex of the Segment
func (gs *GLSegment) StartVert() uint32 { return gs.Start }

//EndVert End Vertex of the Segment
func (gs *GLSegment) EndVert() uint32 { return gs.End }

// LineDef linedef the Segment belongs to
func (gs *GLSegment) LineDef() int16 { return gs.Linedef }

// Direction direction of the Segment
func (gs *GLSegment) Direction() int16 { return gs.Dir }

// PartnerSeg the partnerseg, not set for classic segs
func (gs *GLSegment) PartnerSeg() int32 { return gs.Partnerseg }

func newSegmentsFromLump(lump *wad.Lump) ([]Segment, error) {
	var segs = make([]Segment, (lump.Size)/segSize)
	r := bytes.NewBuffer(lump.Data)

	for i := 0; i < (lump.Size)/segSize; i++ {
		dseg := ClassicSegment{}
		if err := binary.Read(r, binary.LittleEndian, &dseg); err != nil {
			return nil, err
		}
		segs[i] = &dseg
	}
	return segs, nil
}

func newGLSegmentsFromLump(lump *wad.Lump) ([]Segment, error) {
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
