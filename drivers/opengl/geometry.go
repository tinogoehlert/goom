package opengl

import (
	"time"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/go-gl/gl/v2.1/gl"
)

type glWorldGeometry struct {
	data     []float32
	vao      uint32
	texture  []*glTexture
	light    float32
	position mgl32.Vec3
	seqTime  time.Time
	isSky    bool
}

func addGlWorldutils(dst []*glWorldGeometry, src *glWorldGeometry) []*glWorldGeometry {
	if src == nil {
		return dst
	}
	return append(dst, src)
}

func newGlWorldutils(data []float32, light float32, texture []*glTexture) *glWorldGeometry {
	m := &glWorldGeometry{
		data:    data,
		light:   light,
		texture: texture,
		seqTime: time.Now(),
	}
	m.generateGLBuffers()
	return m
}

func (m *glWorldGeometry) pos() mgl32.Vec3 {
	return m.position
}

func (m *glWorldGeometry) Draw(method uint32) {

	if m == nil {
		return
	}
	gl.BindTexture(gl.TEXTURE_2D, m.texture[0].ID)
	gl.BindVertexArray(m.vao)
	gl.DrawArrays(method, 0, int32(len(m.data)/5))
}

func (m *glWorldGeometry) generateGLBuffers() {
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
