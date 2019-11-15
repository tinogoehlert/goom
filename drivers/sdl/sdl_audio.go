package sdl

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/tinogoehlert/goom/audio/music"
	"github.com/tinogoehlert/goom/audio/sfx"

	"github.com/tinogoehlert/go-sdl2/mix"

	"github.com/tinogoehlert/go-sdl2/sdl"
)

type SoundManager struct {
	sounds           sfx.Sounds
	chunks           map[string]*mix.Chunk
	currentTrackName string
	currentTrack     *mix.Music
}

func InitSDLAudio() error {
	return sdl.InitSubSystem(sdl.INIT_AUDIO)
}

func QuitSDLAudio() {
	sdl.QuitSubSystem(sdl.INIT_AUDIO)
}

func NewSoundManager(sounds sfx.Sounds) (*SoundManager, error) {
	if _, err := mix.OpenAudioDevice(22050, mix.DEFAULT_FORMAT, 2, 4096, "", sdl.AUDIO_ALLOW_ANY_CHANGE); err != nil {
		log.Println(err)
		return nil, err
	}
	return &SoundManager{
		sounds: sounds,
		chunks: make(map[string]*mix.Chunk),
	}, nil
}

func (sm *SoundManager) PlayMusic(m *music.Track) error {
	sm.currentTrackName = m.Name + ".mid"
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

	if err := sm.currentTrack.FadeIn(-1, 1000); err != nil {
		return err
	}
	for mix.PlayingMusic() {
		sdl.Delay(40)
	}
	return nil
}

func (sm *SoundManager) Play(name string) error {
	chunk, err := sm.getChunk(name)
	if err != nil {
		return err
	}
	_, err = chunk.Play(-1, 0)
	for mix.Playing(-1) > 0 {
		sdl.Delay(60)
	}

	return err
}

func (sm *SoundManager) PlayAtPosition(name string, distance float32, angle int16) error {
	chunk, err := sm.getChunk(name)
	if err != nil {
		return err
	}
	channel, err := chunk.Play(-1, 0)
	mix.SetPosition(channel, angle, uint8(distance))
	return err
}

func (sm *SoundManager) getChunk(name string) (*mix.Chunk, error) {
	if chunk, ok := sm.chunks[name]; ok {
		return chunk, nil
	}
	return sm.createChunk(name)
}

func (sm *SoundManager) createChunk(name string) (*mix.Chunk, error) {
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

func (sm *SoundManager) Close() {
	mix.CloseAudio()
	sdl.AudioQuit()
}

/*
package sdl

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/tinogoehlert/goom/audio/sfx"

	"github.com/tinogoehlert/go-sdl2/mix"

	"github.com/tinogoehlert/go-sdl2/mix"
)

type SoundManager struct {
	sounds sfx.Sounds
	chunks
}

func InitSDLAudio() error {
	return sdl.InitSubSystem(sdl.INIT_AUDIO)
}

func QuitSDLAudio() {
	sdl.QuitSubSystem(sdl.INIT_AUDIO)
}

func NewSoundManager(sounds sfx.Sounds) (*SoundManager, error) {
	if err := mix.OpenAudio(11025, mix.DEFAULT_FORMAT, 1, 8); err != nil {
		log.Println(err)
		return nil, err
	}

	bytes, _ := ioutil.ReadFile("DSPSITOL.WAV")
	//bytes2, _ := ioutil.ReadFile("DSPISTOL")

	chunk, err := mix.QuickLoadWAV(bytes)
	log.Println(err)

	chunk2, err := mix.QuickLoadRAW(&bytes[0], uint32(len(bytes)))
	log.Println(err)

	fmt.Println(chunk.Play(-1, 0))
	// Wait until it finishes playing
	sdl.Delay(120)
	fmt.Println(chunk2.Play(-1, 0))
	for mix.Playing(-1) > 1 {
		sdl.Delay(40)
	}

}

func (sm *SoundManager) Close() {
	mix.CloseAudio()
	sdl.AudioQuit()
}
*/
