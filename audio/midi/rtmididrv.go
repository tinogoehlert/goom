package midi

import (
	"gitlab.com/gomidi/rtmididrv"
)

func init() {
	drvInit[RTMidi] = rtmididrv.New
}
