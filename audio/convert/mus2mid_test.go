package convert_test

import (
	"fmt"
	"path"
	"testing"

	"github.com/tinogoehlert/goom/audio/convert"
	"github.com/tinogoehlert/goom/audio/files"
	gmidi "github.com/tinogoehlert/goom/audio/midi"
	"github.com/tinogoehlert/goom/audio/music"
	"github.com/tinogoehlert/goom/goom"
	"github.com/tinogoehlert/goom/test"
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/mid"
)

func loadMus(name string, t *testing.T) *music.Track {
	gd, err := goom.GetWAD(path.Join("..", "..", "DOOM1"), "")
	test.Check(err, t)
	return gd.Music.Track(name)
}

func TestMus2Mid(t *testing.T) {
	f, err := files.LoadFile("..", "files", "SLADE_E1M1.mid")
	test.Check(err, t)
	e1m1 := loadMus("E1M1", t)

	var events []string
	fn := func(pos *mid.Position, msg midi.Message) {
		events = append(events, msg.String())
	}
	gmidi.Process(f.Data, fn)
	expected := events

	events = make([]string, 0)
	gmidi.Process(e1m1.MidiStream.Bytes(), fn)
	observed := events

	test.Assert(len(expected) > 0, "expected is empty", t)
	test.Assert(len(observed) > 0, "observed is empty", t)

	maxErrors := 10
	numErrors := 0

	for i := 0; i < len(expected) && i < len(observed) && numErrors < maxErrors; i++ {
		e := expected[i]
		o := observed[i]
		if e != o {
			numErrors++
			fmt.Println(i, "expected", o)
			fmt.Println(i, "observed", e)
		}
	}

	test.Assert(numErrors == 0, "midi output bytes mismatch", t)
}

func TestMus2MidDump(t *testing.T) {
	e1mid, err := files.LoadFile("..", "files", "SLADE_E1M1.mid")
	test.Check(err, t)
	e1mus := loadMus("E1M1", t)
	mi, err := convert.Mus2Mid(e1mus.MusStream)
	test.Check(err, t)
	f2 := files.NewBinFile("..", "files", "GOOM_E1M1.mid")
	f2.Data = mi.Data.Bytes()
	f2.Dump()
	fmt.Println("saving", f2.Path)
	f2.Save()
	fmt.Printf("skipping %s <--> %s binary compare\n",
		e1mid.Path, f2.Path)

	// test.Assert(e1mid.Compare(f2) == 0, "invalid MIDI output", t)
}
