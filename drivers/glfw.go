package drivers

import (
	glfw "github.com/tinogoehlert/goom/drivers/glfw"
)

func init() {
	gw := &glfw.Window{}
	WindowDrivers[GlfwWindow] = gw
	InputDrivers[GlfwInput] = gw
	TimerFuncs[GlfwTimer] = glfw.GetTime
}
