package sdl

import (
	"fmt"
	"log"

	"github.com/tinogoehlert/go-sdl2/sdl"
)

type window struct {
	window        *sdl.Window
	width         int
	height        int
	fbWidth       int
	fbHeight      int
	glContext     sdl.GLContext
	secsPerUpdate float64
	fbSizeChanged func(width int, height int)
	shouldClose   bool
}

// NewWindow creates a new sdl window with GL context
func NewWindow(title string, width, height int) (*window, error) {
	if err := initVideo(); err != nil {
		log.Println(err)
		return nil, err
	}

	sdl.GLSetAttribute(sdl.GL_DOUBLEBUFFER, 2)
	sdl.GLSetAttribute(sdl.GL_DEPTH_SIZE, 32)
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
		log.Println(err)
		return nil, err
	}

	fbWidth, fbHeight := sdlwin.GLGetDrawableSize()
	var win = &window{
		window:        sdlwin,
		width:         width,
		height:        height,
		secsPerUpdate: float64(1) / 60,
		fbWidth:       int(fbWidth),
		fbHeight:      int(fbHeight),
	}

	if win.glContext, err = sdlwin.GLCreateContext(); err != nil {
		sdlwin.Destroy()
		log.Println(err)
		return nil, err
	}

	return win, nil
}

// Size Returns the current size of the Window
func (w *window) Size() (int, int) {
	return w.width, w.height
}

// GetSize Returns the current size of the Window
func (w *window) GetSize() (int, int) {
	fbWidth, fbHeight := w.window.GLGetDrawableSize()
	return int(fbWidth * 2), int(fbHeight * 2)
}

// ShouldClose determines if the window should close
func (w *window) ShouldClose() bool {
	return w.shouldClose
}

// RunGame runs the game loop
func (w *window) RunGame(input func(), update func(), render func(float64)) {
	var (
		previous         = float64(sdl.GetTicks()) / 1000
		lag              = float64(0)
		elapsed, current float64
	)

	fmt.Println(w.secsPerUpdate)
	for {
		current = float64(sdl.GetTicks()) / 1000
		elapsed = current - previous
		previous = current

		lag += elapsed

		sdl.PollEvent()

		for lag >= w.secsPerUpdate {
			lag -= w.secsPerUpdate
			input()
			update()
		}
		// TODO: This tells the renderer close we are to the next tick, so if we
		//       are between two ticks we can display (as an example) the movement
		//       of a projectile by and additional 0.8 units instead of fixed 1 unit.
		//       The todo is to implement this in the renderer ^^.

		render(lag / w.secsPerUpdate)
		w.window.GLSwap()
	}
}

// Close closes the window
func (w *window) Close() {
	w.window.Destroy()
}
