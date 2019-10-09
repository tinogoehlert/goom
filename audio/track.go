package audio

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	midi "github.com/tinogoehlert/goom/audio/midi"
	mus "github.com/tinogoehlert/goom/audio/mus"
	"github.com/tinogoehlert/goom/wad"
)

func ellips(data []byte, limit int) string {
	if len(data) <= limit {
		return fmt.Sprintf("%x", data)
	}
	return fmt.Sprintf("%x...", data[:limit])
}

func head(data []byte) string { return ellips(data, 100) }

func saveTestFile(file string, data []byte) (string, error) {
	dir := os.Getenv("DOOM_TEST")
	if dir == "" {
		dir = "."
	}
	f := filepath.Join(dir, file)
	os.Remove(f)
	return f, ioutil.WriteFile(f, data, 0644)
}

// MusicTrack contains a playable Music track.
type MusicTrack struct {
	wad.Lump
	MidiData *midi.Data
	MusData  *mus.Data
}

// Play plays the MusicTrack.
func (*MusicTrack) Play() {}

// Loop plays the MusicTrack forever.
func (*MusicTrack) Loop() {}

// Stop stops playing the MusicTrack.
func (*MusicTrack) Stop() {}

// SaveMus saves Track Lump as a MUS file.
func (t *MusicTrack) SaveMus() error {
	name := strings.ReplaceAll(t.Name, " ", "_")
	musfile := fmt.Sprintf("test_%s.mus", name)
	data := t.Data
	f, err := saveTestFile(musfile, data)
	if err != nil {
		return err
	}
	fmt.Printf("MUS %s: %s    written as &%s\n", t.Name, head(t.Data), f)
	return nil
}

// SaveMidi saves Track MIDI data as a MID file.
func (t *MusicTrack) SaveMidi() error {
	name := strings.ReplaceAll(t.Name, " ", "_")
	midfile := fmt.Sprintf("test_%s.mid", name)
	data := t.MidiData.Bytes()
	f, err := saveTestFile(midfile, data)
	if err != nil {
		return err
	}
	fmt.Printf("MID %s: %s\n    written as file %s\n", name, head(data), f)
	return nil
}

// Validate checks the track for errors.
func (t *MusicTrack) Validate() error {
	return t.MusData.Validate()
}
