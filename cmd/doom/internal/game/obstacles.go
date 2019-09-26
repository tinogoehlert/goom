package game

import (
	"time"

	"github.com/tinogoehlert/goom"
)

// Obstacle - a static thing in the game
type Obstacle struct {
	position     [2]float32
	angle        float32
	height       float32
	spriteName   string
	sequence     []byte
	currentFrame int
	lastTick     time.Time
}

func (o *Obstacle) Position() [2]float32 {
	return o.position
}

func (o *Obstacle) Angle() float32 {
	return o.angle
}

func (o *Obstacle) Sequence() []byte {
	return o.sequence
}

func (o *Obstacle) SpriteName() string {
	return o.spriteName
}

func (o *Obstacle) Height() float32 {
	return o.height
}

func (o *Obstacle) SetHeight(height float32) {
	o.height = height
}

func (o *Obstacle) CurrentFrame(angle int) (byte, byte) {
	if len(o.sequence) > 0 {
		if time.Now().Sub(o.lastTick) >= 180*time.Millisecond {
			if o.currentFrame+1 >= len(o.sequence) {
				o.currentFrame = 0
			} else {
				o.currentFrame++
			}
			o.lastTick = time.Now()
		}
	}
	return '0', o.sequence[o.currentFrame]
}

func NewObstacle(t *goom.Thing) *Obstacle {
	switch t.Type {
	case 48:
		return buildObstacle("ELEC", []byte{'A'}, t)
	case 2028:
		return buildObstacle("COLU", []byte{'A'}, t)
	case 10:
		return buildObstacle("PLAY", []byte{'W'}, t)
	case 12:
		return buildObstacle("PLAY", []byte{'W'}, t)
	case 15:
		return buildObstacle("PLAY", []byte{'N'}, t)
	case 2035:
		return buildObstacle("BAR1", []byte{'A', 'B'}, t)
	case 55:
		return buildObstacle("SMBT", []byte{'A'}, t)
	case 2019:
		return buildObstacle("ARM2", []byte{'A', 'B'}, t)
	case 2018:
		return buildObstacle("ARM1", []byte{'A', 'B'}, t)
	}

	return nil
}

func buildObstacle(name string, sequence []byte, t *goom.Thing) *Obstacle {
	return &Obstacle{
		spriteName: name,
		sequence:   sequence,
		position:   [2]float32{t.X, t.Y},
		angle:      t.Angle,
	}
}
