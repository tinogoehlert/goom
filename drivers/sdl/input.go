package sdl

import (
	"github.com/tinogoehlert/go-sdl2/sdl"
	shared "github.com/tinogoehlert/goom/drivers/pkg"
)

// Input is the SDL input driver.
type Input struct{}

// IsPressed returns if a the corresponding Key is pressed.
func (id *Input) IsPressed(k shared.Keycode) bool {
	key, ok := sdlDriversKeyMap[k]
	if !ok {
		return false
	}

	states := sdl.GetKeyboardState()
	scanCode := sdl.GetScancodeFromKey(key)
	return states[scanCode] != 0
}

/*
var sdlKeyMap = map[sdl.Keycode]shared.Keycode{
	sdl.K_UP:       shared.KeyUp,
	sdl.K_LEFT:     shared.KeyLeft,
	sdl.K_RIGHT:    shared.KeyRight,
	sdl.K_DOWN:     shared.KeyDown,
	sdl.K_SPACE:    shared.KeySpace,
	sdl.K_KP_ENTER: shared.KeyEnter,
	sdl.K_0:        shared.Key0,
	sdl.K_1:        shared.Key1,
	sdl.K_2:        shared.Key2,
	sdl.K_3:        shared.Key3,
	sdl.K_4:        shared.Key4,
	sdl.K_5:        shared.Key5,
	sdl.K_6:        shared.Key6,
	sdl.K_7:        shared.Key7,
	sdl.K_8:        shared.Key8,
	sdl.K_9:        shared.Key9,
	sdl.K_a:        shared.KeyA,
	sdl.K_b:        shared.KeyB,
	sdl.K_c:        shared.KeyC,
	sdl.K_d:        shared.KeyD,
	sdl.K_e:        shared.KeyE,
	sdl.K_f:        shared.KeyF,
	sdl.K_g:        shared.KeyG,
	sdl.K_h:        shared.KeyH,
	sdl.K_i:        shared.KeyI,
	sdl.K_j:        shared.KeyJ,
	sdl.K_k:        shared.KeyK,
	sdl.K_l:        shared.KeyL,
	sdl.K_m:        shared.KeyM,
	sdl.K_n:        shared.KeyN,
	sdl.K_o:        shared.KeyO,
	sdl.K_p:        shared.KeyP,
	sdl.K_q:        shared.KeyQ,
	sdl.K_r:        shared.KeyR,
	sdl.K_s:        shared.KeyS,
	sdl.K_t:        shared.KeyT,
	sdl.K_u:        shared.KeyU,
	sdl.K_v:        shared.KeyV,
	sdl.K_w:        shared.KeyW,
	sdl.K_x:        shared.KeyX,
	sdl.K_y:        shared.KeyY,
	sdl.K_z:        shared.KeyZ,
}
*/

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
}
