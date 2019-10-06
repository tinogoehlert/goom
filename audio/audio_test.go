package audio

import (
	"encoding/binary"
	"fmt"
	"strings"
	"testing"

	mus "github.com/tinogoehlert/goom/audio/mus"
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

		80  # ReleaseNote                      Delay = 0 (inital value)       High-bit is set, delay byte follows
		10  # ReleaseNote payload (note 1)     Delay = 0
		82  # First byte of delay              Delay = 2 (0x82 & 0x7F)        High-bit is set, delay byte follows
		05  # Second byte of delay             Delay = Delay * 128 + 5 = 261  High-bit is unset, delay is complete
		80  # RelaseNote delayed by 261 ticks  Delay = 0 (value reset)
		20  # ReleaseNote payload (note 2)     Delay = 0
		05  # First byte of second delay       Delay = Delay * 128 + 5 = 5    High-bit is unset, delay is complete
		60  # ScoreEnd delayed by 5 ticks

	*/

	scores := []byte("\x80\x10\x82\x05\x80\x20\x05\x60")
	inst1 := 1
	inst2 := 2
	numInst := 2

	data := concat(
		// Value             Description            Bytes  Offset
		[]byte(mus.LumpID), // MUS ID                 4       0
		b16(len(scores)),   // score size             2       4
		b16(20),            // score offset           2       6   <--.
		b16(1),             // primary channels       2       8      |
		b16(0),             // secndary channels      2      10      |
		b16(numInst),       // number of instruments  2      12      |
		b16(0),             // dummy                  2      14      |
		b16(inst1),         // instrument 1           2      16      |
		b16(inst2),         // instrument 2           2      18      |
		scores,             // scores bytes           8      20   ---'
		// EOF                                      0      28
	)

	hex := fmt.Sprintf("%x", data)

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
		"80108205" + // ReleaseNote 1, 261 ticks Delay
		"802005" + // ReleaseNote 2, 5 ticks Delay
		"60") // ScoresEnd

	if hex != expected {
		t.Errorf("invalid sample.\nobserved: %s\nexpected: %s", hex, expected)
	}
	return data
}

func doomSample(t *testing.T) []byte {
	// First bytes of the scores of a track from DOOM1.
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

	header := []byte("\x4d\x55\x53\x1a" + // MUS ID
		"\x0c\x00" + //     12 len
		"\x12\x00" + //     18 offset
		"\x01\x00" + //      1 prim ch
		"\x00\x00" + //      0 sec ch
		"\x02\x00" + //      1 num inst
		"\x00\x00" + //      0 dummy
		"\x01\x00", //       1 inst 1
	)

	// encode partial song and append a ScoresEnd event
	song := "\x40\x00\x1e\x40\x03\x64\x40\x04\x22\x10\xa8\x77" + "\x60"
	return []byte(append(header, song...))
}

func doomMidi() string {
	mid := `
4d 54 68 64 00 00 00 06 00 01 00 05 00 59 4d 54
72 6b 00 00 00 2b 00 ff 02 16 51 4d 55 53 32 4d
49 44 20 28 43 29 20 53 2e 42 41 43 51 55 45 54
00 ff 59 02 00 00 00 ff 51 03 09 a3 1a 00 ff 2f
00 4d 54 72 6b 00 00 15 f2 00 ff 03 1b 51 75 69
63 6b 20 4d 55 53 2d 3e 4d 49 44 20 21 20 66 6f
72 20 57 69 6e 54 65 78 00 c0 1e 00 b0 07 64 00
0a 18 00 07 64 00 0a 18 89 45 90 28 6c 06 28 00
0d 28 72 14 34 72 02 28 00 10 34 00 01 28 6e 06
28 00 0d 28 6e 13 32 6f 02 28 00 0f 32 00 02 28
68 06 28 00 0d 28 6f 13 30 74 00 28 00 12 30 00
01 28 6c 06 28 00 0d 28 6e 13 2e 6e 02 28 00 11
2e 00 00 28 6e 08 28 00 0c 28 6f 13 2f 6e 12 2f
00 01 30 75 0d 28 00 04 30 00 02 28 6f 04 28 00
0f 28 6f 13 34 72 04 28 00 0d 34 00 02 28 6a 06
28 00 0d 28 6c 13 32 6a 03 28 00 10 32 00 00 28
72 06 28 00 0d 28 6a 14 30 75 00 28 00 0f 30 00
04 28 6f 07 28 00 0c 28 72 13 2e 72 08 e0 00 3e
02 00 41 01 00 44 01 00 47 01 00 4b 01 00 4d 01
00 4c 01 00 46 01 00 40 0a 00 43 02 00 44 01 00
`
	return strings.ReplaceAll(strings.ReplaceAll(mid, " ", ""), "\n", "")
}

func TestParseScores(t *testing.T) {
	type Case struct {
		Data   []byte
		Length int
	}

	cases := []Case{
		Case{doomSample(t)[18:], 5},
		Case{[]byte("\x80\x10\x82\x05\x80\x20\x05\x60"), 3},
	}

	for _, c := range cases {
		data := c.Data
		scores, err := mus.ParseScores(data)
		test.Check(err, t)
		if len(scores) != c.Length {
			t.Errorf("invalid scores: %+v, expected: %d scores", scores, 3)
		}
	}
}

func TestMusLoading(t *testing.T) {
	data := sampleTrack(t)
	md, err := NewMusData(data)
	test.Check(err, t)

	type Case struct {
		Index int
		Type  mus.Event
		Delay int
		Data  string
	}

	cases := []Case{
		Case{0, mus.RelaseNote, 261, "\x10"},
		Case{1, mus.RelaseNote, 5, "\x20"},
		Case{2, mus.ScoreEnd, 0, ""},
	}

	for _, c := range cases {
		s := md.Scores[c.Index]
		// fmt.Printf("comparing score %+v with test case: %+v", s, c)
		if s.Type != c.Type {
			t.Errorf("invalid mus type %d, expected %d", s.Type, c.Type)
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
	for i, d := range [][]byte{sampleTrack(t), doomSample(t)} {
		mid, err := NewMidiData(d)
		test.Check(err, t)
		name := fmt.Sprintf("D_TEST%d", i)
		track := MusicTrack{wad.Lump{Name: name, Data: d}, mid}
		track.Play()
		defer track.Stop()

		test.Check(track.SaveMus(), t)
		test.Check(track.SaveMidi(), t)
		// fmt.Println(mus.Info())
	}
	fmt.Printf("EXPECTED\nMID D_E1M1: %s\n", doomMidi()[:200])
}

func TestMus2Mid(t *testing.T) {
	for i, d := range [][]byte{sampleTrack(t), doomSample(t)} {
		md, err := NewMusData(d)
		test.Check(err, t)
		mid := Mus2Mid(md)
		fmt.Printf("MID M2M_%d: %x\n", i, mid.Data)
	}
}
