package glfw

import (
	"fmt"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/tinogoehlert/goom/drivers"
)

// Window GLFW implementation for Window
type Window struct {
	window        *glfw.Window
	width         int
	height        int
	fbWidth       int
	fbHeight      int
	secsPerUpdate float64
	fbSizeChanged func(width int, height int)
	inputDrv      *InputDriver
}

// NewWindow creates a new GLFW Window
func NewWindow(title string, width, height int) (*Window, error) {
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfwWin, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create GLFW Window: %s", err.Error())
	}
	glfwWin.MakeContextCurrent()
	win := &Window{
		window:        glfwWin,
		width:         width,
		height:        height,
		secsPerUpdate: 1.0 / 120.0,
		inputDrv:      NewInputDriver(),
	}

	win.fbWidth, win.fbHeight = glfwWin.GetFramebufferSize()

	glfwWin.SetFramebufferSizeCallback(func(w *glfw.Window, width int, height int) {
		win.fbWidth = width
		win.fbHeight = height
		if win.fbSizeChanged != nil {
			win.fbSizeChanged(width, height)
		}
	})

	glfwWin.SetSizeCallback(func(w *glfw.Window, width int, height int) {
		win.width = width
		win.height = height
	})

	return win, nil
}

func (w *Window) Input() drivers.InputDriver {
	return w.inputDrv
}

// Size Returns the current size of the Window
func (w *Window) Size() (int, int) {
	return w.width, w.height
}

// Size Returns the current size of the Window
func (w *Window) FrameBufferSize() (int, int) {
	return w.fbWidth, w.fbHeight
}

func (w *Window) onKeyCallback(win *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	w.inputDrv.keyStates <- drivers.Key{
		Keycode: w.inputDrv.NormalizeKeyCode(key),
		State:   drivers.KeyState(action),
	}
}

// Run runs the window loop
func (w *Window) Run(loop func(elapsed float32)) {
	var (
		previous = glfw.GetTime()
	)

	w.window.SetKeyCallback(w.onKeyCallback)

	for !w.window.ShouldClose() {
		var (
			frameTime = glfw.GetTime() - previous
			waitTime  = w.secsPerUpdate - frameTime
		)
		loop(float32(waitTime))
		if waitTime > 0 {
			glfw.WaitEventsTimeout(waitTime)
		}

		w.window.SwapBuffers()
		glfw.PollEvents()
		previous = glfw.GetTime()
	}
}

// Close closes the window
func (w *Window) Close() {
	w.window.Destroy()
}
