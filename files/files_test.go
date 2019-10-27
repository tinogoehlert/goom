package files_test

import (
	"fmt"
	"testing"

	"github.com/tinogoehlert/goom/files"
	"github.com/tinogoehlert/goom/test"
)

var (
	doomDir = ".."
)

type Case struct {
	name string
	len  int
}

func TestHexConversion(t *testing.T) {
	for _, c := range []Case{
		Case{"D_INTROA.mid", 871},
		Case{"D_INTROA.mus", 631},
	} {
		name := c.name
		n := c.len
		f1 := files.NewBinFile(doomDir, "files", name)
		err := f1.Load()
		test.Check(err, t)
		h1 := f1.Hex()
		test.Assert(len(h1) != n, "hex data has wrong length", t)

		fmt.Println(f1.Path)
		fmt.Println(h1)

		f2 := files.NewBinFile(doomDir, "files", name)
		err = f2.FromHex(h1)
		test.Check(err, t)
		test.Assert(f1.Hex() == f2.Hex(), "hex data mismatch", t)
	}
}
