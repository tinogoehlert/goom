package sdl

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/tinogoehlert/go-sdl2/mix"
	"github.com/tinogoehlert/go-sdl2/sdl"

	"github.com/tinogoehlert/goom/audio/music"
	"github.com/tinogoehlert/goom/audio/sfx"
)

// Audio is the SDL audio driver.
type Audio struct {
	sounds           *sfx.Sounds
	chunks           map[string]*mix.Chunk
	currentTrackName string
	currentTrack     *mix.Music
	tempFolder       string
	test             bool
}

// NewAudio returns an SDL audio driver.
func NewAudio(sounds *sfx.Sounds, tempFolder string) (*Audio, error) {
	err := initAudio()
	if err != nil {
		return nil, fmt.Errorf("failed to init SDL subsystem: %s", err.Error())
	}

	if _, err := mix.OpenAudioDevice(22050, mix.DEFAULT_FORMAT, 2, 4096, "", sdl.AUDIO_ALLOW_ANY_CHANGE); err != nil {
		return nil, fmt.Errorf("failed to open audio device: %s", err.Error())
	}

	os.MkdirAll(tempFolder, 0700)

	a := &Audio{
		sounds:     sounds,
		chunks:     make(map[string]*mix.Chunk),
		tempFolder: tempFolder,
	}

	return a, nil
}

// TestMode silences all sounds and music and sets all delays to 0 for testing.
func (a *Audio) TestMode() {
	a.test = true
}

// PlayMusic plays a MUS track.
// For SDL playback the MUS track is converted to a MID file and
// stored in a temp dir unless the target MID file is already present.
func (a *Audio) PlayMusic(track *music.Track) error {
	if track == nil {
		return nil
	}

	a.currentTrackName = path.Join(a.tempFolder, track.Name+".mid")
	if _, err := os.Stat(a.currentTrackName); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		ioutil.WriteFile(a.currentTrackName, track.MidiStream.Bytes(), 0644)
	}

	var err error
	if a.currentTrack, err = mix.LoadMUS(a.currentTrackName); err != nil {
		return err
	}

	return a.currentTrack.FadeIn(-1, 1000)
}

// Play simply plays an audio chunk with the given name
func (a *Audio) Play(name string) error {
	chunk, err := a.getChunk(name)
	if err != nil {
		return err
	}
	_, err = chunk.Play(-1, 0)
	return err
}

// PlayAtPosition plays a sound in 2D virtual space using a given distance and angle.
func (a *Audio) PlayAtPosition(name string, distance float32, angle int16) error {
	chunk, err := a.getChunk(name)
	if err != nil {
		return err
	}
	if distance > 255 {
		distance = 255
	}
	channel, err := chunk.Play(-1, 0)
	mix.SetPosition(channel, angle, uint8(distance))
	return err
}

func (a *Audio) getChunk(name string) (*mix.Chunk, error) {
	if chunk, ok := a.chunks[name]; ok {
		return chunk, nil
	}
	return a.createChunk(name)
}

func (a *Audio) createChunk(name string) (*mix.Chunk, error) {
	sound, ok := sfx.Sounds(*a.sounds)[name]
	if !ok {
		return nil, fmt.Errorf("%s not found", name)
	}
	rwOps, err := sdl.RWFromMem(sound.ToWAV())
	if err != nil {
		return nil, err
	}
	chunk, err := mix.LoadWAVRW(rwOps, false)
	// chunk, err := mix.QuickLoadWAV(sound.ToWAV())
	if err != nil {
		return nil, fmt.Errorf("could not load WAV: %s", err.Error())
	}
	a.chunks[name] = chunk
	return chunk, nil
}

// Close closes the mixer and quits the SDL audio driver.
func (a *Audio) Close() {
	defer mix.CloseAudio()
	defer sdl.AudioQuit()
	t := time.Now()
	fadeOutDur := time.Second

	fmt.Printf("waiting for audio channels to stop: #0")
	for {
		n := mix.Playing(-1)
		// Wait up to 500 ms for playing channels when in non-test mode.
		if time.Now().Sub(t) > fadeOutDur || n == 0 || a.test {
			fmt.Printf(", OK\n")
			break
		}
		fmt.Printf("\b%d", n)
		time.Sleep(fadeOutDur / 10)
	}
}
