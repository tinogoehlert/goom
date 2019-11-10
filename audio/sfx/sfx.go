package sfx

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"path"
	"regexp"
	"time"

	pa "github.com/gordonklaus/portaudio"
	"github.com/tinogoehlert/goom/wad"
)

// SoundID defines the first 4 bytes of a DOOM sound.
const SoundID = "0300"

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

// Get returns a Sound sample by name.
func (s Sounds) Get(name string) *Sound {
	return s["DS"+name]
}

// LoadWAD loads the sound data from the WAD.
func (s Sounds) LoadWAD(w *wad.WAD) error {
	sfxRegex := regexp.MustCompile(`^DS`)
	for _, l := range w.Lumps() {
		if sfxRegex.Match([]byte(l.Name)) {
			if hex.EncodeToString(l.Data[:2]) != SoundID {
				return fmt.Errorf("invalid DS header for LUMP %s: %x", l.Name, l.Data[:4])
			}
			snd := &Sound{l}
			s[l.Name] = snd
			if len(snd.Data) < 32 {
				return fmt.Errorf("too few bytes for %s: %x", snd.Name, len(snd.Data))
			}
		}
	}
	return nil
}

// SampleRate returns the bits per sample of the Sound.
func (s *Sound) SampleRate() int {
	return 8
}

// SampleFreq returns the sample frequency in Hz.
func (s *Sound) SampleFreq() int {
	return int(binary.LittleEndian.Uint16(s.Data[2:]))
}

// NumSamples returns the number of PCM samples that define the Sound.
func (s *Sound) NumSamples() int {
	return len(s.SampleBytes()) * 8 / s.SampleRate()
}

// SampleBytes returns the PCM bytes without the header.
func (s *Sound) SampleBytes() []byte {
	return s.Data[8:]
}

// Duration returns the duration of the sound.
func (s *Sound) Duration() time.Duration {
	numMillis := int64(1000 * float32(s.NumSamples()) / float32(11025))
	return time.Millisecond * time.Duration(numMillis)
}

// Info decsribes the Sound.
func (s *Sound) Info() string {
	head := fmt.Sprintln(hex.Dump(s.Data[:32]))
	dur := s.Duration()
	return fmt.Sprintf(
		"Sound(name=%s, bits=%d, num=%d, size=%d, dur=%s)\n%s",
		s.Name, s.SampleRate(), s.NumSamples(), len(s.SampleBytes()), dur,
		head)
}

// Play plays a Sound using portaudio.
func Play(s *Sound) error {
	fmt.Println("playing sound:", s.Name, s.Size)
	fmt.Println(s.Info())

	if err := pa.Initialize(); err != nil {
		return err
	}
	defer pa.Terminate()
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

	/*
		size := len(data)

		readBytes := func(in, out []byte) {
			// head := hex.Dump(in[:16])
			// fmt.Printf("reading %d bytes from pos=%d:\n%s\n", len(out), chunkStart, head)
			for i := range out {
				if i < len(in) {
					// copy PCM byte
					out[i] = in[i]
				} else {
					// fill buffer with silence
					out[i] = 0x7f
				}
			}
		}
	*/

	bufSize := len(data) // 8192
	// out := make([]byte, bufSize)
	if test {
		for i := range data {
			data[i] = 127 // silence
		}
	}

	freq := float64(s.SampleFreq())

	stream, err := pa.OpenDefaultStream(0, 1, freq, bufSize, &data)
	if err == pa.DeviceUnavailable {
		fmt.Println("skipping unavailable device:", dev.Name)
		return nil
	}
	if err != nil {
		return err
	}

	stream.Start()
	defer stream.Stop()
	defer stream.Close()

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
		w, err := wad.NewWADFromFile(path.Join("..", "..", "DOOM1.WAD"))
		if err != nil {
			return err
		}
		sounds = make(Sounds)
		if err := sounds.LoadWAD(w); err != nil {
			return err
		}
		if len(sounds) == 0 {
			return fmt.Errorf("no sounds loaded")
		}
		fmt.Printf("loaded %d sounds to test cache\n", len(sounds))
	}

	for _, n := range names {
		if err := Play(sounds.Get(n)); err != nil {
			return err
		}
	}
	return nil
}
