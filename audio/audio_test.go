package audio_test

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
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
		Example Track

		Based on example from: http://www.shikadi.net/moddingwiki/MUS_Format

		80  # ReleaseNote                   Delay = 0 (inital value)       High-bit is set, delay byte follows
		10  # ReleaseNote payload (note 16) Delay = 0
		82  # First byte of delay           Delay = 2 (0x82 & 0x7F)        High-bit is set, delay byte follows
		05  # Second byte of delay          Delay = Delay * 128 + 5 = 261  High-bit is unset, delay is complete
		80  # RelaseNote after 261 ticks    Delay = 0 (value reset)
		20  # ReleaseNote payload (note 32) Delay = 0
		05  # First byte of second delay    Delay = Delay * 128 + 5 = 5    High-bit is unset, delay is complete
		60  # ScoreEnd after 5 ticks        Delay = 0                      High-bin must be unset
	*/

	events := []byte("\x80\x10\x82\x05\x80\x20\x05\x60")
	inst1 := 1
	inst2 := 2
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
		"0800" + //     8 len
		"1400" + //    20 offset
		"0100" + //     1 prim ch
		"0000" + //     0 sec ch
		"0200" + //     2 num inst
		"0000" + //     0 dummy
		"0100" + //     1 inst 1
		"0200" + //     2 inst 2
		"80108205" + // ReleaseNote 16, 261 ticks Delay
		"802005" + // ReleaseNote 32, 5 ticks Delay
		"60") // ScoreEnd

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

	header := []byte("\x4d\x55\x53\x1a" + // MUS ID
		"\x0e\x00" + //     14 len
		"\x12\x00" + //     18 offset
		"\x01\x00" + //      1 prim ch
		"\x00\x00" + //      0 sec ch
		"\x02\x00" + //      1 num inst
		"\x00\x00" + //      0 dummy
		"\x01\x00", //       1 inst 1
	)

	// encode partial song and append a ScoreEnd event
	song := "\x40\x00\x1e\x40\x03\x64\x40\x04\x22" +
		"\x10\xa8\x77" + "\x00\x28" + "\x60"
	return []byte(append(header, song...))
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
		Case{sampleTrack(t), 3},
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
		Data  string
	}

	cases := []Case{
		Case{0, mus.RelaseNote, 16, 261, "\x10"},
		Case{1, mus.RelaseNote, 32, 5, "\x20"},
		Case{2, mus.ScoreEnd, 0, 0, ""},
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
		if c.Data == "" && s.Data != nil {
			t.Errorf("Data %x is not nil", c.Data)
		}
		if c.Data != "" && s.Data == nil {
			t.Errorf("missing Data %x", c.Data)
		}
		if s.Data != nil {
			if string(s.Data) != c.Data {
				t.Errorf("invalid Data %x, expected %s", s.Data, c.Data)
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
