package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/tinogoehlert/goom/drivers"
	"github.com/tinogoehlert/goom/drivers/glfw"
	"github.com/tinogoehlert/goom/drivers/opengl"
	"github.com/tinogoehlert/goom/game"
	"github.com/tinogoehlert/goom/goom"
	"github.com/tinogoehlert/goom/graphics"
	"github.com/tinogoehlert/goom/level"
	"github.com/tinogoehlert/goom/utils"
)

var logger = utils.GoomConsole

func main() {
	if err := glfw.Init(); err != nil {
		logger.Fatalf(err.Error())
	}

	win, err := glfw.NewWindow("GOOM", 800, 600)
	if err != nil {
		logger.Fatalf(err.Error())
	}

	if err := opengl.Init(); err != nil {
		logger.Fatalf(err.Error())
	}

	gameData, err := goom.LoadWAD("DOOM1", "")
	if err != nil {
		logger.Red("could not load WAD data: %s", err.Error())
	}

	renderer, err := opengl.NewRenderer(gameData)
	if err := renderer.LoadShaderProgram("main", "resources/shaders/main.vert", "resources/shaders/main.frag"); err != nil {
		logger.Red("could not init GL: %s", err.Error())
	}

	mission := strings.ToUpper("E1M1")
	m := gameData.Level(mission)

	renderer.LoadLevel(m, gameData)
	renderer.SetShaderProgram("main")
	world := game.NewWorld(m, game.NewDefStore("resources/defs.yaml"), gameData)
	player := world.Me()
	renderer.Camera().SetCamera(player.Position(), player.Direction(), player.Height())
	renderer.SetViewPort(win.FrameBufferSize())
	rs := &renderStats{lastUpdate: time.Now()}
	win.Run(func(elapsed float32) {
		renderer.RenderNewFrame(elapsed)
		renderer.SetViewPort(win.FrameBufferSize())

		m.WalkBsp(func(i int, n *level.Node, b level.BBox) {
			renderer.DrawSubSector(i)
		})

		renderer.DrawThings(world.Things())
		renderer.DrawHUD(world.Me())

		ssect, err := m.FindPositionInBsp(level.GLNodesName, player.Position()[0], player.Position()[1])
		if err != nil {
			logger.Print("could not find GLnode for pos %v", player.Position())
		} else {
			var sector = m.SectorFromSSect(ssect)
			player.SetSector(sector)
			player.Lift(sector.FloorHeight()+40, float32(elapsed))
		}
		renderer.Camera().SetCamera(player.Position(), player.Direction(), player.Height())
		handleInput(player, win.Input(), elapsed)
		player.Update(elapsed)
		player.Turn(turnvel, elapsed)

		rs.showStats(gameData, renderer)
		rs.countedFrames++
		rs.accumulatedTime += time.Duration(elapsed)
	})
}

// DrawText draws a string on the screen
func drawText(fonts graphics.FontBook, fontName graphics.FontName, text string, xpos, ypos float32, scaleFactor float32, gr *opengl.GLRenderer) {
	font := fonts[fontName]
	spacing := float32(font.GetSpacing()) * scaleFactor

	// currently, we only know uppercase glyphs
	text = strings.ToUpper(text)

	for _, r := range text {
		if r == ' ' {
			xpos -= spacing
			continue
		}

		glyph := font.GetGlyph(r)
		if glyph == nil {
			xpos -= spacing
			continue
		}

		gr.DrawHUdElement(glyph.GetName(), xpos, ypos, scaleFactor)
		xpos -= spacing + float32(glyph.Width())*scaleFactor
	}
}

func (rs *renderStats) showStats(gd *goom.GameData, gr *opengl.GLRenderer) {
	t1 := time.Now()
	if t1.Sub(rs.lastUpdate) >= time.Second {
		rs.fps = rs.countedFrames
		rs.meanFrameTime = (float32(rs.accumulatedTime) / float32(rs.countedFrames)) / float32(time.Millisecond)
		rs.countedFrames = 0
		rs.accumulatedTime = time.Duration(0)
		rs.lastUpdate = t1
	}

	fpsText := fmt.Sprintf("FPS: %d", rs.fps)
	drawText(gd.Fonts, graphics.FnCompositeRed, fpsText, 800, 600, 0.6, gr)
	ftimeText := fmt.Sprintf("frame time: %.3f ms", rs.meanFrameTime)
	drawText(gd.Fonts, graphics.FnCompositeRed, ftimeText, 800, 580, 0.6, gr)
}

type renderStats struct {
	countedFrames   int
	accumulatedTime time.Duration
	fps             int
	meanFrameTime   float32
	lastUpdate      time.Time
}

var fwdvel = float32(0)
var turnvel = float32(0)

func handleInput(player *game.Player, input drivers.InputDriver, t float32) {
	select {
	case k := <-input.KeyStates():
		handleKey(player, k, t)
	default:
	}
}

func handleKey(player *game.Player, k drivers.Key, t float32) float32 {
	switch k.Keycode {
	case drivers.KeyUp:
		if k.State == drivers.KeyPressed {
			fwdvel = 520
			player.SmoothWalk(520, t)
		}
		if k.State == drivers.KeyReleased {
			fwdvel = 0
			player.SmoothWalk(0, t)
		}
	case drivers.KeyDown:
		if k.State == drivers.KeyPressed {
			fwdvel = -520
			player.SmoothWalk(-520, t)
		}
		if k.State == drivers.KeyReleased {
			fwdvel = 0
		}

	case drivers.KeyLeft:
		if k.State == drivers.KeyPressed {
			turnvel = -300
		}
		if k.State == drivers.KeyReleased {
			turnvel = 0
		}
	case drivers.KeyRight:
		if k.State == drivers.KeyPressed {
			turnvel = 300
		}
		if k.State == drivers.KeyReleased {
			turnvel = 0
		}
	}
	return fwdvel
}
