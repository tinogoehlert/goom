package opengl

import (
	"github.com/tinogoehlert/goom/graphics"
)

type sprite struct {
	mesh   *Mesh
	width  float32
	height float32
	median float32
}

type spriteList map[string]sprite

func BuildSpritesFromGfx(sprites graphics.SpriteStore, palette graphics.Palette) spriteList {
	sl := make(spriteList)

	for key, sprite := range sprites {
		sl[key] = makeSprite(&sprite, palette)
	}

	return sl
}

func makeSprite(gs *graphics.Sprite, palette graphics.Palette) sprite {
	s := sprite{}

	first := gs.FirstFrame()

	var (
		w = float32(first.Image().Width())
		h = float32(first.Image().Height())
	)
	verts := []float32{
		-w, -h, 0, 0.0, 1.0,
		-w, h, 0, 0.0, 0.0,
		w, h, 0, 1.0, 0.0,

		-w, -h, 0, 0.0, 1.0,
		w, h, 0, 1.0, 0.0,
		w, -h, 0, 1.0, 1.0,
	}
	s.mesh = NewMesh(verts, 0, first.Image(), palette, first.Name())
	gs.Frames(func(frame *graphics.SpriteFrame) {
		s.mesh.AddTexture(frame.Image(), &palette, frame.Name())
	})
	s.height = float32(first.Image().Height())
	s.width = float32(first.Image().Width())
	s.median = s.height / 2
	return s
}
