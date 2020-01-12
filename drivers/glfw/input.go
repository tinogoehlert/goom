package glfw

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	shared "github.com/tinogoehlert/goom/drivers/pkg"
)

// IsPressed returns if a key is pressed.
func (in *Window) IsPressed(k shared.Keycode) bool {
	key, ok := glfwDriversKeyMap[k]
	if !ok {
		return false
	}
	return in.window.GetKey(key) == glfw.Press
}

// GetCursorPos returns the last reported position of the cursor.
func (in *Window) GetCursorPos() (xpos, ypos float64) {
	x, y := in.window.GetCursorPos()

	// // reset the mouse pos (only useful for camera positioning to not "run out of space")
	// // needs to be somewhere else once we actually want to use a cursor
	// in.window.SetCursorPos(0, 0)

	return x, y
}

// IsMousePressed returns the corresponding Button is pressed.
func (in *Window) IsMousePressed(b shared.MouseButton) bool {
	btn, ok := glfwMouseButtonMap[b]
	if !ok {
		return false
	}
	return in.window.GetMouseButton(btn) == glfw.Press
}

// SetMouseCameraEnabled enables or disables mouse camera control.
func (in *Window) SetMouseCameraEnabled(en bool) {
	if !en {
		in.window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
		return
	}

	in.window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	if glfw.RawMouseMotionSupported() {
		in.window.SetInputMode(glfw.RawMouseMotion, glfw.True)
	}
}

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
	shared.KeyF5:     glfw.KeyF5,
	shared.KeyF6:     glfw.KeyF6,
}

var glfwMouseButtonMap = map[shared.MouseButton]glfw.MouseButton{
	shared.MouseLeft:   glfw.MouseButton1,
	shared.MouseMiddle: glfw.MouseButton3,
	shared.MouseRight:  glfw.MouseButton2,
	shared.Mouse4:      glfw.MouseButton4,
	shared.Mouse5:      glfw.MouseButton5,
	shared.Mouse6:      glfw.MouseButton6,
	shared.Mouse7:      glfw.MouseButton7,
	shared.Mouse8:      glfw.MouseButton8,
}
