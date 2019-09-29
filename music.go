package goom

import (
	"encoding/binary"
	"fmt"
	"regexp"
	"strings"
)

// MusID identifies MUS data.
var MusID = "MUS\x1a"

// MusEvent defines the type of the event
type MusEvent int

// MusEvent types.
const (
	RelaseNote = iota
	PlayNote
	PitchWheel
	SystemEvent
	ChangeController
	Unknown5
	ScoreEnd
	Unknown7
)

// MusData represents the header of a MUS track.
type MusData struct {
	ID          []byte // 4-byte Music identifier "MUS" 0x1A
	scoreLen    uint16 // size of the MUS body
	scoreStart  uint16 // start of the MUS body
	channels    uint16 // Number of primary channels (excl. percussion channel 15)
	secChannels uint16 // Number of secondary channels
	numInstr    uint16 // Number of instruments
	dummy       uint16 // ???
	instruments []uint16
	scores      []byte
}

// Info returns the header information as string.
func (h *MusData) Info() string {
	c := MusData(*h)
	c.instruments = nil
	c.scores = nil
	return fmt.Sprintf("%+v (%d scores bytes)", c, len(h.scores))
}

// MusicTrack contains a playable Music track.
type MusicTrack struct {
	Lump
	MusData *MusData
}

// NewMusData creates a MusHeader from the given WAD bytes.
func NewMusData(data []byte) (*MusData, error) {
	if data == nil {
		return &MusData{ID: []byte(MusID)}, nil
	}
	if len(data) < 16 {
		return nil, fmt.Errorf("failed to load bytes '%s' as MUS", data)
	}
	h := MusData{
		ID:          data[:4],
		scoreLen:    binary.LittleEndian.Uint16(data[4:]),
		scoreStart:  binary.LittleEndian.Uint16(data[6:]),
		channels:    binary.LittleEndian.Uint16(data[8:]),
		secChannels: binary.LittleEndian.Uint16(data[10:]),
		numInstr:    binary.LittleEndian.Uint16(data[12:]),
		dummy:       binary.LittleEndian.Uint16(data[14:]),
		instruments: nil,
		scores:      nil,
	}
	lastInst := int(16+2*h.numInstr) - 2
	for i := 16; i <= lastInst; i += 2 {
		h.instruments = append(h.instruments, binary.LittleEndian.Uint16(data[i:]))
	}
	h.scores = data[h.scoreStart:]
	if len(h.scores) != int(h.scoreLen) {
		return nil, fmt.Errorf(
			"wrong MUS data size: len(h.scores) = %d, h.scoreLen = %d",
			len(h.scores), h.scoreLen,
		)
	}
	return &h, nil
}

// MusicSuite is a suite of named MusicTracks.
type MusicSuite map[string]*MusicTrack

// LoadMusic loads the music data from the WAD and returns it
// as playble music tracks.
func (wm *WadManager) LoadMusic() (MusicSuite, error) {
	var (
		midiRegex = regexp.MustCompile(`^D_`)
	)
	suite := make(MusicSuite)
	for _, w := range wm.wads {
		for i := 0; i < len(w.lumps); i++ {
			l := w.lumps[i]
			switch {
			case midiRegex.Match([]byte(l.Name)):
				m, err := NewMusData(l.Data)
				t := &MusicTrack{l, m}
				if err != nil || !t.isMusic() {
					GoomConsole.Red("failed to load MUS track: %s, err: %s", t.Name, err)
				}
				suite[l.Name] = t
			}
		}
	}
	return suite, nil
}

// Play plays the MusicTrack.
func (*MusicTrack) Play() {}

// Loop plays the MusicTrack forever.
func (*MusicTrack) Loop() {}

// Stop stops playing the MusicTrack.
func (*MusicTrack) Stop() {}

func (t MusicTrack) isMusic() bool {
	return len(t.Data) > 0 && string(t.Data[:4]) == MusID
}

// Info shows a summary of the loaded tracks.
func (suite MusicSuite) Info() string {
	var text []string
	for _, t := range suite {
		text = append(text, fmt.Sprintf("%s (%d): %v", t.Name, t.Size, t.MusData.Info()))
	}
	return strings.Join(text, "\n")
}
