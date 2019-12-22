package drivers

import "github.com/tinogoehlert/goom/drivers/sdl"

func init() {
	noopAudio := &sdl.Audio{}
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
