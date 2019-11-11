package sdl

import (
	"github.com/veandco/go-sdl2/sdl"
)

// Init inits SDL
func Init() error {
	return sdl.Init(sdl.INIT_EVERYTHING)
}

// Destroy destroys sdl
func Destroy() {
	sdl.Quit()
}
