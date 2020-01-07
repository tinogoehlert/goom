package drivers

import (
	"fmt"

	"gitlab.com/gomidi/midi/mid"
	"gitlab.com/ubunatic/portmididrv"

	gomididrv "github.com/tinogoehlert/goom/drivers/gomidi"
)

func init() {
	MusicDrivers[PortMidiMusic] = gomididrv.NewMidiPlayer(PortMidiDriver())
}

// PortMidiDriver safely returns the rtmidi driver or nil.
func PortMidiDriver() mid.Driver {
	drv, err := portmididrv.New()
	if err != nil {
		fmt.Printf("failed to load portmididrv: %s", err.Error())
		return nil
	}
	return drv
}
