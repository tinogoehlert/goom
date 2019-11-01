package game

type Projectile struct {
	*Movable
	damage   int
	maxRange int
}

func NewProjectile(pos, dir [2]float32, damage, maxRange int) *Projectile {
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
