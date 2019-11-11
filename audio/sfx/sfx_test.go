package sfx_test

import (
	"testing"

	"github.com/tinogoehlert/goom/audio/sfx"
	"github.com/tinogoehlert/goom/test"
)

func TestPlaySound(t *testing.T) {
	sfx.TestMode()
	// 22 kHz
	test.Check(sfx.PlaySounds("ITMBK"), t)
	// 11 kHz
	test.Check(sfx.PlaySounds("PISTOL", "OOF", "SHOTGN"), t)
}
