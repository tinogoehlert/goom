package convert_test

import (
	"fmt"
	"path"
	"testing"

	"github.com/tinogoehlert/goom/audio/convert"
	"github.com/tinogoehlert/goom/audio/files"
	"github.com/tinogoehlert/goom/audio/music"
	"github.com/tinogoehlert/goom/goom"
	"github.com/tinogoehlert/goom/test"
)

func loadMus(name string, t *testing.T) *music.Track {
	gd, err := goom.GetWAD(path.Join("..", "..", "DOOM1"), "")
	test.Check(err, t)
	return gd.Music.Track(name)
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
	fmt.Printf("skipping %s <--> %s binary compare\n", e1mid.Path, f2.Path)

	// test.Assert(e1mid.Compare(f2) == 0, "invalid MIDI output", t)
}
