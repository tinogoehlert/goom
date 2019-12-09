package run

import (
	"path"

	"github.com/tinogoehlert/goom/drivers"
	"github.com/tinogoehlert/goom/drivers/opengl"
	"github.com/tinogoehlert/goom/game"
	"github.com/tinogoehlert/goom/goom"
	"github.com/tinogoehlert/goom/utils"
)

// Runner stores game data and drivers.
type Runner struct {
	gameData *goom.GameData
	world    *game.World
	window   drivers.Window
	renderer *opengl.GLRenderer
	gameDir  string
}

var (
	testRunner *Runner
	logger     = utils.GoomConsole
)

// TestRunner return a non initalized runner for headless testing
func TestRunner(gameDir ...string) *Runner {
	if testRunner == nil {
		r := &Runner{gameDir: path.Join(gameDir...)}
		iwad := path.Join(r.gameDir, "DOOM1")
		defs := path.Join(r.gameDir, "resources/defs.yaml")
		r.InitWAD(iwad, "", defs)
		r.InitAudio(drivers.AudioDrivers[drivers.SdlAudio])
		testRunner = r
	}
	return testRunner
}

// Window returns the game window.
func (r *Runner) Window() drivers.Window {
	return r.window
}

// World returns the game world.
func (r *Runner) World() *game.World {
	return r.world
}

// Renderer returns the game world.
func (r *Runner) Renderer() *opengl.GLRenderer {
	return r.renderer
}

// GameData returns the game world.
func (r *Runner) GameData() *goom.GameData {
	return r.gameData
}

// InitWAD loads the game data.
func (r *Runner) InitWAD(iwadfile, pwadfile, gameDefs string) {
	logger.Green("loading %s", iwadfile)
	var err error
	r.gameData, err = goom.LoadWAD(iwadfile, pwadfile)
	if err != nil {
		logger.Red("failed to load WAD data: %s", err.Error())
	}
	r.world = game.NewWorld(r.gameData, game.NewDefStore(gameDefs))
}

// InitAudio starts the audio driver.
func (r *Runner) InitAudio(newAudio drivers.AudioCreator) {
	audio, err := newAudio(&r.GameData().Sounds, path.Join(r.gameDir, "temp", "music"))
	if err != nil {
		logger.Red("failed to init audio system: %s", err.Error())
	}

	r.world.SetAudioDriver(audio)
}

// InitRenderer starts the window driver and GL renderer.
func (r *Runner) InitRenderer(newWindow drivers.WindowCreator, w, h int) error {
	var err error
	r.window, err = newWindow("GOOM", w, h)
	if err != nil {
		return err
	}

	err = opengl.Init()
	if err != nil {
		return err
	}

	r.renderer, err = opengl.NewRenderer(r.gameData)
	if err != nil {
		return err
	}

	err = r.renderer.LoadShaderProgram("main", "resources/shaders/main.vert", "resources/shaders/main.frag")
	if err != nil {
		logger.Red("failed to load shaders: %s", err.Error())
	}

	err = r.renderer.SetShaderProgram("main")
	if err != nil {
		logger.Red("failed to init shaders: %s", err.Error())
	}
	return nil
}
