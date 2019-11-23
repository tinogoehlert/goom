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
		secsPerUpdate: 1.0 / 60.0,
		inputDrv:      newInputDriver(glfwWin),
	}

	win.fbWidth, win.fbHeight = glfwWin.GetSize()

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

func (w *Window) GetInput() drivers.InputDriver {
	return w.inputDrv
}

// Size Returns the current size of the Window
func (w *Window) Size() (int, int) {
	return w.width, w.height
}

// GetSize returns the current size of the Window
func (w *Window) GetSize() (int, int) {
	return w.fbWidth, w.fbHeight
}

// RunGame runs the game loop
func (w *Window) RunGame(input func(), update func(), render func(float64)) {
	var (
		previous         = glfw.GetTime()
		lag              = float64(0)
		elapsed, current float64
	)

	for !w.window.ShouldClose() {
		current = glfw.GetTime()
		elapsed = current - previous
		previous = current

		lag += elapsed

		w.inputDrv.poll()

		for lag >= w.secsPerUpdate {
			lag -= w.secsPerUpdate
			input()
			update()
		}

		// TODO: This tells the renderer how close we are to the next tick, so if we
		//       are between two ticks we can display (as an example) the movement
		//       of a projectile by and additional 0.8 units instead of fixed 1 unit.
		//       The todo is to implement this in the renderer ^^.
		render(lag / w.secsPerUpdate)
		w.window.SwapBuffers()
	}
}

// Close closes the window
func (w *Window) Close() {
	w.window.Destroy()
}

func GetGameTime() float64 {
	return glfw.GetTime()
}
