package drivers

import (
	"github.com/tinogoehlert/goom/audio/music"
	"github.com/tinogoehlert/goom/audio/sfx"
)

type audioDriver string

const (
	// NoopAudio is a stubbed audio driver
	NoopAudio audioDriver = "noop"

	// SdlAudio is the SDL audio driver
	SdlAudio audioDriver = "Sdl"
)

// Audio interface
type Audio interface {
	PlayMusic(m *music.Track) error
	Play(name string) error
	PlayAtPosition(name string, distance float32, angle int16) error
	Close()
}

// audioCreator is a function that creates a Audio driver
type audioCreator func(sounds *sfx.Sounds, tempFolder string) (Audio, error)

// AudioDrivers contains all available audio drivers
var AudioDrivers = make(map[audioDriver]audioCreator)
