package sdl

import (
	"github.com/tinogoehlert/goom/drivers"
	"github.com/veandco/go-sdl2/sdl"
)

// InputDriver handles GLFW Input Events
type InputDriver struct {
	keyStates chan drivers.Key
	win       *sdl.Window
}

// NewInputDriver creates a new GLFW Input Driver
func newInputDriver(win *sdl.Window) *InputDriver {
	return &InputDriver{
		keyStates: make(chan drivers.Key, 2),
		win:       win,
	}
}

// KeyStates returns a channel where key state changes will be published
func (id *InputDriver) KeyStates() chan drivers.Key {
	return id.keyStates
}

// IsPressed is keycode pressed? -.^
func (id *InputDriver) IsPressed(keycode drivers.Keycode) bool {
	return false
}
