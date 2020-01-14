package sdl

import (
	"github.com/tinogoehlert/go-sdl2/sdl"
	shared "github.com/tinogoehlert/goom/drivers/pkg"
)

// IsPressed returns if the corresponding Key is pressed.
func (w *Window) IsPressed(k shared.Keycode) bool {
	key, ok := sdlDriversKeyMap[k]
	if !ok {
		return false
	}

	states := sdl.GetKeyboardState()
	scanCode := sdl.GetScancodeFromKey(key)
	return states[scanCode] != 0
}

// GetCursorDelta returns the mouse movement since last call.
func (w *Window) GetCursorDelta() (float64, float64) {
	if !w.mouseCameraEnabled {
		return 0, 0
	}

	x, y, _ := sdl.GetRelativeMouseState()

	return float64(x), float64(y)
}

// IsMousePressed returns the corresponding Button is pressed.
func (w *Window) IsMousePressed(b shared.MouseButton) bool {
	button, ok := sdlMouseButtonMap[b]
	if !ok {
		return false
	}

	_, _, state := sdl.GetMouseState()
	if state&sdl.Button(button) != 0 {
		return true
	}

	return false
}

// SetMouseCameraEnabled enables or disables mouse camera control.
func (w *Window) SetMouseCameraEnabled(en bool) {
	w.mouseCameraEnabled = en
	sdl.SetRelativeMouseMode(en)
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
	shared.KeyE:      sdl.K_e,
	shared.KeyF:      sdl.K_f,
	shared.KeyG:      sdl.K_g,
	shared.KeyH:      sdl.K_h,
	shared.KeyI:      sdl.K_i,
	shared.KeyJ:      sdl.K_j,
	shared.KeyK:      sdl.K_k,
	shared.KeyL:      sdl.K_l,
	shared.KeyM:      sdl.K_m,
	shared.KeyN:      sdl.K_n,
	shared.KeyO:      sdl.K_o,
	shared.KeyP:      sdl.K_p,
	shared.KeyQ:      sdl.K_q,
	shared.KeyR:      sdl.K_r,
	shared.KeyS:      sdl.K_s,
	shared.KeyT:      sdl.K_t,
	shared.KeyU:      sdl.K_u,
	shared.KeyV:      sdl.K_v,
	shared.KeyW:      sdl.K_w,
	shared.KeyX:      sdl.K_x,
	shared.KeyY:      sdl.K_y,
	shared.KeyZ:      sdl.K_z,
	shared.KeyF5:     sdl.K_F5,
	shared.KeyF6:     sdl.K_F6,
	shared.KeyF7:     sdl.K_F7,
	shared.KeyF8:     sdl.K_F8,
}

var sdlMouseButtonMap = map[shared.MouseButton]uint32{
	shared.MouseLeft:   sdl.BUTTON_LEFT,
	shared.MouseMiddle: sdl.BUTTON_MIDDLE,
	shared.MouseRight:  sdl.BUTTON_RIGHT,
	shared.Mouse4:      sdl.BUTTON_X1,
	shared.Mouse5:      sdl.BUTTON_X2,
}
