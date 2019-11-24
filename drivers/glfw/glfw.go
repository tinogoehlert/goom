package glfw

import (
	"runtime"

	"github.com/go-gl/glfw/v3.2/glfw"
)

var videoInitialized bool

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
}
