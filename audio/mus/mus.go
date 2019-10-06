package audio

import (
	"encoding/binary"
	"fmt"
)

// LumpID identifies MUS data.
var LumpID = "MUS\x1a"

// PercussionChannel is the channel used for perscussion events.
const PercussionChannel = 15

// Event defines the type of a MUS event.
type Event byte

// Long is a unit16.
type Long uint16

// Event types.
const (
	// Event                    Byte 1         Byte 2     Description
	RelaseNote  Event = iota // [0,note]                  Stops playing the note on a channel.
	PlayNote                 // [vol?,note]    [0,volume] Play note and optionally set the volume if vol? bit is 1.
	PitchBend                // [bend amount]             Bend all notes on a channel by -1(0), -½(64), 0(128), +½(192), +1(255) tones.
	SystemEvent              // [0,controller]            Used for OPL2 (see: http://www.shikadi.net/moddingwiki/MUS_Format)
	Controller               // [0,controller] [0,value]  Change controller for channel (skipped if SystemEvent is used)
	MeasureEnd               //                           End current musical measure reached (does not affect playback).
	ScoreEnd                 //                           Last event in a song.
	Unused                   // [empty]                   Not used.
)

// Score describes a Musical Scores to play.
type Score struct {
	Type    Event  // MUS event type
	Index   int    // position of the score bytes in the track
	Channel int    // Channel number
	Delay   int    // computed delay in ticks
	Data    []byte // 0-2 payload bytes
}

// Data represents the header of a MUS track.
type Data struct {
	ID          string  // 4-byte Music identifier "MUS" 0x1A
	ScoreLen    Long    // size of the MUS body
	ScoreStart  Long    // start of the MUS body
	Channels    Long    // Number of primary channels (excl. percussion channel 15)
	SecChannels Long    // Number of secondary channels
	NumInstr    Long    // Number of instruments
	Dummy       Long    // Separator between header and instruments list
	Instruments []Long  // List of used instruments (can be used to load sound patches, etc.)
	Scores      []Score // The actual music notes, pauses, etc.
}

// NewLong converts the first two bytes from the given data to a unit16 (long) value.
func NewLong(data []byte) Long {
	return Long(binary.LittleEndian.Uint16(data))
}

// Info returns summarized header information as string.
func (h *Data) Info() string {
	// create dummy copy to safely remove not-logged data
	c := Data(*h)
	c.Instruments = nil
	c.Scores = nil
	events := make([]Event, len(h.Scores))
	for i, s := range h.Scores {
		events[i] = s.Type
	}
	return fmt.Sprintf("mus.Data (summary): %+v (%d scores: %v)", c, len(h.Scores), events)
}

// NextIndex returns the byte index of the next score or delay byte,
// based on this Score's byte index.
func (s *Score) NextIndex() int {
	return s.Index + len(s.Data) + 1
}
