package main

import (
	"flag"
	"fmt"
	"math"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/go-gl/glfw/v3.2/glfw"

	"github.com/tinogoehlert/goom"
	"github.com/tinogoehlert/goom/cmd/doom/internal/game"
	"github.com/tinogoehlert/goom/cmd/doom/internal/game/monsters"
	"github.com/tinogoehlert/goom/cmd/doom/internal/opengl"
	"github.com/ttacon/chalk"
)

var (
	shaderDir = "resources/shaders"
)

func dist(x1, y1, x2, y2 float32) float32 {
	first := math.Pow(float64(x1-x2), 2)
	second := math.Pow(float64(y1-y2), 2)
	return float32(math.Sqrt(first + second))
}

func tryMove(m *goom.Map, from, to mgl32.Vec3, angleX, angleY float64) mgl32.Vec3 {
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
	wadfile := flag.String("wad", "DOOM1", "wad file to load (without extension)")
	flag.Parse()
	fmt.Println(chalk.Green.Color("GOOM - DOOM clone written in go"))
	fmt.Println(chalk.Green.Color(fmt.Sprintf("load %s", *wadfile)))

	renderer, err := opengl.NewRenderer()
	if err != nil {
		fmt.Printf(chalk.Red.Color("could not init GL: %s\n"), err.Error())
	}

	doomWAD := goom.NewWadManager()
	if err := doomWAD.LoadFile(*wadfile + ".wad"); err != nil {
		fmt.Printf(chalk.Red.Color("could not load WAD: %s\n"), err.Error())
	}
	doomGfx, err := doomWAD.LoadGraphics()
	if err != nil {
		fmt.Printf(chalk.Red.Color("could not load gfx: %s\n"), err.Error())
	}

	if err := doomWAD.LoadFile(*wadfile + ".gwa"); err != nil {
		fmt.Printf(chalk.Red.Color("could not load gwa: %s\n"), err.Error())
	}

	doomMaps, err := doomWAD.LoadMaps()
	if err != nil {
		fmt.Printf(chalk.Red.Color("could not load maps: %s\n"), err.Error())
	}
	if err := renderer.CreateWindow(800, 600, "GOOM"); err != nil {
		fmt.Printf(chalk.Red.Color("could not load maps: %s\n"), err.Error())
	}
	if err := renderer.LoadShaderProgram("main", shaderDir+"/main.vert", shaderDir+"/main.frag"); err != nil {
		fmt.Printf(chalk.Red.Color("could not init GL: %s\n"), err.Error())
	}
	if err := renderer.LoadShaderProgram("red", shaderDir+"/main.vert", shaderDir+"/simpleRed.frag"); err != nil {
		fmt.Printf(chalk.Red.Color("could not init GL: %s\n"), err.Error())
	}

	renderer.SetShaderProgram("main")

	m := &doomMaps[0]
	fmt.Println(len(m.Nodes(goom.GLNodesName)))
	renderer.BuildLevel(m, doomGfx)
	renderer.BuildSprites(doomGfx)

	playerPos := m.Things[0]
	var player = game.NewPlayer(float32(playerPos.X), float32(playerPos.Y), 45, float32(playerPos.Angle))
	renderer.Camera().SetCamera(player.Position(), player.Direction(), player.Height())

	var things = []game.DoomThing{}

	for _, t := range m.Things {
		if obstacle := game.NewObstacle(&t); obstacle != nil {
			things = appendDoomThing(things, obstacle, m)
		}
		if monster := monsters.NewMonster(&t); monster != nil {
			things = appendDoomThing(things, monster, m)
		}
	}
	//os.Exit(0)
	renderer.Loop(30, func() {
		renderer.DrawThings(things)
	}, func(w *glfw.Window) {
		playerInput(m, renderer.Camera(), player, w)
	})
}

func appendDoomThing(dst []game.DoomThing, src game.DoomThing, m *goom.Map) []game.DoomThing {
	var ssect, err = m.FindPositionInBsp(goom.GLNodesName, src.Position()[0], src.Position()[1])
	if err != nil {
		fmt.Println("could not find GLnode for pos", src.Position())
	} else {
		var sector = m.SectorFromSSect(ssect)
		src.SetHeight(sector.FloorHeight())
	}

	return append(dst, src)
}

func playerInput(m *goom.Map, cam *opengl.Camera, player *game.Player, w *glfw.Window) {
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

	var ssect, err = m.FindPositionInBsp(goom.GLNodesName, player.Position()[0], player.Position()[1])
	if err != nil {
		fmt.Println("could not find GLnode for pos", player.Position())
	} else {
		var sector = m.SectorFromSSect(ssect)
		player.SetHeight(sector.FloorHeight() + 50)
	}
	cam.SetCamera(player.Position(), player.Direction(), player.Height())
}
