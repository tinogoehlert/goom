package glfw

import (
	"fmt"
	"runtime"

	"github.com/go-gl/glfw/v3.2/glfw"
)

// Init initialize glfw
func Init() error {
	runtime.LockOSThread()

	if err := glfw.Init(); err != nil {
		return fmt.Errorf("could not initialize GLFW: %s", err.Error())
	}
	return nil
}

// Destroy terminates the GLFW Driver
func Destroy() {
	glfw.Terminate()
}
