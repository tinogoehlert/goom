package glfw

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/tinogoehlert/goom/drivers"
)

// InputDriver handles GLFW Input Events
type InputDriver struct {
	keyStates chan drivers.Key
	win       *glfw.Window
}

// NewInputDriver creates a new GLFW Input Driver
func newInputDriver(win *glfw.Window) *InputDriver {
	return &InputDriver{
		keyStates: make(chan drivers.Key, 2),
		win:       win,
	}
}

// Poll polls for events
func (id *InputDriver) poll() {
	glfw.PollEvents()
}

// KeyStates returns a channel where key state changes will be published
func (id *InputDriver) KeyStates() chan drivers.Key {
	return id.keyStates
}

// IsPressed is keycode pressed? -.^
func (id *InputDriver) IsPressed(keycode drivers.Keycode) bool {
	if k, ok := driversKeyMap[keycode]; ok {
		return id.win.GetKey(k) == glfw.Press
	}
	return false
}

// NormalizeKeyCode returns drivers Keycode for equal glfw keycode
func (id *InputDriver) NormalizeKeyCode(key glfw.Key) drivers.Keycode {
	if normKey, ok := glfwKeyMap[key]; ok {
		return normKey
	}
	return 0
}
