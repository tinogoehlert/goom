package main

import (
	"github.com/tinogoehlert/goom/drivers"
	"github.com/tinogoehlert/goom/drivers/glfw"
)

func init() {
	if err := glfw.Init(); err != nil {
		logger.Fatalf(err.Error())
	}
}

func initWindow(title string, width, height int) (drivers.Window, error) {
	return glfw.NewWindow(title, width, height)
}

func getTime() float64 {
	return float64(glfw.GetGameTime())
}
