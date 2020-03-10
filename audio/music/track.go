package music

import (
	"fmt"

	"github.com/tinogoehlert/goom/audio/convert"
	midi "github.com/tinogoehlert/goom/audio/midi"
	mus "github.com/tinogoehlert/goom/audio/mus"
	"github.com/tinogoehlert/goom/wad"
)

// ellips dumps the first `limit` bytes of the data in hex format.
func ellips(data []byte, limit int) string {
	if len(data) <= limit {
		return fmt.Sprintf("%x", data)
	}
	return fmt.Sprintf("%x...", data[:limit])
}

// head dumps the first 100 bytes of the data in hex format.
func head(data []byte) string { return ellips(data, 100) }

// Track contains a playable Music track.
type Track struct {
	wad.Lump
	MidiStream *midi.Stream
	MusStream  *mus.Stream
}

// NewTrack loads MUS bytes as music.Track.
func NewTrack(lump wad.Lump) (*Track, error) {
	header := lump.Data[0:4]
	var (
		mi  *midi.Stream
		mu  *mus.Stream
		err error
	)

	switch string(header) {
	case "MUS\x1a":
		mu, err = mus.NewMusStream(lump.Data)
		if err != nil {
			return nil, err
		}
		mi, err = convert.Mus2Mid(mu)
		if err != nil {
			return nil, err
		}
	case "MThd":
		mi = midi.NewStreamFromBytes(lump.Data)
	}

	return &Track{lump, mi, mu}, nil
}

// Validate checks the track for errors.
func (t *Track) Validate() error {
	return t.MusStream.Simulate()
}
