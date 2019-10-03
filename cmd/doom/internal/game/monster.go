package game

// State state
type State int

const (
	// StateLurking monster is lurking around.
	StateLurking State = 0
	// StateAttacking mosnter attacks something
	StateAttacking State = 1
)

// Monster A DOOM Monster
type Monster struct {
	*Movable
	health int
	state  int
}

// MonsterFromDef creates monster from definition
func MonsterFromDef(x, y, height, angle float32, def *MonsterDef) *Monster {
	var m = NewMonster(x, y, height, angle, def.Sprite)
	for k, v := range def.Animations {
		m.animations[k] = []byte(v)
	}
	m.currentAnimation = m.animations["walk"]
	return m
}

// NewMonster converts ID to name and sequence
func NewMonster(x, y, height, angle float32, sprite string) *Monster {
	return &Monster{
		Movable: NewMovable(x, y, height, angle, sprite),
	}
}

/*
// NewMonster converts ID to name and sequence
func NewMonster(t *level.Thing) Monster {
	switch t.Type {
	case 68:
		return NewDummyMonster(t, "BSPI")
	case 64:
		return NewDummyMonster(t, "VILE")
	case 3003:
		return NewDummyMonster(t, "BOSS")
	case 3005:
		return NewDummyMonster(t, "HEAD")
	case 65:
		return NewDummyMonster(t, "CPOS")
	case 72:
		return NewDummyMonster(t, "KEEN")
	case 16:
		return NewDummyMonster(t, "CYBR")
	case 3002:
		return NewDummyMonster(t, "SARG")
	case 69:
		return NewDummyMonster(t, "BOS2")
	case 3001:
		return NewDummyMonster(t, "TROO")
	case 3006:
		return NewDummyMonster(t, "SKUL")
	case 67:
		return NewDummyMonster(t, "FATT")
	case 71:
		return NewDummyMonster(t, "PAIN")
	case 66:
		return NewDummyMonster(t, "SKEL")
	case 7:
		return NewDummyMonster(t, "SPID")
	case 84:
		return NewDummyMonster(t, "SSWV")
	}

	return nil
}
*/
