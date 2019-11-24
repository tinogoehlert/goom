package sdl

import (
	"github.com/tinogoehlert/go-sdl2/sdl"
)

var videoInitialized bool
var audioInitialized bool

func initVideo() error {
	if videoInitialized {
		return nil
	}

	err := sdl.InitSubSystem(sdl.INIT_VIDEO)
	if err != nil {
		return err
	}

	videoInitialized = true
	return nil
}

func initAudio() error {
	if audioInitialized {
		return nil
	}

	err := sdl.InitSubSystem(sdl.INIT_AUDIO)
	if err != nil {
		return err
	}

	audioInitialized = true
	return nil
}

// Destroy terminates the SDL driver
func Destroy() {
	sdl.QuitSubSystem(sdl.INIT_VIDEO)
	sdl.QuitSubSystem(sdl.INIT_AUDIO)
}
