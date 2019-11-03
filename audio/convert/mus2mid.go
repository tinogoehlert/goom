package convert

import (
	"fmt"

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

// ClampVolume limits MUS volume values to 127 and logs all clamps.
func ClampVolume(vol uint8) uint8 {
	if vol > 127 {
		fmt.Printf("clamping MUS volume = %d to max MIDI volume = 127", vol)
		vol = 127
	}
	return vol
}

// Mus2Mid converst MUS data to MIDI data.
func Mus2Mid(in *mus.Stream) (*midi.Stream, error) {
	p := midi.NewStream(numChans)

	if len(in.Events) == 0 {
		fmt.Println("skiping to parse empty MUS stream:", in.ID)
		return nil, nil
	}

	var ev mus.Event

	for _, ev = range in.Events {
		ch := p.GetChannel(int(ev.Channel))

		switch ev.Type {
		case mus.RelaseNote:
			p.Add(midi.ReleaseKey, ch, ev.GetNote(), 0)
		case mus.PlayNote:
			var vol byte
			if ev.HasVolume() {
				// use event volume
				vol = ev.GetVolume()
				// set it as new channel volume for notes without volume
				p.SetVelocity(ch, vol)
			} else {
				// use the last volume for playing the note
				vol = p.GetVelocity(ch)
			}
			vol = ClampVolume(vol)
			p.Add(midi.PressKey, ch, ev.GetNote(), vol)
		case mus.PitchBend:
			bend := ev.GetBend()
			// Allowed MUS Bend Values:
			//   0  one tone down
			//  64  half-tone down
			// 128  normal (no bend)
			// 192  half-tone up
			// 255  one tone up

			// Scale up to MIDI Bend Range: [0:16384]
			wheel := uint16(bend) * 64
			// encode LSB and MSB of pitch value
			mid1 := byte(wheel & 127)
			mid2 := byte(wheel >> 7 & 127)
			if w := uint16(mid2)<<7 | uint16(mid1); w != wheel {
				return nil, fmt.Errorf("invalid wheel=%d, expected=%d", w, wheel)
			}
			p.Add(midi.PitchWheel, ch, mid1, mid2)
		case mus.System:
			ctrl := ev.GetController()
			mid1 := byte(midControls[ctrl])
			mid2 := byte(0)
			if mus.Control(ctrl) == mus.MonoOn {
				mid2 = byte(in.Channels)
			}
			p.Add(midi.ChangeController, ch, mid1, mid2)
		case mus.Controller:
			ctrl := ev.GetController()
			mctrl := byte(midControls[ctrl])
			val := ev.GetControllerValue()
			if mus.Control(ctrl) == mus.Volume {
				// track the channel volumes as given
				p.SetVelocity(ch, val)
				// only clamp volumes when sending them to the MIDI out
				val = ClampVolume(val)
			}
			if mus.Control(ctrl) == mus.ChangeInstr {
				p.Add(midi.ChangePatch, ch, val)
				break
			}
			p.Add(midi.ChangeController, ch, mctrl, val)
		case mus.MeasureEnd:
		case mus.ScoreEnd:
			p.CompleteTrack()
		default:
			return nil, fmt.Errorf("uknown event: %s", ev.Info())
		}
		p.SetTime(int(ev.Delay))
	}

	if ev.Type != mus.ScoreEnd {
		fmt.Println("MUS stream does not end with ScoreEnd, but with:", ev.Info())
		fmt.Println("forcefully completing track")
		p.CompleteTrack()
	}
	return p, nil
}
