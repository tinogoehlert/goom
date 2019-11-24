package drivers

type windowDriver string

const (
	// GlfwWindow is the GLFW window driver
	GlfwWindow windowDriver = "glfw"

	// SdlWindow is the SDl window driver
	SdlWindow windowDriver = "sdl"
)

// Window interface for the DOOM engine
type Window interface {
	Close()
	GetSize() (width, height int)
	RunGame(input func(), update func(), render func(nextFrameDelta float64))
}

// windowCreator is a function that creates a Window
type windowCreator func(title string, width, height int) (Window, error)

// WindowMakers contains all available window creators
var WindowMakers = make(map[windowDriver]windowCreator)
