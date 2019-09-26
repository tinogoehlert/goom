package monsters

import (
	"math"
	"time"

	"github.com/tinogoehlert/goom"
)

// HumanTrooper Former Human Trooper (POSS)
type HumanTrooper struct {
	position     [2]float32
	angle        float32
	sequence     []byte
	spriteName   string
	height       float32
	thing        *goom.Thing
	walkSequence []byte
	currentFrame int
	lastTick     time.Time
}

func NewTrooper(t *goom.Thing, name string) *HumanTrooper {
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

func (ht *HumanTrooper) CurrentFrame(angle int) (byte, byte) {
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
	return calcAngle(angle, int(ht.angle)), ht.walkSequence[ht.currentFrame]
}

func calcAngle(pa, oa int) byte {

	diff := math.Abs(float64((pa - oa)))

	//fmt.Println(diff)
	if diff >= 140 {
		return '1'
	}

	if diff >= 100 && diff <= 140 {
		return '2'
	}

	if diff >= 80 && diff <= 100 {
		return '3'
	}

	if diff >= 40 && diff <= 80 {
		return '4'
	}

	return '5'
}
