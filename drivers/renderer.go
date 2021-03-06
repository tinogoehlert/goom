package drivers

import (
	"github.com/tinogoehlert/goom/goom"
	"github.com/tinogoehlert/goom/level"
)

// Renderer interface for the DOOM engine
type Renderer interface {
	LoadLevel(m *level.Level, gd *goom.GameData)
}
