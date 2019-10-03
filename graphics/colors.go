package graphics

import (
	"image/color"

	"github.com/tinogoehlert/goom/wad"
)

const (
	numPalettes = 14
)

// Palette DOOM color palette
type Palette struct {
	Colors [256]color.RGBA
}

var defaultPalette Palette

// Palettes one or more palettes
type Palettes [numPalettes]Palette

// NewPalettes one or more palettes
func NewPalettes(w *wad.WAD) (*Palettes, error) {
	var palettes = Palettes{}
	for _, lump := range w.Lumps() {
		if lump.Name == "PLAYPAL" {
			for i := 0; i < numPalettes; i++ {
				p := Palette{}
				for ci := 0; ci < 256*3; ci += 3 {
					p.Colors[ci/3].R = lump.Data[ci]
					p.Colors[ci/3].G = lump.Data[ci+1]
					p.Colors[ci/3].B = lump.Data[ci+2]
					p.Colors[ci/3].A = 255
				}
				if i == 0 {
					defaultPalette = p
				}
				palettes[i] = p
			}
			return &palettes, nil
		}
	}
	return &palettes, nil
}

func DefaultPalette() Palette {
	return defaultPalette
}
