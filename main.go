package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/tinogoehlert/goom/drivers"
	"github.com/tinogoehlert/goom/drivers/opengl"
	"github.com/tinogoehlert/goom/game"
	"github.com/tinogoehlert/goom/goom"
	"github.com/tinogoehlert/goom/graphics"
	"github.com/tinogoehlert/goom/level"
	"github.com/tinogoehlert/goom/run"
	"github.com/tinogoehlert/goom/utils"
)

var (
	logger = utils.GoomConsole

	// flags
	iwadfile     = flag.String("iwad", "DOOM1", "IWAD file to load (without extension)")
	pwadfile     = flag.String("pwad", "", "PWAD file to load (without extension)")
	levelName    = flag.String("level", "E1M1", "Level to start e.g. E1M1")
	fpsMax       = flag.Int("fpsmax", 0, "Limit FPS")
	windowHeight = 600
	windowWidth  = 800
	gameDefs     = "resources/defs.yaml"

	// engine functions
	newWindow = drivers.WindowMakers[drivers.GlfwWindow]
	newAudio  = drivers.AudioDrivers[drivers.SdlAudio]
	getTime   = drivers.Timers[drivers.SdlTimer]
	newInput  = drivers.InputProviders[drivers.GlfwInput]
)

func main() {
	flag.Parse()

	logger.Green("GOOM - DOOM clone written in Go")
	logger.Green("Press Q to exit GOOM.")

	e := newEngine()

	inputDriver := newInput(e.Window())

	inputFunc := func() {
		input(inputDriver, e.World().Me())
	}

	e.Window().RunGame(inputFunc, e.World().Update, e.render)
}

type engine struct {
	*run.Runner
	stats *renderStats
}

func newEngine() *engine {
	var err error
	e := &engine{
		&run.Runner{},
		&renderStats{lastUpdate: time.Now()},
	}

	// init all subsystems
	e.InitWAD(*iwadfile, *pwadfile, gameDefs)
	e.InitAudio(newAudio)
	err = e.InitRenderer(newWindow, windowWidth, windowHeight)
	if err != nil {
		logger.Redf("failed to init renderer %s", err.Error())
	}

	// load mission
	mission := e.GameData().Level(strings.ToUpper(*levelName))
	e.Renderer().LoadLevel(mission, e.GameData())
	e.World().LoadLevel(mission)
	player := e.World().Me()

	ssect, err := mission.FindPositionInBsp(level.GLNodesName, player.Position()[0], player.Position()[1])
	if err != nil {
		logger.Print("could not find GLnode for pos %v", player.Position())
	} else {
		sector := mission.SectorFromSSect(ssect)
		player.SetSector(sector)
	}

	e.Renderer().Camera().SetCamera(player.Position(), player.Direction(), player.Height())
	e.Renderer().SetViewPort(e.Window().GetSize())

	return e
}

func (e *engine) render(interpolTime float64) {
	started := getTime()
	e.Renderer().RenderNewFrame()
	e.Renderer().SetViewPort(e.Window().GetSize())

	mission := e.World().GetLevel()

	mission.WalkBsp(func(i int, n *level.Node, b level.BBox) {
		e.Renderer().DrawSubSector(i)
	})

	player := e.World().Me()

	e.Renderer().DrawThings(e.World().Things())
	e.Renderer().DrawHUD(player, interpolTime)

	ssect, err := mission.FindPositionInBsp(level.GLNodesName, player.Position()[0], player.Position()[1])
	if err != nil {
		logger.Print("could not find GLnode for pos %v", player.Position())
	} else {
		var sector = mission.SectorFromSSect(ssect)
		player.SetSector(sector)
		player.Lift(sector.FloorHeight())
	}
	e.Renderer().Camera().SetCamera(player.Position(), player.Direction(), player.Height())

	e.stats.showStats(e.GameData(), e.Renderer())
	e.stats.countedFrames++
	ft := getTime() - started
	e.stats.accumulatedTime += time.Duration(ft * float64(time.Second))

	if *fpsMax > 0 {
		time.Sleep(time.Second / time.Duration(*fpsMax))
	}
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

type renderStats struct {
	countedFrames   int
	accumulatedTime time.Duration
	fps             int
	meanFrameTime   float32
	lastUpdate      time.Time
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
	// TODO: position realtively to the window size
	drawText(gd.Fonts, graphics.FnCompositeRed, fpsText, 800, 600, 0.6, gr)
	ftimeText := fmt.Sprintf("frame time: %.6f ms", rs.meanFrameTime)
	// TODO: position realtively to the window size
	drawText(gd.Fonts, graphics.FnCompositeRed, ftimeText, 800, 580, 0.6, gr)
}

func input(id drivers.Input, player *game.Player) {
	if id.IsPressed(drivers.KeyUp) || id.IsPressed(drivers.KeyW) {
		player.Forward(1)
	}
	if id.IsPressed(drivers.KeyDown) || id.IsPressed(drivers.KeyS) {
		player.Forward(-1)
	}

	if id.IsPressed(drivers.KeyA) {
		player.Strafe(-1)
	}
	if id.IsPressed(drivers.KeyD) {
		player.Strafe(1)
	}

	if id.IsPressed(drivers.KeyLeft) {
		player.Turn(-1.5)
	}
	if id.IsPressed(drivers.KeyRight) {
		player.Turn(1.5)
	}

	if id.IsPressed(drivers.KeyLShift) {
		player.FireWeapon()
	}

	if id.IsPressed(drivers.KeyQ) {
		os.Exit(0)
	}
}
