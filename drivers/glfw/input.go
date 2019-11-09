package glfw

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/tinogoehlert/goom/drivers"
)

type Keycode int

// InputDriver handles GLFW Input Events
type InputDriver struct {
	keyStates chan drivers.Key
}

// NewInputDriver creates a new GLFW Input Driver
func NewInputDriver() *InputDriver {
	return &InputDriver{
		keyStates: make(chan drivers.Key, 2),
	}
}

// Poll polls for events
func (id *InputDriver) poll(win *Window) {
	glfw.PollEvents()
}

// KeyStates returns a channel where key state changes will be published
func (id *InputDriver) KeyStates() chan drivers.Key {
	return id.keyStates
}

func (id *InputDriver) NormalizeKeyCode(key glfw.Key) drivers.Keycode {
	switch key {
	case glfw.KeyUp:
		return drivers.KeyUp
	case glfw.KeyDown:
		return drivers.KeyDown
	case glfw.KeyLeft:
		return drivers.KeyLeft
	case glfw.KeyRight:
		return drivers.KeyRight
	default:
		return 0
	}
}
