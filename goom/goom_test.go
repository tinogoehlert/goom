package goom_test

import (
	"os"
	"testing"

	"github.com/tinogoehlert/goom/goom"
	"github.com/tinogoehlert/goom/test"
)

const testWad = "DOOM1"

// loads test WAD if present, otherwise skips the test.
func loadTestWAD(t *testing.T) *goom.GameData {
	gd, err := goom.GetWAD(testWad, "")
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
	track := gd.Music.Track("INTRO")
	test.Assert(track != nil, "track not found: INTRO", t)
	test.Check(track.Validate(), t)
}
