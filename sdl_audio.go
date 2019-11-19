//+build ignore

package main

import (
	"os"

	"github.com/tinogoehlert/goom/drivers/sdl"
	"github.com/tinogoehlert/goom/game"
)

func initAudio(world *game.World) {
	os.MkdirAll("temp/music/", 0700)
	sm, err := sdl.NewAudioDriver(world.Data().Sounds, "temp/music")
	if err != nil {
		logger.Fatalf("could not load sounds: %s", err.Error())
	}
	world.SetAudioDriver(sm)
}
