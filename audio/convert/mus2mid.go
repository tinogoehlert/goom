package convert

import (
	midi "github.com/tinogoehlert/goom/audio/midi"
	mus "github.com/tinogoehlert/goom/audio/mus"
)

const (
	numChans = 16
)

// MIDI control values for MUS control keys
var midControls = []midi.Control{
	// MidController          MIDI value       MUS value
	midi.Undefined,       //  N/A              0
	midi.BankSelect,      //  0 or 32          1
	midi.ModulationWheel, //  1                2
	midi.Volume,          //  7                3
	midi.PanPot,          //  10               4
	midi.ExpressionCtrl,  //  11               5
	midi.ReverbDepth,     //  91               6
	midi.ChorusDepth,     //  93               7
	midi.DamperPedal,     //  64               8
	midi.SoftPedal,       //  67               9
	midi.AllSoundsOff,    //  120              10
	midi.AllNotesOff,     //  123              11
	midi.MonoOn,          //  126              12
	midi.PolyOn,          //  127              13
	midi.ResetAllCtrl,    //  121              14
}

// MUS/MID file specs:
// http://www.shikadi.net/moddingwiki/MID_Format
// http://www.shikadi.net/moddingwiki/MUS_Format
// https://github.com/AyrA/WADex/blob/master/WADex/MUS2MID.cs
// https://github.com/sirjuddington/SLADE/tree/master/thirdparty/mus2mid
// https://github.com/madame-rachelle/qzdoom/blob/newmaster/src/sound/midisources/midisource_mus.cpp

// Mus2Mid converst MUS data to MIDI data.
func Mus2Mid(in *mus.Data) *midi.Data {
	p := midi.NewParser(numChans)

	for _, ev := range in.Events {
		ch := p.GetChannel(int(ev.Channel))

		// fmt.Printf("converting event: %s", ev.Info())
		// continue

		switch ev.Type {
		case mus.RelaseNote:
			p.Add(midi.ReleaseKey, ch, ev.GetNote())
		case mus.PlayNote:
			p.Add(midi.PressKey, ch, ev.GetNote(), ev.GetVolume())
		case mus.PitchBend:
			v := ev.GetBend()
			mid1 := (v & 1) << 6
			mid2 := (v >> 1) & 0x7F
			p.Add(midi.PitchWheel, ch, mid1, mid2)
		case mus.System:
			ctrl := ev.GetController()
			mid1 := uint8(midControls[ctrl])
			mid2 := uint8(0)
			if ctrl == 12 {
				mid2 = uint8(in.Channels)
			}
			p.Add(midi.ChangeController, ch, mid1, mid2)
		case mus.Controller:
			ctrl := ev.GetController()
			val := ev.GetControllerValue()
			if ctrl == 0 {
				p.Add(midi.ChangePatch, ch, val)
				break
			}
			midCtrl := midControls[ctrl]
			if midCtrl == midi.Volume {
				// TODO: Check if clamp volume to 127 is required.
				// see https://github.com/madame-rachelle/qzdoom/blob/5fa7520fd7b499aa3b7e3b939deb72920a294a6b/src/sound/midisources/midisource_mus.cpp#L324
			}
			p.Add(midi.ChangeController, ch, uint8(midCtrl), val)

		case mus.MeasureEnd:
		case mus.ScoreEnd:
		}
		p.SetTime(int(ev.Delay))
	}
	p.CompleteTrack()

	return &p.Data
}
