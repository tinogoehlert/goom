package glfw

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

// InputDriver handles GLFW Input Events
type InputDriver struct {
	win    *glfw.Window
	mapper func(keycode interface{}) (glfw.Key, bool)
}

// NewInputDriver creates a new GLFW Input Driver
func NewInputDriver(win *Window, mapper func(keycode interface{}) (glfw.Key, bool)) *InputDriver {
	return &InputDriver{
		win:    win.window,
		mapper: mapper,
	}
}

// Poll polls for events
func (id *InputDriver) poll() {
	glfw.PollEvents()
}

// IsPressed is keycode pressed? -.^
func (id *InputDriver) IsPressed(keycode interface{}) bool {
	key, ok := id.mapper(keycode)
	if !ok {
		return false
	}

	return id.win.GetKey(key) == glfw.Press
}
