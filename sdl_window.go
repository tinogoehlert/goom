//+build ignore

package main

import (
	sdl_native "github.com/tinogoehlert/go-sdl2/sdl"
	"github.com/tinogoehlert/goom/drivers"
	"github.com/tinogoehlert/goom/drivers/sdl"
)

func init() {
	if err := sdl.InitSDLVideo(); err != nil {
		logger.Fatalf(err.Error())
	}
}

func initWindow(title string, width, height int) (drivers.Window, error) {
	return sdl.NewGLWindow(title, width, height)
}

func getTime() float64 {
	return float64(sdl_native.GetTicks())
}
