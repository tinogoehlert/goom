package wad

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/tinogoehlert/goom/utils"
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
			Name:     utils.WadString(buff[8:16]),
		}
		l.dataFromBuff(data)
		w.lumps[i] = l
	}

	return nil
}
