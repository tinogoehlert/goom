package drivers

import (
	"github.com/tinogoehlert/goom/audio/music"
	"github.com/tinogoehlert/goom/audio/sfx"
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
	InitMusic(tracks *music.TrackStore, tempFolder string) error
	PlayMusic(m *music.Track) error
	Close()
}

// Audio driver interface
type Audio interface {
	TestDriver
	// Init starts the driver.
	InitAudio(sounds *sfx.Sounds) error
	Play(name string) error
	PlayAtPosition(name string, distance float32, angle int16) error
	Close()
}
