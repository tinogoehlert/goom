package glfw

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

// GetTime provides the game time in seconds
func GetTime() float64 {
	return glfw.GetTime()
}
