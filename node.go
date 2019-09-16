package goom

import (
	"encoding/binary"
)

const (
	nodeSize   = 28
	glNodeSize = 32
)

// BBox describe a rectangle which is the area covered by each of the two child nodes respectively.
// A bounding box consists of four short values (top, bottom, left and right)
// giving the upper and lower bounds of the y coordinate and the lower and upper bounds of the x coordinate (in that order).
type BBox [4]float32

// NewBBoxFromi16 creates new BBox from i16 buffer
func NewBBoxFromi16(buff []byte) BBox {
	return BBox{
		i16Tof(buff[0:2]),
		i16Tof(buff[2:4]),
		i16Tof(buff[4:6]),
		i16Tof(buff[6:8]),
	}
}

func (b BBox) Top() float32    { return b[0] }
func (b BBox) Bottom() float32 { return b[1] }
func (b BBox) Left() float32   { return b[2] }
func (b BBox) Right() float32  { return b[3] }

func (b BBox) PosInBox(x, y float32) bool {
	return x > b[2] && x < b[3] && y > b[1] && y < b[0]
}

// NodeChild can be a Sector or a Node
type NodeChild uint32

// IsSubSector Checks if NodeChild is Subsector
func (n *NodeChild) IsSubSector() bool {
	return (*n & (1 << 31)) > 0
}

// Num gets the number without magic
func (n *NodeChild) Num() uint32 {
	return uint32(*n &^ (1 << 31))
}

// Node The nodes lump constitutes a binary space partition of the level.
type Node struct {
	position  Point
	diff      Point
	RightBBox BBox
	LeftBBox  BBox
	Right     uint32
	Left      uint32
}

func newNodesFromLump(lump *Lump) ([]Node, error) {
	var (
		nodeCount = len(lump.Data) / nodeSize
		nodes     = make([]Node, nodeCount)
	)
	for i := 0; i < nodeCount; i++ {
		vb := lump.Data[(i * nodeSize) : (i*nodeSize)+nodeSize]
		nodes[i] = Node{
			position:  Point{X: i16Tof(vb[0:2]), Y: i16Tof(vb[2:4])},
			diff:      Point{X: i16Tof(vb[4:6]), Y: i16Tof(vb[6:8])},
			RightBBox: NewBBoxFromi16(vb[8:16]),
			LeftBBox:  NewBBoxFromi16(vb[16:24]),
			Right:     uint32(binary.LittleEndian.Uint16(vb[24:26])),
			Left:      uint32(binary.LittleEndian.Uint16(vb[26:28])),
		}
	}
	return nodes, nil
}

func newGLNodesFromLump(lump *Lump) ([]Node, error) {
	var (
		nodeCount = len(lump.Data) / glNodeSize
		nodes     = make([]Node, nodeCount)
	)
	for i := 0; i < nodeCount; i++ {
		vb := lump.Data[(i * glNodeSize) : (i*glNodeSize)+glNodeSize]
		nodes[i] = Node{
			position:  Point{X: i16Tof(vb[0:2]), Y: i16Tof(vb[2:4])},
			diff:      Point{X: i16Tof(vb[4:6]), Y: i16Tof(vb[6:8])},
			RightBBox: NewBBoxFromi16(vb[8:16]),
			LeftBBox:  NewBBoxFromi16(vb[16:24]),
			Right:     uint32(binary.LittleEndian.Uint32(vb[24:28])),
			Left:      uint32(binary.LittleEndian.Uint32(vb[28:32])),
		}
	}
	return nodes, nil
}
