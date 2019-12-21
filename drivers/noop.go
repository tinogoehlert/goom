package drivers

import (
	"github.com/tinogoehlert/goom/drivers/noop"
)

func init() {
	AudioDrivers[NoopAudio] = &noop.Audio{}
}
