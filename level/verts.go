package level

import (
	"encoding/binary"

	"github.com/tinogoehlert/goom/geometry"
	"github.com/tinogoehlert/goom/wad"
)

const (
	vertSize     = 4
	vertSizeGlV5 = 8
)

const (
	glMagicV5 = "gNd5"
)

// NewVerticesFromLump loads vertices from Lump
func newVerticesFromLump(lump *wad.Lump) ([]geometry.Vec2, error) {
	var verts []geometry.Vec2
	switch string(lump.Data[0:4]) {
	case glMagicV5:
		verts = readGLVertsV5(lump.Data[4:])
	default:
		verts = readNormalVerts(lump.Data)
	}

	return verts, nil
}

func readNormalVerts(buff []byte) []geometry.Vec2 {
	var (
		vertCount = len(buff) / vertSize
		verts     = make([]geometry.Vec2, vertCount)
	)
	for i := 0; i < vertCount; i++ {
		vb := buff[(i * vertSize) : (i*vertSize)+vertSize]
		verts[i] = geometry.V2(
			float32(int16(binary.LittleEndian.Uint16(vb[0:2]))),
			float32(int16(binary.LittleEndian.Uint16(vb[2:4]))),
		)
	}
	return verts
}

func readGLVertsV5(buff []byte) []geometry.Vec2 {
	var (
		vertCount = (len(buff) / vertSizeGlV5)
		verts     = make([]geometry.Vec2, vertCount)
	)
	for i := 0; i < vertCount; i++ {
		vb := buff[(i * vertSizeGlV5) : (i*vertSizeGlV5)+vertSizeGlV5]

		verts[i] = geometry.V2(
			float32(int32(binary.LittleEndian.Uint32(vb[0:4])))/65536.0,
			float32(int32(binary.LittleEndian.Uint32(vb[4:8])))/65536.0,
		)
	}
	return verts
}
