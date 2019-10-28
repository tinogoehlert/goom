package mus

import (
	"fmt"
	"strings"
)

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
		Channel: uint8(b & 0x0f),
		Byte:    b,
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
		ev.DelayBytes = data[index : index+nd]
		index += nd
	}
	ev.NextIndex = index

	// fmt.Println(ev.Info())

	err = ev.Validate()
	if err != nil {
		return nil, err
	}

	return &ev, nil
}

// ParsePayload reads the payload bytes of an event.
func (ev *Event) ParsePayload(data []byte) error {
	switch ev.Type {
	case RelaseNote, PitchBend, System:
		ev.Data = data[0:1]
	case PlayNote:
		if data[0]>>7 == 0 {
			// has no volume flag and thus no volume byte
			ev.Data = data[0:1]
		} else {
			// has no volume flag and a volume byte
			ev.Data = data[0:2]
		}
	case Controller:
		ev.Data = data[0:2]
	case ScoreEnd, MeasureEnd, Unused:
		// payload is empty
	default:
		return fmt.Errorf("invalid event: %d", ev)
	}
	// fmt.Printf("parsed %s payload: %x\n", ev.Name(), ev.Data)
	return nil
}

// ParseDelay reads delay bytes and computes the number of delay ticks.
func (ev *Event) ParseDelay(data []byte) (numDelayBytes int, err error) {
	delay := uint16(0)
	for i := 0; i < len(data); i++ {
		b := data[i]
		delay = delay*128 + uint16(b&0x7f)
		if (b >> 7) == 0 {
			ev.Delay = delay
			return i + 1, nil
		}
	}
	return 0, fmt.Errorf("invalid delay bytes in MUS data")
}

// Validate validates the parsed event values.
func (ev *Event) Validate() error {
	var errs []string
	err := func(format string, v ...interface{}) {
		errs = append(errs, fmt.Sprintf(format, v...))
	}
	check127 := func(i int, msg string) {
		if ev.Data[i] > 127 {
			err("%s: ev.Data[%d] = %x", msg, i, ev.Data[i])
		}
	}

	switch ev.Type {
	case RelaseNote:
		check127(0, "invalid release note")
	case PlayNote:
		if ev.HasVolume() {
			if ev.Data[0] < 127 {
				err("invalid play note (with volume): ev.Data[%d] = %x", 0, ev.Data[0])
			}
			check127(1, "invalid volume")
		} else {
			check127(0, "invalid play note (without volume)")
		}
	case System:
		check127(0, "invalid system controller number")
		ctrl := uint8(ev.Data[0])
		if ctrl < 10 || ctrl > 15 {
			err("invalid system controller: %d", ctrl)
		}
	case Controller:
		ctrl := uint8(ev.Data[0])
		check127(0, "invalid controller number")
		if ctrl != 3 {
			// allow bigger vaules for volume controller
			check127(1, "invalid controller value")
		}
		if ctrl > 15 {
			err("invalid controller: %d", ctrl)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("Invalid %s\n%s", ev.Info(), strings.Join(errs, "\n"))
	}
	return nil
}
