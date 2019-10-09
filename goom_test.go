package goom

import (
	"os"
	"testing"

	"github.com/tinogoehlert/goom/test"
)

const testWad = "DOOM1"

// loads test WAD if present, otherwise skips the test.
func loadTestWAD(t *testing.T) *GameData {
	gd, err := GetWAD(testWad, "")
	if os.IsNotExist(err) {
		t.Logf("skipping WAD test for missing %s WAD.", testWad)
		t.Skip()
	}
	test.Check(err, t)
	return gd
}

// Test loading the DOOM1.wad and gwa files
func TestWAD(t *testing.T) {
	if gd := loadTestWAD(t); gd == nil {
		t.Fail()
	}
}

// Test loading and playing music.
func TestMusic(t *testing.T) {
	gd := loadTestWAD(t)
	music, ok := gd.Music["D_E1M1"]
	test.Assert(ok, "track not found: D_E1M1", t)
	music.Play()
	defer music.Stop()

	test.Check(music.SaveMus(), t)
	test.Check(music.SaveMidi(), t)
	test.Check(music.Validate(), t)
}
