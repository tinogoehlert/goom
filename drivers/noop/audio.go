package noop

import (
	"github.com/tinogoehlert/goom/audio/music"
	"github.com/tinogoehlert/goom/audio/sfx"
	"github.com/tinogoehlert/goom/drivers"
)

// Audio stub audio driver
type Audio struct{}

func init() {
	drivers.AudioDrivers[drivers.NoopAudio] = newAudio
}

func newAudio(sounds *sfx.Sounds, tempFolder string) (drivers.Audio, error) {
	return drivers.Audio(Audio{}), nil
}

// PlayMusic does nothing
func (a Audio) PlayMusic(m *music.Track) error { return nil }

// Play does nothing
func (a Audio) Play(name string) error { return nil }

// PlayAtPosition does nothing
func (a Audio) PlayAtPosition(name string, distance float32, angle int16) error { return nil }

// Close does nothing
func (a Audio) Close() {}
