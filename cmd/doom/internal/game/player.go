package game

import (
	"fmt"
	"math"
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

// Player DOOM SLAYER!!
type Player struct {
	*Movable
	weaponBag map[string]*Weapon
	weapon    *Weapon
	lastTick  time.Time
}

// NewPlayer creates a new player with the given values
func NewPlayer(x, y, height, angle float32) *Player {
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
		weaponBag: make(map[string]*Weapon),
	}
	return p
}

func (p *Player) Walk(steps, passedTime float32) {
	p.Movable.Walk(steps, passedTime)
	p.weapon.bobbing(passedTime)
	p.lastTick = time.Now()
}

// AddWeapon add a new weapon into player's bag or adds ammo
// if he weapon is already in the bag
func (p *Player) AddWeapon(weapon *Weapon) {
	if w, ok := p.weaponBag[weapon.Name]; ok {
		w.ammo += 20
		return
	}
	fmt.Println(len(weapon.Animations["fire"]))
	p.weaponBag[weapon.Name] = weapon
	p.weapon = weapon
}

// CollidedWithThing something collides with thing
func (p *Player) CollidedWithThing(thing *DoomThing) {
	fmt.Println("i collided with", thing.SpriteName())
}

// SwitchWeapon switches
func (p *Player) SwitchWeapon(name string) {
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
