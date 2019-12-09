package sfx_test

import (
	"fmt"
	"testing"

	"github.com/tinogoehlert/goom/run"
	"github.com/tinogoehlert/goom/test"
)

func TestPlaySound(t *testing.T) {
	r := run.TestRunner("..", "..")
	fmt.Sprintf("TestRunner: %v", r)
	drv := r.World().AudioDriver()
	drv.TestMode()
	// 22 kHz
	test.Check(drv.Play("DSITMBK"), t)
	// 11 kHz
	test.Check(drv.Play("DSPISTOL"), t)
	drv.Close()
}
