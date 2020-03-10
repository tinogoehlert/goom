package sfx

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"regexp"
	"time"

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
