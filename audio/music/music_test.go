package music_test

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"path"
	"strings"
	"testing"

	gmidi "github.com/tinogoehlert/goom/audio/midi"
	gmus "github.com/tinogoehlert/goom/audio/mus"
	mus "github.com/tinogoehlert/goom/audio/mus"
	"github.com/tinogoehlert/goom/audio/music"
	"github.com/tinogoehlert/goom/goom"
	"github.com/tinogoehlert/goom/test"
	"github.com/tinogoehlert/goom/wad"
)

func b16(v int) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(v))
	return b
}

func concat(values ...[]byte) []byte {
	var data []byte
	for _, b := range values {
		data = append(data, b...)
	}
	return data
}

func sampleMus(t *testing.T) []byte {
	/*
		Artificial Example Track

		Composed using ModdingWiki MUS spec: http://www.shikadi.net/moddingwiki/MUS_Format

		30  # System Controller     Delay = 0, Channel = 0
		0A  # AllNotesOff

		40  # Controller            Delay = 0, Channel = 0
		30  # Volume
		64  # VolumeLevel           Volume = 100

		90  # PlayNote              Delay = 0, Channel = 0         High-bit is set,   delay byte follows
		10  # PlayNote Payload      Note  = 16
		0F  # First byte of delay   Delay = Delay * 128 + 15 = 15  High-bit is unset, delay is complete
			# 15 Ticks delay
		90  # PlayNote              Delay = 0, Channel = 0
		20  # PlayNote Payload      Note  = 32
		0F  # First byte of delay   Delay = Delay * 128 + 15 = 15  High-bit is unset, delay is complete
			# 15 Ticks delay
		80  # ReleaseNote           Delay = 0, Channel = 0         High-bit is set,   delay byte follows
		10  # ReleaseNote payload   Note  = 16
		82  # First byte of delay   Delay = 2 (0x82 & 0x7F)        High-bit is set,   delay byte follows
		05  # Second byte of delay  Delay = Delay * 128 + 5 = 261  High-bit is unset, delay is complete
			# 261 Ticks delay
		80  # RelaseNote            Delay = 0, Channel = 0         High-bit is set,   delay byte follows
		20  # ReleaseNote payload   Note  = 32
		05  # First byte of delay   Delay = Delay * 128 + 5 = 5    High-bit is unset, delay is complete
			# 5 Ticks delay
		60  # ScoreEnd              Delay = 0                      High-bit must be unset

		Channels:       0
		NumChannels:    1
		Instruments:    16, 32
		NumInstruments: 2
		NumScores:      7
		NumBytes:       19

		This track should be fully MIDI-convertible.
		The result should be playable using any MID-file player.
	*/

	eventsHex := strings.Join([]string{
		"300a",     // all notes off
		"400364",   // volume 100
		"90100f",   // play 16 + delay
		"90200f",   // play 32 + delay
		"80108205", // release 16 + delay
		"802005",   // release 32 + delay
		"60",       // end
	}, "")
	events, _ := hex.DecodeString(eventsHex)
	inst1 := 16
	inst2 := 32
	numInst := 2

	// create header and append events
	data := concat(
		// Value               Description            Bytes  ByteIndex
		[]byte(gmus.LumpID), // MUS ID                 4       0
		b16(len(events)),    // score size             2       4
		b16(20),             // score offset           2       6   <--.
		b16(1),              // primary channels       2       8      |
		b16(0),              // secondary channels     2      10      |
		b16(numInst),        // number of instruments  2      12      |
		b16(0),              // dummy                  2      14      |
		b16(inst1),          // instrument 1           2      16      |
		b16(inst2),          // instrument 2           2      18      |
		events,              // event bytes            8      20   ---'
		// EOF                                         0      28
	)

	hx := hex.EncodeToString(data)

	// fmt.Println("sampleData:", hex)

	expected := ("4d55531a" + // MUS ID
		"1300" + //    19 len
		"1400" + //    20 offset
		"0100" + //     1 prim ch
		"0000" + //     0 sec ch
		"0200" + //     2 num inst
		"0000" + //     0 dummy
		"1000" + //    16 first inst
		"2000" + //    32 second inst
		eventsHex)

	if hx != expected {
		t.Errorf("invalid sample.\nobserved: %s\nexpected: %s", hx, expected)
	}
	return data
}

func introMid(t *testing.T) {
	/*
		bytes total:  2002
		bytes header:   22
		bytes track:  1980

		4D 54 68 64 00 00 00 06    Mthd + header size
		00 00 00 01 00 46          number of tracks and resolution
		4D 54 72 6B 00 00 07 BC    Mtrk + track size

		00 C0 00
		00 B0 07 64
		00 B0 0A 40
		00 90 08 7F
		00 C1 2F
		00 B1 07 64
		00 B1 0A 40
		00 91 20 7F
		00 C2 33
		00 B2 07 64
		00 B2 0A 0E
		00 92 20 7F
		00 B2 07 64
		00 B2 0A 0E
		00 C3 33
		00 B3 07 64
		00 B3 0A 72
		00 93 26 7F
		00 B3 07 64
		00 B3 0A 72
		00 C4 5E
		00 B4 07 64
		00 B4 0A 0E
		00 94 2F 7C
		00 B4 07 64
		00 B4 0A 0E
		00 C5 5E
		00 B5 07 64
		00 B5 0A 72
		00 95 35 5C
		00 B5 07 64
		00 B5 0A 72
		00 C6 75
		00 B6 07 7F
		00 B6 0A 40
		00 96 26 7F
		00 C7 66
		00 B7 07 64
		00 B7 0A 40
		00 97 26 57
		00 B7 07 64
		00 B7 0A 40
		00 C9 00
		00 B9 07 64
		00 B9 0A 40
		00 99 23 7F
		00 99 28 7F
		00 C8 12
		00 B8 07 64
		00 B8 0A 7C
		00 98 11 7F
		00 CA 0B
		00 BA 07 64
		...
	*/
}

func e1m1Mid() {
	/*
		4D 54 68 64 00 00 00 06 00 00 00 01 00
		46 4D 54 72 6B 00 00 5B 04

		00 C0 1E         ChangeInst 0  30
		00 B0 07 64      ChangeCtrl 0   7  100
		00 B0 0A 18      ChangeCtrl 0  10   24
		00 B0 07 64      ChangeCtrl 0   7  100
		00 B0 0A 18      ChangeCtrl 0  10   24

		00 C1 1D         ChangeInst 1  29
		00 B1 07 64      ChangeCtrl 1   7  100
		00 B1 0A 68      ChangeCtrl 1  10  104
		00 B1 07 64      ChangeCtrl 1   7  100
		00 B1 0A 68      ChangeCtrl 1  10  104

		00 C2 22         ChangeInst 2  34
		00 B2 07 64      ChangeCtrl 2   7  100
		00 B2 0A 40      ChangeCtrl 2  10   64

		00 92 28 78      NoteOn     2  40  120

		00 C9 00         ChangeInst 9   0
		00 B9 07 64      ChangeCtrl 9   7  100
		00 99 24 7B      NoteOn     9  36  123
		00 99 28 73      NoteOn     9  40  115
		00 99 29 78      NoteOn     9  41  120

		00 91 28 6C
		06 81 28 00
		07 82 28 00
		04 89 24 00
		02 91 28 72
		01 89 28 00
		01 89 29 00
		11 91 34 72
		03 81 28 00
		10 81 34 00
		01 91 28 6E
		06 81 28 00
	*/
}

// Hand-coded, playable MIDI file, tested using aplaymidi and Timidity.
// This MIDI file encodes the same notes as sampleMus.
// Converting sampleMus should produce this MIDI file.
func sampleMid(t *testing.T) []byte {
	/*
		                                    	MIDI  MIDI
			MUS                                 DELAY EVENT
			"300a",     // all notes off        00    B0 78 00  // controller 120, notes off
			"400364",   // volume 100           00    B0 07 64  // controller 7,
			"90100f",   // play 16 + delay      00    90 10 64  // play note 16 with volume 100
			"90200f",   // play 32 + delay      0f       20 64  // play note 32 with volume 100 after 15 ticks
			"80108205", // release 16 + delay   0f       10 00  // play note 16 with volume 0   after 15 ticks
			"802005",   // release 32 + delay   82 05    20 00  // play note 31 with volume 0   after 261 ticks
			"60",       // end                  00    FF 2F 00  // end track
	*/
	h := gmidi.MidHeader()
	data, _ := hex.DecodeString(strings.Join([]string{
		"00b07800", // controller 120, value = 0
		"00b00764", // controller 7, volume = 100
		"00901064", // play 16, volume = 100
		"0f2064",   // play 32, volume = 100 after 15 ticks
		"0f1000",   // release 16 after 15 ticks
		"82052000", // release 32 after 261 ticks
		"00ff2f00", // end track
	}, ""))
	tl := gmidi.TrackLength(data)
	return append(append(h, tl...), data...)
}

func doomSample(t *testing.T) []byte {
	// First bytes of the event data of a track from DOOM1.
	//
	// 40 00 1e 40 03 64 40 04 22 10 a8 77
	//
	// byte  bits         description
	// 40    0100 0000    controller (4) -> 2 bytes will follow
	// 00    0000 0000      number = 0,  change instrument
	// 1e    0001 1110      value  = 30, instrument number
	//
	// 40    0100 0000    controller (4) -> 2 bytes will follow
	// 03    0000 0011      number = 3,   set volume
	// 64    0110 0100      value  = 100, volume level
	//
	// 40    0100 0000    controller (4) -> 2 bytes will follow
	// 04    0000 0100      number = 4,  set balance
	// 22    0010 0010      value  = 34, betweem left(0) and center(64)

	// 10    0001 0000    play note (ch:0, delay:0)
	// a8    1010 1000      vol? = 1, note = 40,  play note 40 and change volume
	// 77    0111 0111      volume = 119,         set volume to 119  (max is 127)

	// To create a valid sample track, we must turn off note 40 again:
	// 00    0000 0000    release note (ch:0, delay:0)
	// 28    0010 1000      note = 40,   release note 40

	v := []string{
		"4d55531a", // MUS ID
		"0e00",     // 14 len
		"1200",     // 18 offset
		"0100",     //  1 prim ch
		"0000",     //  0 sec ch
		"0200",     //  1 num inst
		"0000",     //  0 dummy
		"0100",     //  1 inst 1
		"40001e",   // ChangeInst 30
		"400364",   // Volume 100
		"400422",   // Balance 22
		"10a877",   // Play Note 40, Volume 119
		"0028",     // Release Note 40
		"60",       // ScoreEnd
	}
	data, _ := hex.DecodeString(strings.Join(v, ""))
	return data
}

func loadMus(name string, t *testing.T) []byte {
	gd, err := goom.GetWAD(path.Join("..", "..", "DOOM1"), "")
	test.Check(err, t)
	return gd.Music.Track(name).Data
}

type Mus struct {
	Name   string
	Data   []byte
	Length int
}

func allMus(t *testing.T) []Mus {
	return []Mus{
		{"SAMPLE", sampleMus(t), 7},
		{"DOOM", doomSample(t), 6},
		{"INTROA", loadMus("INTROA", t), 214},
		{"INTRO", loadMus("INTRO", t), 498},
		{"E1M1", loadMus("E1M1", t), 5826},
	}
}

func TestParseEvents(t *testing.T) {
	for _, m := range allMus(t) {
		dmus, err := mus.NewMusStream(m.Data)
		test.Check(err, t)
		if len(dmus.Events) != m.Length {
			t.Errorf("track %s: invalid number of MUS events: %d, expected: %d events",
				m.Name, len(dmus.Events), m.Length)
		}
	}
}

func TestMusLoading(t *testing.T) {
	data := sampleMus(t)
	md, err := mus.NewMusStream(data)
	test.Check(err, t)

	type Case struct {
		Index int
		Type  gmus.EventType
		Note  uint8
		Delay uint16
		Hex   string
	}

	cases := []Case{
		{0, gmus.System, 10, 0, "0a"},
		{1, gmus.Controller, 3, 0, "0364"},
		{2, gmus.PlayNote, 16, 15, "10"},
		{3, gmus.PlayNote, 32, 15, "20"},
		{4, gmus.RelaseNote, 16, 261, "10"},
		{5, gmus.RelaseNote, 32, 5, "20"},
		{6, gmus.ScoreEnd, 0, 0, ""},
	}

	for _, c := range cases {
		s := md.Events[c.Index]
		// fmt.Printf("comparing score %+v with test case: %+v\n", s, c)
		if s.Type != c.Type {
			t.Errorf("invalid mus type %d, expected %d.", s.Type, c.Type)
		}
		if s.Type == gmus.RelaseNote && s.GetNote() != c.Note {
			t.Errorf("invalid note %d, expected %d.", s.GetNote(), c.Note)
		}
		if s.Delay != c.Delay {
			t.Errorf("wrong delay %d, expected %d", s.Delay, c.Delay)
		}
		if c.Hex == "" && s.Data != nil {
			t.Errorf("Data %x is not nil", c.Hex)
		}
		if c.Hex != "" && s.Data == nil {
			t.Errorf("missing Data %x", c.Hex)
		}
		if s.Data != nil {
			hx := hex.EncodeToString(s.Data)
			if hx != c.Hex {
				t.Errorf("invalid Data %x, expected %s", hx, c.Hex)
			}
		}
	}
}

func TestTrackLoading(t *testing.T) {
	for _, d := range allMus(t) {
		track, err := music.NewTrack(wad.Lump{Name: d.Name, Data: d.Data})
		test.Check(err, t)
		test.Check(track.Validate(), t)
		mu := track.MusStream
		ev := mu.Events[0]

		// test.Check(track.SaveMus(), t)
		// test.Check(track.SaveMidi(), t)

		// test the info methods
		test.Assert(mu.Info()[0:8] == "mus.Data", "invalid mus info", t)
		test.Assert(ev.Info()[0:9] == "mus.Event", "invalid event info", t)
	}
}

func TestPlaybackProviders(t *testing.T) {
	allProviders := []gmidi.Provider{gmidi.RTMidi, gmidi.PortMidi, gmidi.Any}
	songs := allMus(t)[:3]
	gmidi.TestMode()

	for _, provider := range allProviders {
		fmt.Println("starting player with provider: ", provider)
		p, err := gmidi.NewPlayer(provider)
		test.Check(err, t)
		test.Assert(p != nil, "no midi device found cannot test playback", t)
		defer p.Close()

		for _, song := range songs {
			track, err := music.NewTrack(wad.Lump{Data: song.Data})
			test.Check(err, t)
			fmt.Println("playing song", track.Name, "using provider", provider)
			p.Play(track.MidiStream)
		}
		p.Close()
	}
}
