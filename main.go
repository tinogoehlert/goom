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
	"github.com/tinogoehlert/goom/utils"
)

var (
	logger = utils.GoomConsole

	// flags
	iwadfile  = flag.String("iwad", "DOOM1", "IWAD file to load (without extension)")
	pwadfile  = flag.String("pwad", "", "PWAD file to load (without extension)")
	levelName = flag.String("level", "E1M1", "Level to start e.g. E1M1")
	fpsMax    = flag.Int("fpsmax", 0, "Limit FPS")

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

	r := &runner{
		stats: &renderStats{lastUpdate: time.Now()},
	}

	r.initialize()

	inputDriver := newInput(r.window)

	inputFunc := func() {
		input(inputDriver, r.world.Me())
	}

	r.window.RunGame(inputFunc, r.world.Update, r.render)
}

type runner struct {
	gameData *goom.GameData
	world    *game.World
	window   drivers.Window
	renderer *opengl.GLRenderer
	stats    *renderStats
}

func (r *runner) initialize() {
	var err error

	logger.Green("loading %s", *iwadfile)
	r.gameData, err = goom.LoadWAD(*iwadfile, *pwadfile)
	if err != nil {
		logger.Red("could not load WAD data: %s", err.Error())
	}

	mission := r.gameData.Level(strings.ToUpper(*levelName))
	r.world = game.NewWorld(r.gameData, game.NewDefStore("resources/defs.yaml"))

	audio, err := newAudio(&r.world.Data().Sounds, "temp/music")
	if err != nil {
		logger.Fatalf(err.Error())
	}

	r.world.SetAudioDriver(audio)

	r.window, err = newWindow("GOOM", 800, 600)
	if err != nil {
		logger.Fatalf(err.Error())
	}

	err = opengl.Init()
	if err != nil {
		logger.Fatalf(err.Error())
	}

	r.renderer, err = opengl.NewRenderer(r.gameData)
	if err != nil {
		logger.Red("could not init GL: %s", err.Error())
	}

	err = r.renderer.LoadShaderProgram("main", "resources/shaders/main.vert", "resources/shaders/main.frag")
	if err != nil {
		logger.Red("could not init GL: %s", err.Error())
	}

	r.renderer.LoadLevel(mission, r.gameData)

	err = r.renderer.SetShaderProgram("main")
	if err != nil {
		logger.Red("could not init GL: %s", err.Error())
	}

	r.world.LoadLevel(mission)

	player := r.world.Me()

	ssect, err := mission.FindPositionInBsp(level.GLNodesName, player.Position()[0], player.Position()[1])
	if err != nil {
		logger.Print("could not find GLnode for pos %v", player.Position())
	} else {
		sector := mission.SectorFromSSect(ssect)
		player.SetSector(sector)
	}

	r.renderer.Camera().SetCamera(player.Position(), player.Direction(), player.Height())
	r.renderer.SetViewPort(r.window.GetSize())
}

func (r *runner) render(interpolTime float64) {
	started := getTime()
	r.renderer.RenderNewFrame()
	r.renderer.SetViewPort(r.window.GetSize())

	mission := r.world.GetLevel()

	mission.WalkBsp(func(i int, n *level.Node, b level.BBox) {
		r.renderer.DrawSubSector(i)
	})

	player := r.world.Me()

	r.renderer.DrawThings(r.world.Things())
	r.renderer.DrawHUD(player, interpolTime)

	ssect, err := mission.FindPositionInBsp(level.GLNodesName, player.Position()[0], player.Position()[1])
	if err != nil {
		logger.Print("could not find GLnode for pos %v", player.Position())
	} else {
		var sector = mission.SectorFromSSect(ssect)
		player.SetSector(sector)
		player.Lift(sector.FloorHeight())
	}
	r.renderer.Camera().SetCamera(player.Position(), player.Direction(), player.Height())

	r.stats.showStats(r.gameData, r.renderer)
	r.stats.countedFrames++
	ft := getTime() - started
	r.stats.accumulatedTime += time.Duration(ft * float64(time.Second))

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
