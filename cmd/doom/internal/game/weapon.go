package game

import (
	"math"
	"time"
)

const (
	bobSpeed float64 = 0.4
	bobPow   float64 = 2
)

// Weapon weapon
type Weapon struct {
	Name       string `yaml:"name"`
	Sprite     string `yaml:"sprite"`
	FireSprite string `yaml:"fireSprite"`
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
func (w *Weapon) Fire() {
	if w.state == 0 {
		w.state = 1
		w.lastTick = time.Now()
	}
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
	w.offset[1] += frameTime / 2.5
	if frameTime >= 0 && w.offset[1] > 200 {
		w.pulledDown()
		w.state = 0
	}
	if frameTime < 0 && w.offset[1] <= 0 {
		w.offset[1] = 0
		w.state = 0
	}
}

func (w *Weapon) bobbing(passedTime float32) {
	if w.state != 0 {
		return
	}
	w.bobPhase += float64(passedTime) * math.Pi * 2 * bobSpeed
	x := math.Sin(w.bobPhase/2.0) * bobSpeed * 25.5 * bobPow
	y := math.Abs((math.Sin(w.bobPhase/2.0) * bobSpeed * 12.5 * bobPow))
	w.offset[0] = float32(math.Round(x))
	w.offset[1] = float32(math.Round(y))
}

// NextFrames gets weapon and fire frame. if no fire, value will be 255
func (w *Weapon) NextFrames(frameTime float32) (byte, byte) {
	if w.state == 0 {
		return 'A', 255
	}
	var (
		fire  byte = 255
		frame      = w.Animations["shoot"][w.currentFrame]
	)
	if w.state == 1 && w.currentFrame+1 <= len(w.Animations["fire"]) {
		fire = w.Animations["fire"][w.currentFrame]
	}
	if w.state == 2 {
		w.pull((frameTime))
	}
	if w.state == 3 {
		w.pull(-(frameTime))

	}
	if w.state == 1 && time.Now().Sub(w.lastTick) >= 70*time.Millisecond {
		if w.currentFrame+1 < len(w.Animations["shoot"]) {
			w.currentFrame++
		} else {
			w.currentFrame = 0
			w.state = 0
		}
		w.lastTick = time.Now()
	}
	return frame, fire
}
