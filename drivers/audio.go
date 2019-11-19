package drivers

import "github.com/tinogoehlert/goom/audio/music"

// AudioDriver interface
type AudioDriver interface {
	PlayMusic(m *music.Track) error
	Play(name string) error
	PlayAtPosition(name string, distance float32, angle int16) error
}

// NOPlayer stub audio driver
type NOPlayer struct{}

// PlayMusic does nothing
func (nop *NOPlayer) PlayMusic(m *music.Track) error { return nil }

// Play does nothing
func (nop *NOPlayer) Play(name string) error { return nil }

// PlayAtPosition does nothing
func (nop *NOPlayer) PlayAtPosition(name string, distance float32, angle int16) error { return nil }
