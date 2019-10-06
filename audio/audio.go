package audio

import (
	"fmt"

	mus "github.com/tinogoehlert/goom/audio/mus"
)

// NewMusData creates a MusHeader from the given WAD bytes.
func NewMusData(data []byte) (*mus.Data, error) {
	if data == nil {
		return &mus.Data{ID: mus.LumpID}, nil
	}
	id := string(data[:4])
	if len(data) < 16 || id != mus.LumpID {
		return nil, fmt.Errorf("failed to load bytes '%s' as MUS", data)
	}

	h := mus.Data{
		ID:          string(data[:4]),
		ScoreLen:    mus.NewLong(data[4:]),
		ScoreStart:  mus.NewLong(data[6:]),
		Channels:    mus.NewLong(data[8:]),
		SecChannels: mus.NewLong(data[10:]),
		NumInstr:    mus.NewLong(data[12:]),
		Dummy:       mus.NewLong(data[14:]),
		Instruments: nil,
		Scores:      nil,
	}
	scoreStart := int(h.ScoreStart)
	for i := 16; i < scoreStart; i += 2 {
		h.Instruments = append(h.Instruments, mus.NewLong(data[i:]))
	}
	// fmt.Printf("header: %s\n", h.Info())

	scores, err := LoadScores(data[h.ScoreStart:])
	if err != nil {
		return nil, err
	}
	h.Scores = scores

	return &h, nil
}

// LoadScores parses the given bytes and converts them to a slice of MusScores.
func LoadScores(data []byte) ([]mus.Score, error) {
	scores := make([]mus.Score, 0, len(data)/2)
	// last := len(data)
	// if last > 10 { last = 10 }
	// fmt.Printf("loading MUS Scores: %x...\n", data[:last])
	for i := 0; i < len(data); {
		b := data[i] // read score byte that describes the following event

		// bits      int  purpose
		// 01110000  112  MusType bit mask (requires shifting by 4 bits afterwards)
		// 00001111  15   Channel bit mask
		// 01111111  127  delay bit mask used for delay bytes
		delayBit := b >> 7 // delayBit shows if event is followed by delay bytes
		s := mus.Score{
			Type:    mus.Event((b & 112) >> 4),
			Index:   i,
			Channel: int(b & 15),
		}
		i++

		// read 0-2 subsequent payload bytes
		payload, err := ReadPayload(s.Type, data[i:])
		if err != nil {
			return nil, err
		}
		s.Data = payload
		i += len(s.Data)

		// read 0-n subsequent delay bytes
		if delayBit == 1 {
			delay, nd, err := ReadDelay(data[i:])
			if err != nil {
				return nil, err
			}
			s.Delay = delay
			// advance index by number of delay bytes
			i += nd
		}

		// fmt.Printf("append data[%d...] as score[%d]: %+v\n", i, scoreNum, s)
		scores = append(scores, s)
	}
	// fmt.Printf("scores: %+v\n", scores)
	return scores, nil
}

// ReadPayload reads the payload bytes of an event.
func ReadPayload(ev mus.Event, data []byte) (payload []byte, err error) {
	// if len(data) > 2 { data = data[:2] }
	// fmt.Printf("ReadPayload (%d): \\x%x", ev, sample)

	switch ev {
	case mus.RelaseNote, mus.PitchBend, mus.SystemEvent:
		payload = data[0:1]
	case mus.PlayNote:
		if data[0]>>7 == 0 {
			// has no volume flag and thus no volume byte
			payload = data[0:1]
		}
		payload = data[0:2]
	case mus.Controller:
		payload = data[0:2]
	case mus.ScoreEnd, mus.MeasureEnd, mus.Unused:
		// payload is empty
	default:
		err = fmt.Errorf("invalid event: %d", ev)
	}
	return
}

// ReadDelay reads delay bytes and computes the number of delay ticks.
func ReadDelay(data []byte) (value, numDelayBytes int, err error) {
	delay := 0
	for i := 0; i < len(data); i++ {
		b := data[i]
		delay = delay*128 + int(b&127)
		if (b >> 7) == 0 {
			return delay, i + 1, nil
		}
	}
	return 0, 0, fmt.Errorf("invalid delay bytes in MUS data")
}
