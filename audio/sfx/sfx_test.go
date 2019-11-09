package sfx_test

import (
	"testing"

	"github.com/tinogoehlert/goom/audio/sfx"
	"github.com/tinogoehlert/goom/test"
)

func TestPlaySound(t *testing.T) {
	sfx.TestMode()
	test.Check(sfx.PlaySounds("PISTOL", "OOF", "SHOTGN"), t)
}
