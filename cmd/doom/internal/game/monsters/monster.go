package monsters

import (
	"github.com/tinogoehlert/goom"
	"github.com/tinogoehlert/goom/cmd/doom/internal/game"
)

// Monster A DOOM Monster
type Monster interface {
	game.DoomThing
}

// NewMonster converts ID to name and sequence
func NewMonster(t *goom.Thing) Monster {
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
	case 3004:
		return NewTrooper(t, "POSS")
	case 9:
		return NewTrooper(t, "SPOS")
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
