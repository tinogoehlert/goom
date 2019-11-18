package drivers

import "github.com/tinogoehlert/goom/audio/music"

type AudioDriver interface {
	PlayMusic(m *music.Track) error
	Play(name string) error
	PlayAtPosition(name string, distance float32, angle int16) error
}
