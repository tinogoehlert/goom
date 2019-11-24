package drivers

import (
	"github.com/tinogoehlert/goom/audio/sfx"
	"github.com/tinogoehlert/goom/drivers/noop"
)

func init() {
	AudioDrivers[NoopAudio] = newNoopAudio
}

func newNoopAudio(sounds *sfx.Sounds, tempFolder string) (Audio, error) {
	audio, err := noop.NewAudio(sounds, tempFolder)
	return Audio(audio), err
}
