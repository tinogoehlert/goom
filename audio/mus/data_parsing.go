package mus

import (
	"encoding/binary"
	"fmt"
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
