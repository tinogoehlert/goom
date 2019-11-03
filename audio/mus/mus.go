package mus

import (
	"encoding/binary"
	"fmt"
)

// LumpID identifies MUS data.
var LumpID = "MUS\x1a"

// PercussionChannel is the channel used for perscussion events.
const PercussionChannel = 15

// EventType defines the type of a MUS event.
type EventType byte

// Event types.
const (
	// Event                    Byte 1         Byte 2     Description
	RelaseNote EventType = 0 // [0,note]                  Stops playing the note on a channel.
	PlayNote   EventType = 1 // [vol?,note]    [0,volume] Play note and optionally set the volume if vol? bit is 1.
	PitchBend  EventType = 2 // [bend amount]             Bend all notes on a channel by -1(0), -½(64), 0(128), +½(192), +1(255) tones.
	System     EventType = 3 // [0,controller]            Used for OPL2 (see: http://www.shikadi.net/moddingwiki/MUS_Format)
	Controller EventType = 4 // [0,controller] [0,value]  Change controller for channel (skipped if SystemEvent is used)
	MeasureEnd EventType = 5 //                           End current musical measure reached (does not affect playback).
	ScoreEnd   EventType = 6 //                           Last event in a song.
	Unused     EventType = 7 // [empty]                   Not used.
)

// Event names.
var eventNames = map[EventType]string{
	RelaseNote: "RelaseNote",
	PlayNote:   "PlayNote",
	PitchBend:  "PitchBend",
	System:     "System",
	Controller: "Controller",
	MeasureEnd: "MeasureEnd",
	ScoreEnd:   "ScoreEnd",
	Unused:     "Unused",
}

// Control defines a MIDI Controller number used by MUS events.
type Control uint8

// MUS Controller Numbers
const (
	//                           //  MIDI Controller
	ChangeInstr     Control = 0  //  N/A (Event 0xC0)
	BankSelect      Control = 1  //  0 or 32
	ModulationWheel Control = 2  //  1
	Volume          Control = 3  //  7
	PanPot          Control = 4  //  10
	ExpressionCtrl  Control = 5  //  11
	ReverbDepth     Control = 6  //  91
	ChorusDepth     Control = 7  //  93
	DamperPedal     Control = 8  //  64
	SoftPedal       Control = 9  //  67
	AllSoundsOff    Control = 10 //  120
	AllNotesOff     Control = 11 //  123
	MonoOn          Control = 12 //  126
	PolyOn          Control = 13 //  127
	ResetAllCtrl    Control = 14 //  121
)

// HeaderStart find the MUS header bytes.
// This is required since some MUS files start with garbage.
func HeaderStart(data []byte) int {
	if len(data) < 32 {
		return 0
	}
	for i := 0; i < len(data)-4; i++ {
		if string(data[i:4]) == LumpID {
			return i
		}
	}
	return 0
}

// ParseInt converts the first two bytes from the given data to a unit16 (long) value.
func ParseInt(data []byte) int {
	return int(binary.LittleEndian.Uint16(data))
}

// ParseInstruments reads `numInstr` words as `[]int`.
func ParseInstruments(data []byte, numInstr int) (inst []int, err error) {
	if len(data) < numInstr {
		return nil, fmt.Errorf("too few instrument bytes: %x", data)
	}
	for i := 0; i < numInstr; i += 2 {
		inst = append(inst, ParseInt(data[i:]))
	}
	return inst, nil
}
