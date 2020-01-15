package music

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/tinogoehlert/goom/wad"
)

// TrackStore is a suite of named Tracks.
type TrackStore map[string]*Track

// NewTrackStore creates a new music.TrackStore.
func NewTrackStore() TrackStore {
	return make(TrackStore)
}

// LoadWAD loads the music data from the WAD and returns it
// as playble music tracks.
func (s TrackStore) LoadWAD(w *wad.WAD) {
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
				fmt.Printf("failed to load track: %s, err: %s\n", l.Name, err)
			}
			s[l.Name] = t
		}
	}
}

// Info shows a summary of the loaded tracks.
func (s TrackStore) Info() string {
	var text []string
	for _, t := range s {
		text = append(text, fmt.Sprintf("%s (%d): %v", t.Name, t.Size, t.MidiStream.Info()))
	}
	return strings.Join(text, "\n")
}

// Track returns a specific MusicTrack.
func (s TrackStore) Track(name string) *Track {
	if t, ok := s["D_"+name]; ok {
		return t
	}
	fmt.Println("invalid music track", name)
	return nil
}
