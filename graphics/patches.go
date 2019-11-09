package graphics

import (
	"github.com/tinogoehlert/goom/internal/utils"
)

const (
	patchSize = 10
)

// Patch fragment of a wall texture
type Patch struct {
	originX       int
	originY       int
	pictureID     int16
	stepDirection int16
	colorMap      int16
	*DoomPicture
	name string
}

// NewPatch create a new patch
func NewPatch(buff []byte) *Patch {
	if len(buff) != patchSize {
		panic("patch size missmatch")
	}
	p := &Patch{
		originX:       utils.Int(buff[0:2]),
		originY:       utils.Int(buff[2:4]),
		pictureID:     utils.I16(buff[4:6]),
		stepDirection: utils.I16(buff[6:8]),
		colorMap:      utils.I16(buff[8:10]),
	}
	return p
}

func (p *Patch) loadPicture(buff []byte) {
	p.DoomPicture = NewDoomPicture(buff)
}
