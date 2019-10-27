package mus

import (
	"encoding/hex"
	"fmt"
)

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

// Event describes a generic music event.
type Event struct {
	Type       EventType // MUS event type
	Index      int       // position of the event's bytes in the track
	NextIndex  int       // position of the next event's bytes in the track
	Channel    uint8     // Channel number
	Delay      uint16    // computed delay in ticks
	Byte       byte      // source byte of the event
	Data       []byte    // 0-2 payload bytes
	DelayBytes []byte    // 0-n delay bytes
}

// Name returns the name of the event.
func (ev *Event) Name() string {
	return eventNames[ev.Type]
}

// Hex returns the source bytes of the event in hex-format.
func (ev *Event) Hex() string {
	hx := hex.EncodeToString(append([]byte{ev.Byte}, ev.Data...))
	if len(ev.DelayBytes) > 0 {
		hx = hx + " " + hex.EncodeToString(ev.DelayBytes)
	}
	return hx
}

// GetNote returns the Note value (0-127) for Play and Release events.
func (ev *Event) GetNote() uint8 {
	return ev.Data[0] & 0x7F
}

// HasVolume checks the note byte of a PlayNote event for the Volume flag.
func (ev *Event) HasVolume() bool {
	return ev.Data[0]>>7 == 1
}

// GetVolume returns the volume value (0-127).
func (ev *Event) GetVolume() uint8 {
	if ev.HasVolume() {
		return ev.Data[1]
	}
	return 0
}

// GetBend returns the bend value (0-255).
func (ev *Event) GetBend() uint8 {
	return ev.Data[0]
}

// GetController returns the controller number (0-127).
func (ev *Event) GetController() uint8 {
	return ev.Data[0]
}

// GetControllerValue returns the controller value (0-127).
func (ev *Event) GetControllerValue() uint8 {
	return ev.Data[1]
}

// Info descibes the event.
func (ev *Event) Info() string {
	var id, val uint8
	valueBits := 0
	switch ev.Type {
	case RelaseNote:
		id = ev.GetNote()
	case PlayNote:
		id = ev.GetNote()
		val = ev.GetVolume()
		valueBits = 7
	case PitchBend:
		val = ev.GetBend()
		valueBits = 8
	case System:
		id = ev.GetController()
	case Controller:
		id = ev.GetController()
		val = ev.GetControllerValue()
		valueBits = 7
	}
	return fmt.Sprintf("mus.Event[%d](%d, %s, %d, val=%d, bits=%d, ch=%d, delay=%d, hex=%s)",
		ev.Index, ev.Type, ev.Name(), id, val, valueBits, ev.Channel, ev.Delay, ev.Hex())
}
