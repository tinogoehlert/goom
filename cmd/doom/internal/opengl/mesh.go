package opengl

import (
	"fmt"
	"image"
	"time"

	"github.com/tinogoehlert/goom/graphics"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/go-gl/gl/v2.1/gl"
)

type glTexture struct {
	name   string
	ID     uint32
	width  float32
	height float32
}

type Mesh struct {
	data     []float32
	vao      uint32
	textures map[string]*glTexture
	light    float32
	position mgl32.Vec3
	firstTex *glTexture
	seqTime  time.Time
}

func AddMesh(dst []*Mesh, src *Mesh) []*Mesh {
	if src == nil {
		return dst
	}
	return append(dst, src)
}

func NewMesh(data []float32, light float32, texture graphics.Image, palette graphics.Palette, frameID string) *Mesh {
	if texture == nil {
		return nil
	}
	m := &Mesh{
		data:     data,
		light:    light,
		textures: make(map[string]*glTexture),
		seqTime:  time.Now(),
	}
	m.generateGLBuffers()
	if texture != nil {
		m.generateGLTexture(texture.ToRGBA(palette.Colors), frameID)
		m.firstTex = m.textures[frameID]
	}
	return m
}

func (m *Mesh) Pos() mgl32.Vec3 {
	return m.position
}

func (m *Mesh) DrawMesh(method uint32) {
	if m == nil {
		return
	}
	m.DrawWithTexture(method, m.firstTex)
}

func (m *Mesh) GetTexture(name string) *glTexture {
	if m == nil {
		return nil
	}
	return m.textures[name]
}

func (m *Mesh) DrawWithTexture(method uint32, texture *glTexture) {
	if texture != nil {
		gl.BindTexture(gl.TEXTURE_2D, texture.ID)
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

func (m *Mesh) AddTexture(tex graphics.Image, palette *graphics.Palette, frameID string) {
	m.generateGLTexture(tex.ToRGBA(palette.Colors), frameID)
}

func (m *Mesh) generateGLTexture(tex *image.RGBA, frameID string) {
	if tex == nil {
		fmt.Println("no texture found")
		return
	}
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
		gl.Ptr(tex.Pix))

	m.textures[frameID] = &glTexture{
		width:  float32(tex.Rect.Size().X),
		height: float32(tex.Rect.Size().Y),
		ID:     texID,
		name:   frameID,
	}
}
