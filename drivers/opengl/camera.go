package opengl

import (
	"github.com/go-gl/mathgl/mgl32"
)

//Camera OpenGL camera
type Camera struct {
	position  mgl32.Vec2
	velocity  float32
	direction mgl32.Vec3
	angle     float32
	height    float32
	view      mgl32.Mat4
}

// NewCamera creates a new camera at the given position
func NewCamera() *Camera {
	return &Camera{
		view:   mgl32.Ident4(),
		height: 45,
	}
}

// SetCamera set the cam position
func (cam *Camera) SetCamera(pos [2]float32, dir [3]float32, height float32) {
	cam.position = pos
	cam.direction = dir
	cam.height = height
}

// ViewMat4 Get current direction
func (cam *Camera) ViewMat4() mgl32.Mat4 {
	return mgl32.LookAt(
		-cam.position.X(),
		cam.height,
		cam.position.Y(),
		-cam.position.X()+cam.direction.X(),
		cam.height+cam.direction.Z(),
		cam.position.Y()+cam.direction.Y(),
		0.0, 1.0, -0.0)
}
