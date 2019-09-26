package game

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

// Player DOOM SLAYER!!
type Player struct {
	position     [2]float32
	dir          [2]float32
	height       float32
	angle        float32
	posChanged   func(xy [2]float32, height float32)
	angleChanged func(angle float32, dir [2]float32)
}

func NewPlayer(x, y, height, angle float32) *Player {
	dy, dx := math.Sincos(float64(angle) * math.Pi / 180)
	p := &Player{
		position: [2]float32{x, y},
		height:   height,
		angle:    angle,
		dir:      mgl32.Vec2{float32(dx), float32(dy)},
	}
	return p
}

// Walk move player x steps back or forth
func (p *Player) Walk(steps float32) {
	p.position[0] += (-p.dir[0] * steps)
	p.position[1] += (p.dir[1] * steps)
}

// Strafe move player x steps to the side
func (p *Player) Strafe(steps float32) {
	p.position[0] += steps
	p.position[1] += (p.dir[1] * steps)
}

// Lift set players height
func (p *Player) Lift(height float32) {
	p.height += height
}

// Turn player
func (p *Player) Turn(angle float32) {
	p.angle += angle
	y, x := math.Sincos(float64(p.angle) * math.Pi / 180)
	p.dir = mgl32.Vec2{float32(x), float32(y)}
}

// Position get XY position
func (p *Player) Position() [2]float32 {
	return p.position
}

// Direction get XY direction
func (p *Player) Direction() [2]float32 {
	return p.dir
}

// Height get players height
func (p *Player) Height() float32 {
	return p.height
}

func (p *Player) SetHeight(height float32) {
	p.height = height
}

// PositionChanged callback for pos change
func (p *Player) PositionChanged(fn func(xy [2]float32, height float32)) {
	p.posChanged = fn
}

// AngleChanged callback for angle change
func (p *Player) AngleChanged(fn func(angle float32, dir [2]float32)) {
	p.angleChanged = fn
}
