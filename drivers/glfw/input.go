package glfw

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

// InputDriver handles GLFW Input Events
type InputDriver struct {
	win *glfw.Window
}

// NewInputDriver creates a new GLFW Input Driver
func NewInputDriver(win *Window) *InputDriver {
	return &InputDriver{
		win: win.window,
	}
}

// Poll polls for events
func (id *InputDriver) poll() {
	glfw.PollEvents()
}

// IsPressed is keycode pressed? -.^
func (id *InputDriver) IsPressed(keycode uint16) bool {
	if keycode < 31 {
		return false
	}

	return id.win.GetKey(glfw.Key(keycode)) == glfw.Press
}
