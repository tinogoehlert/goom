package engine

import "github.com/tinogoehlert/goom/game"

// NullAudio is a silent audio engine
type NullAudio struct{}

func (n NullAudio) initialize(world *game.World) error { return nil }
