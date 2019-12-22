package drivers

import "github.com/tinogoehlert/goom/drivers/noop"

func init() {
	noopAudio := &noop.Audio{}
	AudioDrivers[NoopAudio] = noopAudio
	MusicDrivers[NoopMusic] = noopAudio
}

// NoopDrivers returns all Noop drivers.
func NoopDrivers() *Drivers {
	return &Drivers{
		Audio: AudioDrivers[SdlAudio],
		Music: MusicDrivers[SdlMusic],
	}
}
