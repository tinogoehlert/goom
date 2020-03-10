package level

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/tinogoehlert/goom/utils"
	"github.com/tinogoehlert/goom/wad"
)

const (
	GLNodesName  = "GL_NODES"
	NodesName    = "NODES"
	GLSsectsName = "GL_SSECT"
	SSectsName   = "SSECTORS"
	GLSegsName   = "GL_SEGS"
	SegsName     = "SEGS"
)

// Level - A map in Doom is made up of several lumps,
// each containing specific data required to construct and execute the map.
type Level struct {
	Name      string
	Things    []Thing
	LinesDefs []LineDef
	Sectors   []Sector
	SideDefs  []SideDef
	Walls     []Wall
	// private pools
	vertexPool map[string][]utils.Vec2
	segPool    map[string][]Segment
	ssectPool  map[string][]SubSector
	nodePool   map[string][]Node
}

// Store stores map of levels
type Store map[string]*Level

// NewStore creates new level store
func NewStore() Store {
	return make(Store)
}

// LoadWAD loads wad into store
func (s Store) LoadWAD(w *wad.WAD) error {
	var (
		nameRegex   = regexp.MustCompile(`^E\dM\d|^MAP\d\d`)
		glNameRegex = regexp.MustCompile(`^GL_E\dM\d|^GL_MAP\d\d`)
	)
	for i := 0; i < len(w.Lumps()); i++ {
		lump := w.Lumps()[i]
		switch {
		case nameRegex.Match([]byte(lump.Name)):
			l, err := NewLevel(w.Lumps()[i+1 : i+9])
			if err != nil {
				fmt.Printf("ERROR %s: %s\n", lump.Name, err.Error())
				continue
			}
			l.Name = lump.Name
			s[l.Name] = l
			i += 7
		case glNameRegex.Match([]byte(lump.Name)):
			appendGLNodes(s[lump.Name[3:]], w.Lumps()[i+1:i+5])
			i += 3
		}
	}
	return nil
}

// NewLevel Loads a level from a list of lumps
func NewLevel(lumps []wad.Lump) (l *Level, err error) {
	l = &Level{
		vertexPool: make(map[string][]utils.Vec2),
		segPool:    make(map[string][]Segment),
		ssectPool:  make(map[string][]SubSector),
		nodePool:   make(map[string][]Node),
	}
	l.Things, err = loadThingsFromLump(&lumps[0])
	if err != nil {
		return nil, fmt.Errorf("could not read things from WAD: %s", err.Error())
	}
	l.LinesDefs, err = newLinedefsFromLump(&lumps[1])
	if err != nil {
		return nil, fmt.Errorf("could not read linedefs from WAD: %s", err.Error())
	}
	l.SideDefs, err = newSidesDefFromLump(&lumps[2])
	if err != nil {
		return nil, fmt.Errorf("could not read sidedefs from WAD: %s", err.Error())
	}
	l.vertexPool[lumps[3].Name], err = newVerticesFromLump(&lumps[3])
	if err != nil {
		return nil, fmt.Errorf("could not read vertices from WAD: %s", err.Error())
	}
	l.segPool[lumps[4].Name], err = newSegmentsFromLump(&lumps[4])
	if err != nil {
		return nil, fmt.Errorf("could not read segs from WAD: %s", err.Error())
	}
	segs := l.segPool[lumps[4].Name]
	l.ssectPool[lumps[5].Name], err = newSSectsFromLump(&lumps[5], segs)

	if err != nil {
		return nil, fmt.Errorf("could not read subsectors from WAD: %s", err.Error())
	}
	l.nodePool[lumps[6].Name], err = newNodesFromLump(&lumps[6])
	if err != nil {
		return nil, fmt.Errorf("could not read nodes from WAD: %s", err.Error())
	}
	l.Sectors, err = newSectorsFromLump(&lumps[7])
	if err != nil {
		return nil, fmt.Errorf("could not read sectors from WAD: %s", err.Error())
	}

	for _, line := range l.LinesDefs {
		l.Walls = append(l.Walls, NewWall(&line, l))
	}

	return l, nil
}

func appendGLNodes(l *Level, lumps []wad.Lump) (err error) {
	if l == nil {
		panic("level not found")
	}
	l.vertexPool[lumps[0].Name], err = newVerticesFromLump(&lumps[0])
	if err != nil {
		return fmt.Errorf("could not load GL_VERT: %s", err.Error())
	}
	l.segPool[lumps[1].Name], err = newGLSegmentsFromLump(&lumps[1])
	if err != nil {
		return fmt.Errorf("could not load GL_SEGS: %s", err.Error())
	}
	segs := l.segPool[lumps[1].Name]
	l.ssectPool[lumps[2].Name], err = newGLSSectsV5FromLump(&lumps[2], segs)
	if err != nil {
		return fmt.Errorf("could not load GL_SEGS: %s", err.Error())
	}
	l.nodePool[lumps[3].Name], err = newGLNodesFromLump(&lumps[3])
	if err != nil {
		return fmt.Errorf("could not read GL_NODES from WAD: %s", err.Error())
	}
	return nil
}

func (l *Level) buildWalls() error {
	if len(l.LinesDefs) == 0 {
		return errors.New("no linedefs given")
	}

	return nil
}

// Vert gets a vert
func (l *Level) Vert(id uint32) utils.Vec2 {
	if utils.MagicU32(id).MagicBit() {
		return l.vertexPool["GL_VERT"][utils.MagicU32(id).Uint32()]
	}
	return l.vertexPool["VERTEXES"][id]
}

// Segments gets segs (SEGS or GL_SEGS)
func (l *Level) Segments(name string) []Segment {
	return l.segPool[name]
}

// Nodes gets BSP Nodes (NODES or GL_NODES)
func (l *Level) Nodes(name string) []Node {
	return l.nodePool[name]
}

// SubSectors gets Subsectors (GL_SSECT or SSECT)
func (l *Level) SubSectors(name string) []SubSector {
	return l.ssectPool[name]
}

// OtherSide gets the opposite side of a side / seg
func (l *Level) OtherSide(line *LineDef, seg Segment) *SideDef {
	if seg.Direction() == 0 {
		if line.Left == -1 {
			return nil
		}
		return &l.SideDefs[line.Left]
	}
	return &l.SideDefs[line.Right]
}

// SectorFromSSect gets the sector from a subsector
func (l *Level) SectorFromSSect(ssect *SubSector) *Sector {
	var (
		fseg   = ssect.Segments()[0]
		line   = l.LinesDefs[fseg.LineDef()]
		side   = l.SideDefs[line.Right]
		sector = l.Sectors[side.Sector]
	)

	if fseg.Direction() == 1 {
		side = l.SideDefs[line.Left]
		sector = l.Sectors[side.Sector]
	}
	return &sector
}

// WalkBsp walks through the node tree
func (l *Level) WalkBsp(fn func(index int, n *Node, b BBox)) error {
	nodes, ok := l.nodePool[GLNodesName]
	if !ok {
		return fmt.Errorf("could not find %s", GLNodesName)
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
func (l *Level) FindPositionInBsp(nodeType string, x, y float32) (*SubSector, error) {
	ssects := l.SubSectors(SSectsName)
	if nodeType == GLNodesName {
		ssects = l.SubSectors(GLSsectsName)
	}
	nodes, ok := l.nodePool[nodeType]
	if !ok {
		return nil, fmt.Errorf("could not find %s", nodeType)
	}
	n := nodes[len(nodes)-1]
	for i := 0; i < len(nodes); i++ {
		if x*n.direction.X()+y*n.direction.Y() > n.dirDeg {
			if n.Left.IsSubSector() {
				return &ssects[n.Left.Num()], nil
			}
			n = nodes[n.Left.Num()]
		} else {
			if n.Right.IsSubSector() {
				return &ssects[n.Right.Num()], nil
			}
			n = nodes[n.Right.Num()]
		}
	}
	return nil, fmt.Errorf("not found")
}
