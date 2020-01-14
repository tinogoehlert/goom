package game

import (
	"time"

	"github.com/tinogoehlert/goom/utils"
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
func NewPlayer(x, y, hAngle, vAngle float32, w *World) *Player {
	p := &Player{
		Movable: &Movable{
			DoomThing: &DoomThing{
				position: [2]float32{x, y},
				hAngle:   hAngle,
				vAngle:   vAngle,
			},
		},
		world:     w,
		weaponBag: make(map[string]*Weapon),
		maxSpeed:  0.5,
	}

	p.SetDirectionAngles(hAngle, vAngle)

	return p
}

// Forward sets the player to a forward moving state using the given speed.
func (p *Player) Forward(speed float32) {
	p.weapon.bobbing()
	p.currSpeed = utils.Clamp(p.currSpeed+speed, -p.maxSpeed, p.maxSpeed)
	p.velocityX += p.currSpeed
}

// Strafe sets the player to a sideways moving state using the given speed.
func (p *Player) Strafe(speed float32) {
	p.weapon.bobbing()
	p.currSpeed = utils.Clamp(p.currSpeed+speed, -p.maxSpeed, p.maxSpeed)
	p.velocityY += p.currSpeed
}

// Lift sets the player to an upward moving state towards the given height.
func (p *Player) Lift(height float32) {
	if p.targetHeight == height {
		return
	}
	p.targetHeight = height
	p.velocityZ = (p.targetHeight - p.height) / 4
}

// Stop sets the player's velocity to 0.
func (p *Player) Stop() {
	p.velocityX = 0
}

// Height returns the player's height.
func (p *Player) Height() float32 {
	// the players viewport height is 41 according to https://eev.ee/blog/2016/10/10/doom-scale/
	return p.DoomThing.height + 41
}

// Update updates all velocities and to deacclerate all types of movement.
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

// FireWeapon --> BOOOM!
func (p *Player) FireWeapon() {
	if p.weapon.Fire() {
		// TODO: play correct sounds for other weapons
		// TODO: reuse playback device instead of naive playback
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
