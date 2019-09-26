package game

// DoomThing a doom thing
type DoomThing interface {
	Position() [2]float32
	SpriteName() string
	Sequence() []byte
	Angle() float32
	Height() float32
	SetHeight(float32)
	CurrentFrame(angle int) (byte, byte)
}
