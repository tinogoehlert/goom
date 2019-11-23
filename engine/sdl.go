package engine

import (
	"os"

	sdl_native "github.com/tinogoehlert/go-sdl2/sdl"
	"github.com/tinogoehlert/goom/drivers"
	"github.com/tinogoehlert/goom/drivers/sdl"
	"github.com/tinogoehlert/goom/game"
)

// SdlGraphics is a graphics engine.
type SdlGraphics struct{}

var sdlInitialized bool

func (s SdlGraphics) initialize() error {

	if sdlInitialized {
		return nil
	}

	return sdl.InitVideo()
}

// GetWindow creates a new window
func (s SdlGraphics) GetWindow(title string, width, height int) (drivers.Window, error) {
	if err := s.initialize(); err != nil {
		return nil, err
	}

	return sdl.NewGLWindow(title, width, height)
}

// GetTime provides the game time
func (s SdlGraphics) GetTime() float64 {
	return float64(sdl_native.GetTicks())
}

// SdlAudio is an audio engine
type SdlAudio struct{}

// Initialize the sdl audio driver
func (s SdlAudio) Initialize(world *game.World) error {
	os.MkdirAll("temp/music/", 0700)
	sm, err := sdl.NewAudioDriver(world.Data().Sounds, "temp/music")
	if err != nil {
		return err
	}

	world.SetAudioDriver(sm)

	return nil
}
