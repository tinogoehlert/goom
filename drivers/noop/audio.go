package noop

import (
	"github.com/tinogoehlert/goom/audio/music"
	"github.com/tinogoehlert/goom/audio/sfx"
)

// Audio is the dummy audio and music driver
type Audio struct{}

// InitAudio does nothing.
func (a *Audio) InitAudio(sounds *sfx.Sounds) error {
	return nil
}

// InitMusic does nothing.
func (a *Audio) InitMusic(tracks *music.TrackStore) error {
	return nil
}

// TestMode does nothing.
func (a *Audio) TestMode() {
}

// PlayMusic does nothing
func (a *Audio) PlayMusic(m *music.Track) error {
	return nil
}

// Play does nothing
func (a *Audio) Play(name string) error { return nil }

// PlayAtPosition does nothing
func (a *Audio) PlayAtPosition(name string, distance float32, angle int16) error { return nil }

// Close does nothing
func (a *Audio) Close() {}
