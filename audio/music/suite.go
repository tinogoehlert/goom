package music

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/tinogoehlert/goom/wad"
)

// Suite is a suite of named Tracks.
type Suite map[string]*Track

// NewSuite creates a new music.Suite.
func NewSuite() Suite {
	return make(Suite)
}

// LoadWAD loads the music data from the WAD and returns it
// as playble music tracks.
func (suite Suite) LoadWAD(w *wad.WAD) error {
	var (
		midiRegex = regexp.MustCompile(`^D_`)
		lumps     = w.Lumps()
	)
	for i := 0; i < len(lumps); i++ {
		l := lumps[i]
		switch {
		case midiRegex.Match([]byte(l.Name)):
			t, err := NewTrack(l)
			if err != nil {
				fmt.Printf("failed to load track: %s, err: %s\n", t.Name, err)
			}
			suite[l.Name] = t
		}
	}
	return nil
}

// Info shows a summary of the loaded tracks.
func (suite Suite) Info() string {
	var text []string
	for _, t := range suite {
		text = append(text, fmt.Sprintf("%s (%d): %v", t.Name, t.Size, t.MidiStream.Info()))
	}
	return strings.Join(text, "\n")
}

// Track returns a specific MusicTrack.
func (suite Suite) Track(name string) *Track {
	if t, ok := suite["D_"+name]; ok {
		return t
	}
	fmt.Println("invalid music track", name)
	return nil
}
