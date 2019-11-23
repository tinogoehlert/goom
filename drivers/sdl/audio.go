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

type AudioDriver struct {
	sounds           sfx.Sounds
	chunks           map[string]*mix.Chunk
	currentTrackName string
	currentTrack     *mix.Music
	tempFolder       string
}

func NewAudioDriver(sounds sfx.Sounds, tempFolder string) (*AudioDriver, error) {
	if err := sdl.InitSubSystem(sdl.INIT_AUDIO); err != nil {
		return nil, err
	}

	if _, err := mix.OpenAudioDevice(22050, mix.DEFAULT_FORMAT, 2, 4096, "", sdl.AUDIO_ALLOW_ANY_CHANGE); err != nil {
		log.Println(err)
		return nil, err
	}
	return &AudioDriver{
		sounds:     sounds,
		chunks:     make(map[string]*mix.Chunk),
		tempFolder: tempFolder,
	}, nil

}

func (sm *AudioDriver) PlayMusic(m *music.Track) error {
	if m == nil {
		return nil
	}
	sm.currentTrackName = sm.tempFolder + "/" + m.Name + ".mid"
	if _, err := os.Stat(sm.currentTrackName); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		ioutil.WriteFile(sm.currentTrackName, m.MidiStream.Bytes(), 0644)
	}

	var err error
	if sm.currentTrack, err = mix.LoadMUS(sm.currentTrackName); err != nil {
		return err
	}

	return sm.currentTrack.FadeIn(-1, 1000)
}

// Play simply plays an audio chunk with the given name
func (sm *AudioDriver) Play(name string) error {
	chunk, err := sm.getChunk(name)
	if err != nil {
		return err
	}
	_, err = chunk.Play(-1, 0)
	return err
}

func (sm *AudioDriver) PlayAtPosition(name string, distance float32, angle int16) error {
	chunk, err := sm.getChunk(name)
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

func (sm *AudioDriver) getChunk(name string) (*mix.Chunk, error) {
	if chunk, ok := sm.chunks[name]; ok {
		return chunk, nil
	}
	return sm.createChunk(name)
}

func (sm *AudioDriver) createChunk(name string) (*mix.Chunk, error) {
	sound, ok := sm.sounds[name]
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
	sm.chunks[name] = chunk
	return chunk, nil
}

func (sm *AudioDriver) Close() {
	mix.CloseAudio()
	sdl.AudioQuit()
	sdl.QuitSubSystem(sdl.INIT_AUDIO)
}
