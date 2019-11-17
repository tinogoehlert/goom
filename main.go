package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/tinogoehlert/goom/drivers/sdl"

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
	iwadfile := flag.String("iwad", "DOOM1", "IWAD file to load (without extension)")
	pwadfile := flag.String("pwad", "", "PWAD file to load (without extension)")
	levelName := flag.String("level", "E1M1", "Level to start e.g. E1M1")
	test := flag.Bool("test", false, "Exit GOOM after loading all data.")

	flag.Parse()
	logger.Green("GOOM - DOOM clone written in Go")
	logger.Green("loading %s", *iwadfile)

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

	logger.Green("load WAD: %s", *iwadfile)
	gameData, err := goom.LoadWAD(*iwadfile, *pwadfile)
	if err != nil {
		logger.Red("could not load WAD data: %s", err.Error())
	}
	mission := strings.ToUpper(*levelName)

	sm, err := sdl.NewAudioDriver(gameData.Sounds)
	if err != nil {
		logger.Fatalf("could not load sounds: %s", err.Error())
	}
	if err := sm.PlayMusic(gameData.Music.Track(mission)); err != nil {
		logger.Fatalf("could not play music: %s", err.Error())
	}

	renderer, err := opengl.NewRenderer(gameData)
	if err := renderer.LoadShaderProgram("main", "resources/shaders/main.vert", "resources/shaders/main.frag"); err != nil {
		logger.Red("could not init GL: %s", err.Error())
	}

	m := gameData.Level(mission)

	renderer.LoadLevel(m, gameData)
	renderer.SetShaderProgram("main")
	world := game.NewWorld(m, game.NewDefStore("resources/defs.yaml"), gameData, sm)
	player := world.Me()
	ssect, err := m.FindPositionInBsp(level.GLNodesName, player.Position()[0], player.Position()[1])
	if err != nil {
		logger.Print("could not find GLnode for pos %v", player.Position())
	} else {
		var sector = m.SectorFromSSect(ssect)
		player.SetSector(sector)
	}

	logger.Green("Press Q to exit GOOM.")

	renderer.Camera().SetCamera(player.Position(), player.Direction(), player.Height())
	renderer.SetViewPort(win.FrameBufferSize())
	rs := &renderStats{lastUpdate: time.Now()}
	if *test {
		logger.Green("Test run finished. Exiting GOOM.")
		os.Exit(0)
	}
	win.Run(func(elapsed float32) {
		renderer.RenderNewFrame(elapsed)
		renderer.SetViewPort(win.FrameBufferSize())

		m.WalkBsp(func(i int, n *level.Node, b level.BBox) {
			renderer.DrawSubSector(i)
		})

		renderer.DrawThings(world.Things())
		renderer.DrawHUD(world.Me(), elapsed)

		ssect, err := m.FindPositionInBsp(level.GLNodesName, player.Position()[0], player.Position()[1])
		if err != nil {
			logger.Print("could not find GLnode for pos %v", player.Position())
		} else {
			var sector = m.SectorFromSSect(ssect)
			player.SetSector(sector)
			player.Lift(sector.FloorHeight(), float32(elapsed))
		}
		world.Update(elapsed / float32(time.Second))

		renderer.Camera().SetCamera(player.Position(), player.Direction(), player.Height())
		input(win.Input(), player, elapsed)
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

func input(id drivers.InputDriver, player *game.Player, delta float32) {
	if id.IsPressed(drivers.KeyUp) {
		player.Forward(100, delta)
	}
	if id.IsPressed(drivers.KeyDown) {
		player.Forward(-100, delta)
	}
	if id.IsPressed(drivers.KeyLeft) {
		player.Turn(-130, delta)
	}
	if id.IsPressed(drivers.KeyRight) {
		player.Turn(130, delta)
	}
	if id.IsPressed(drivers.KeyLShift) {
		player.FireWeapon()
	}
}
