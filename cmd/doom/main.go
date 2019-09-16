package main

import (
	"fmt"
	"math"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/go-gl/glfw/v3.2/glfw"

	"github.com/tinogoehlert/goom"
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
	fmt.Println(chalk.Green.Color("GOOM - DOOM clone written in go"))

	renderer, err := opengl.NewRenderer()
	if err != nil {
		fmt.Printf(chalk.Red.Color("could not init GL: %s\n"), err.Error())
	}

	doomWAD := goom.NewWadManager()
	if err := doomWAD.LoadFile("DOOM1.WAD"); err != nil {
		fmt.Printf(chalk.Red.Color("could not load WAD: %s\n"), err.Error())
	}
	doomGfx, err := doomWAD.LoadGraphics()
	if err != nil {
		fmt.Printf(chalk.Red.Color("could not load gfx: %s\n"), err.Error())
	}

	if err := doomWAD.LoadFile("DOOM1.gwa"); err != nil {
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

	m := &doomMaps[0]
	renderer.BuildLevel(m, doomGfx)
	renderer.BuildSprites(doomGfx)

	playerPos := m.Things[0]
	var player = NewPlayer(float32(playerPos.X), float32(playerPos.Y), 45, float32(playerPos.Angle))
	renderer.Camera().SetCamera(player.Position(), player.dir, player.Height())

	renderer.Loop(30, func(w *glfw.Window) {
		playerInput(m, renderer.Camera(), player, w)
	})
}

func playerInput(m *goom.Map, cam *opengl.Camera, player *Player, w *glfw.Window) {
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

	var height = float32(0)
	for _, node := range m.Nodes("GL_NODES") {
		if goom.MagicU32(node.Right).MagicBit() {
			if node.RightBBox.PosInBox(player.position[0], player.position[1]) {
				var (
					ssect  = m.SubSectors("GL_SSECT")[int(goom.MagicU32(node.Right).Uint32())]
					fseg   = ssect.Segments()[0]
					line   = m.LinesDefs[fseg.GetLineDef()]
					side   = m.SideDefs[line.Right]
					sector = m.Sectors[side.Sector]
				)

				if fseg.GetDirection() == 1 {
					side = m.SideDefs[line.Left]
					sector = m.Sectors[side.Sector]
				}

				height = sector.FloorHeight()
			}
		}
		if goom.MagicU32(node.Left).MagicBit() {
			if node.LeftBBox.PosInBox(player.position[0], player.position[1]) {
				var (
					ssect  = m.SubSectors("GL_SSECT")[int(goom.MagicU32(node.Left).Uint32())]
					fseg   = ssect.Segments()[0]
					line   = m.LinesDefs[fseg.GetLineDef()]
					side   = m.SideDefs[line.Right]
					sector = m.Sectors[side.Sector]
				)

				if fseg.GetDirection() == 1 {
					side = m.SideDefs[line.Left]
					sector = m.Sectors[side.Sector]
				}

				height = sector.FloorHeight()
			}
		}
	}
	cam.SetCamera(player.Position(), player.dir, height+50)
}
