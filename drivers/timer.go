package drivers

type timer string

const (
	// GlfwTimer is the GLFW game time provider
	GlfwTimer timer = "glfw"

	// SdlTimer is the SDL game time provider
	SdlTimer timer = "sdl"
)

// TimerFunction is the function that returns the game time
type TimerFunction func() float64

// Timers contains all available game time providers
var Timers = make(map[timer]TimerFunction)
