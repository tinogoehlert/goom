package audio

import (
	"encoding/binary"
	"fmt"

	mus "github.com/tinogoehlert/goom/audio/mus"
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
	BankSelect      = Control(0x00)
	ModulationWheel = Control(0x01)
	Volume          = Control(0x07)
	PanPot          = Control(0x0A)
	ExpressionCtrl  = Control(0x0B)
	ReverbDepth     = Control(0x5B)
	ChorusDepth     = Control(0x5D)
	DamperPedal     = Control(0x40)
	SoftPedal       = Control(0x43)
	AllSoundsOff    = Control(0x78)
	AllNotesOff     = Control(0x7B)
	MonoOn          = Control(0x7E)
	PolyOn          = Control(0x7F)
	ResetAllCtrl    = Control(0x79)
	Undefined       = Control(0x0F)
)

// Event defines a MIDI event.
type Event struct {
	DeltaTime uint32   // MIDI ticks between previous and current event
	StreamID  uint32   // Reserved (0).
	Event     uint32   // Encoded event including event code and parameters
	Params    []uint32 // Additional parameters (not used)
}

// Channel is the integer number of a MIDI channel.
type Channel int

// PercussionChannel is the number of the MIDI channel
// used for percussion events.
const PercussionChannel = Channel(9)

// Data stores parsed MIDI events.
type Data struct {
	Delay  int
	Events []Event
}

// Parser describes MIDI data.
type Parser struct {
	Data
	// Channels provides an index for the used MIDI channels.
	channels   map[int]Channel
	velocities map[Channel]byte
	volumes    map[Channel]byte
	time       uint32
}

// NewParser returns a MIDI parser.
func NewParser(numChannels int) *Parser {
	return &Parser{
		channels:   make(map[int]Channel, numChannels),
		velocities: make(map[Channel]byte, numChannels),
		volumes:    make(map[Channel]byte, numChannels),
	}
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

// Bytes returns the MIDI bytes for a mid file.
func (d *Data) Bytes() []byte {
	// TODO: convert events to bytes
	var data []byte
	for _, ev := range d.Events {
		md := make([]byte, 12)
		binary.LittleEndian.PutUint32(md[0:], ev.DeltaTime)
		binary.LittleEndian.PutUint32(md[4:], ev.StreamID)
		binary.LittleEndian.PutUint32(md[8:], ev.Event)
		data = append(data, md...)
	}
	return append(append(MidHeader(), TrackLength(data)...), data...)
}

// Info returns summarized header information as string.
func (d *Data) Info() string {
	n := len(d.Events)
	if n > 10 {
		n = 10
	}
	return fmt.Sprintf("midi Events: %x", d.Events[:n])
}

// InitChan initializes the given MIDI channel
// and resets all instruments.
func (p *Parser) InitChan(ch Channel) {
	p.velocities[ch] = 100
	p.volumes[ch] = 127
	// turn all notes off on channel, write 0x7b, 0
	p.Add(ChangeController, ch, byte(AllNotesOff))
}

// GetChannel acquires and returns a (new) MIDI channel.
// If the channel is not used yet, the channel is initialized.
func (p *Parser) GetChannel(num int) Channel {
	ch, ok := p.channels[num]
	if !ok {
		if num == mus.PercussionChannel {
			ch = PercussionChannel
		} else {
			ch = Channel(len(p.channels))
		}
		p.channels[num] = ch
		p.InitChan(ch)
	}
	return ch
}

// GetVelocity returns the velocity byte for a channel.
func (p *Parser) GetVelocity(ch Channel) byte {
	return p.velocities[ch]
}

// SetVelocity sets the velocity byte for a channel.
func (p *Parser) SetVelocity(ch Channel, vel byte) {
	p.velocities[ch] = vel
}

// SetTime sets the delay for the next event.
func (p *Parser) SetTime(time int) {
	p.time = uint32(time)
}

// CompleteTrack completes a track using a Noop event if required.
func (p *Parser) CompleteTrack() {
	if p.time > 0 {
		p.Events = append(p.Events, Event{
			DeltaTime: p.time,
			Event:     uint32(Noop) << 24,
		})
	}
	p.time = 0
}

// Add encodes and adds a MIDI event to the MIDI data.
// Supports upto two payload values `mid1` (required) and `mid2` (optional).
// Additional values are ignored.
func (p *Parser) Add(ev EventType, ch Channel, mid1 uint8, mid2 ...uint8) {
	event := uint32(ev) | uint32(ch) | uint32(mid1)<<8
	if len(mid2) > 0 {
		event |= uint32(mid2[0]) << 16
	}
	p.Events = append(p.Events, Event{
		DeltaTime: p.time,
		Event:     event,
	})
	p.time = 0
}

/*
// WriteData append the data to the MIDI data.
func (md *Data) WriteData(data []byte) {
	md.Data = append(md.Data, data...)
}

// WriteByte writes a byte.
func (md *Data) WriteByte(b byte) {
	md.Data = append(md.Data, b)
}

// WriteEventByte writes an events byte.
func (md *Data) WriteEventByte(ev EventType, ch Channel) {
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
*/
