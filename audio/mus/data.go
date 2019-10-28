package mus

import (
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

// Info returns summarized header information as string.
func (h *Data) Info() string {
	// create dummy copy to safely remove not-logged data
	c := Data(*h)
	c.Events = nil
	events := make([]EventType, len(h.Events))
	for i, s := range h.Events {
		events[i] = s.Type
	}
	return fmt.Sprintf("mus.Data: %+v (%d events: %v)", c, len(h.Events), events)
}

// Simulate simulates Playback and checks the validity
// of all played Events.
func (h *Data) Simulate() error {
	var on map[uint8]bool
	var errs []string
	// turn off all notes
	off := func() { on = make(map[uint8]bool, 127) }
	off()
	err := func(ev Event, format string, v ...interface{}) {
		err := fmt.Sprintf("%s[%d](%s): %s", ev.Name(), ev.Index, ev.Hex(), fmt.Sprintf(format, v...))
		errs = append(errs, err)
	}
	for _, ev := range h.Events {
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
