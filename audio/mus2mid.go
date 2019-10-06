package audio

import (
	midi "github.com/tinogoehlert/goom/audio/midi"
	mus "github.com/tinogoehlert/goom/audio/mus"
)

const (
	numChans = 16
)

// MIDI control values for MUS control keys
var midControls = []midi.Control{
	// MidControl        MIDI controller  MUS value
	midi.Undefined,       //  N/A              0
	midi.BankSelect,      //  0 or 32          1
	midi.ModulationWheel, //  1                2
	midi.Volume,          //  7                3
	midi.PanPot,          //  10               4
	midi.ExpressionPot,   //  11               5
	midi.ReverbPot,       //  91               6
	midi.ChorusDepth,     //  93               7
	midi.SustainPedal,    //  64               8
	midi.SoftPedal,       //  67               9
	midi.AllSoundsOff,    //  120              10
	midi.AllNotesOff,     //  123              11
	midi.Mono,            //  126              12
	midi.Poly,            //  127              13
	midi.Reset,           //  121              14
}

// Mus2Mid converst MUS data to MIDI data.
func Mus2Mid(in *mus.Data) *midi.Data {
	out := midi.NewData(numChans)

	for _, s := range in.Scores {
		ch := out.GetChannel(s.Channel)
		vel := out.GetVelocity(ch)

		switch s.Type {
		case mus.RelaseNote:
			out.WriteReleaseKey(ch, s.Data[0])
		case mus.PlayNote:
			if len(s.Data) == 2 {
				vel = s.Data[1] & 0x7f
				out.SetVelocity(ch, vel)
			}
			out.WritePressKey(ch, s.Data[0], vel)
		case mus.PitchBend:
			out.WritePitchWheel(ch, byte(s.Data[0]*64))
		case mus.SystemEvent:
			ctrl := s.Data[0]
			if ctrl >= 10 && ctrl <= 14 {
				out.WriteController(ch, midControls[ctrl], 0)
			}
		case mus.Controller:
			ctrl := s.Data[0]
			val := s.Data[1]
			if ctrl == 0 {
				out.WriteChangePatch(ch, val)
				break
			}
			if ctrl >= 1 && ctrl <= 9 {
				out.WriteController(ch, midControls[ctrl], val)
			}
		case mus.MeasureEnd:
		case mus.ScoreEnd:
		}
		if s.Delay > 0 || s.Type == mus.ScoreEnd {
			// adjust and write delay and complete track
			out.Delay += s.Delay
			out.WriteEndTrack()
		}
	}

	return out
}
