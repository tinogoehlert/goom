package drivers

import (
	"github.com/tinogoehlert/goom/audio/music"
	"github.com/tinogoehlert/goom/audio/sfx"
)

// AudioDriver defines an audio driver by name.
type AudioDriver string

// MusicDriver defines a music driver by name.
type MusicDriver string

const (
	// NoopAudio is a dummy audio driver.
	NoopAudio AudioDriver = "Noop"

	// NoopMusic is a dummy music driver.
	NoopMusic MusicDriver = "Noop"

	// SdlAudio is the SDL audio driver.
	SdlAudio AudioDriver = "Sdl"

	// SdlMusic is the SDL audio driver.
	SdlMusic MusicDriver = "Sdl"

	// PortMidiMusic is the PortMIDI music driver.
	PortMidiMusic MusicDriver = "PortMidi"
)

// TestDriver inferface
type TestDriver interface {
	// TestMode sets the driver into test mode with silenced sounds and music and zero delays.
	// Call this function before playing music in unit test.
	TestMode()
}

// Music driver interface
type Music interface {
	TestDriver
	// Init starts the driver.
	Init(tracks *music.TrackStore, tempFolder string) error
	PlayMusic(m *music.Track) error
}

// Audio driver interface
type Audio interface {
	TestDriver
	// Init starts the driver.
	Init(sounds *sfx.Sounds) error
	Play(name string) error
	PlayAtPosition(name string, distance float32, angle int16) error
	Close()
}

// AudioDrivers contains all available audio drivers.
var AudioDrivers = make(map[AudioDriver]Audio)

// MusicDrivers contains all available music drivers.
var MusicDrivers = make(map[MusicDriver]Music)
