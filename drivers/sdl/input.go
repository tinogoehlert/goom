package sdl

import (
	"github.com/tinogoehlert/go-sdl2/sdl"
)

type inputDriver struct{}

// NewInputDriver creates a new input driver
func NewInputDriver() *inputDriver {
	return &inputDriver{}
}

// IsPressed is keycode pressed? -.^
func (id *inputDriver) IsPressed(keycode uint16) bool {
	states := sdl.GetKeyboardState()
	scanCode := sdl.GetScancodeFromKey(sdl.Keycode(keycode))
	return states[scanCode] != 0
}
