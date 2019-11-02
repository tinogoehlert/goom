package mus

import (
	"fmt"
	"strings"
)

// Stream represents the header of a MUS track.
type Stream struct {
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

// NewMusStream creates a MUS data from the given WAD bytes.
func NewMusStream(data []byte) (*Stream, error) {
	if data == nil {
		return &Stream{ID: LumpID}, nil
	}
	data = data[HeaderStart(data):]
	id := string(data[:4])
	if len(data) < 16 || id != LumpID {
		return nil, fmt.Errorf("failed to load bytes '%s' as MUS", data)
	}

	md := &Stream{
		ID:          string(data[:4]),
		ScoreLen:    ParseInt(data[4:]),
		ScoreStart:  ParseInt(data[6:]),
		Channels:    ParseInt(data[8:]),
		SecChannels: ParseInt(data[10:]),
		NumInstr:    ParseInt(data[12:]),
		Dummy:       ParseInt(data[14:]),
		Instruments: nil,
		Events:      nil,
	}
	inst, err := ParseInstruments(data[16:], md.NumInstr)
	if err != nil {
		return nil, err
	}
	md.Instruments = inst

	events, err := ParseEvents(data[md.ScoreStart:])
	if err != nil {
		return nil, err
	}
	md.Events = events

	return md, nil
}

// ParseEvents parses the given bytes and converts them to a slice of MusScores.
func ParseEvents(data []byte) ([]Event, error) {
	events := make([]Event, 0, len(data)/2)
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

// Info returns summarized header information as string.
func (s *Stream) Info() string {
	// create shallow copy to safely remove unlogged data
	c := Stream(*s)
	c.Events = nil
	events := make([]EventType, len(s.Events))
	for i, s := range s.Events {
		events[i] = s.Type
	}
	return fmt.Sprintf("mus.Data: %+v (%d events: %v)", c, len(c.Events), events)
}

// Simulate simulates Playback and checks the validity
// of all played Events.
func (s *Stream) Simulate() error {
	var on map[uint8]bool
	var errs []string
	// turn off all notes
	off := func() { on = make(map[uint8]bool, 127) }
	off()
	err := func(ev Event, format string, v ...interface{}) {
		err := fmt.Sprintf("%s[%d](%s): %s", ev.Name(), ev.Index, ev.Hex(), fmt.Sprintf(format, v...))
		errs = append(errs, err)
	}
	for _, ev := range s.Events {
		// fmt.Println(ev.Info())
		if verr := ev.Validate(); verr != nil {
			err(ev, "validation failed: %s", verr.Error())
		}
		switch ev.Type {
		case RelaseNote:
			on[ev.GetNote()] = false
		case PlayNote:
			on[ev.GetNote()] = true
		case System:
			ctrl := ev.GetController()
			switch ctrl {
			case 10, 11:
				off()
			case 12, 13, 14, 15:
				// add more valueless controller tests here
			}
		case Controller:
			// ctrl := ev.GetController()
			// add more valued controller tests here
		}
	}
	for note, state := range on {
		if state == true {
			errs = append(errs, fmt.Sprintf("note not turned off: %d", note))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("Invalid MUS events:\n%v", strings.Join(errs, "\n"))
	}
	return nil
}
