package opengl

import (
	"image"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/tinogoehlert/goom/graphics"
)

type glTexture struct {
	image graphics.Image
	ID    uint32
}

type glTextureStore map[string][]*glTexture

func newGLTextureStore() glTextureStore {
	var ts = make(glTextureStore)
	ts.initTexture("null", 1)
	//ts["null"][0] = makeNoTexture()
	return ts
}

func (ts glTextureStore) initTexture(name string, count int) {
	ts[name] = make([]*glTexture, count)
}

func (ts glTextureStore) addTexture(name string, idx int, img graphics.Image) {
	ts[name][idx] = makeGLTexture(img)
}

func (ts glTextureStore) Get(name string, idx int) *glTexture {
	if tex, ok := ts[name]; ok {
		if tex[idx] == nil && idx != 0 {
			return tex[0]
		}
		return ts[name][idx]
	}
	return ts["null"][0]
}

func makeGLTexture(img graphics.Image) *glTexture {
	if img == nil {
		return nil
	}
	tex := img.ToRGBA(graphics.DefaultPalette().Colors)

	return &glTexture{
		ID:    genGLTexture(tex),
		image: img,
	}
}

func makeNoTexture() *glTexture {
	tex := image.NewRGBA(image.Rect(0, 0, 64, 64))
	return &glTexture{
		ID:    genGLTexture(tex),
		image: nil,
	}
}

func genGLTexture(tex *image.RGBA) uint32 {
	var texID uint32
	gl.GenTextures(1, &texID)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texID)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(tex.Rect.Size().X),
		int32(tex.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(tex.Pix),
	)
	return texID
}
