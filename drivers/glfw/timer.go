package glfw

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

// GetTime provides the game time
func GetTime() float64 {
	return glfw.GetTime()
}
