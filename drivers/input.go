package drivers

type KeyState int
type inputDriver string

const (
	KeyReleased KeyState = 0
	KeyPressed  KeyState = 1
	KeyRepeated KeyState = 2

	// GlfwInput is the GLFW input driver
	GlfwInput inputDriver = "glfw"

	// SdlInput is the SDL input driver
	SdlInput inputDriver = "sdl"
)

type Key struct {
	Keycode Keycode
	State   KeyState
}

type Input interface {
	IsPressed(keycode interface{}) bool
}

type inputProvider func(Window) Input

// InputProviders contains all available input providers
var InputProviders = make(map[inputDriver]inputProvider)
