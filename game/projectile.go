package game

// Projectile is a Movable projectile with a specific range and damage.
type Projectile struct {
	*Movable
	damage   int
	maxRange int
}

// NewProjectile returns a new projectile.
func NewProjectile(pos [2]float32, dir [3]float32, damage, maxRange int) *Projectile {
	return &Projectile{
		Movable: &Movable{
			DoomThing: &DoomThing{
				position:  pos,
				direction: dir,
			},
		},
		maxRange: maxRange,
		damage:   damage,
	}
}
