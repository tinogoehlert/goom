package files_test

import (
	"fmt"
	"testing"

	"github.com/tinogoehlert/goom/audio/files"
	"github.com/tinogoehlert/goom/test"
)

type Case struct {
	name string
	len  int
}

func TestHexConversion(t *testing.T) {
	name := "SLADE_E1M1.mid"
	n := 23322
	f1 := files.NewBinFile(name)
	err := f1.Load()
	test.Check(err, t)
	h1 := f1.Hex()
	test.Assert(len(h1) != n, "hex data has wrong length", t)

	fmt.Println(f1.Path)
	// fmt.Println(h1)

	f2 := files.NewBinFile(name)
	err = f2.FromHex(h1)
	test.Check(err, t)
	test.Assert(f1.Hex() == f2.Hex(), "hex data mismatch", t)
}
