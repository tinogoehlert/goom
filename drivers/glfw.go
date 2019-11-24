package drivers

import (
	"github.com/go-gl/glfw/v3.2/glfw"

	glfw_internal "github.com/tinogoehlert/goom/drivers/glfw"
)

func init() {
	WindowMakers[GlfwWindow] = newGlfwWindow
	Timers[GlfwTimer] = glfw_internal.GetTime
	InputProviders[GlfwInput] = newGlfwInput
}

func newGlfwWindow(title string, width, height int) (Window, error) {
	win, err := glfw_internal.NewWindow(title, width, height)

	return Window(win), err
}

func newGlfwInput(win Window) Input {
	drv := glfw_internal.NewInputDriver(
		win.(*glfw_internal.Window),
		mapGlfwKey,
	)
	return Input(drv)
}

func mapGlfwKey(keycode uint16) (glfw.Key, bool) {
	key, ok := glfwDriversKeyMap[Keycode(keycode)]
	return key, ok
}

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

var glfwDriversKeyMap = map[Keycode]glfw.Key{
	KeyUp:     glfw.KeyUp,
	KeyLeft:   glfw.KeyLeft,
	KeyRight:  glfw.KeyRight,
	KeyDown:   glfw.KeyDown,
	KeySpace:  glfw.KeySpace,
	KeyEnter:  glfw.KeyEnter,
	KeyRShift: glfw.KeyRightShift,
	KeyLShift: glfw.KeyLeftShift,
	Key0:      glfw.Key0,
	Key1:      glfw.Key1,
	Key2:      glfw.Key2,
	Key3:      glfw.Key3,
	Key4:      glfw.Key4,
	Key5:      glfw.Key5,
	Key6:      glfw.Key6,
	Key7:      glfw.Key7,
	Key8:      glfw.Key8,
	Key9:      glfw.Key9,
	KeyA:      glfw.KeyA,
	KeyB:      glfw.KeyB,
	KeyC:      glfw.KeyC,
	KeyD:      glfw.KeyD,
	KeyE:      glfw.KeyE,
	KeyF:      glfw.KeyF,
	KeyG:      glfw.KeyG,
	KeyH:      glfw.KeyH,
	KeyI:      glfw.KeyI,
	KeyJ:      glfw.KeyJ,
	KeyK:      glfw.KeyK,
	KeyL:      glfw.KeyL,
	KeyM:      glfw.KeyM,
	KeyN:      glfw.KeyN,
	KeyO:      glfw.KeyO,
	KeyP:      glfw.KeyP,
	KeyQ:      glfw.KeyQ,
	KeyR:      glfw.KeyR,
	KeyS:      glfw.KeyS,
	KeyT:      glfw.KeyT,
	KeyU:      glfw.KeyU,
	KeyV:      glfw.KeyV,
	KeyW:      glfw.KeyW,
	KeyX:      glfw.KeyX,
	KeyY:      glfw.KeyY,
	KeyZ:      glfw.KeyZ,
}
