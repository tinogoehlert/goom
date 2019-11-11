package sdl

import (
	"fmt"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/tinogoehlert/goom/drivers"
	"github.com/veandco/go-sdl2/sdl"
)

type GLWindow struct {
	window        *sdl.Window
	width         int
	height        int
	fbWidth       int
	fbHeight      int
	glContext     sdl.GLContext
	secsPerUpdate float64
	fbSizeChanged func(width int, height int)
	inputDrv      *InputDriver
	shouldClose   bool
}

// NewGLWindow creates a new sdl window with GL context
func NewGLWindow(title string, width, height int) (*GLWindow, error) {
	sdl.GLSetAttribute(sdl.GL_DOUBLEBUFFER, 1)
	sdl.GLSetAttribute(sdl.GL_DEPTH_SIZE, 16)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 2)

	sdlwin, err := sdl.CreateWindow(
		title,
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		int32(width),
		int32(height),
		sdl.WINDOW_OPENGL|sdl.WINDOW_RESIZABLE,
	)
	if err != nil {
		return nil, err
	}
	fbWidth, fbHeight := sdlwin.GLGetDrawableSize()
	var win = &GLWindow{
		window:   sdlwin,
		width:    width,
		height:   height,
		fbWidth:  int(fbWidth),
		fbHeight: int(fbHeight),
	}

	if win.glContext, err = sdlwin.GLCreateContext(); err != nil {
		sdlwin.Destroy()
		return nil, err
	}

	fmt.Println(sdlwin.GLMakeCurrent(win.glContext))

	return win, nil
}

func (w *GLWindow) Input() drivers.InputDriver {
	return w.inputDrv
}

// Size Returns the current size of the Window
func (w *GLWindow) Size() (int, int) {
	return w.width, w.height
}

// FrameBufferSize Returns the current size of the Window
func (w *GLWindow) FrameBufferSize() (int, int) {
	fbWidth, fbHeight := w.window.GLGetDrawableSize()
	return int(fbWidth * 2), int(fbHeight * 2)
}

// ShouldClose determines if the window should close
func (w *GLWindow) ShouldClose() bool {
	return w.shouldClose
}

// Run runs the window loop
func (w *GLWindow) Run(loop func(elapsed float32)) {
	var (
		previous = float64(sdl.GetTicks())
	)

	for !w.shouldClose {
		var (
			frameTime = float64(sdl.GetTicks()) - previous
			waitTime  = w.secsPerUpdate - frameTime
		)
		loop(float32(waitTime))
		if waitTime > 0 {
			glfw.WaitEventsTimeout(waitTime)
		}

		w.window.GLSwap()
		sdl.PollEvent()
		//w.window.UpdateSurface()
		previous = float64(sdl.GetTicks())
	}
}

// Close closes the window
func (w *GLWindow) Close() {
	w.window.Destroy()
}
