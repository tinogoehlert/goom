package gomididrv

import (
	"bytes"
	"fmt"
	"regexp"
	"time"

	gomidi "gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/gm"
	"gitlab.com/gomidi/midi/mid"
	"gitlab.com/gomidi/midi/midimessage/channel"
	"gitlab.com/gomidi/midi/midimessage/meta"

	"github.com/tinogoehlert/goom/audio/midi"
)

// DefaultTicksPerSecond defines the default number of MIDI ticks per second used for playback.
// Most DOOM titles use 140Hz. Raptor uses 70Hz.
// The value can be adjusted by setting `Player.TicksPerSecond` before playback.
const DefaultTicksPerSecond = 140

// MessageHandler processes MIDI messages.
type MessageHandler func(*mid.Position, gomidi.Message)

// Process plays the stream using the given message handler.
func Process(data []byte, fn MessageHandler) {
	b := bytes.NewReader(data)
	rd := mid.NewReader(mid.NoLogger())
	rd.Msg.Each = fn
	rd.ReadAllSMF(b)
}

// InitMidiOutput looks for common MIDI devices
// and initalizes a MIDI device that is suitable for playback.
// Make sure to `defer m.Close()` the device later.
func (p *MidiPlayer) InitMidiOutput() error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("MIDI Player panicking:", r)
		}
	}()

	expPortTimidity := regexp.MustCompile(`TiMidity.*port.*[0-9]`)
	expPortMacSynth := regexp.MustCompile(`SimpleSynth virtual.*`)

	outs, err := p.driver.Outs()
	if err != nil {
		return err
	}

	var out mid.Out

	for _, o := range outs {
		name := []byte(o.String())
		mtmy := expPortTimidity.Find(name)
		mmac := expPortMacSynth.Find(name)
		if mtmy != nil || mmac != nil {
			fmt.Printf("using device #%d: %s\n", o.Number(), o.String())
			out = o
			break
		}
		fmt.Printf("skipping unknown MIDI device #%d: %s\n", o.Number(), o.String())
	}

	if out == nil {
		fmt.Println("no known devices found")
	}

	if out == nil && len(outs) > 0 {
		i := 0
		if len(outs) > 1 {
			i = 1
		}
		out = outs[i]
		fmt.Printf("fallback to output #%d, MIDI device #%d: %s\n", i, out.Number(), out.String())
	}

	if out == nil {
		p.driver.Close()
		return fmt.Errorf("no playback device found")
	}

	if err := out.Open(); err != nil {
		return err
	}

	p.writer = mid.ConnectOut(out)
	p.out = out

	return nil
}

// Play plays a track once.
func (p *MidiPlayer) Play(s *midi.Stream) {
	p.Reset()
	Process(s.Bytes(), p.HandleMessage)
}

// Loop plays a song repeatedly until the player closes.
func (p *MidiPlayer) Loop(s *midi.Stream) {
	for p.out.IsOpen() {
		p.Play(s)
	}
}

// ProgramChange implements workaround for ProgramChange messages
// for broken drivers.
func (p *MidiPlayer) ProgramChange(msg channel.ProgramChange) error {
	switch {
	case p.test:
		return nil
	default:
		return gm.WriteGMProgram(p.writer, msg.Channel(), msg.Program())
	}
}

// Reset send common reset messages to the MIDI device.
func (p *MidiPlayer) Reset() {
	fmt.Println("resetting MIDI player")
	if err := gm.WriteReset(p.writer, 0, 0); err != nil {
		fmt.Printf("failed to GM-reset player: %s", err)
	}
}

// Off silences all MIDI channels.
func (p *MidiPlayer) Off() {
	if p.test {
		return
	}
	fmt.Println("muting all channels")
	for ch := 0; ch < 16; ch++ {
		if err := p.writer.Silence(int8(ch), false); err != nil {
			fmt.Printf("failed to silence channel %d: %s\n", ch, err.Error())
		}
	}
}

// HandleMessage plays MIDI messages.
func (p *MidiPlayer) HandleMessage(pos *mid.Position, msg gomidi.Message) {
	if !p.out.IsOpen() {
		return
	}
	if p.test {
		time.Sleep(time.Nanosecond)
	} else {
		delay := time.Second / time.Duration(p.TicksPerSecond) * time.Duration(pos.DeltaTicks)
		time.Sleep(delay)
	}
	var err error
	switch msg.(type) {
	case channel.ProgramChange:
		err = p.ProgramChange(msg.(channel.ProgramChange))
	default:
		if msg == meta.EndOfTrack {
			p.Off()
		}
		if p.test {
			return
		}
		err = p.writer.Write(msg)
	}
	if err != nil {
		p.Off()
		panic(err)
	}
}
