package game

import (
	"math"
	"time"
)

const (
	bobSpeed float64 = 1.7
	bobPow   float64 = 1.3
)

// Weapon weapon
type Weapon struct {
	Name       string `yaml:"name"`
	Sprite     string `yaml:"egoSprite"`
	FireSprite string `yaml:"fireSprite"`
	Damage     int    `yaml:"damage"`
	Range      int    `yaml:"range"`
	Sound      string `yaml:"sound"`
	FireOffset struct {
		X float32 `yaml:"x"`
		Y float32 `yaml:"y"`
	} `yaml:"fire_offset"`
	ammo         int
	Animations   map[string]string `yaml:"anim"`
	state        int
	lastTick     time.Time
	offset       [2]float32
	pulledDown   func()
	currentFrame int
	bobPhase     float64
}

func (w *Weapon) Offset() [2]float32 {
	return w.offset
}

// Fire fires the weapon
func (w *Weapon) Fire() bool {
	if w.state == 0 {
		w.state = 1
		w.lastTick = time.Now()
		return true
	}
	return false
}

// PutDown puts the weapon down
func (w *Weapon) PutDown(fin func()) {
	w.pulledDown = fin
	if w.state == 0 {
		w.state = 2
	}
}

// PutUp puts the weapon up
func (w *Weapon) PutUp() {
	w.state = 3
}

func (w *Weapon) pull(frameTime float32) {
	w.offset[1] += 450 * frameTime
	if w.offset[1] > 200 {
		w.pulledDown()
		w.state = 0
	}
	if w.offset[1] <= 0 {
		w.offset[1] = 0
		w.state = 0
	}
}

func (w *Weapon) bobbing(passedTime float32) {
	if w.state != 0 {
		return
	}
	w.bobPhase += float64(passedTime) * math.Pi * 0.7 * bobSpeed
	x := math.Cos(w.bobPhase) * bobSpeed * 20.5 * bobPow
	y := math.Abs((math.Sin(w.bobPhase) * bobSpeed * 17 * bobPow))
	w.offset[0] = float32(math.Round(x))
	w.offset[1] = float32(math.Round(y))
}

// NextFrames gets weapon and fire frame. if no fire, value will be 255
func (w *Weapon) NextFrames(frameTime float32) (byte, byte) {
	var (
		anim      = w.Animations["idle"]
		fire byte = 255
	)
	if w.state == 1 && w.currentFrame+1 <= len(w.Animations["fire"]) {
		fire = w.Animations["fire"][w.currentFrame]
	}
	if w.state == 1 {
		anim = w.Animations["shoot"]
	}
	var frame = anim[w.currentFrame]
	if w.state == 2 {
		w.pull((frameTime))
		return frame, fire
	}
	if w.state == 3 {
		w.pull(-(frameTime))
		return frame, fire
	}
	if time.Now().Sub(w.lastTick) >= 100*time.Millisecond {
		if w.currentFrame+1 < len(anim) {
			w.currentFrame++
		} else {
			w.currentFrame = 0
			w.state = 0
		}
		w.lastTick = time.Now()
	}
	return frame, fire
}
