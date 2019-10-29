package audio

import (
	"fmt"

	"github.com/tinogoehlert/goom/audio/convert"
	midi "github.com/tinogoehlert/goom/audio/midi"
	mus "github.com/tinogoehlert/goom/audio/mus"
)

// NewMusData creates a MUS data from the given WAD bytes.
func NewMusData(data []byte) (*mus.Data, error) {
	if data == nil {
		return &mus.Data{ID: mus.LumpID}, nil
	}
	data = data[mus.HeaderStart(data):]
	id := string(data[:4])
	if len(data) < 16 || id != mus.LumpID {
		return nil, fmt.Errorf("failed to load bytes '%s' as MUS", data)
	}

	md := &mus.Data{
		ID:          string(data[:4]),
		ScoreLen:    mus.ParseInt(data[4:]),
		ScoreStart:  mus.ParseInt(data[6:]),
		Channels:    mus.ParseInt(data[8:]),
		SecChannels: mus.ParseInt(data[10:]),
		NumInstr:    mus.ParseInt(data[12:]),
		Dummy:       mus.ParseInt(data[14:]),
		Instruments: nil,
		Events:      nil,
	}
	inst, err := mus.ParseInstruments(data[16:], md.NumInstr)
	if err != nil {
		return nil, err
	}
	md.Instruments = inst

	events, err := mus.ParseEvents(data[md.ScoreStart:])
	if err != nil {
		return nil, err
	}
	md.Events = events

	return md, nil
}

// NewMidiData creates MIDI data from the given WAD bytes.
func NewMidiData(data []byte) (*midi.Data, error) {
	md, err := NewMusData(data)
	if err != nil {
		return nil, err
	}
	return convert.Mus2Mid(md), nil
}