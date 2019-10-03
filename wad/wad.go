package wad

import (
	"encoding/binary"
	"fmt"
	"os"
)

const (
	wadHeaderSize = 12
	lumpSize      = 16
)

// Type determines the type of a WAD file (IWAD or PWAD)
type Type string

const (
	// TypeInternal refers to a WAD file which contains all of the game data for a complete game.
	TypeInternal Type = "IWAD"
	// TypePatch containing lumps of data created by a user as an add-on.
	TypePatch Type = "PWAD"
)

// WadManager holds list of wads
type WadManager struct {
	wads []*WAD
}

// NewWadManager creates a new WAD Manager
func NewWadManager() *WadManager {
	return &WadManager{}
}

// LoadFile processes a WAD file
func (wm *WadManager) LoadFile(file string) error {
	w, err := NewWADFromFile(file)
	if err != nil {
		return err
	}
	wm.wads = append(wm.wads, w)
	return nil
}

// WAD (which, according to the Doom Bible, is an acrostic for "Where's All the Data?")
// is the file format used by Doom and all Doom-engine-based games for storing data.
// A WAD file consists of a header, a directory, and the data lumps that make up the resources stored within the file.
type WAD struct {
	// Type Defines whether the WAD is an IWAD or a Pgoom.
	Type Type
	// An integer specifying the number of lumps in the goom.
	NumLumps int
	// An integer holding a pointer to the location of the directory.
	InfoTableOFS int
	lumps        []Lump
}

// NewWADFromFile Loads WAD from the given file
func NewWADFromFile(file string) (*WAD, error) {
	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	// load header
	header := make([]byte, wadHeaderSize)
	sz, err := fd.Read(header)
	if err != nil {
		return nil, fmt.Errorf("could read WAD header: %s", err.Error())
	}
	if sz != wadHeaderSize {
		return nil, fmt.Errorf("could read WAD header: wrong size (%d)", sz)
	}
	wt := Type(header[:4])
	if wt != TypeInternal && wt != TypePatch {
		return nil, fmt.Errorf("unsupported WAD type: " + string(wt))
	}

	wad := &WAD{
		Type:         wt,
		NumLumps:     int(binary.LittleEndian.Uint32(header[4:8])),
		InfoTableOFS: int(binary.LittleEndian.Uint32(header[8:12])),
	}

	// load header
	data := make([]byte, wad.InfoTableOFS-wadHeaderSize)
	sz, err = fd.Read(data)
	if err != nil {
		return nil, fmt.Errorf("could read WAD data: %s", err.Error())
	}

	return wad, wad.loadLumpsFromFD(fd, data)
}

// Lumps gets slice of lumps
func (w *WAD) Lumps() []Lump {
	return w.lumps
}

// Lump get lump by name
func (w *WAD) Lump(name string) *Lump {
	for _, l := range w.lumps {
		if name == l.Name {
			return &l
		}
	}
	return nil
}
