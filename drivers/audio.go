package drivers

import (
	"github.com/tinogoehlert/goom/audio/music"
	"github.com/tinogoehlert/goom/audio/sfx"
)

// AudioDriver defines an audio driver by name.
type AudioDriver string

const (
	// NoopAudio is a stubbed audio driver.
	NoopAudio AudioDriver = "noop"

	// SdlAudio is the SDL audio driver.
	SdlAudio AudioDriver = "Sdl"
)

// Audio interface
type Audio interface {
	PlayMusic(m *music.Track) error
	Play(name string) error
	PlayAtPosition(name string, distance float32, angle int16) error
	Close()

	// TestMode sets the driver into test mode with silenced sounds and delays set to 0.
	// Call this function before playing sounds in unit test.
	TestMode()
}

// AudioCreator is a function that creates an Audio driver.
type AudioCreator func(sounds *sfx.Sounds, tempFolder string) (Audio, error)

// AudioDrivers contains all available audio drivers.
var AudioDrivers = make(map[AudioDriver]AudioCreator)
