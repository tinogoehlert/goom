package sfx

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"regexp"
	"time"

	pa "github.com/gordonklaus/portaudio"
	"github.com/tinogoehlert/goom/wad"
)

// SoundID defines the first 4 bytes of a DOOM sound.
const SoundID = "0300"

// sfxPrefix is the prefix used for sound lump names.
var sfxPrefix string

// sfxRegex is the regex used to indentify sound lumps.
var sfxRegex *regexp.Regexp

// SetPrefix sets the sound lump prefix.
// Use this function to define different prefixes for different games.
func SetPrefix(p string) {
	sfxPrefix = p
	sfxRegex = regexp.MustCompile("^" + p)
}

func init() {
	SetPrefix("DS")
}

// vars for testing and development
var test bool
var sounds Sounds

// TestMode silences all sound and remove all delays.
// Use this for testing.
func TestMode() {
	test = true
}

// Sound stores PCM sound bytes.
type Sound struct {
	wad.Lump
}

// Sounds is a suite of named Tracks.
type Sounds map[string]*Sound

// Get returns a Sound sample by name for the current game.
func (s Sounds) Get(name string) *Sound {
	return s[sfxPrefix+name]
}

// GetByID returns a Sound sample by its full lump name.
func (s Sounds) GetByID(key string) *Sound {
	return s[key]
}

// LoadWAD loads the sound data from the WAD.
func (s Sounds) LoadWAD(w *wad.WAD) error {
	for _, l := range w.Lumps() {
		if sfxRegex.Match([]byte(l.Name)) {
			if hex.EncodeToString(l.Data[:2]) != SoundID {
				return fmt.Errorf("invalid sound header for LUMP %s: %x", l.Name, l.Data[:4])
			}
			s[l.Name] = &Sound{l}
		}
	}
	sounds = s
	return nil
}

// BitDepth returns the bits per sample of the Sound.
func (s *Sound) BitDepth() int {
	return 8
}

// SampleRate returns the sample frequency in Hz.
func (s *Sound) SampleRate() int {
	return int(binary.LittleEndian.Uint16(s.Data[2:]))
}

// NumSamples returns the number of PCM samples that define the Sound.
func (s *Sound) NumSamples() int {
	return len(s.SampleBytes()) * 8 / s.BitDepth()
}

// ToWAV returns sampleBytes with RIFF/WAV header
func (s *Sound) ToWAV() []byte {
	var (
		wavBuff = new(bytes.Buffer)
	)

	wavBuff.WriteString("RIFF")
	binary.Write(wavBuff, binary.LittleEndian, uint32(44+len(s.SampleBytes())-8))
	wavBuff.WriteString("WAVE")
	wavBuff.WriteString("fmt ")
	binary.Write(wavBuff, binary.LittleEndian, uint32(16))
	binary.Write(wavBuff, binary.LittleEndian, uint16(1))
	binary.Write(wavBuff, binary.LittleEndian, uint16(1))
	binary.Write(wavBuff, binary.LittleEndian, uint32(s.SampleRate()))
	binary.Write(wavBuff, binary.LittleEndian, uint32(s.SampleRate()))
	binary.Write(wavBuff, binary.LittleEndian, uint16(1))
	binary.Write(wavBuff, binary.LittleEndian, uint16(8))
	wavBuff.WriteString("data")
	binary.Write(wavBuff, binary.LittleEndian, uint32(len(s.SampleBytes())))
	wavBuff.Write(s.SampleBytes())

	// ioutil.WriteFile("new.wav", wavBuff.Bytes(), 0644)
	return wavBuff.Bytes()
}

// SampleBytes returns the PCM bytes without the header.
func (s *Sound) SampleBytes() []byte {
	return s.Data[8:]
}

// Duration returns the duration of the sound.
func (s *Sound) Duration() time.Duration {
	sampleTime := time.Second / time.Duration(s.SampleRate())
	return time.Duration(s.NumSamples()) * sampleTime
}

// Info decsribes the Sound.
func (s *Sound) Info() string {
	head := fmt.Sprintln(hex.Dump(s.Data[:32]))
	dur := s.Duration()
	return fmt.Sprintf(
		"Sound(name=%s, bits=%d, rate=%d, num=%d, size=%d, dur=%s)\n%s",
		s.Name, s.BitDepth(), s.SampleRate(), s.NumSamples(), len(s.SampleBytes()), dur,
		head)
}

func init() {
	if err := pa.Initialize(); err != nil {
		panic(err)
	}
}

// Play plays a Sound using portaudio.
func Play(s *Sound) error {
	fmt.Println("playing sound:", s.Name, s.Size)
	fmt.Println(s.Info())

	//defer pa.Terminate()
	api, err := pa.DefaultHostApi()
	if err != nil {
		return err
	}
	fmt.Println("using portaudio api:", api.Type)

	/*
		devs, err := pa.Devices()
		if err != nil {
			return err
		}
		for devNum, dev := range devs {
			fmt.Printf("found portaudio device %d: %s\n", devNum, dev.Name)
		}
	*/

	dev := api.DefaultOutputDevice
	fmt.Println("using portaudio device:", dev.Name)

	data := s.SampleBytes()

	if test {
		for i := range data {
			data[i] = 127 // silence
		}
	}

	rate := float64(s.SampleRate())
	stream, err := pa.OpenDefaultStream(0, 1, rate, pa.FramesPerBufferUnspecified, &data)
	if err == pa.DeviceUnavailable {
		fmt.Println("skipping unavailable device:", dev.Name)
		return nil
	}
	if err != nil {
		return err
	}

	stream.Start()
	//defer stream.Stop()
	//defer stream.Close()

	// fmt.Println("started playback on device:", dev.Name)
	// t := time.Now()
	if err := stream.Write(); err != nil && err != io.EOF {
		return err
	}
	if !test {
		time.Sleep(s.Duration())
	}
	// d := time.Now().Sub(t)
	// fmt.Printf("finished playback on device: %s, after: %s\n", dev.Name, d)

	return nil
}

// PlaySounds plays all given sounds.
func PlaySounds(names ...string) error {
	if len(sounds) == 0 {
		fmt.Println("no sounds loaded")
		return fmt.Errorf("no sounds loaded")
	}
	for _, n := range names {
		if err := Play(sounds.Get(n)); err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}
