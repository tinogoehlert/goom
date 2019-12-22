package gomididrv

import (
	"github.com/tinogoehlert/goom/audio/music"
	"gitlab.com/gomidi/midi/mid"
)

// MidiPlayer is the rtmidi music driver.
type MidiPlayer struct {
	tracks         *music.TrackStore
	driver         mid.Driver
	writer         *mid.Writer
	out            mid.Out
	test           bool
	TicksPerSecond int
}

// NewMidiPlayer initalizes a MidiPlayer with a specific midi driver.
func NewMidiPlayer(drv mid.Driver) *MidiPlayer {
	return &MidiPlayer{
		driver: drv,
	}
}

// InitMusic initializes the MIDI device.
func (p *MidiPlayer) InitMusic(tracks *music.TrackStore, tempDir string) error {
	p.TicksPerSecond = DefaultTicksPerSecond
	if err := p.InitMidiOutput(); err != nil {
		return err
	}
	p.tracks = tracks
	return nil
}

// TestMode slicences all music and sets delays to 0.
func (p *MidiPlayer) TestMode() {
	p.test = true
}

// PlayMusic loops the music in the background.
func (p *MidiPlayer) PlayMusic(m *music.Track) error {
	go p.Play(m.MidiStream)
	return nil
}

// Close does nothing
func (p *MidiPlayer) Close() {
	defer p.Off()
	defer p.out.Close()
	defer p.driver.Close()
}
