package goom

import (
	"fmt"
	"regexp"
)

const (
	GLNodesName  = "GL_NODES"
	NodesName    = "NODES"
	GLSsectsName = "GL_SSECT"
	SSectsName   = "SSECT"
	GLSegsName   = "GL_SEGS"
	SegsName     = "SEGS"
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

// LoadMaps loads all the maps
func (wm *WadManager) LoadMaps() (maps []Map, err error) {
	var (
		nameRegex   = regexp.MustCompile(`^E\dM\d|^MAP\d\d`)
		glNameRegex = regexp.MustCompile(`^GL_E\dM\d|^GL_MAP\d\d`)
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

// Vert gets a vert
func (m *Map) Vert(id uint32) *Vertex {
	if MagicU32(id).MagicBit() {
		return &m.vertexPool["GL_VERT"][MagicU32(id).Uint32()]
	}
	return &m.vertexPool["VERTEXES"][id]
}

// Segments gets segs (SEGS or GL_SEGS)
func (m *Map) Segments(name string) []Segment {
	return m.segPool[name]
}

// Nodes gets BSP Nodes (NODES or GL_NODES)
func (m *Map) Nodes(name string) []Node {
	return m.nodePool[name]
}

// SubSectors gets Subsectors (GL_SSECT or SSECT)
func (m *Map) SubSectors(name string) []SubSector {
	return m.ssectPool[name]
}

// OtherSide gets the opposite side of a side / seg
func (m *Map) OtherSide(line *LineDef, seg Segment) *SideDef {
	if seg.GetDirection() == 0 {
		if line.Left == -1 {
			return nil
		}
		return &m.SideDefs[line.Left]
	}
	return &m.SideDefs[line.Right]
}

// SectorFromSSect gets the sector from a subsector
func (m *Map) SectorFromSSect(ssect *SubSector) *Sector {
	var (
		fseg   = ssect.Segments()[0]
		line   = m.LinesDefs[fseg.GetLineDef()]
		side   = m.SideDefs[line.Right]
		sector = m.Sectors[side.Sector]
	)

	if fseg.GetDirection() == 1 {
		side = m.SideDefs[line.Left]
		sector = m.Sectors[side.Sector]
	}
	return &sector
}

// WalkBsp walks through the node tree
func (m *Map) WalkBsp(nodeType string, fn func(index int, n *Node, b BBox)) error {
	nodes, ok := m.nodePool[nodeType]
	if !ok {
		return fmt.Errorf("could not find %s", nodeType)
	}
	for i := len(nodes) - 1; i >= 0; i-- {
		if nodes[i].Right.IsSubSector() {
			fn(int(nodes[i].Right.Num()), &nodes[i], nodes[i].RightBBox)
		}
		if nodes[i].Left.IsSubSector() {
			fn(int(nodes[i].Left.Num()), &nodes[i], nodes[i].LeftBBox)
		}
	}
	return nil
}

// FindPositionInBsp finds a position in the nodeTree
func (m *Map) FindPositionInBsp(nodeType string, x, y float32) (*SubSector, error) {
	ssects := m.SubSectors(SSectsName)
	if nodeType == GLNodesName {
		ssects = m.SubSectors(GLSsectsName)
	}
	nodes, ok := m.nodePool[nodeType]
	if !ok {
		return nil, fmt.Errorf("could not find %s", nodeType)
	}
	var (
		lastSub    *SubSector
		lastHeight = float32(-10000)
	)
	for i := 0; i < len(nodes); i++ {
		if nodes[i].Right.IsSubSector() {
			if nodes[i].RightBBox.PosInBox(x, y) {
				height := m.SectorFromSSect(&ssects[nodes[i].Right.Num()]).FloorHeight()
				if height > lastHeight {
					lastSub = &ssects[nodes[i].Right.Num()]
					lastHeight = height
				}
			}
		}
		if nodes[i].Left.IsSubSector() {
			if nodes[i].LeftBBox.PosInBox(x, y) {
				height := m.SectorFromSSect(&ssects[nodes[i].Left.Num()]).FloorHeight()
				if height > lastHeight {
					lastSub = &ssects[nodes[i].Left.Num()]
					lastHeight = height
				}
			}
		}
	}

	if lastSub == nil {
		return nil, fmt.Errorf("not found")
	}
	return lastSub, nil
}
