package drivers

// WindowDriver defines a window driver by name.
type WindowDriver string

// InputDriver defines an input by name.
type InputDriver string

// Timer defins a timer by name.
type Timer string

// AudioDriver defines an audio driver by name.
type AudioDriver string

// MusicDriver defines a music driver by name.
type MusicDriver string

// setup driver maps to allow dynamic access to actual drivers
var (
	WindowDrivers = make(map[WindowDriver]Window)
	AudioDrivers  = make(map[AudioDriver]Audio)
	MusicDrivers  = make(map[MusicDriver]Music)
	InputDrivers  = make(map[InputDriver]Input)
	TimerFuncs    = make(map[Timer]TimerFunc)
)

// Drivers stores the engine drivers.
type Drivers struct {
	Window  Window
	Audio   Audio
	Music   Music
	Input   Input
	GetTime TimerFunc
}

// Define common and specific driver names by name.
const (
	Noop      = "Noop"
	NoopAudio = AudioDriver(Noop)
	NoopMusic = MusicDriver(Noop)

	Glfw       = "Glfw"
	GlfwWindow = WindowDriver(Glfw)
	GlfwTimer  = Timer(Glfw)
	GlfwInput  = InputDriver(Glfw)

	Sdl       = "Sdl"
	SdlWindow = WindowDriver(Sdl)
	SdlTimer  = Timer(Sdl)
	SdlInput  = InputDriver(Sdl)
	SdlAudio  = AudioDriver(Sdl)
	SdlMusic  = MusicDriver(Sdl)

	PortMidiMusic = MusicDriver("PortMidi")
)
