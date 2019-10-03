package audio

import (
	"encoding/binary"
	"fmt"
	"regexp"
	"strings"

	"github.com/tinogoehlert/goom/wad"
)

// MusID identifies MUS data.
var MusID = "MUS\x1a"

// MusEvent defines the type of the event
type MusEvent int

// MusEvent types.
const (
	// Event              Byte 1         Byte 2     Description
	RelaseNote  = iota // [0,note]                  Stops playing the note on a channel.
	PlayNote           // [vol?,note]    [0,volume] Play note and optionally set the volume if vol? bit is 1.
	PitchBend          // [bend amount]             Bend all notes on a channel by -1(0), -½(64), 0(128), +½(192), +1(255) tones.
	SystemEvent        // [0,controller]            Used for OPL2 (see: http://www.shikadi.net/moddingwiki/MUS_Format)
	Controller         // [0,controller] [0,value]  Change controller for channel (skipped if SystemEvent is used)
	MeasureEnd         //                           End current musical measure reached (does not affect playback).
	ScoreEnd           //                           Last event in a song.
	Unused             // [empty]                   Not used.
)

// MusScore describes a Musical Scores to play.
type MusScore struct {
	Type    MusEvent // MusEvent type
	Channel int      // Channel number
	Delay   int      // computed delay in ticks
	Byte1   byte     // first payload byte for the event
	Byte2   byte     // second payload byte for the event
}

// MusData represents the header of a MUS track.
type MusData struct {
	ID          []byte     // 4-byte Music identifier "MUS" 0x1A
	scoreLen    uint16     // size of the MUS body
	scoreStart  uint16     // start of the MUS body
	channels    uint16     // Number of primary channels (excl. percussion channel 15)
	secChannels uint16     // Number of secondary channels
	numInstr    uint16     // Number of instruments
	dummy       uint16     // Separator between header and instruments list
	instruments []uint16   // List of used instruments (ca be used to load sound patches, etc.)
	scores      []MusScore // The actual music notes, pauses, etc.
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
	wad.Lump
	MusData *MusData
}

// NewMusData creates a MusHeader from the given WAD bytes.
func NewMusData(data []byte) (*MusData, error) {
	if data == nil {
		return &MusData{ID: []byte(MusID)}, nil
	}
	id := string(data[:4])
	if len(data) < 16 || id != MusID {
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
	scores, err := LoadScores(data[h.scoreStart:])
	if err != nil {
		return nil, err
	}
	h.scores = scores
	return &h, nil
}

// LoadScores parses the given bytes and converts them to a slice of MusScores.
func LoadScores(data []byte) ([]MusScore, error) {
	scores := make([]MusScore, 0, len(data))
	scoreNum := 0
	for i := 0; i < len(data); i++ {
		// bits      int  purpose
		// 01110000  112  MusType bit mask (requires shifting by 4 bits afterwards)
		// 00001111  15   Channel bit mask
		// 01111111  127  delay bit mask used for delay bytes
		b := data[i]
		mtype := (b & 128) >> 4 // shift and mask mus type bits
		channel := b & 15       // mask channel
		last := b >> 7          // get delay bytes flag
		delay := 0

		// TODO: read payload bytes

		// read the subsequent delay bytes
		if last == 1 {
			var err error
			numDelayBytes := 0
			delay, numDelayBytes, err = ReadDelay(data[i+1:])
			if err != nil {
				return nil, err
			}
			i = i + numDelayBytes
		}

		s := MusScore{
			Type:    MusEvent(mtype),
			Channel: int(channel),
			Delay:   delay,
		}
		scores = append(scores, s)
	}
	return scores[:scoreNum], nil
}

// ReadDelay reads delay bytes and computes the number of delay ticks.
func ReadDelay(data []byte) (value, numDelayBytes int, err error) {
	delay := 0
	for i := 0; i < len(data); i++ {
		b := data[i]
		delay = delay*128 + int(b&127)
		if (b >> 7) == 0 {
			return delay, i + 1, nil
		}
	}
	return 0, 0, fmt.Errorf("invalid delay bytes in MUS data")
}

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

// Play plays the MusicTrack.
func (*MusicTrack) Play() {}

// Loop plays the MusicTrack forever.
func (*MusicTrack) Loop() {}

// Stop stops playing the MusicTrack.
func (*MusicTrack) Stop() {}

// Info shows a summary of the loaded tracks.
func (suite MusicSuite) Info() string {
	var text []string
	for _, t := range suite {
		text = append(text, fmt.Sprintf("%s (%d): %v", t.Name, t.Size, t.MusData.Info()))
	}
	return strings.Join(text, "\n")
}
