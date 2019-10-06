package audio

import (
	mus "github.com/tinogoehlert/goom/audio/mus"
	"github.com/tinogoehlert/goom/wad"
)

// MusicTrack contains a playable Music track.
type MusicTrack struct {
	wad.Lump
	MusData *mus.Data
}

// Play plays the MusicTrack.
func (*MusicTrack) Play() {}

// Loop plays the MusicTrack forever.
func (*MusicTrack) Loop() {}

// Stop stops playing the MusicTrack.
func (*MusicTrack) Stop() {}
