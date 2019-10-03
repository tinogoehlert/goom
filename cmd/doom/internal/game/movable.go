package game

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

// Movable is a movable Thing. eg. Monsters
type Movable struct {
	*DoomThing
}

// NewMovable creates a new movable thing
func NewMovable(x, y, height, angle float32, sprite string) *Movable {
	var m = &Movable{
		DoomThing: NewDoomThing(x, y, height, angle, sprite),
	}
	return m
}

// Lift set players height
func (m *Movable) Lift(height float32) {
	m.height = height
}

// Turn player
func (m *Movable) Turn(angle float32) {
	m.angle += angle
	y, x := math.Sincos(float64(m.angle) * math.Pi / 180)
	m.direction = mgl32.Vec2{float32(x), float32(y)}
}
