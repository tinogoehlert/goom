package goom

import (
	"os"
	"testing"
)

const testWad = "DOOM1"

// loads test WAD if present, otherwise skips the test.
func loadTestWAD(t *testing.T) *GameData {
	gd, err := GetWAD(testWad, "")
	if os.IsNotExist(err) {
		t.Logf("skipping WAD test for missing %s WAD.", testWad)
		t.Skip()
	}
	if err != nil {
		t.Error(err)
	}
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
	if !ok {
		t.Fail()
	}
	music.Play()
	defer music.Stop()
}
