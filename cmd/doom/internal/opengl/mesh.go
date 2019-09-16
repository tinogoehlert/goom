package opengl

import (
	"fmt"
	"image"
	"time"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/tinogoehlert/goom"

	"github.com/go-gl/gl/v2.1/gl"
)

type Mesh struct {
	data       []float32
	vao        uint32
	textureIDs []uint32
	light      float32
	position   mgl32.Vec3
	currentTex int
	seqTime    time.Time
}

func AddMesh(dst []*Mesh, src *Mesh) []*Mesh {
	if src == nil {
		return dst
	}
	return append(dst, src)
}

func NewMesh(data []float32, light float32, texture goom.DoomTex, palette *goom.Palette) *Mesh {
	if texture == nil {
		return nil
	}
	m := &Mesh{
		data:       data,
		light:      light,
		textureIDs: make([]uint32, 0),
		seqTime:    time.Now(),
	}
	m.generateGLBuffers()
	if texture != nil && palette != nil {
		m.generateGLTexture(texture.ToRGBA(palette))
	}
	return m
}

func (m *Mesh) Pos() mgl32.Vec3 {
	return m.position
}

func (m *Mesh) DrawMesh(method uint32) {
	if len(m.textureIDs) > 0 {
		if time.Now().Sub(m.seqTime) >= 180*time.Millisecond {
			if m.currentTex+1 >= len(m.textureIDs) {
				m.currentTex = 0
			} else {
				m.currentTex++
			}
			m.seqTime = time.Now()
		}
		gl.BindTexture(gl.TEXTURE_2D, m.textureIDs[m.currentTex])
	}
	gl.BindVertexArray(m.vao)
	gl.DrawArrays(method, 0, int32(len(m.data)/5))
}

func (m *Mesh) generateGLBuffers() {
	var fVBO uint32
	gl.GenVertexArrays(1, &m.vao)
	gl.GenBuffers(1, &fVBO)
	// bind the Vertex Array Object first, then bind and set vertex buffer(s), and then configure vertex attributes(s).
	gl.BindVertexArray(m.vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, fVBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(m.data)*4, gl.Ptr(m.data), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)
}

func (m *Mesh) AddTexture(tex goom.DoomTex, palette *goom.Palette) {
	m.generateGLTexture(tex.ToRGBA(palette))
}

func (m *Mesh) generateGLTexture(tex *image.RGBA) {
	if tex == nil {
		fmt.Println("no texture found")
		return
	}
	var texID uint32
	gl.GenTextures(1, &texID)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texID)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
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
		gl.Ptr(tex.Pix))
	m.textureIDs = append(m.textureIDs, texID)
}
