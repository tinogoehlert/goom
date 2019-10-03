package opengl

import (
	"github.com/go-gl/gl/v2.1/gl"
)

type glSpriter struct {
	vao      uint32
	meshSize int
}

func NewSpriter() *glSpriter {
	verts := []float32{
		-60, -60, 0, 0.0, 1.0,
		-60, 60, 0, 0.0, 0.0,
		60, 60, 0, 1.0, 0.0,

		-60, -60, 0, 0.0, 1.0,
		60, 60, 0, 1.0, 0.0,
		60, -60, 0, 1.0, 1.0,
	}
	gls := &glSpriter{
		meshSize: len(verts) / 5,
	}

	var fVBO uint32
	gl.GenVertexArrays(1, &gls.vao)
	gl.GenBuffers(1, &fVBO)
	// bind the Vertex Array Object first, then bind and set vertex buffer(s), and then configure vertex attributes(s).
	gl.BindVertexArray(gls.vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, fVBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(verts)*4, gl.Ptr(verts), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)
	return gls
}

func (gls *glSpriter) Draw(method uint32, tex *glTexture) {
	gl.BindTexture(gl.TEXTURE_2D, tex.ID)
	gl.BindVertexArray(gls.vao)
	gl.DrawArrays(method, 0, int32(gls.meshSize))
}
