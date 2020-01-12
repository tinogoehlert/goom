package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/tinogoehlert/goom/drivers"
	"github.com/tinogoehlert/goom/drivers/opengl"
	drvShared "github.com/tinogoehlert/goom/drivers/pkg"
	"github.com/tinogoehlert/goom/goom"
	"github.com/tinogoehlert/goom/graphics"
	"github.com/tinogoehlert/goom/level"
	"github.com/tinogoehlert/goom/run"
	"github.com/tinogoehlert/goom/utils"
)

var (
	logger = utils.GoomConsole

	midiOptions = strings.Join([]string{
		string(drivers.PortMidiMusic),
		string(drivers.RtMidiMusic),
		string(drivers.SdlMusic),
		string(drivers.Noop),
	}, "|")

	// flags
	iwadfile     = flag.String("iwad", "DOOM1", "IWAD file to load (without extension)")
	pwadfile     = flag.String("pwad", "", "PWAD file to load (without extension)")
	levelName    = flag.String("level", "E1M1", "Level to start e.g. E1M1")
	fpsMax       = flag.Int("fpsmax", 0, "Limit FPS")
	midiDrv      = flag.String("mididrv", "noop", "MIDI driver name ("+midiOptions+")")
	winDrv       = flag.String("windowdrv", "glfw", "Window and Input driver name")
	windowHeight = 600
	windowWidth  = 800
	gameDefs     = "resources/defs.yaml"
)

func main() {
	flag.Parse()

	mainDrivers := drivers.Drivers{
		Window:  drivers.WindowDrivers[drivers.WindowDriver(strings.ToLower(*winDrv))],
		Audio:   drivers.AudioDrivers[drivers.SdlAudio],
		Music:   drivers.MusicDrivers[drivers.MusicDriver(strings.ToLower(*midiDrv))],
		Input:   drivers.InputDrivers[drivers.InputDriver(strings.ToLower(*winDrv))],
		GetTime: drivers.TimerFuncs[drivers.SdlTimer],
	}

	logger.Green("GOOM - DOOM clone written in Go")
	logger.Green("Press Q to exit GOOM.")

	logger.Green("MUSIC: %T", mainDrivers.Music)

	e := newEngine(&mainDrivers)

	inputFunc := func() {
		input(e)
	}

	e.Window().RunGame(inputFunc, e.World().Update, e.render)
}

type engine struct {
	*run.Runner
	stats      *renderStats
	xpos, ypos float64
}

func newEngine(drivers *drivers.Drivers) *engine {
	var err error
	e := &engine{
		&run.Runner{Drivers: drivers},
		&renderStats{lastUpdate: time.Now()},
		0, 0,
	}

	// init all subsystems
	e.InitWAD(*iwadfile, *pwadfile, gameDefs)
	e.InitAudio()
	err = e.InitRenderer(windowWidth, windowHeight)
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

	e.Renderer().SetViewPort(e.Window().GetSize())

	return e
}

func (e *engine) render(interpolTime float64) {
	started := e.GetTime()
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
	e.Renderer().SetPlayerPosition(mgl32.Vec3{
		-player.Position()[0],
		player.Height(),
		player.Position()[1],
	})
	e.stats.showStats(e.GameData(), e.Renderer())
	e.stats.countedFrames++
	ft := e.GetTime() - started
	e.stats.accumulatedTime += time.Duration(ft * float64(time.Second))

	if *fpsMax > 0 {
		time.Sleep(time.Second / time.Duration(*fpsMax))
	}
}

// DrawText draws a string on the screen
func drawText(fonts graphics.FontBook, fontName graphics.FontName, text string, xpos, ypos, scaleFactor float32, gr *opengl.GLRenderer) {
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
	// TODO: position relatively to the window size
	drawText(gd.Fonts, graphics.FnCompositeRed, fpsText, 800, 600, 0.6, gr)
	ftimeText := fmt.Sprintf("frame time: %.6f ms", rs.meanFrameTime)
	// TODO: position relatively to the window size
	drawText(gd.Fonts, graphics.FnCompositeRed, ftimeText, 800, 580, 0.6, gr)
}

func input(e *engine) {
	in := e.Input
	player := e.World().Me()

	if in.IsPressed(drvShared.KeyW) {
		player.Forward(1)
	}
	if in.IsPressed(drvShared.KeyS) {
		player.Forward(-1)
	}

	if in.IsPressed(drvShared.KeyA) {
		player.Strafe(-1)
	}
	if in.IsPressed(drvShared.KeyD) {
		player.Strafe(1)
	}

	if in.IsPressed(drvShared.KeyLeft) {
		player.Turn(-1.5)
	}
	if in.IsPressed(drvShared.KeyRight) {
		player.Turn(1.5)
	}

	if in.IsPressed(drvShared.KeyUp) {
		player.Pitch(1.5)
	}
	if in.IsPressed(drvShared.KeyDown) {
		player.Pitch(-1.5)
	}

	if in.IsPressed(drvShared.KeyLShift) || in.IsMousePressed(drvShared.MouseLeft) {
		player.FireWeapon()
	}

	if in.IsPressed(drvShared.KeyQ) {
		os.Exit(0)
	}

	if in.IsPressed(drvShared.KeyF5) {
		in.SetMouseCameraEnabled(true)
	}

	if in.IsPressed(drvShared.KeyF6) {
		in.SetMouseCameraEnabled(false)
	}
}
