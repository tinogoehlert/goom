package drivers

import (
	"fmt"

	"gitlab.com/gomidi/midi/mid"
	"gitlab.com/gomidi/rtmididrv"

	gomididrv "github.com/tinogoehlert/goom/drivers/gomidi"
)

func init() {
	MusicDrivers[RtMidiMusic] = gomididrv.NewMidiPlayer(RtMidiDriver())
}

// RtMidiDriver safely returns the rtmidi driver or nil.
func RtMidiDriver() mid.Driver {
	drv, err := rtmididrv.New()
	if err != nil {
		fmt.Printf("failed to load rtmididrv: %s", err.Error())
		return nil
	}
	return drv
}
