package opengl

import (
	"fmt"

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

	first := gs.FirstFrame()
	if first == nil {
		fmt.Println(gs.Name)
	}
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
	gs.Frames(func(frame *goom.SpriteFrame) {
		s.mesh.AddTexture(frame.Image(), palette, frame.Name())
	})
	s.height = first.Image().Height()
	s.width = first.Image().Width()
	s.median = s.height / 2
	return s
}
