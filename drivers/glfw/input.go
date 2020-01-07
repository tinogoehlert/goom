package glfw

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	shared "github.com/tinogoehlert/goom/drivers/pkg"
)

// Input handling is done via the GLFW window, thus we can implement
// IsPressed as part of the Window.

// IsPressed returns if a key is pressed.
func (in *Window) IsPressed(k shared.Keycode) bool {
	key, ok := glfwDriversKeyMap[k]
	if !ok {
		return false
	}
	return in.window.GetKey(key) == glfw.Press
}

/*
var glfwKeyMap = map[glfw.Key]Keycode{
	glfw.KeyUp:    KeyUp,
	glfw.KeyLeft:  KeyLeft,
	glfw.KeyRight: KeyRight,
	glfw.KeyDown:  KeyDown,
	glfw.KeySpace: KeySpace,
	glfw.KeyEnter: KeyEnter,
	glfw.Key0:     Key0,
	glfw.Key1:     Key1,
	glfw.Key2:     Key2,
	glfw.Key3:     Key3,
	glfw.Key4:     Key4,
	glfw.Key5:     Key5,
	glfw.Key6:     Key6,
	glfw.Key7:     Key7,
	glfw.Key8:     Key8,
	glfw.Key9:     Key9,
	glfw.KeyA:     KeyA,
	glfw.KeyB:     KeyB,
	glfw.KeyC:     KeyC,
	glfw.KeyD:     KeyD,
	glfw.KeyE:     KeyE,
	glfw.KeyF:     KeyF,
	glfw.KeyG:     KeyG,
	glfw.KeyH:     KeyH,
	glfw.KeyI:     KeyI,
	glfw.KeyJ:     KeyJ,
	glfw.KeyK:     KeyK,
	glfw.KeyL:     KeyL,
	glfw.KeyM:     KeyM,
	glfw.KeyN:     KeyN,
	glfw.KeyO:     KeyO,
	glfw.KeyP:     KeyP,
	glfw.KeyQ:     KeyQ,
	glfw.KeyR:     KeyR,
	glfw.KeyS:     KeyS,
	glfw.KeyT:     KeyT,
	glfw.KeyU:     KeyU,
	glfw.KeyV:     KeyV,
	glfw.KeyW:     KeyW,
	glfw.KeyX:     KeyX,
	glfw.KeyY:     KeyY,
	glfw.KeyZ:     KeyZ,
}
*/

var glfwDriversKeyMap = map[shared.Keycode]glfw.Key{
	shared.KeyUp:     glfw.KeyUp,
	shared.KeyLeft:   glfw.KeyLeft,
	shared.KeyRight:  glfw.KeyRight,
	shared.KeyDown:   glfw.KeyDown,
	shared.KeySpace:  glfw.KeySpace,
	shared.KeyEnter:  glfw.KeyEnter,
	shared.KeyRShift: glfw.KeyRightShift,
	shared.KeyLShift: glfw.KeyLeftShift,
	shared.Key0:      glfw.Key0,
	shared.Key1:      glfw.Key1,
	shared.Key2:      glfw.Key2,
	shared.Key3:      glfw.Key3,
	shared.Key4:      glfw.Key4,
	shared.Key5:      glfw.Key5,
	shared.Key6:      glfw.Key6,
	shared.Key7:      glfw.Key7,
	shared.Key8:      glfw.Key8,
	shared.Key9:      glfw.Key9,
	shared.KeyA:      glfw.KeyA,
	shared.KeyB:      glfw.KeyB,
	shared.KeyC:      glfw.KeyC,
	shared.KeyD:      glfw.KeyD,
	shared.KeyE:      glfw.KeyE,
	shared.KeyF:      glfw.KeyF,
	shared.KeyG:      glfw.KeyG,
	shared.KeyH:      glfw.KeyH,
	shared.KeyI:      glfw.KeyI,
	shared.KeyJ:      glfw.KeyJ,
	shared.KeyK:      glfw.KeyK,
	shared.KeyL:      glfw.KeyL,
	shared.KeyM:      glfw.KeyM,
	shared.KeyN:      glfw.KeyN,
	shared.KeyO:      glfw.KeyO,
	shared.KeyP:      glfw.KeyP,
	shared.KeyQ:      glfw.KeyQ,
	shared.KeyR:      glfw.KeyR,
	shared.KeyS:      glfw.KeyS,
	shared.KeyT:      glfw.KeyT,
	shared.KeyU:      glfw.KeyU,
	shared.KeyV:      glfw.KeyV,
	shared.KeyW:      glfw.KeyW,
	shared.KeyX:      glfw.KeyX,
	shared.KeyY:      glfw.KeyY,
	shared.KeyZ:      glfw.KeyZ,
}
