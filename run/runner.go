package run

import (
	"fmt"
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

var logger = utils.GoomConsole

// TestRunner return a non initalized runner for headless testing
func TestRunner(gameDir ...string) *Runner {
	r := &Runner{gameDir: path.Join(gameDir...)}
	iwad := path.Join(r.gameDir, "DOOM1")
	defs := path.Join(r.gameDir, "resources", "defs.yaml")
	r.InitWAD(iwad, "", defs)
	r.InitAudio(
		drivers.AudioDrivers[drivers.SdlAudio],
		drivers.MusicDrivers[drivers.SdlMusic],
	)

	return r
}

// Window returns the game window.
func (r *Runner) Window() drivers.Window {
	if r == nil {
		return nil
	}
	return r.window
}

// World returns the game world.
func (r *Runner) World() *game.World {
	if r == nil {
		return nil
	}
	return r.world
}

// Renderer returns the game world.
func (r *Runner) Renderer() *opengl.GLRenderer {
	if r == nil {
		return nil
	}
	return r.renderer
}

// GameData returns the game world.
func (r *Runner) GameData() *goom.GameData {
	if r == nil {
		return nil
	}
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
func (r *Runner) InitAudio(audio drivers.Audio, music drivers.Music) {
	err := audio.Init(&r.GameData().Sounds)
	if err != nil {
		logger.Red("failed to init audio system: %s", err.Error())
	}
	r.world.SetAudioDriver(audio)

	err := music.Init(&r.GameData().Music, path.Join(r.gameDir, "temp", "music"))
	if err != nil {
		logger.Red("failed to init music system: %s", err.Error())
	}
	r.world.SetMusicDriver(music)
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

	err = r.renderer.LoadShaderProgram(
		"main",
		path.Join("resources", "shaders", "main.vert"),
		path.Join("resources", "shaders", "main.frag"),
	)
	if err != nil {
		return fmt.Errorf("failed to load shaders: %w", err)
	}

	err = r.renderer.SetShaderProgram("main")
	if err != nil {
		return fmt.Errorf("failed to init shaders: %w", err)
	}

	return nil
}
