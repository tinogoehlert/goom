// +build !windows

package midi

import (
	"gitlab.com/gomidi/midi/mid"
	"gitlab.com/ubunatic/portmididrv"
)

func init() {
	drvInit[PortMidi] = new
}

func new() (mid.Driver, error) {
	return portmididrv.New()
}
