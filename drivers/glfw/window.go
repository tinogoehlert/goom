package glfw

import (
	"fmt"

	"github.com/go-gl/glfw/v3.3/glfw"
)

// Window GLFW implementation for Window
type Window struct {
	window             *glfw.Window
	width              int
	height             int
	fbWidth            int
	fbHeight           int
	secsPerUpdate      float64
	fbSizeChanged      func(width int, height int)
	mouseCameraEnabled bool
	lastMouseX         float64
	lastMouseY         float64
}

// Open creates a new GLFW Window.
func (w *Window) Open(title string, width, height int) error {
	w.width = width
	w.height = height
	w.secsPerUpdate = 1.0 / 60.0

	if err := initVideo(); err != nil {
		return err
	}

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	gw, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		return fmt.Errorf("could not create GLFW Window: %s", err.Error())
	}
	w.window = gw

	gw.MakeContextCurrent()

	w.fbWidth, w.fbHeight = gw.GetFramebufferSize()

	gw.SetFramebufferSizeCallback(func(gw *glfw.Window, width int, height int) {
		w.fbWidth = width
		w.fbHeight = height
		if w.fbSizeChanged != nil {
			w.fbSizeChanged(width, height)
		}
	})

	gw.SetSizeCallback(func(gw *glfw.Window, width int, height int) {
		w.width = width
		w.height = height
	})

	return nil
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

		glfw.PollEvents()

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
