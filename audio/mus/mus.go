package mus

import (
	"encoding/binary"
	"fmt"
	"strings"
)

// LumpID identifies MUS data.
var LumpID = "MUS\x1a"

// PercussionChannel is the channel used for perscussion events.
const PercussionChannel = 15

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
	Events      []Event // The actual music notes, pauses, etc.
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

// ParseEvents parses the given bytes and converts them to a slice of MusScores.
func ParseEvents(data []byte) ([]Event, error) {
	events := make([]Event, 0, len(data)/2)
	// last := len(data)
	// if last > 10 { last = 10 }
	// fmt.Printf("loading MUS events: %x...\n", data[:last])
	for i := 0; i < len(data); {
		ev, err := NewEvent(i, data)
		if err != nil {
			return nil, err
		}
		events = append(events, *ev)
		i = ev.NextIndex
	}

	// fmt.Printf("events: %+v\n", events)
	return events, nil
}

// NewEvent creates a new events from the given MUS bytes.
func NewEvent(index int, data []byte) (*Event, error) {
	// Read event byte that describes the music event and how
	// to interpret subsequent payload and delay bytes.
	b := data[index]

	// Byte   Bits      Shift  Mask  Purpose
	// event  01110000  >> 4   0x07  MusType bit mask
	// event  00001111         0x0F  Channel bit mask
	// delay  01111111         0x7F  Delay   bit mask
	ev := Event{
		Type:    EventType((b >> 4) & 0x07),
		Index:   index,
		Channel: int(b & 0x0F),
	}
	delayBit := b >> 7 // delayBit shows if event is followed by delay bytes
	index++

	// read 0-2 subsequent payload bytes
	err := ev.ParsePayload(data[index:])
	if err != nil {
		return nil, err
	}
	index += len(ev.Data)

	// read 0-n subsequent delay bytes
	if delayBit == 1 {
		nd, err := ev.ParseDelay(data[index:])
		if err != nil {
			return nil, err
		}
		// advance index by number of delay bytes
		index += nd
	}
	ev.NextIndex = index
	return &ev, nil
}

// ParsePayload reads the payload bytes of an event.
func (ev *Event) ParsePayload(data []byte) error {
	// if len(data) > 2 { data = data[:2] }
	// fmt.Printf("ReadPayload (%d): \\x%x", ev, sample)

	switch ev.Type {
	case RelaseNote, PitchBend, System:
		ev.Data = data[0:1]
	case PlayNote:
		if data[0]>>7 == 0 {
			// has no volume flag and thus no volume byte
			ev.Data = data[0:1]
		}
		ev.Data = data[0:2]
	case Controller:
		ev.Data = data[0:2]
	case ScoreEnd, MeasureEnd, Unused:
		// payload is empty
	default:
		return fmt.Errorf("invalid event: %d", ev)
	}
	return nil
}

// ParseDelay reads delay bytes and computes the number of delay ticks.
func (ev *Event) ParseDelay(data []byte) (numDelayBytes int, err error) {
	delay := 0
	for i := 0; i < len(data); i++ {
		b := data[i]
		delay = delay*128 + int(b&127)
		if (b >> 7) == 0 {
			ev.Delay = delay
			return i + 1, nil
		}
	}
	return 0, fmt.Errorf("invalid delay bytes in MUS data")
}

// Info returns summarized header information as string.
func (h *Data) Info() string {
	// create dummy copy to safely remove not-logged data
	c := Data(*h)
	c.Instruments = nil
	c.Events = nil
	events := make([]EventType, len(h.Events))
	for i, s := range h.Events {
		events[i] = s.Type
	}
	return fmt.Sprintf("mus.Data (summary): %+v (%d events: %v)", c, len(h.Events), events)
}

// Validate simulates Playback and checks the validity
// of all played Events.
func (h *Data) Validate() error {
	var on map[uint8]bool
	var errors []error
	off := func() {
		on = make(map[uint8]bool, 127)
	}
	err := func(ev Event, format string, v ...interface{}) {
		v = append([]interface{}{ev.Name()}, v...)
		errors = append(errors, fmt.Errorf("%s: "+format, v...))
	}
	check127 := func(ev Event, pos int) {
		b := ev.Data[pos]
		if b > 127 {
			err(ev, "invalid byte: %x", b)
		}
	}
	off()
	for _, ev := range h.Events {
		fmt.Println(ev.Info())
		switch ev.Type {
		case RelaseNote:
			on[ev.GetNote()] = false
			check127(ev, 0)
		case PlayNote:
			on[ev.GetNote()] = true
			if ev.HasVolume() {
				check127(ev, 1)
			}
		case System:
			ctrl := ev.GetController()
			switch ctrl {
			case 10, 11:
				off()
			case 12, 13, 14, 15:
			default:
				err(ev, "invalid controller: %d", ctrl)
			}
		case Controller:
			check127(ev, 0)
			check127(ev, 1)
			ctrl := ev.GetController()
			switch {
			case ctrl > 15:
				err(ev, "invalid controller: %d", ctrl)
			}
		}
	}
	for note, state := range on {
		if state == true {
			errors = append(errors, fmt.Errorf("note not turned off: %d", note))
		}
	}
	if len(errors) > 0 {
		var texts []string
		for _, err := range errors {
			texts = append(texts, err.Error())
		}
		return fmt.Errorf("Invalid MUS events:\n%v", strings.Join(texts, "\n"))
	}
	return nil
}
