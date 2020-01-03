package gomididrv_test

import (
	"fmt"
	"path"
	"testing"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/mid"

	"github.com/tinogoehlert/goom/audio/files"
	"github.com/tinogoehlert/goom/audio/music"
	gomididrv "github.com/tinogoehlert/goom/drivers/gomidi"
	"github.com/tinogoehlert/goom/goom"
	"github.com/tinogoehlert/goom/test"
)

func loadMus(name string, t *testing.T) *music.Track {
	gd, err := goom.GetWAD(path.Join("..", "..", "DOOM1"), "")
	test.Check(err, t)
	return gd.Music.Track(name)
}

func TestMidiProcessing(t *testing.T) {
	f, err := files.LoadFile("..", "..", "audio", "files", "SLADE_E1M1.mid")
	test.Check(err, t)
	e1m1 := loadMus("E1M1", t)

	var events []string
	fn := func(pos *mid.Position, msg midi.Message) {
		events = append(events, msg.String())
	}
	gomididrv.Process(f.Data, fn)
	expected := events

	events = make([]string, 0)
	gomididrv.Process(e1m1.MidiStream.Bytes(), fn)
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
