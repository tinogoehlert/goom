package level

import (
	"encoding/binary"

	"github.com/tinogoehlert/goom/wad"
)

const (
	vertSize     = 4
	vertSizeGlV5 = 8
)

const (
	glMagicV5 = "gNd5"
)

// Vertex XY coordinates
type Vertex struct {
	x float32
	y float32
}

// X return X coord
func (v *Vertex) X() float32 { return v.x }

// Y return Y coord
func (v *Vertex) Y() float32 { return v.y }

// NewVerticesFromLump loads vertices from Lump
func newVerticesFromLump(lump *wad.Lump) ([]Vertex, error) {
	var verts []Vertex
	switch string(lump.Data[0:4]) {
	case glMagicV5:
		verts = readGLVertsV5(lump.Data[4:])
	default:
		verts = readNormalVerts(lump.Data)
	}

	return verts, nil
}

func readNormalVerts(buff []byte) []Vertex {
	var (
		vertCount = len(buff) / vertSize
		verts     = make([]Vertex, vertCount)
	)
	for i := 0; i < vertCount; i++ {
		vb := buff[(i * vertSize) : (i*vertSize)+vertSize]
		verts[i] = Vertex{
			x: float32(int16(binary.LittleEndian.Uint16(vb[0:2]))),
			y: float32(int16(binary.LittleEndian.Uint16(vb[2:4]))),
		}
	}
	return verts
}

func readGLVertsV5(buff []byte) []Vertex {
	var (
		vertCount = (len(buff) / vertSizeGlV5)
		verts     = make([]Vertex, vertCount)
	)
	for i := 0; i < vertCount; i++ {
		vb := buff[(i * vertSizeGlV5) : (i*vertSizeGlV5)+vertSizeGlV5]

		verts[i] = Vertex{
			x: float32(int32(binary.LittleEndian.Uint32(vb[0:4]))) / 65536.0,
			y: float32(int32(binary.LittleEndian.Uint32(vb[4:8]))) / 65536.0,
		}
	}
	return verts
}
