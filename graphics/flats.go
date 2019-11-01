package graphics

import (
	"regexp"

	"github.com/tinogoehlert/goom/wad"
)

// Flat floor / ceiling texture
type Flat struct {
	*DoomPicture
	name string
}

// NewFlat create a new flat
func NewFlat(name string, buff []byte) *Flat {
	f := &Flat{
		name: name,
	}

	f.DoomPicture = newDummyPicture(64, 64)

	for i, b := range buff {
		f.DoomPicture.data[i] = b
	}

	return f
}

// FlatStore stores map of flats
type FlatStore map[string][]*Flat

func NewFlatStore() FlatStore {
	return make(FlatStore)
}

func (fs FlatStore) LoadWAD(w *wad.WAD) {
	var (
		flatStartRegex = regexp.MustCompile(`^F?_START`)
		flatEndRegex   = regexp.MustCompile(`^F?_END`)
		lumps          = w.Lumps()
	)

	for i := 0; i < len(lumps); i++ {
		lump := &lumps[i]
		if flatStartRegex.Match([]byte(lump.Name)) {

			for {
				lump := &lumps[i]
				if lump.Size > 0 {
					fs.Append(lump.Name, NewFlat(lump.Name, lump.Data))
				}
				if flatEndRegex.Match([]byte(lump.Name)) {
					break
				}
				i++
			}
		}
	}
}

func (fs FlatStore) Append(name string, flat *Flat) {
	if byte(name[len(name)-1]) >= 60 && byte(name[len(name)-1]) <= 71 {
		name = name[:len(name)-1]
	}
	if _, ok := fs[name]; !ok {
		fs[name] = make([]*Flat, 0, 1)
	}
	fs[name] = append(fs[name], flat)
}
