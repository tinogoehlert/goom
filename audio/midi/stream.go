package midi

import (
	"fmt"

	"github.com/tinogoehlert/goom/audio/mus"
)

// Data stores parsed MIDI events.
type Data struct {
	rawBytes []byte
	Events   []Event
}

// Stream describes MIDI data stream.
type Stream struct {
	Data
	// Channels provides an index for the used MIDI channels.
	channels   map[int]Channel
	velocities map[Channel]byte
	volumes    map[Channel]byte
	time       uint32
}

// NewStream returns a MIDI parser.
func NewStream(numChannels int) *Stream {
	return &Stream{
		channels:   make(map[int]Channel, numChannels),
		velocities: make(map[Channel]byte, numChannels),
		volumes:    make(map[Channel]byte, numChannels),
	}
}

func NewStreamFromBytes(bytes []byte) *Stream {
	return &Stream{
		Data: Data{
			rawBytes: bytes,
		},
	}
}

// Bytes returns the MIDI bytes for a mid file.
func (d *Data) Bytes() []byte {
	if d.rawBytes != nil {
		return d.rawBytes
	}

	var data []byte
	for _, ev := range d.Events {
		data = append(data, ev.Bytes()...)
	}

	d.rawBytes = append(append(MidHeader(), TrackLength(data)...), data...)
	return d.rawBytes
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
func (s *Stream) InitChan(ch Channel) {
	s.velocities[ch] = 100
	s.volumes[ch] = 127
	// turn all notes off on channel, write 0x7b, 0
	// s.Add(ChangeController, ch, byte(AllNotesOff), 0)
}

// GetChannel acquires and returns a (new) MIDI channel.
// If the channel is not used yet, the channel is initialized.
func (s *Stream) GetChannel(num int) Channel {
	ch, ok := s.channels[num]
	if !ok {
		if num == mus.PercussionChannel {
			ch = PercussionChannel
		} else {
			ch = Channel(len(s.channels))
		}
		s.channels[num] = ch
		s.InitChan(ch)
	}
	return ch
}

// GetVelocity returns the velocity byte for a channel.
func (s *Stream) GetVelocity(ch Channel) byte {
	return s.velocities[ch]
}

// SetVelocity sets the velocity byte for a channel.
func (s *Stream) SetVelocity(ch Channel, vel byte) {
	s.velocities[ch] = vel
}

// SetTime sets the delay for the next event.
func (s *Stream) SetTime(time int) {
	s.time = uint32(time)
}

// CompleteTrack completes a track using a Noop event if required.
func (s *Stream) CompleteTrack() {
	s.Events = append(s.Events, Event{
		Delay: s.time,
		Data:  []byte{0xff, 0x2f, 0x00},
	})
	s.time = 0
}

// Add encodes and adds a MIDI event to the MIDI data.
// Supports upto two payload values `mid1` (required) and `mid2` (optional).
// Additional values are ignored.
func (s *Stream) Add(ev EventType, ch Channel, mid1 byte, mid2 ...byte) {
	data := append([]byte{byte(ev) | byte(ch), mid1}, mid2...)
	s.Events = append(s.Events, Event{
		Delay: s.time,
		Data:  data,
	})
	s.time = 0
}
