package noop

import (
	"github.com/tinogoehlert/goom/audio/music"
	"github.com/tinogoehlert/goom/audio/sfx"
)

// Audio stub audio driver
type Audio struct{}

func NewAudio(sounds *sfx.Sounds, tempFolder string) (*Audio, error) {
	return &Audio{}, nil
}

// PlayMusic does nothing
func (a Audio) PlayMusic(m *music.Track) error {
	return nil
}

// Play does nothing
func (a Audio) Play(name string) error { return nil }

// PlayAtPosition does nothing
func (a Audio) PlayAtPosition(name string, distance float32, angle int16) error { return nil }

// Close does nothing
func (a Audio) Close() {}
