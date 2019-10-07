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
	ScoreLen    int     // size of the MUS body
	ScoreStart  int     // start of the MUS body
	Channels    int     // Number of primary channels (excl. percussion channel 15)
	SecChannels int     // Number of secondary channels
	NumInstr    int     // Number of instruments
	Dummy       int     // Separator between header and instruments list
	Instruments []int   // List of used instruments (can be used to load sound patches, etc.)
	Scores      []Score // The actual music notes, pauses, etc.
}

// ParseInt converts the first two bytes from the given data to a unit16 (long) value.
func ParseInt(data []byte) int {
	return int(binary.LittleEndian.Uint16(data))
}

// ParseIntruments reads `numInstr` words as `[]int`.
func ParseIntruments(data []byte, numInstr int) (inst []int, err error) {
	if len(data) < numInstr {
		return nil, fmt.Errorf("too few instrument bytes: %x", data)
	}
	for i := 0; i < numInstr; i += 2 {
		inst = append(inst, ParseInt(data[i:]))
	}
	return inst, nil
}

// ParseScores parses the given bytes and converts them to a slice of MusScores.
func ParseScores(data []byte) ([]Score, error) {
	scores := make([]Score, 0, len(data)/2)
	// last := len(data)
	// if last > 10 { last = 10 }
	// fmt.Printf("loading MUS Scores: %x...\n", data[:last])
	for i := 0; i < len(data); {
		b := data[i] // read score byte that describes the following event

		// bits      int  purpose
		// 01110000  112  MusType bit mask (requires shifting by 4 bits afterwards)
		// 00001111  15   Channel bit mask
		// 01111111  127  delay bit mask used for delay bytes
		delayBit := b >> 7 // delayBit shows if event is followed by delay bytes
		s := Score{
			Type:    Event((b & 112) >> 4),
			Index:   i,
			Channel: int(b & 15),
		}
		i++

		// read 0-2 subsequent payload bytes
		payload, err := ReadPayload(s.Type, data[i:])
		if err != nil {
			return nil, err
		}
		s.Data = payload
		i += len(s.Data)

		// read 0-n subsequent delay bytes
		if delayBit == 1 {
			delay, nd, err := ReadDelay(data[i:])
			if err != nil {
				return nil, err
			}
			s.Delay = delay
			// advance index by number of delay bytes
			i += nd
		}

		// fmt.Printf("append data[%d...] as score[%d]: %+v\n", i, scoreNum, s)
		scores = append(scores, s)
	}
	// fmt.Printf("scores: %+v\n", scores)
	return scores, nil
}

// ReadPayload reads the payload bytes of an event.
func ReadPayload(ev Event, data []byte) (payload []byte, err error) {
	// if len(data) > 2 { data = data[:2] }
	// fmt.Printf("ReadPayload (%d): \\x%x", ev, sample)

	switch ev {
	case RelaseNote, PitchBend, SystemEvent:
		payload = data[0:1]
	case PlayNote:
		if data[0]>>7 == 0 {
			// has no volume flag and thus no volume byte
			payload = data[0:1]
		}
		payload = data[0:2]
	case Controller:
		payload = data[0:2]
	case ScoreEnd, MeasureEnd, Unused:
		// payload is empty
	default:
		err = fmt.Errorf("invalid event: %d", ev)
	}
	return
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
