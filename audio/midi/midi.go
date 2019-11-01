package midi

import (
	"encoding/binary"
)

// EventType defines a MIDI event type.
type EventType byte

// EventTypes:
const (
	Noop              = EventType(0x02)
	ReleaseKey        = EventType(0x80)
	PressKey          = EventType(0x90)
	AfterTouchKey     = EventType(0xA0)
	ChangeController  = EventType(0xB0)
	ChangePatch       = EventType(0xC0)
	AfterTouchChannel = EventType(0xD0)
	PitchWheel        = EventType(0xE0)
)

// Control defines a MIDI controller type.
type Control byte

// MidControl types
// For more details, read the MIDI docs:
// - https://www.midi.org/specifications-old/item/table-3-control-change-messages-data-bytes-2
// - http://www.personal.kent.edu/~sbirch/Music_Production/MP-II/MIDI/midi_control_change_messages.htm
const (
	Undefined       Control = 15
	BankSelect      Control = 0
	ModulationWheel Control = 1
	Volume          Control = 7
	PanPot          Control = 10
	ExpressionCtrl  Control = 11
	ReverbDepth     Control = 91
	ChorusDepth     Control = 93
	DamperPedal     Control = 64
	SoftPedal       Control = 67
	AllSoundsOff    Control = 120
	AllNotesOff     Control = 123
	MonoOn          Control = 126
	PolyOn          Control = 127
	ResetAllCtrl    Control = 121
)

// Event defines a MIDI event.
type Event struct {
	Delay uint32 // MIDI ticks between previous and current event
	Data  []byte // Encoded event including event code and parameters
}

// Bytes returns the event data as MIDI bytes.
func (ev *Event) Bytes() []byte {
	return append(EncodeVarInt(ev.Delay), ev.Data...)
}

// Channel is the integer number of a MIDI channel.
type Channel int

// PercussionChannel is the number of the MIDI channel
// used for percussion events.
const PercussionChannel = Channel(9)

// EncodeVarInt encodes an uint32 as MIDI length.
func EncodeVarInt(v uint32) []byte {
	data := []byte{0, 0, 0, byte(v & 0x7f)}
	i := 3
	for i > 0 {
		if v < 128 {
			break
		}
		i--
		v >>= 7
		data[i] = byte(v&0x7f) | 0x80
	}
	return data[i:]
}

// DecodeVarInt decodes a MIDI length as uint32.
func DecodeVarInt(data []byte) (v uint32) {
	for i := 0; i < 4; i++ {
		v |= uint32(data[i] & 127)
		if data[i]&0x80 == 0 {
			break
		}
		v <<= 7
	}
	return
}

// MidHeader returns the generic MIDI header bytes.
func MidHeader() []byte {
	return []byte("MThd" + // Header start
		"\x00\x00\x00\x06" + // Header size
		"\x00\x00" + // MIDI type (0, single track)
		"\x00\x01" + // Number of tracks
		"\x00\x46" + // Resolution
		"MTrk", // Track start
	)
}

// TrackLength returns the length of a track as MIDI bytes.
func TrackLength(data []byte) []byte {
	tl := make([]byte, 4)
	binary.BigEndian.PutUint32(tl, uint32(len(data)))
	return tl
}
