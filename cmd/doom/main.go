package main

import (
	"flag"
	"os"
	"strings"

	"github.com/go-gl/glfw/v3.2/glfw"

	"github.com/tinogoehlert/goom"
	"github.com/tinogoehlert/goom/cmd/doom/internal/game"
	"github.com/tinogoehlert/goom/cmd/doom/internal/opengl"
	"github.com/tinogoehlert/goom/level"
)

var (
	shaderDir = "resources/shaders"
	log       = goom.GoomConsole
)

func main() {
	iwadfile := flag.String("iwad", "DOOM1", "IWAD file to load (without extension)")
	pwadfile := flag.String("pwad", "", "PWAD file to load (without extension)")
	levelName := flag.String("level", "E1M1", "Level to start e.g. E1M1")
	test := flag.Bool("test", false, "Exit GOOM after loading all data.")

	flag.Parse()
	log.Green("GOOM - DOOM clone written in Go")
	log.Green("loading %s", *iwadfile)

	renderer, err := opengl.NewRenderer()
	if err != nil {
		log.Red("could not init GL: %s", err.Error())
	}

	gameData, err := goom.LoadWAD(*iwadfile, *pwadfile)
	if err != nil {
		log.Red("could not load WAD: %s", err.Error())
	}

	if err := renderer.CreateWindow(800, 600, "GOOM"); err != nil {
		log.Red("could not load maps: %s", err.Error())
	}
	if err := renderer.LoadShaderProgram("main", shaderDir+"/main.vert", shaderDir+"/main.frag"); err != nil {
		log.Red("could not init GL: %s", err.Error())
	}

	renderer.SetShaderProgram("main")

	m := gameData.Level(strings.ToUpper(*levelName))
	renderer.BuildGraphics(gameData)

	renderer.BuildLevel(m, gameData)

	world := game.NewWorld(m, game.NewDefStore("defs.yaml"), gameData)
	log.Green("Press Q to exit GOOM.")

	renderer.SetFPSCap(30)
	renderer.Loop(func() {

		m.WalkBsp(func(i int, n *level.Node, b level.BBox) {
			renderer.DrawSubSector(i)
		})
		renderer.DrawThings(world.Things())

		renderer.DrawHUD(world.Me())
	}, func(w *glfw.Window, frameTime float32) {
		world.Update(frameTime)
		playerInput(m, renderer.Camera(), world.Me(), w, frameTime)
		if *test {
			log.Green("Test run finished. Exiting GOOM.")
			os.Exit(0)
		}
	})
}

var speed = float32(120)

func playerInput(m *level.Level, cam *opengl.Camera, player *game.Player, w *glfw.Window, frameTime float32) {
	if w.GetKey(glfw.KeyW) == glfw.Press || w.GetKey(glfw.KeyUp) == glfw.Press {
		player.Walk(speed, frameTime)
	}
	if w.GetKey(glfw.KeyS) == glfw.Press || w.GetKey(glfw.KeyDown) == glfw.Press {
		player.Walk(-speed, frameTime)
	}
	if w.GetKey(glfw.KeyLeft) == glfw.Press {
		player.Turn(-speed, frameTime)
	}
	if w.GetKey(glfw.KeyRight) == glfw.Press {
		player.Turn(speed, frameTime)
	}
	if w.GetKey(glfw.KeyD) == glfw.Press {
		player.Strafe(speed, frameTime)
	}
	if w.GetKey(glfw.KeyA) == glfw.Press {
		player.Strafe(-speed, frameTime)
	}
	if w.GetKey(glfw.KeySpace) == glfw.Press {
		player.FireWeapon()
	}
	if w.GetKey(glfw.Key1) == glfw.Press {
		player.SwitchWeapon("pistol")
	}
	if w.GetKey(glfw.Key2) == glfw.Press {
		player.SwitchWeapon("shotgun")
	}
	if w.GetKey(glfw.Key8) == glfw.Press {
		player.SwitchWeapon("super-shotgun")
	}
	if w.GetKey(glfw.Key0) == glfw.Press {
		player.SwitchWeapon("chainsaw")
	}
	if w.GetKey(glfw.KeyQ) == glfw.Press {
		os.Exit(0)
	}

	var ssect, err = m.FindPositionInBsp(level.GLNodesName, player.Position()[0], player.Position()[1])
	if err != nil {
		log.Print("could not find GLnode for pos %v", player.Position())
	} else {
		var sector = m.SectorFromSSect(ssect)
		player.SetSector(sector)
		player.Lift(sector.FloorHeight()+50, frameTime)
	}
	cam.SetCamera(player.Position(), player.Direction(), player.Height())
}
