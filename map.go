package goom

import (
	"fmt"
	"regexp"
)

// Map - A map in Doom is made up of several lumps,
// each containing specific data required to construct and execute the map.
type Map struct {
	Name       string
	Things     []Thing
	LinesDefs  []LineDef
	Sectors    []Sector
	SideDefs   []SideDef
	vertexPool map[string][]Vertex
	segPool    map[string][]Segment
	ssectPool  map[string][]SubSector
	nodePool   map[string][]Node
}

func (m *Map) Vert(id uint32) *Vertex {
	if MagicU32(id).MagicBit() {
		return &m.vertexPool["GL_VERT"][MagicU32(id).Uint32()]
	}
	return &m.vertexPool["VERTEXES"][id]
}

func (m *Map) Segments(name string) []Segment {
	return m.segPool[name]
}

func (m *Map) Nodes(name string) []Node {
	return m.nodePool[name]
}

func (m *Map) SubSectors(name string) []SubSector {
	return m.ssectPool[name]
}

func (m *Map) OtherSide(line *LineDef, seg Segment) *SideDef {
	if seg.GetDirection() == 0 {
		if line.Left == -1 {
			return nil
		}
		return &m.SideDefs[line.Left]
	}
	return &m.SideDefs[line.Right]
}

func (wm *WadManager) LoadMaps() (maps []Map, err error) {
	var (
		nameRegex   = regexp.MustCompile(`^E\dM\d|^MAP\d\d`)
		glNameRegex = regexp.MustCompile(`^GL_E\dM\d|^MAP\d\d`)
	)
	maps = make([]Map, 0)
	tmpMaps := make(map[string]*Map)
	for _, w := range wm.wads {
		for i := 0; i < len(w.lumps); i++ {
			l := w.lumps[i]
			switch {
			case nameRegex.Match([]byte(l.Name)):
				m, err := loadMap(w.lumps[i+1 : i+9])
				if err != nil {
					return nil, fmt.Errorf("ERROR %s: %s", l.Name, err.Error())
				}
				m.Name = l.Name
				maps = append(maps, *m)
				tmpMaps[l.Name] = m
				i += 7
			case glNameRegex.Match([]byte(l.Name)):
				appendGLNodes(tmpMaps[l.Name[3:]], w.lumps[i+1:i+5])
				i += 3
			}
		}
	}
	return maps, nil
}

func loadMap(lumps []Lump) (m *Map, err error) {
	m = &Map{
		vertexPool: make(map[string][]Vertex),
		segPool:    make(map[string][]Segment),
		ssectPool:  make(map[string][]SubSector),
		nodePool:   make(map[string][]Node),
	}
	m.Things, err = loadThingsFromLump(&lumps[0])
	if err != nil {
		return nil, fmt.Errorf("could not read things from WAD: %s", err.Error())
	}
	m.LinesDefs, err = newLinedefsFromLump(&lumps[1])
	if err != nil {
		return nil, fmt.Errorf("could not read linedefs from WAD: %s", err.Error())
	}
	m.SideDefs, err = newSidesDefFromLump(&lumps[2])
	if err != nil {
		return nil, fmt.Errorf("could not read sidedefs from WAD: %s", err.Error())
	}
	m.vertexPool[lumps[3].Name], err = newVerticesFromLump(&lumps[3])
	if err != nil {
		return nil, fmt.Errorf("could not read vertices from WAD: %s", err.Error())
	}
	m.segPool[lumps[4].Name], err = newSegmentsFromLump(&lumps[4])
	if err != nil {
		return nil, fmt.Errorf("could not read segs from WAD: %s", err.Error())
	}
	segs := m.segPool[lumps[4].Name]
	m.ssectPool[lumps[5].Name], err = newSSectsFromLump(&lumps[5], segs)
	if err != nil {
		return nil, fmt.Errorf("could not read subsectors from WAD: %s", err.Error())
	}
	m.nodePool[lumps[6].Name], err = newNodesFromLump(&lumps[6])
	if err != nil {
		return nil, fmt.Errorf("could not read nodes from WAD: %s", err.Error())
	}
	m.Sectors, err = newSectorsFromLump(&lumps[7])
	if err != nil {
		return nil, fmt.Errorf("could not read sectors from WAD: %s", err.Error())
	}
	return m, nil
}

func appendGLNodes(m *Map, lumps []Lump) (err error) {
	m.vertexPool[lumps[0].Name], err = newVerticesFromLump(&lumps[0])
	if err != nil {
		return fmt.Errorf("could not load GL_VERT: %s", err.Error())
	}
	m.segPool[lumps[1].Name], err = newGLSegmentsFromLump(&lumps[1])
	if err != nil {
		return fmt.Errorf("could not load GL_SEGS: %s", err.Error())
	}
	segs := m.segPool[lumps[1].Name]
	m.ssectPool[lumps[2].Name], err = newGLSSectsV5FromLump(&lumps[2], segs)
	if err != nil {
		return fmt.Errorf("could not load GL_SEGS: %s", err.Error())
	}
	m.nodePool[lumps[3].Name], err = newGLNodesFromLump(&lumps[3])
	if err != nil {
		return fmt.Errorf("could not read GL_NODES from WAD: %s", err.Error())
	}
	return nil
}

/*


func (m *Map) verticesFromLump(lump Lump) error {
	if lump.Size%vertexSize != 0 {
		return fmt.Errorf("size missmatch")
	}

	var vertCount = lump.Size / vertexSize

	verts := make([]Vertex, vertCount)
	for i := 0; i < vertCount; i++ {
		buff := lump.Data[(i * vertexSize) : (i*vertexSize)+vertexSize]
		v, err := newVertexFromBuffer(buff)
		if err != nil {
			return fmt.Errorf("could not load vertex: %s", err.Error())
		}
		verts[i] = *v
	}
	m.Bsp.vertexPool[lump.Name] = verts
	return nil
}


func (m *Map) sectorsFromLump(lump Lump) error {
	if lump.Size%sectorSize != 0 {
		return fmt.Errorf("size missmatch")
	}

	var sectorCount = lump.Size / sectorSize
	m.Sectors = make([]Sector, sectorCount)
	for i := 0; i < sectorCount; i++ {
		r := bytes.NewBuffer(lump.Data[(i * sectorSize) : (i*sectorSize)+sectorSize])
		s := &Sector{}
		if err := binary.Read(r, binary.LittleEndian, s); err != nil {
			return err
		}
		m.Sectors[i] = *s
	}
	return nil
}

func (m *Map) nodesFromLump(lump Lump) error {
	r := bytes.NewBuffer(lump.Data)
	n := &Node{}
	if err := binary.Read(r, binary.LittleEndian, n); err != nil {
		return err
	}
	m.Nodes = append(m.Nodes, *n)
	return nil
}

// DoomMaps holds Doom Maps (levels)
type DoomMaps []*Map

func LoadMapsFromWAD(wad *WAD) (DoomMaps, error) {

	maps := make(DoomMaps, 0, 8)
	var m *Map
	var err error
	for _, lump := range goom.GetLumps() {
		if nameRegex.Match([]byte(lump.Name)) {
			m = &Map{Name: lump.Name}
			maps = append(maps, m)
		}
		switch lump.Name {
		case "THINGS":
			m.Things, err = loadThingsFromLump(&lump)
			if err != nil {
				return nil, err
			}
			break
		case "LINEDEFS":
			err := m.linesFromLump(lump)
			if err != nil {
				return nil, err
			}
		case "VERTEXES":
			err := m.verticesFromLump(lump)
			if err != nil {
				return nil, err
			}
		case "SIDEDEFS":
			err := m.sidesFromLump(lump)
			if err != nil {
				return nil, fmt.Errorf("could not read sides: %s", err.Error())
			}
		case "SECTORS":
			err := m.sectorsFromLump(lump)
			if err != nil {
				return nil, fmt.Errorf("could not read sectors: %s", err.Error())
			}
		case "NODES":
			err := m.nodesFromLump(lump)
			if err != nil {
				return nil, fmt.Errorf("could not read sectors: %s", err.Error())
			}
		}
	}
	return maps, nil
}
*/
