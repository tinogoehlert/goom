package drivers

import (
	shared "github.com/tinogoehlert/goom/drivers/pkg"
)

// KeyState is the state of the key.
type KeyState int

// Key states and input names
const (
	KeyReleased KeyState = 0
	KeyPressed  KeyState = 1
	KeyRepeated KeyState = 2
)

// Key stores key code and press state.
type Key struct {
	Keycode shared.Keycode
	State   KeyState
}

// Input interface
type Input interface {
	IsPressed(shared.Keycode) bool
	GetCursorPos() (xpos, ypos float64)
	IsMousePressed(shared.MouseButton) bool
	SetMouseCameraEnabled(bool)
}
