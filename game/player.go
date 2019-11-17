package game

import (
	"fmt"
	"math"
	"time"

	"github.com/tinogoehlert/goom/audio/sfx"
	"github.com/tinogoehlert/goom/utils"

	"github.com/go-gl/mathgl/mgl32"
)

// Player DOOM SLAYER!!
type Player struct {
	*Movable
	weaponBag    map[string]*Weapon
	weapon       *Weapon
	lastTick     time.Time
	world        *World
	velocityX    float32
	velocityY    float32
	velocityZ    float32
	maxSpeed     float32
	currSpeed    float32
	targetHeight float32
}

// NewPlayer creates a new player with the given values
func NewPlayer(x, y, height, angle float32, w *World) *Player {
	dy, dx := math.Sincos(float64(angle) * math.Pi / 180)
	p := &Player{
		Movable: &Movable{
			DoomThing: &DoomThing{
				position:  [2]float32{x, y},
				height:    height,
				angle:     angle,
				direction: mgl32.Vec2{float32(dx), float32(dy)},
			},
		},
		world:     w,
		weaponBag: make(map[string]*Weapon),
		maxSpeed:  0.5,
	}
	return p
}

func (p *Player) Forward(speed float32) {
	p.weapon.bobbing()
	p.currSpeed = utils.Clamp(p.currSpeed+speed, -p.maxSpeed, p.maxSpeed)
	p.velocityX += p.currSpeed
}

func (p *Player) Strafe(speed float32) {
	p.weapon.bobbing()
	p.currSpeed = utils.Clamp(p.currSpeed+speed, -p.maxSpeed, p.maxSpeed)
	p.velocityY += p.currSpeed
}

func (p *Player) Lift(height float32) {
	if p.targetHeight == height {
		return
	}
	p.targetHeight = height
	p.velocityZ = (p.targetHeight - p.height) / 4
}

func (p *Player) Stop() {
	p.velocityX = 0
}

func (p *Player) Height() float32 {
	return p.DoomThing.height + 40
}

func (p *Player) Update() {
	p.velocityX *= 0.90
	p.Movable.Walk(p.velocityX)
	p.velocityY *= 0.90
	p.Movable.Strafe(p.velocityY)
	if p.targetHeight != p.height {
		p.SetHeight(p.height + p.velocityZ)
	}
	p.lastTick = time.Now()
}

// AddWeapon add a new weapon into player's bag or adds ammo
// if he weapon is already in the bag
func (p *Player) AddWeapon(weapon *Weapon) {
	if w, ok := p.weaponBag[weapon.Name]; ok {
		w.ammo += 20
		return
	}
	p.weaponBag[weapon.Name] = weapon
	p.SwitchWeapon(weapon.Name)
}

func (p *Player) FireWeapon() {
	if p.weapon.Fire() {
		fmt.Println("FIRE:", p.weapon.Name)
		// TODO: play correct sounds for other weapons
		// TODO: reuse playback device instead of naive playback
		go sfx.PlaySounds(p.weapon.Sound)
		p.world.spawnShot(p)
	}
}

// SwitchWeapon switches
func (p *Player) SwitchWeapon(name string) {
	if p.weapon == nil {
		p.weapon = p.weaponBag[name]
		p.weapon.PutUp()
	}
	if name == p.weapon.Name {
		return
	}
	if w, ok := p.weaponBag[name]; ok {
		p.weapon.PutDown(func() {
			p.weapon = w
			p.weapon.PutUp()
		})
	}
}

// Weapon gets current weapon
func (p *Player) Weapon() *Weapon {
	return p.weapon
}
