package audio_test

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/tinogoehlert/goom/audio"
	mus "github.com/tinogoehlert/goom/audio/mus"
	"github.com/tinogoehlert/goom/files"
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

func sampleTrack(t *testing.T) []byte {
	/*
		Artifical Example Track

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
		Intruments:     16, 32
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
		[]byte(mus.LumpID), // MUS ID                 4       0
		b16(len(events)),   // score size             2       4
		b16(20),            // score offset           2       6   <--.
		b16(1),             // primary channels       2       8      |
		b16(0),             // secndary channels      2      10      |
		b16(numInst),       // number of instruments  2      12      |
		b16(0),             // dummy                  2      14      |
		b16(inst1),         // instrument 1           2      16      |
		b16(inst2),         // instrument 2           2      18      |
		events,             // event bytes            8      20   ---'
		// EOF                                        0      28
	)

	hex := hex.EncodeToString(data)

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

	bf := &files.BinFile{
		"test",
		data,
	}

	if hex != expected {
		t.Errorf("invalid sample.\nobserved: %s\nexpected: %s", hex, expected)
	}
	if bf.Hex() != expected {
		t.Errorf("invalid hex data.\nobserved: %s\nexpected: %s", bf.Hex(), expected)
	}
	return data
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

func doomFile(name string, t *testing.T) *files.BinFile {
	f := files.NewBinFile("..", "files", name)
	test.Check(f.Load(), t)
	return f
}

func introMus(t *testing.T) *files.BinFile  { return doomFile("D_INTRO.mus", t) }
func introaMus(t *testing.T) *files.BinFile { return doomFile("D_INTROA.mus", t) }
func e1m1Mus(t *testing.T) *files.BinFile   { return doomFile("D_E1M1.mus", t) }

func TestParseEvents(t *testing.T) {
	type Case struct {
		Data   []byte
		Length int
	}

	cases := []Case{
		Case{sampleTrack(t), 7},
		Case{doomSample(t), 6},
		Case{introaMus(t).Data, 214},
		Case{introMus(t).Data, 498},
		Case{e1m1Mus(t).Data, 5826},
	}

	for _, c := range cases {
		data := c.Data
		dmus, err := audio.NewMusData(data)
		test.Check(err, t)
		if len(dmus.Events) != c.Length {
			t.Errorf("invalid number of MUS events: %d, expected: %d events", len(dmus.Events), c.Length)
		}
	}
}

func TestMusLoading(t *testing.T) {
	data := sampleTrack(t)
	md, err := audio.NewMusData(data)
	test.Check(err, t)

	type Case struct {
		Index int
		Type  mus.EventType
		Note  uint8
		Delay uint16
		Hex   string
	}

	cases := []Case{
		Case{0, mus.System, 10, 0, "0a"},
		Case{1, mus.Controller, 3, 0, "0364"},
		Case{2, mus.PlayNote, 16, 15, "10"},
		Case{3, mus.PlayNote, 32, 15, "20"},
		Case{4, mus.RelaseNote, 16, 261, "10"},
		Case{5, mus.RelaseNote, 32, 5, "20"},
		Case{6, mus.ScoreEnd, 0, 0, ""},
	}

	for _, c := range cases {
		s := md.Events[c.Index]
		// fmt.Printf("comparing score %+v with test case: %+v\n", s, c)
		if s.Type != c.Type {
			t.Errorf("invalid mus type %d, expected %d.", s.Type, c.Type)
		}
		if s.Type == mus.RelaseNote && s.GetNote() != c.Note {
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
	for i, d := range [][]byte{
		sampleTrack(t),
		doomSample(t),
		introMus(t).Data,
		introaMus(t).Data,
		e1m1Mus(t).Data,
	} {
		musd, err := audio.NewMusData(d)
		test.Check(err, t)
		mid, err := audio.NewMidiData(d)
		test.Check(err, t)
		name := fmt.Sprintf("D_TEST%d", i)
		track := audio.MusicTrack{wad.Lump{Name: name, Data: d}, mid, musd}
		track.Play()
		defer track.Stop()

		test.Check(track.Validate(), t)

		ev := musd.Events[0]

		// test.Check(track.SaveMus(), t)
		// test.Check(track.SaveMidi(), t)

		// test the info methods
		test.Assert(musd.Info()[0:8] == "mus.Data", "invalid mus info", t)
		test.Assert(ev.Info()[0:9] == "mus.Event", "invalid event info", t)
	}
}
