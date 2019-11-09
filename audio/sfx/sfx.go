package sfx

import (
	"encoding/hex"
	"fmt"
	"io"
	"path"
	"time"

	pa "github.com/gordonklaus/portaudio"
	"github.com/tinogoehlert/goom/wad"
)

var test bool

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
	return s[name]
}

// SampleRate returns the bits per sample of the Sound.
func (s *Sound) SampleRate() int {
	return 8
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
	fmt.Println("playing sound:", s.Info())

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

	stream, err := pa.OpenDefaultStream(0, 1, 11025, bufSize, &data)
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
	w, err := wad.NewWADFromFile(path.Join("..", "..", "DOOM1.WAD"))
	if err != nil {
		return err
	}

	for _, n := range names {
		s := &Sound{*w.Lump("DS" + n)}
		if hex.EncodeToString(s.Data[:4]) != "0300112b" {
			return fmt.Errorf("invalid DS header for LUMP: %v", s)
		}
		// fmt.Printf("loaded %d sfx bytes:\n", len(s.Data))
		if err := Play(s); err != nil {
			return err
		}
	}
	return nil
}
