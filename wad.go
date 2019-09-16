package goom

import (
	"encoding/binary"
	"fmt"
	"os"
	"strings"
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

// Lump consists of a number of entries, each with a length of 16 bytes.
type Lump struct {
	Name     string
	Size     int
	Position int
	Data     []byte
}

func (l *Lump) dataFromBuff(data []byte) {
	if l.Size > 0 {
		l.Data = data[l.Position-wadHeaderSize : l.Position-wadHeaderSize+l.Size]
	}
}

type WadManager struct {
	wads []*WAD
}

func NewWadManager() *WadManager {
	return &WadManager{}
}

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
		return nil, fmt.Errorf("could not open WAD file: %s", err.Error())
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

func (w *WAD) loadLumpsFromFD(fd *os.File, data []byte) error {
	buff := make([]byte, lumpSize)
	w.lumps = make([]Lump, w.NumLumps)
	for i := 0; i < w.NumLumps; i++ {
		sz, err := fd.Read(buff)
		if err != nil {
			return fmt.Errorf("could read lump %d: %s", i, err.Error())
		}
		if sz != lumpSize {
			return fmt.Errorf("could read lump %d: %s", i, "wrong size")
		}
		l := Lump{
			Position: int(binary.LittleEndian.Uint32(buff[0:4])),
			Size:     int(binary.LittleEndian.Uint32(buff[4:8])),
			Name:     wadString(buff[8:16]),
		}
		l.dataFromBuff(data)
		w.lumps[i] = l
	}

	return nil
}

func (w *WAD) GetLumps() []Lump {
	return w.lumps
}

func (w *WAD) GetLump(name string) *Lump {
	for _, l := range w.lumps {
		if name == l.Name {
			return &l
		}
	}
	return nil
}

func wadString(buff []byte) string {
	return strings.TrimRight(string(buff), "\x00")
}
