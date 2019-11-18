package level

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/tinogoehlert/goom/utils"

	"github.com/tinogoehlert/goom/wad"
)

const (
	sectorSize  = 26
	ssectSize   = 4
	glSsectSize = 8
)

type Sector struct {
	floorHeight    float32
	ceilingHeight  float32
	floorTexture   string
	ceilingTexture string
	lightLevel     float32
	sectorType     int16
	// A number that makes the sector a target of the action specified by any linedef with the same tag number.
	// Used, for example, to alter the sector's ceiling and/or floor height, lighting level or flats.
	tag int16
}

// LightLevel the amount of light applied to the sector and it's childs
func (s *Sector) LightLevel() float32 { return s.lightLevel }

// Type type of the sector
func (s *Sector) Type() int16 { return s.sectorType }

// CeilHeight gets ceiling height
func (s *Sector) CeilHeight() float32 { return s.ceilingHeight }

// CeilTexture gets ceiling Texture
func (s *Sector) CeilTexture() string { return s.ceilingTexture }

// FloorHeight gets ceiling height
func (s *Sector) FloorHeight() float32 { return s.floorHeight }

// FloorTexture gets ceiling Texture
func (s *Sector) FloorTexture() string { return s.floorTexture }

// Tag sector tags for segs
func (s *Sector) Tag() int16 { return s.tag }

func newSectorsFromLump(lump *wad.Lump) ([]Sector, error) {
	if lump.Size%sectorSize != 0 {
		return nil, fmt.Errorf("size missmatch")
	}

	var sectorCount = lump.Size / sectorSize
	var sectors = make([]Sector, sectorCount)

	for i := 0; i < sectorCount; i++ {
		b := lump.Data[(i * sectorSize) : (i*sectorSize)+sectorSize]
		sectors[i] = Sector{
			floorHeight:    utils.Int16Tof32(b[0:2]),
			ceilingHeight:  utils.Int16Tof32(b[2:4]),
			floorTexture:   strings.TrimRight(string(b[4:12]), "\x00"),
			ceilingTexture: strings.TrimRight(string(b[12:20]), "\x00"),
			lightLevel:     utils.Int16Tof32(b[20:22]),
			sectorType:     utils.I16(b[20:24]),
			tag:            utils.I16(b[20:26]),
		}
	}
	return sectors, nil
}

type SubSector struct {
	Count    uint32
	firstSeg uint32
	segments []Segment
}

func (ssect *SubSector) Segments() []Segment {
	return ssect.segments
}

func newSSectsFromLump(lump *wad.Lump, segs []Segment) ([]SubSector, error) {
	var (
		ssectCount = len(lump.Data) / ssectSize
		subsectors = make([]SubSector, (lump.Size)/ssectSize)
	)
	for i := 0; i < ssectCount; i++ {
		vb := lump.Data[(i * ssectSize) : (i*ssectSize)+ssectSize]
		ssect := SubSector{
			Count:    uint32(int16(binary.LittleEndian.Uint16(vb[0:2]))),
			firstSeg: uint32(int16(binary.LittleEndian.Uint16(vb[2:4]))),
		}
		ssect.segments = segs[ssect.firstSeg : ssect.firstSeg+ssect.Count]
		subsectors[i] = ssect
	}
	return subsectors, nil
}

func newGLSSectsV5FromLump(lump *wad.Lump, segs []Segment) ([]SubSector, error) {
	var (
		ssectCount = len(lump.Data) / glSsectSize
		subsectors = make([]SubSector, (lump.Size)/glSsectSize)
	)
	for i := 0; i < ssectCount; i++ {
		vb := lump.Data[(i * glSsectSize) : (i*glSsectSize)+glSsectSize]
		ssect := SubSector{
			Count:    binary.LittleEndian.Uint32(vb[0:4]),
			firstSeg: binary.LittleEndian.Uint32(vb[4:8]),
		}
		ssect.segments = segs[ssect.firstSeg : ssect.firstSeg+ssect.Count]
		subsectors[i] = ssect
	}
	return subsectors, nil
}
