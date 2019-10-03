package main

import (
	"flag"
	"math"
	"os"
	"strings"

	"github.com/go-gl/mathgl/mgl32"

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

func dist(x1, y1, x2, y2 float32) float32 {
	first := math.Pow(float64(x1-x2), 2)
	second := math.Pow(float64(y1-y2), 2)
	return float32(math.Sqrt(first + second))
}

func tryMove(m *level.Level, from, to mgl32.Vec3, angleX, angleY float64) mgl32.Vec3 {
	for _, line := range m.LinesDefs {
		if line.Left != -1 {
			continue
		}
		ls := m.Vert(uint32(line.Start))
		le := m.Vert(uint32(line.End))

		d1 := dist(to.X(), to.Z(), ls.X(), ls.Y())
		d2 := dist(to.X(), to.Z(), le.X(), le.Y())
		lineLen := dist(ls.X(), ls.Y(), le.X(), le.Y())
		buffer := float32(2) // higher # = less accurate
		if d1+d2 >= lineLen-buffer && d1+d2 <= lineLen+buffer {
			return to
		}
	}
	return to
}

func main() {
	iwadfile := flag.String("iwad", "DOOM1", "IWAD file to load (without extension)")
	pwadfile := flag.String("pwad", "", "PWAD file to load (without extension)")
	level := flag.String("level", "E1M1", "Level to start e.g. E1M1")
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

	m := gameData.Level(strings.ToUpper(*level))
	renderer.BuildGraphics(gameData)

	renderer.BuildLevel(m, gameData)

	playerPos := m.Things[0]
	var player = game.NewPlayer(float32(playerPos.X), float32(playerPos.Y), 45, float32(playerPos.Angle))
	renderer.Camera().SetCamera(player.Position(), player.Direction(), player.Height())
	defstore := game.NewDefStore("defs.yaml")
	var things = []*game.DoomThing{}

	for _, t := range m.Things {
		if obstacleDef := defstore.GetObstacleDef(int(t.Type)); obstacleDef != nil {
			obstacle := game.ThingFromDef(t.X, t.Y, 0, t.Angle, obstacleDef)
			things = appendDoomThing(things, obstacle, m)
		}

		if monsterDef := defstore.GetMonsterDef(int(t.Type)); monsterDef != nil {
			monster := game.MonsterFromDef(t.X, t.Y, 0, t.Angle, monsterDef)
			things = appendDoomThing(things, monster.DoomThing, m)
		}
	}
	music := gameData.Music["D_E1M1"]
	music.Loop()
	defer music.Stop()
	log.Green("Press Q to exit GOOM.")
	renderer.Loop(30, func() {
		renderer.DrawThings(things)
	}, func(w *glfw.Window) {
		playerInput(m, renderer.Camera(), player, w)
		if *test {
			log.Green("Test run finished. Exiting GOOM.")
			os.Exit(0)
		}
	})
}

func appendDoomThing(dst []*game.DoomThing, src *game.DoomThing, m *level.Level) []*game.DoomThing {
	var ssect, err = m.FindPositionInBsp(level.GLNodesName, src.Position()[0], src.Position()[1])
	if err != nil {
		log.Print("could not find GLnode for pos %v", src.Position())
	} else {
		var sector = m.SectorFromSSect(ssect)
		src.SetHeight(sector.FloorHeight())
	}

	return append(dst, src)
}

func playerInput(m *level.Level, cam *opengl.Camera, player *game.Player, w *glfw.Window) {
	if w.GetKey(glfw.KeyUp) == glfw.Press {
		player.Walk(7)
	}
	if w.GetKey(glfw.KeyDown) == glfw.Press {
		player.Walk(-7)
	}
	if w.GetKey(glfw.KeyLeft) == glfw.Press {
		player.Turn(-3)
	}
	if w.GetKey(glfw.KeyRight) == glfw.Press {
		player.Turn(3)
	}
	if w.GetKey(glfw.KeyQ) == glfw.Press {
		os.Exit(0)
	}

	var ssect, err = m.FindPositionInBsp(level.GLNodesName, player.Position()[0], player.Position()[1])
	if err != nil {
		log.Print("could not find GLnode for pos %v", player.Position())
	} else {
		var sector = m.SectorFromSSect(ssect)
		player.Lift(sector.FloorHeight() + 50)
	}
	cam.SetCamera(player.Position(), player.Direction(), player.Height())
}
