package monsters

import (
	"math"
	"time"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/tinogoehlert/goom/level"
)

// HumanTrooper Former Human Trooper (POSS)
type HumanTrooper struct {
	position     [2]float32
	angle        float32
	sequence     []byte
	spriteName   string
	height       float32
	thing        *level.Thing
	walkSequence []byte
	currentFrame int
	flipped      int
	lastTick     time.Time
}

func NewTrooper(t *level.Thing, name string) *HumanTrooper {
	return &HumanTrooper{
		position:     [2]float32{t.X, t.Y},
		angle:        t.Angle,
		thing:        t,
		spriteName:   name,
		walkSequence: []byte{'A', 'B', 'C', 'D'},
	}
}

func (ht *HumanTrooper) Position() [2]float32 {
	return ht.position
}

func (ht *HumanTrooper) Angle() float32 {
	return ht.angle
}

func (ht *HumanTrooper) SpriteName() string {
	return ht.spriteName
}

func (ht *HumanTrooper) Height() float32 {
	return ht.height
}

func (ht *HumanTrooper) SetHeight(height float32) {
	ht.height = height
}

func (ht *HumanTrooper) Sequence() []byte {
	return ht.walkSequence
}

func (ht *HumanTrooper) Flipped() int {
	return ht.flipped
}

func (ht *HumanTrooper) CurrentFrame(playerDir [2]float32) (byte, byte) {
	if len(ht.walkSequence) > 0 {
		if time.Now().Sub(ht.lastTick) >= 180*time.Millisecond {
			if ht.currentFrame+1 >= len(ht.walkSequence) {
				ht.currentFrame = 0
			} else {
				ht.currentFrame++
			}
			ht.lastTick = time.Now()
		}
	}

	angle, flipped := calcAngle(playerDir, ht.position, ht.angle)
	ht.flipped = flipped
	return angle, ht.walkSequence[ht.currentFrame]
}

func calcAngle(playerDir, myDir mgl32.Vec2, origin float32) (byte, int) {

	dist := playerDir.Sub(myDir)
	angle := mgl32.RadToDeg(float32(math.Atan2(float64(dist.Y()), float64(dist.X())))) - origin

	if angle < 0.0 {
		angle += 360
	}
	switch {
	case (angle >= 292.5 && angle < 337.5):
		return '2', 1
	case (angle >= 22.5 && angle < 67.5):
		return '2', 0
	case (angle >= 67.5 && angle < 112.5):
		return '3', 0
	case (angle >= 112.5 && angle < 157.5):
		return '4', 0
	case (angle >= 157.5 && angle < 202.5):
		return '5', 0
	case (angle >= 202.5 && angle < 247.5):
		return '3', 1
	case (angle >= 247.5 && angle < 292.5):
		return '4', 1
	case (angle >= 337.5 || angle < 22.5):
		return '1', 0
	default:
		return '1', 0
	}
}
