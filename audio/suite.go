package audio

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/tinogoehlert/goom/wad"
)

// MusicSuite is a suite of named MusicTracks.
type MusicSuite map[string]*MusicTrack

// NewMusicSuite creates a new MusicStore
func NewMusicSuite() MusicSuite {
	return make(MusicSuite)
}

// LoadWAD loads the music data from the WAD and returns it
// as playble music tracks.
func (suite MusicSuite) LoadWAD(w *wad.WAD) error {
	var (
		midiRegex = regexp.MustCompile(`^D_`)
		lumps     = w.Lumps()
	)
	for i := 0; i < len(lumps); i++ {
		l := lumps[i]
		switch {
		case midiRegex.Match([]byte(l.Name)):
			m, err := NewMusData(l.Data)
			t := &MusicTrack{l, m}
			if err != nil {
				fmt.Printf("failed to load MUS track: %s, err: %s\n", t.Name, err)
			}
			suite[l.Name] = t
		}
	}
	return nil
}

// Info shows a summary of the loaded tracks.
func (suite MusicSuite) Info() string {
	var text []string
	for _, t := range suite {
		text = append(text, fmt.Sprintf("%s (%d): %v", t.Name, t.Size, t.MusData.Info()))
	}
	return strings.Join(text, "\n")
}