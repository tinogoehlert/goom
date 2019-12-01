package sdl

import (
	"fmt"

	"github.com/tinogoehlert/go-sdl2/sdl"
)

type inputDriver struct {
	mapper func(keycode interface{}) (sdl.Keycode, bool)
}

// NewInputDriver creates a new input driver
func NewInputDriver(mapper func(keycode interface{}) (sdl.Keycode, bool)) *inputDriver {
	return &inputDriver{mapper: mapper}
}

// IsPressed is keycode pressed? -.^
func (id *inputDriver) IsPressed(keycode interface{}) bool {
	key, ok := id.mapper(keycode)
	if !ok {
		return false
	}

	states := sdl.GetKeyboardState()
	fmt.Printf("%v\n", states)
	scanCode := sdl.GetScancodeFromKey(key)
	return states[scanCode] != 0
}
