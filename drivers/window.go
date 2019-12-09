package drivers

// WindowDriver defines a window driver by name.
type WindowDriver string

const (
	// GlfwWindow is the GLFW window driver
	GlfwWindow WindowDriver = "glfw"

	// SdlWindow is the SDl window driver
	SdlWindow WindowDriver = "sdl"
)

// Window interface for the DOOM engine
type Window interface {
	Close()
	GetSize() (width, height int)
	RunGame(input func(), update func(), render func(nextFrameDelta float64))
}

// WindowCreator is a function that creates a Window
type WindowCreator func(title string, width, height int) (Window, error)

// WindowMakers contains all available window creators
var WindowMakers = make(map[WindowDriver]WindowCreator)
