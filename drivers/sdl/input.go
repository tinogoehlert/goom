package sdl

import (
	"github.com/tinogoehlert/go-sdl2/sdl"
	shared "github.com/tinogoehlert/goom/drivers/pkg"
)

// Input is the SDL input driver.
type Input struct{}

// IsPressed returns if the corresponding Key is pressed.
func (id *Input) IsPressed(k shared.Keycode) bool {
	key, ok := sdlDriversKeyMap[k]
	if !ok {
		return false
	}

	states := sdl.GetKeyboardState()
	scanCode := sdl.GetScancodeFromKey(key)
	return states[scanCode] != 0
}

// GetCursorPos returns the last reported position of the cursor.
func (id *Input) GetCursorPos() (xpos, ypos float64) {
	// TODO
	return
}

// IsMousePressed returns the corresponding Button is pressed.
func (id *Input) IsMousePressed(b shared.MouseButton) bool {
	// TODO
	return false
}

// SetMouseCameraEnabled enables or disables mouse camera control.
func (id *Input) SetMouseCameraEnabled(bool) {
	// TODO
	return
}

var sdlDriversKeyMap = map[shared.Keycode]sdl.Keycode{
	shared.KeyUp:     sdl.K_UP,
	shared.KeyLeft:   sdl.K_LEFT,
	shared.KeyRight:  sdl.K_RIGHT,
	shared.KeyDown:   sdl.K_DOWN,
	shared.KeySpace:  sdl.K_SPACE,
	shared.KeyEnter:  sdl.K_KP_ENTER,
	shared.KeyRShift: sdl.K_RSHIFT,
	shared.KeyLShift: sdl.K_LSHIFT,
	shared.Key0:      sdl.K_0,
	shared.Key1:      sdl.K_1,
	shared.Key2:      sdl.K_2,
	shared.Key3:      sdl.K_3,
	shared.Key4:      sdl.K_4,
	shared.Key5:      sdl.K_5,
	shared.Key6:      sdl.K_6,
	shared.Key7:      sdl.K_7,
	shared.Key8:      sdl.K_8,
	shared.Key9:      sdl.K_9,
	shared.KeyA:      sdl.K_a,
	shared.KeyB:      sdl.K_b,
	shared.KeyC:      sdl.K_c,
	shared.KeyD:      sdl.K_d,
}
