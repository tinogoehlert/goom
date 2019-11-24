package sdl

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/tinogoehlert/go-sdl2/mix"
	"github.com/tinogoehlert/go-sdl2/sdl"

	"github.com/tinogoehlert/goom/audio/music"
	"github.com/tinogoehlert/goom/audio/sfx"
)

type Audio struct {
	sounds           *sfx.Sounds
	chunks           map[string]*mix.Chunk
	currentTrackName string
	currentTrack     *mix.Music
	tempFolder       string
}

func NewAudio(sounds *sfx.Sounds, tempFolder string) (*Audio, error) {

	err := initAudio()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if _, err := mix.OpenAudioDevice(22050, mix.DEFAULT_FORMAT, 2, 4096, "", sdl.AUDIO_ALLOW_ANY_CHANGE); err != nil {
		log.Println(err)
		return nil, err
	}

	os.MkdirAll(tempFolder, 0700)

	a := &Audio{
		sounds:     sounds,
		chunks:     make(map[string]*mix.Chunk),
		tempFolder: tempFolder,
	}

	return a, nil
}

func (a Audio) PlayMusic(track *music.Track) error {
	if track == nil {
		return nil
	}

	a.currentTrackName = a.tempFolder + "/" + track.Name + ".mid"
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
func (a Audio) Play(name string) error {
	chunk, err := a.getChunk(name)
	if err != nil {
		return err
	}
	_, err = chunk.Play(-1, 0)
	return err
}

func (a Audio) PlayAtPosition(name string, distance float32, angle int16) error {
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

func (a Audio) Close() {
	mix.CloseAudio()
	sdl.AudioQuit()
}
