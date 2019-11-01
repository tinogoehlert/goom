package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/tinogoehlert/goom/audio/midi"

	"github.com/go-gl/glfw/v3.2/glfw"

	"github.com/tinogoehlert/goom"
	"github.com/tinogoehlert/goom/cmd/doom/internal/game"
	"github.com/tinogoehlert/goom/cmd/doom/internal/opengl"
	"github.com/tinogoehlert/goom/graphics"
	"github.com/tinogoehlert/goom/level"
)

var (
	shaderDir = "resources/shaders"
	log       = goom.GoomConsole
)

type renderStats struct {
	countedFrames   int
	accumulatedTime time.Duration
	fps             int
	meanFrameTime   float32
	lastUpdate      time.Time
}

func main() {
	iwadfile := flag.String("iwad", "DOOM1", "IWAD file to load (without extension)")
	pwadfile := flag.String("pwad", "", "PWAD file to load (without extension)")
	levelName := flag.String("level", "E1M1", "Level to start e.g. E1M1")
	test := flag.Bool("test", false, "Exit GOOM after loading all data.")
	showStats := flag.Bool("rstat", false, "Show renderer stats like FPS and frametime")

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
	renderer.BuildGraphics(gameData)

	mission := strings.ToUpper(*levelName)
	m := gameData.Level(mission)

	renderer.BuildLevel(m, gameData)
	world := game.NewWorld(m, game.NewDefStore("defs.yaml"), gameData)

	track := gameData.Music.Track(mission)
	mPlayer, err := midi.NewPlayer(midi.Any)
	if err != nil {
		log.Red("failed to start MIDI player: %s", err.Error())
	}

	if mPlayer != nil && track != nil {
		defer mPlayer.Close()
		go mPlayer.Loop(track.MidiStream)
	}

	log.Green("Press Q to exit GOOM.")

	stats := &renderStats{lastUpdate: time.Now()}

	renderer.SetFPSCap(30)
	renderer.Loop(func() {
		m.WalkBsp(func(i int, n *level.Node, b level.BBox) {
			renderer.DrawSubSector(i)
		})
		renderer.DrawThings(world.Things())

		renderer.DrawHUD(world.Me())

		if *showStats {
			stats.showStats(gameData, renderer)
		}

	}, func(w *glfw.Window, frameTime float32) {
		world.Update(frameTime)
		playerInput(m, renderer.Camera(), world.Me(), w, frameTime)
		if *test {
			log.Green("Test run finished. Exiting GOOM.")
			os.Exit(0)
		}
	}, func(recordedFrameTime time.Duration) {
		if !(*showStats) {
			return
		}
		stats.countedFrames++
		stats.accumulatedTime += recordedFrameTime
	})
}

var speed = float32(420)

func playerInput(m *level.Level, cam *opengl.Camera, player *game.Player, w *glfw.Window, frameTime float32) {
	if w.GetKey(glfw.KeyW) == glfw.Press || w.GetKey(glfw.KeyUp) == glfw.Press {
		player.Walk(speed, frameTime)
	}
	if w.GetKey(glfw.KeyS) == glfw.Press || w.GetKey(glfw.KeyDown) == glfw.Press {
		player.Walk(-speed, frameTime)
	}
	if w.GetKey(glfw.KeyLeft) == glfw.Press {
		player.Turn(-speed/2, frameTime)
	}
	if w.GetKey(glfw.KeyRight) == glfw.Press {
		player.Turn(speed/2, frameTime)
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
