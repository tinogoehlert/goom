package game

import (
	"log"
	"math"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/tinogoehlert/goom/level"
)

// Thingable is a type that behaves like a thing
type Thingable interface {
	GetID() int
	CollidedWithThing(thing *DoomThing)
	Position() [2]float32
	Direction() [2]float32
	Height() float32
	NextFrame() byte
	SetHeight(height float32)
	CalcAngle(origin mgl32.Vec2) (int, int)
	SpriteName() string
	IsShown() bool
	GetSector() *level.Sector
	SetSector(sector *level.Sector)
}

// DoomThing a doom thing
type DoomThing struct {
	id               int
	angle            float32
	height           float32
	sprite           string
	direction        [2]float32
	position         [2]float32
	animations       map[string][]byte
	currentAnimation []byte
	currentFrame     int
	lastTick         time.Time
	hasAngles        bool
	freeze           bool
	currentSector    *level.Sector
}

// ThingFromDef creates thing from definition
func ThingFromDef(x, y, height, angle float32, def *ThingDef) *DoomThing {
	var m = NewDoomThing(x, y, height, angle, def.Sprite, false)
	m.animations["idle"] = []byte(def.Animation)
	m.currentAnimation = m.animations["idle"]
	m.id = def.ID
	return m
}

// NewDoomThing creates a new DOOM Thing.
func NewDoomThing(x, y, height, angle float32, sprite string, hasAngles bool) *DoomThing {
	return &DoomThing{
		position:   mgl32.Vec2{x, y},
		height:     height,
		angle:      angle,
		sprite:     sprite,
		animations: make(map[string][]byte),
		hasAngles:  hasAngles,
	}
}

func appendDoomThing(dst []Thingable, src Thingable, m *level.Level) []Thingable {
	var ssect, err = m.FindPositionInBsp(level.GLNodesName, src.Position()[0], src.Position()[1])
	if err != nil {
		log.Printf("could not find GLnode for pos %v\n", src.Position())
	} else {
		var sector = m.SectorFromSSect(ssect)
		src.SetSector(sector)
		src.SetHeight(sector.FloorHeight())
	}

	return append(dst, src)
}

// GetSector returns the current sector
func (dt *DoomThing) GetSector() *level.Sector {
	return dt.currentSector
}

// SetSector returns the current sector
func (dt *DoomThing) SetSector(sector *level.Sector) {
	dt.currentSector = sector
}

// GetID returns the THing ID
func (dt *DoomThing) GetID() int {
	return dt.id
}

// CollidedWithThing something collides with thing
func (dt *DoomThing) CollidedWithThing(thing *DoomThing) {}

// Position get XY position
func (dt *DoomThing) Position() [2]float32 {
	return dt.position
}

// Direction get XY direction
func (dt *DoomThing) Direction() [2]float32 {
	return dt.direction
}

// IsShown determines if consumable was consumed
func (dt *DoomThing) IsShown() bool {
	return true
}

// Height get Things height
func (dt *DoomThing) Height() float32 {
	return dt.height
}

// EnterSector thing enters this sector
func (dt *DoomThing) EnterSector(sector *level.Sector) {
	dt.currentSector = sector
}

// SetHeight set things height
func (dt *DoomThing) SetHeight(height float32) {
	dt.height = height
}

// SpriteName get things sprite name
func (dt *DoomThing) SpriteName() string {
	return dt.sprite
}

// NextFrame gets the next frame of the current animation
func (dt *DoomThing) NextFrame() byte {
	if dt.freeze {
		return dt.currentAnimation[dt.currentFrame]
	}
	if time.Now().Sub(dt.lastTick) >= 120*time.Millisecond {
		if dt.currentFrame+1 >= len(dt.currentAnimation) {
			dt.currentFrame = 0
		} else {
			dt.currentFrame++
		}
		dt.lastTick = time.Now()
	}
	return dt.currentAnimation[dt.currentFrame]
}

func (dt *DoomThing) CalcAngle(origin mgl32.Vec2) (int, int) {
	if !dt.hasAngles {
		return 0, 0
	}
	dist := origin.Sub(dt.position)
	angle := mgl32.RadToDeg(float32(math.Atan2(float64(dist.Y()), float64(dist.X())))) - dt.angle

	if angle < 0.0 {
		angle += 360
	}
	switch {
	case (angle >= 292.5 && angle < 337.5):
		return 2, 1
	case (angle >= 22.5 && angle < 67.5):
		return 2, 0
	case (angle >= 67.5 && angle < 112.5):
		return 3, 0
	case (angle >= 112.5 && angle < 157.5):
		return 4, 0
	case (angle >= 157.5 && angle < 202.5):
		return 5, 0
	case (angle >= 202.5 && angle < 247.5):
		return 3, 1
	case (angle >= 247.5 && angle < 292.5):
		return 4, 1
	case (angle >= 337.5 || angle < 22.5):
		return 1, 0
	default:
		return 1, 0
	}
}
