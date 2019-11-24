package sdl

import (
	"github.com/tinogoehlert/go-sdl2/sdl"
)

// GetTime provides the game time
func GetTime() float64 {
	return float64(sdl.GetTicks())
}
