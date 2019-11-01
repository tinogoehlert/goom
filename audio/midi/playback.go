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
	"gitlab.com/gomidi/portmididrv"
	"gitlab.com/gomidi/rtmididrv"
)

// TicksPerSecond defines the default number of MIDI ticks per second used for playback.
// Most DOOM titles use 140Hz. Raptor uses 70Hz.
// The value can be adjusted by setting `Player.TicksPerSecond` before playback.
const TicksPerSecond = 140

// FixedMessage wraps a MIDI message and fixes bytes for portmididrv.
type FixedMessage struct {
	midi.Message
}

// Raw returns the fixed bytes.
func (m FixedMessage) Raw() []byte {
	b := m.Message.Raw()
	if len(b) == 2 {
		fmt.Println("fixing ProgramChange message for portmididrv")
		b = append(b, 0)
	}
	return b
}

func (m FixedMessage) String() string {
	return m.Message.String()
}

// Player defines a MIDI port and driver ana allows
// closing these.
type Player struct {
	out              mid.Out
	drv              mid.Driver
	wr               *mid.Writer
	TicksPerSecond   int
	useFixedMessages bool
}

var test = false

// Provider defines MIDI driver types.
type Provider string

// MIDI driver types
const (
	PortMidi = "gomidi/portmididrv"
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

func (d *Player) initDriver(providers ...Provider) error {
	var errors []string

	if len(providers) == 0 {
		providers = []Provider{Any}
	}

	for _, p := range providers {
		switch {
		case p.match(RTMidi):
			fmt.Println("trying MIDI driver: gomidi/rtmididrv")
			if drv, err := rtmididrv.New(); err != nil {
				errors = append(errors, err.Error())
			} else {
				d.drv = drv
				return nil
			}
		case p.match(PortMidi):
			fmt.Println("trying MIDI driver: gomidi/portmididrv")
			if drv, err := portmididrv.New(); err != nil {
				errors = append(errors, err.Error())
			} else {
				d.drv = drv
				d.useFixedMessages = true
				return nil
			}
		}
	}

	return fmt.Errorf("no driver found:\n%s", strings.Join(errors, "\n"))
}

// NewPlayer looks for common MIDI devices
// and returns an opened MIDI output suitable for playback.
// Make sure to `defer d.Close()` the device later.
func NewPlayer(providers ...Provider) (*Player, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("MIDI Player panicking:", r)
		}
	}()

	expPortTimidity := regexp.MustCompile(`TiMidity.*port.*[0-9]`)
	expPortMacSynth := regexp.MustCompile(`SimpleSynth virtual.*`)

	d := &Player{
		TicksPerSecond: TicksPerSecond,
	}
	if err := d.initDriver(providers...); err != nil {
		return nil, err
	}

	outs, err := d.drv.Outs()
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
		d.drv.Close()
		return nil, fmt.Errorf("no playback device found")
	}

	if err := out.Open(); err != nil {
		return nil, err
	}

	d.wr = mid.ConnectOut(out)
	d.out = out

	return d, nil
}

// Close closes the underlying MIDI port and driver.
func (d *Player) Close() {
	defer d.Off()
	defer d.out.Close()
	defer d.drv.Close()
}

// Off silences all MIDI channels.
func (d *Player) Off() {
	if test {
		return
	}
	fmt.Println("muting all channels")
	for ch := 0; ch < 16; ch++ {
		if err := d.wr.Silence(int8(ch), false); err != nil {
			fmt.Printf("failed to silence channel %d: %s\n", ch, err.Error())
		}
	}
}

// HandleMessage plays MIDI messages.
func (d *Player) HandleMessage(pos *mid.Position, msg midi.Message) {
	if !d.out.IsOpen() {
		return
	}
	if test {
		time.Sleep(time.Nanosecond)
	} else {
		delay := time.Second / time.Duration(d.TicksPerSecond) * time.Duration(pos.DeltaTicks)
		time.Sleep(delay)
	}
	var err error
	switch msg.(type) {
	case channel.ProgramChange:
		err = d.ProgramChange(msg.(channel.ProgramChange))
	default:
		if msg == meta.EndOfTrack {
			d.Off()
		}
		if test {
			return
		}
		err = d.wr.Write(msg)
	}
	if err != nil {
		d.Off()
		panic(err)
	}
}

// ProgramChange implements workaround for ProgramChange messages
// for broken drivers.
func (d *Player) ProgramChange(msg channel.ProgramChange) error {
	switch {
	case test:
		return nil
	case d.useFixedMessages:
		fmt.Println("TODO: fix portmididrv to 2-byte ProgramChange message")
		msg := FixedMessage{msg}
		return d.wr.Write(msg)
	default:
		return gm.WriteGMProgram(d.wr, msg.Channel(), msg.Program())
	}
}

// Reset send common reset messages to the MIDI device.
func (d *Player) Reset() {

	switch {
	case d.useFixedMessages:
		fmt.Println("TODO: fix portmididrv to allow Reset")
	default:
		fmt.Println("resetting MIDI player")
		if err := gm.WriteReset(d.wr, 0, 0); err != nil {
			fmt.Printf("failed to GM-reset player: %s", err)
		}
	}
}

// Play plays a song.
func (d *Player) Play(stream *Stream) {
	d.Reset()
	b := bytes.NewReader(stream.Bytes())
	rd := mid.NewReader(mid.NoLogger())
	rd.Msg.Each = d.HandleMessage
	rd.ReadAllSMF(b)
}

// Loop plays a song repeatedly until the player closes.
func (d *Player) Loop(stream *Stream) {
	for d.out.IsOpen() {
		d.Play(stream)
	}
}
