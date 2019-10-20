package level

import (
	"encoding/binary"

	"github.com/tinogoehlert/goom/internal/utils"

	"github.com/tinogoehlert/goom/geometry"
	"github.com/tinogoehlert/goom/wad"
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
		utils.Int16Tof32(buff[0:2]),
		utils.Int16Tof32(buff[2:4]),
		utils.Int16Tof32(buff[4:6]),
		utils.Int16Tof32(buff[6:8]),
	}
}

// Top top
func (b BBox) Top() float32 { return b[0] }

// Bottom Bottom
func (b BBox) Bottom() float32 { return b[1] }

// Left Left
func (b BBox) Left() float32 { return b[2] }

// Right Right
func (b BBox) Right() float32 { return b[3] }

// PosInBox is position in box?
func (b BBox) PosInBox(x, y float32) bool {
	return x > b[2] && x < b[3] && y > b[1] && y < b[0]
}

// NodeChild can be a Sector or a Node
type NodeChild uint32

func nodeChildI16(v int) NodeChild {
	n := NodeChild(v &^ (1 << 15))
	if v > 3000 {
		n |= (1 << 31)
	}
	return n
}

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
	position  geometry.Vec2
	diagonal  geometry.Vec2
	direction geometry.Vec2
	dirDeg    float32
	RightBBox BBox
	LeftBBox  BBox
	Right     NodeChild
	Left      NodeChild
}

func newNodesFromLump(lump *wad.Lump) ([]Node, error) {
	var (
		nodeCount = len(lump.Data) / nodeSize
		nodes     = make([]Node, nodeCount)
	)
	for i := 0; i < nodeCount; i++ {
		vb := lump.Data[(i * nodeSize) : (i*nodeSize)+nodeSize]
		nodes[i] = Node{
			position:  geometry.V2(utils.Int16Tof32(vb[0:2]), utils.Int16Tof32(vb[2:4])),
			diagonal:  geometry.V2(utils.Int16Tof32(vb[4:6]), utils.Int16Tof32(vb[6:8])),
			RightBBox: NewBBoxFromi16(vb[8:16]),
			LeftBBox:  NewBBoxFromi16(vb[16:24]),
			Right:     nodeChildI16(int(binary.LittleEndian.Uint16(vb[24:26]))),
			Left:      nodeChildI16(int(binary.LittleEndian.Uint16(vb[26:28]))),
		}
		nodes[i].direction = nodes[i].diagonal.CrossVec2().Normalize()
		nodes[i].dirDeg = nodes[i].position.Dot(nodes[i].direction)
	}
	return nodes, nil
}

func newGLNodesFromLump(lump *wad.Lump) ([]Node, error) {
	var (
		nodeCount = lump.Size / glNodeSize
		nodes     = make([]Node, nodeCount)
	)
	for i := 0; i < nodeCount; i++ {
		vb := lump.Data[(i * glNodeSize) : (i*glNodeSize)+glNodeSize]
		nodes[i] = Node{
			position:  geometry.V2(utils.Int16Tof32(vb[0:2]), utils.Int16Tof32(vb[2:4])),
			diagonal:  geometry.V2(utils.Int16Tof32(vb[4:6]), utils.Int16Tof32(vb[6:8])),
			RightBBox: NewBBoxFromi16(vb[8:16]),
			LeftBBox:  NewBBoxFromi16(vb[16:24]),
			Right:     NodeChild(binary.LittleEndian.Uint32(vb[24:28])),
			Left:      NodeChild(binary.LittleEndian.Uint32(vb[28:32])),
		}
		nodes[i].direction = nodes[i].diagonal.CrossVec2().Normalize()
		nodes[i].dirDeg = nodes[i].position.Dot(nodes[i].direction)
	}
	return nodes, nil
}
