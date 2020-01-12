package game

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

// Movable is a movable Thing. eg. Monsters
type Movable struct {
	*DoomThing
	speed       float32
	collisionCB func(me *DoomThing, to mgl32.Vec2) mgl32.Vec2
}

// NewMovable creates a new movable thing
func NewMovable(x, y, angle float32, sprite string) *Movable {
	var m = &Movable{
		DoomThing: NewDoomThing(x, y, angle, sprite, true),
	}
	m.Turn(0)
	m.Pitch(0)
	return m
}

// SetCollision sets the collision callback.
func (m *Movable) SetCollision(cb func(thing *DoomThing, to mgl32.Vec2) mgl32.Vec2) {
	m.collisionCB = cb
}

// Walk move player x steps back or forth
func (m *Movable) Walk(steps float32) {
	if m.collisionCB != nil {
		tmpPos := m.position
		tmpPos[0] += -m.direction[0] * steps
		tmpPos[1] += m.direction[1] * steps
		m.position = m.collisionCB(m.DoomThing, tmpPos)
		return
	}
	m.position[0] += -m.direction[0] * steps
	m.position[1] += m.direction[1] * steps
}

// Strafe move player x steps left or right
func (m *Movable) Strafe(steps float32) {
	if m.collisionCB != nil {
		tmpPos := m.position
		tmpPos[0] += m.direction[1] * steps
		tmpPos[1] += m.direction[0] * steps
		m.position = m.collisionCB(m.DoomThing, tmpPos)
		return
	}
	m.position[0] += m.direction[1] * steps
	m.position[1] += m.direction[0] * steps
}

// Lift set players height
func (m *Movable) Lift(steps float32, timePassed float32) {
	m.height += steps * timePassed
}

// Turn player
func (m *Movable) Turn(angle float32) {
	if angle == 0 {
		return
	}

	m.hAngle += angle

	m.updateDirection()
}

// Pitch is looking up and down
func (m *Movable) Pitch(angle float32) {
	if angle == 0 {
		return
	}

	m.vAngle += angle
	if m.vAngle > 89 {
		m.vAngle = 89
	}

	if m.vAngle < -89 {
		m.vAngle = -89
	}

	m.updateDirection()
}

// ResetPitch makes the player look to the horizon
func (m *Movable) ResetPitch() {
	m.vAngle = 0

	m.updateDirection()
}

func (m *Movable) updateDirection() {
	y, x := math.Sincos(float64(m.hAngle) * math.Pi / 180)
	z := math.Pi * m.vAngle / 90

	m.direction = mgl32.Vec3{float32(x), float32(y), float32(z)}
}
