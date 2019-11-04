package midi

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"time"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/gm"
	"gitlab.com/gomidi/midi/mid"
	"gitlab.com/gomidi/midi/midimessage/channel"
	"gitlab.com/gomidi/midi/midimessage/meta"
	"gitlab.com/gomidi/rtmididrv"
	"gitlab.com/ubunatic/portmididrv"
)

// TicksPerSecond defines the default number of MIDI ticks per second used for playback.
// Most DOOM titles use 140Hz. Raptor uses 70Hz.
// The value can be adjusted by setting `Player.TicksPerSecond` before playback.
const TicksPerSecond = 140

// Player defines a MIDI port and driver ana allows
// closing these.
type Player struct {
	out            mid.Out
	drv            mid.Driver
	wr             *mid.Writer
	TicksPerSecond int
}

var test = false

// Provider defines MIDI driver types.
type Provider string

// MIDI driver types
const (
	PortMidi = "ubunatic/portmididrv"
	RTMidi   = "gomidi/rtmididrv"
	Any      = "any"
)

func (p Provider) match(provider Provider) bool {
	return p == Any || p == provider
}

// TestMode disables all delays and sounds for unit testing.
func TestMode() {
	test = true
}

func (p *Player) initDriver(providers ...Provider) error {
	var errors []string

	if len(providers) == 0 {
		providers = []Provider{Any}
	}

	for _, pr := range providers {
		switch {
		case pr.match(RTMidi):
			fmt.Println("trying MIDI driver: gomidi/rtmididrv")
			if drv, err := rtmididrv.New(); err != nil {
				errors = append(errors, err.Error())
			} else {
				p.drv = drv
				return nil
			}
		case pr.match(PortMidi):
			fmt.Println("trying MIDI driver: ubunatic/portmididrv")
			if drv, err := portmididrv.New(); err != nil {
				errors = append(errors, err.Error())
			} else {
				p.drv = drv
				return nil
			}
		}
	}

	return fmt.Errorf("no driver found:\n%s", strings.Join(errors, "\n"))
}

// NewPlayer looks for common MIDI devices
// and returns an opened MIDI output suitable for playback.
// Make sure to `defer p.Close()` the device later.
func NewPlayer(providers ...Provider) (*Player, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("MIDI Player panicking:", r)
		}
	}()

	expPortTimidity := regexp.MustCompile(`TiMidity.*port.*[0-9]`)
	expPortMacSynth := regexp.MustCompile(`SimpleSynth virtual.*`)

	p := &Player{
		TicksPerSecond: TicksPerSecond,
	}
	if err := p.initDriver(providers...); err != nil {
		return nil, err
	}

	outs, err := p.drv.Outs()
	if err != nil {
		return nil, err
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
		p.drv.Close()
		return nil, fmt.Errorf("no playback device found")
	}

	if err := out.Open(); err != nil {
		return nil, err
	}

	p.wr = mid.ConnectOut(out)
	p.out = out

	return p, nil
}

// Close closes the underlying MIDI port and driver.
func (p *Player) Close() {
	defer p.Off()
	defer p.out.Close()
	defer p.drv.Close()
}

// Off silences all MIDI channels.
func (p *Player) Off() {
	if test {
		return
	}
	fmt.Println("muting all channels")
	for ch := 0; ch < 16; ch++ {
		if err := p.wr.Silence(int8(ch), false); err != nil {
			fmt.Printf("failed to silence channel %d: %s\n", ch, err.Error())
		}
	}
}

// HandleMessage plays MIDI messages.
func (p *Player) HandleMessage(pos *mid.Position, msg midi.Message) {
	if !p.out.IsOpen() {
		return
	}
	if test {
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
		if test {
			return
		}
		err = p.wr.Write(msg)
	}
	if err != nil {
		p.Off()
		panic(err)
	}
}

// ProgramChange implements workaround for ProgramChange messages
// for broken drivers.
func (p *Player) ProgramChange(msg channel.ProgramChange) error {
	switch {
	case test:
		return nil
	default:
		return gm.WriteGMProgram(p.wr, msg.Channel(), msg.Program())
	}
}

// Reset send common reset messages to the MIDI device.
func (p *Player) Reset() {
	fmt.Println("resetting MIDI player")
	if err := gm.WriteReset(p.wr, 0, 0); err != nil {
		fmt.Printf("failed to GM-reset player: %s", err)
	}
}

// Play plays a song.
func (p *Player) Play(stream *Stream) {
	p.Reset()
	Process(stream.Bytes(), p.HandleMessage)
}

// Loop plays a song repeatedly until the player closes.
func (p *Player) Loop(stream *Stream) {
	for p.out.IsOpen() {
		p.Play(stream)
	}
}

// MessageHandler processes MIDI messages.
type MessageHandler func(*mid.Position, midi.Message)

// Process plays the stream using the given message handler.
func Process(data []byte, fn MessageHandler) {
	b := bytes.NewReader(data)
	rd := mid.NewReader(mid.NoLogger())
	rd.Msg.Each = fn
	rd.ReadAllSMF(b)
}
