package drivers

import (
	"github.com/tinogoehlert/goom/drivers/sdl"
)

func init() {
	sdlAudio := &sdl.Audio{}
	AudioDrivers[SdlAudio] = sdlAudio
	MusicDrivers[SdlMusic] = sdlAudio
	WindowDrivers[SdlWindow] = &sdl.Window{}
	InputDrivers[SdlInput] = &sdl.Input{}
	TimerFuncs[SdlTimer] = sdl.GetTime
}

// SdlDrivers returns all SDL drivers.
func SdlDrivers() *Drivers {
	return &Drivers{
		Window:  WindowDrivers[SdlWindow],
		Audio:   AudioDrivers[SdlAudio],
		Music:   MusicDrivers[SdlMusic],
		Input:   InputDrivers[SdlInput],
		GetTime: TimerFuncs[SdlTimer],
	}
}
