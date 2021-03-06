package game

import (
	"time"
)

// State state
type State int

const (
	// StateLurking monster is lurking around.
	StateLurking State = 0
	// StateAttacking mosnter attacks something
	StateAttacking State = 1
	// StateHurt monster was hurt
	StateHurt State = 2
	// StateDying monster is dying
	StateDying State = 3
	// StateSplash monster is dying so hard that it splashes into bloody chunks of hellish meat
	StateSplash State = 4
	// StateDead monster was hurt
	StateDead State = 5
)

// Monster A DOOM Monster
type Monster struct {
	*Movable
	health     int
	state      State
	sizeX      float32
	sizeY      float32
	lastTick   time.Time
	lastChange time.Time
	sounds     map[State]string
}

// MonsterFromDef creates monster from definition
func MonsterFromDef(x, y, sx, sy, angle float32, def *MonsterDef) *Monster {
	var m = NewMonster(x, y, angle, def.Sprite)
	m.health = def.Health
	m.sizeX = sx
	m.sizeY = sy
	for k, v := range def.Sounds {
		switch k {
		case "hit":
			m.sounds[StateHurt] = v
		case "die":
			m.sounds[StateDying] = v
		}
	}
	m.sounds[StateSplash] = "DSSLOP"

	for k, v := range def.Animations {
		m.animations[k] = []byte(v)
	}
	m.currentAnimation = m.animations["walk"]
	return m
}

// NewMonster converts ID to name and sequence
func NewMonster(x, y, angle float32, sprite string) *Monster {
	return &Monster{
		Movable: NewMovable(x, y, angle, sprite),
		sounds:  make(map[State]string),
	}
}

// IsCorpse == DEADBEEF?
func (m *Monster) IsCorpse() bool {
	return m.state == StateDead
}

// Update updates the monster's state, making it Lurk after being hit
// or makining it die after being killed.
func (m *Monster) Update() {
	if m.state == StateHurt {
		if time.Now().Sub(m.lastChange) > 150*time.Millisecond {
			m.Lurk()
		}
	}
	if m.state == StateDying || m.state == StateSplash {
		if m.currentFrame == len(m.currentAnimation)-1 {
			m.state = StateDead
			m.freeze = true
			m.hasAngles = false
		}
	}
	m.lastTick = time.Now()
}

// Hit monster got hit by something
func (m *Monster) Hit(damage int, distance float32) State {
	m.health -= damage - (int(distance) / 100)
	if m.state == StateHurt || m.state == StateDying || m.state == StateDead {
		return -1
	}
	if m.health < 0 {
		m.currentAnimation = m.animations["die"]
		m.state = StateDying
		if damage > 20 && distance < 100 {
			m.currentAnimation = m.animations["splash"]
			m.state = StateSplash
		}
		m.currentFrame = 0
		m.lastChange = time.Now()
		return m.state
	}
	m.currentAnimation = m.animations["hurt"]
	m.currentFrame = 0
	m.state = StateHurt
	m.lastChange = time.Now()
	return StateHurt
}

// Lurk GRRR!
func (m *Monster) Lurk() {
	if m.state != StateLurking {
		m.currentAnimation = m.animations["walk"]
		m.currentFrame = 0
		m.state = StateLurking
		m.lastChange = time.Now()
	}
}

// Think What?
func (m *Monster) Think(player *Player) {
	//m.Walk(12, frameTime)
	//m.Turn(12, frameTime)
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
