package audio

import (
	mus "github.com/tinogoehlert/goom/audio/mus"
)

// Event defines a MIDI event type.
type Event byte

// Event types:
const (
	ReleaseKey        = Event(0x80)
	PressKey          = Event(0x90)
	AfterTouchKey     = Event(0xA0)
	ChangeController  = Event(0xB0)
	ChangePatch       = Event(0xC0)
	AfterTouchChannel = Event(0xD0)
	PitchWheel        = Event(0xE0)
)

// Control defines a MIDI controller type.
type Control byte

// MidControl types
// For more details, read the MIDI docs:
// - https://www.midi.org/specifications-old/item/table-3-control-change-messages-data-bytes-2
// - http://www.personal.kent.edu/~sbirch/Music_Production/MP-II/MIDI/midi_control_change_messages.htm
const (
	BankSelect      = Control(0x00)
	ModulationWheel = Control(0x01)
	Volume          = Control(0x07)
	PanPot          = Control(0x0A)
	ExpressionPot   = Control(0x0B)
	ReverbPot       = Control(0x5B)
	ChorusDepth     = Control(0x5D)
	SustainPedal    = Control(0x40)
	SoftPedal       = Control(0x43)
	AllSoundsOff    = Control(0x78)
	AllNotesOff     = Control(0x7B)
	Mono            = Control(0x7E)
	Poly            = Control(0x7F)
	Reset           = Control(0x79)
	Undefined       = Control(0x0F)
)

// Channel is the integer number of a MIDI channel.
type Channel int

// PercussionChannel is the number of the MIDI channel
// used for percussion events.
const PercussionChannel = Channel(9)

// Data describes MIDI data.
type Data struct {
	// Channels provides an index for the used MIDI channels.
	channels   map[int]Channel
	velocities map[Channel]byte
	Delay      int
	Data       []byte
}

// NewData creates a MIDI data stub.
func NewData(numChannels int) *Data {

	midiHeader := []byte("MThd" + // Header start
		"\x00\x00\x00\x06" + // Header size
		"\x00\x00" + // MIDI type (0, single track)
		"\x00\x01" + // Number of tracks
		"\x00\x46" + // Resolution
		"MTrk" + // Track start
		"\x00\x00\x00\x00", // Track length placeholder
	)

	return &Data{
		Data:       midiHeader,
		channels:   make(map[int]Channel, numChannels),
		velocities: make(map[Channel]byte, numChannels),
	}
}

// InitChan initializes the given MIDI channel
// and resets all instruments.
func (md *Data) InitChan(ch Channel) {
	md.velocities[ch] = 127
	// turn all notes off, write 0x7b, 0
	md.WriteController(ch, AllNotesOff, 0)
}

// GetChannel acquires and returns a (new) MIDI channel.
// If the channel is not used yet, the channel is initialized.
func (md *Data) GetChannel(num int) Channel {
	ch, ok := md.channels[num]
	if !ok {
		if num == mus.PercussionChannel {
			ch = PercussionChannel
		} else {
			ch = Channel(len(md.channels))
		}
		md.channels[num] = ch
		md.InitChan(ch)
	}
	return ch
}

// GetVelocity returns the velocity byte for a channel.
func (md *Data) GetVelocity(ch Channel) byte {
	return md.velocities[ch]
}

// SetVelocity sets the velocity byte for a channel.
func (md *Data) SetVelocity(ch Channel, vel byte) {
	md.velocities[ch] = vel
}

// WriteData append the data to the MIDI data.
func (md *Data) WriteData(data []byte) {
	md.Data = append(md.Data, data...)
}

// WriteByte writes a byte.
func (md *Data) WriteByte(b byte) {
	md.Data = append(md.Data, b)
}

// WriteEventByte writes an events byte.
func (md *Data) WriteEventByte(ev Event, ch Channel) {
	md.WriteByte(byte(ev) | byte(ch))
}

// WriteDataByte writes a data byte.
func (md *Data) WriteDataByte(data byte, mod byte) {
	if mod == 0 {
		mod = 0x7F
	}
	md.WriteByte(data & mod)
}

// WriteTime writes delay bytes and resets the running delay.
func (md *Data) WriteTime(time int) (resetTime bool, bytesWritten int) {
	buffer := byte(time) & 0x7F
	var writeval byte

	for {
		time >>= 7
		if time == 0 {
			break
		}
		buffer <<= 8
		buffer |= byte((time & 0x7F) | 0x80)
	}

	for {
		writeval = byte(buffer & 0xFF)

		md.WriteByte(writeval)

		if (buffer & 0x80) != 0 {
			buffer >>= 8
		} else {
			md.Delay = 0
			return
		}
	}
	return
}

// WriteReleaseKey writes a release note.
func (md *Data) WriteReleaseKey(ch Channel, key byte) {
	md.WriteTime(md.Delay)
	md.WriteEventByte(ReleaseKey, ch)
	md.WriteDataByte(key, 0)
	md.WriteByte(0)
}

// WriteController changes a controller.
func (md *Data) WriteController(ch Channel, ctrl Control, val byte) {
	md.WriteTime(md.Delay)
	md.WriteEventByte(ChangeController, ch)
	md.WriteDataByte(byte(ctrl), 0)
	if (val & 0x80) != 0 {
		val = 0x7F
	}
	md.WriteByte(val)

}

// WritePressKey writes a note.
func (md *Data) WritePressKey(ch Channel, velocity, key byte) {
	md.WriteTime(md.Delay)
	md.WriteEventByte(PressKey, ch)
	md.WriteDataByte(key, 0)
	md.WriteDataByte(md.GetVelocity(ch), 0)
}

// WritePitchWheel writes a wheel pitch.
func (md *Data) WritePitchWheel(ch Channel, wheel byte) {
	md.WriteTime(md.Delay)
	md.WriteEventByte(PitchWheel, ch)
	md.WriteDataByte(wheel, 0)
	md.WriteByte(wheel >> 7)
}

// WriteChangePatch changes the patch.
func (md *Data) WriteChangePatch(ch Channel, patch byte) {
	md.WriteTime(md.Delay)
	md.WriteEventByte(ChangePatch, ch)
	md.WriteDataByte(patch, 0)
}

// WriteEndTrack writes the "end track" bytes and updates the tracklength
func (md *Data) WriteEndTrack() {
	endtrack := []byte{0xFF, 0x2F, 0x00}
	md.WriteTime(md.Delay)
	md.WriteData(endtrack)
	n := len(md.Data)
	md.Data[18+0] = byte(n>>24) & 0xff
	md.Data[18+1] = byte(n>>16) & 0xff
	md.Data[18+2] = byte(n>>8) & 0xff
	md.Data[18+3] = byte(n) & 0xff
}
