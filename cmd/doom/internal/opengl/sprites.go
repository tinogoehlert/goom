package opengl

import (
	"github.com/tinogoehlert/goom"
)

type sprite struct {
	mesh   *Mesh
	width  float32
	height float32
	median float32
}

type spriteList map[string]sprite

func BuildSpritesFromGfx(gfx *goom.Graphics) spriteList {
	sl := make(spriteList)

	for key, sprite := range gfx.GetSprites() {
		sl[key] = makeSprite(&sprite, gfx.Palette(0))
	}

	return sl
}

func makeSprite(gs *goom.Sprite, palette *goom.Palette) sprite {
	s := sprite{}

	var (
		first = gs.GetHeadAngle()[1]
		w     = float32(first.Image().Width())
		h     = float32(first.Image().Height())
	)
	verts := []float32{
		-w, -h, 0, 0.0, 1.0,
		-w, h, 0, 0.0, 0.0,
		w, h, 0, 1.0, 0.0,

		-w, -h, 0, 0.0, 1.0,
		w, h, 0, 1.0, 0.0,
		w, -h, 0, 1.0, 1.0,
	}
	s.mesh = NewMesh(verts, 0, first.Image(), palette)
	for _, frame := range gs.GetHeadAngle() {
		s.mesh.AddTexture(frame.Image(), palette)
	}
	s.height = first.Image().Height()
	s.width = first.Image().Width()
	s.median = s.height / 2
	return s
}
