package engine

import (
	"github.com/tinogoehlert/goom/drivers"
	"github.com/tinogoehlert/goom/drivers/glfw"
)

var glfwInitialized bool

// GlfwGraphics is a graphics engine.
type GlfwGraphics struct{}

func (g GlfwGraphics) initialize() error {
	if glfwInitialized {
		return nil
	}

	glfwInitialized = true
	return glfw.InitVideo()
}

// GetWindow creates a new window
func (g GlfwGraphics) GetWindow(title string, width, height int) (drivers.Window, error) {
	if err := g.initialize(); err != nil {
		return nil, err
	}

	return glfw.NewWindow(title, width, height)
}

// GetTime provides the game time
func (g GlfwGraphics) GetTime() float64 {
	return float64(glfw.GetGameTime())
}
