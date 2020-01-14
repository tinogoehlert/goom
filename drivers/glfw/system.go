package glfw

import (
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
)

var videoInitialized bool

// InitVideo inits Glfw.
func initVideo() error {
	if videoInitialized {
		return nil
	}

	runtime.LockOSThread()
	err := glfw.Init()

	if err == nil {
		videoInitialized = true
	}

	return err
}

// Destroy terminates the GLFW driver
func Destroy() {
	glfw.Terminate()
	runtime.UnlockOSThread()
}
