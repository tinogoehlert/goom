package monsters

import (
	"github.com/tinogoehlert/goom/level"
)

var (
	walkingSeq = []byte{'A', 'B', 'C', 'D'}
)

// DummyMonster dummy stub
type DummyMonster struct {
	position   [2]float32
	angle      float32
	sequence   []byte
	spriteName string
	thing      *level.Thing
	height     float32
}

func NewDummyMonster(t *level.Thing, name string) *DummyMonster {
	return &DummyMonster{
		position:   [2]float32{t.X, t.Y},
		angle:      t.Angle,
		thing:      t,
		spriteName: name,
	}
}

func (ht *DummyMonster) Position() [2]float32 {
	return ht.position
}

func (ht *DummyMonster) Angle() float32 {
	return ht.angle
}

func (ht *DummyMonster) Sequence() []byte {
	return walkingSeq
}

func (ht *DummyMonster) SpriteName() string {
	return ht.spriteName
}

func (ht *DummyMonster) Height() float32 {
	return ht.height
}

func (ht *DummyMonster) SetHeight(height float32) {
	ht.height = height
}

func (ht *DummyMonster) Flipped() int {
	return 0
}

func (ht *DummyMonster) CurrentFrame(dir [2]float32) (byte, byte) {
	return '1', walkingSeq[0]
}
