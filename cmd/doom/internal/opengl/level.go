package opengl

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/tinogoehlert/goom"
)

type level struct {
	name       string
	walls      []*Mesh
	subSectors []*subSector
	mapRef     *goom.Map
}

type subSector struct {
	floors   []*Mesh
	ceilings []*Mesh
	walls    []*Mesh
	sector   goom.Sector
	ref      goom.SubSector
}

func RegisterMap(m *goom.Map, gfx *goom.Graphics) *level {
	l := level{
		name:       m.Name,
		mapRef:     m,
		subSectors: make([]*subSector, 0, len(m.SubSectors("GL_SSECT"))),
	}

	for _, ssect := range m.SubSectors("GL_SSECT") {
		var s = &subSector{ref: ssect}
		s.addFlats(m, gfx)
		s.addWalls(m, gfx)
		l.subSectors = append(l.subSectors, s)
	}
	return &l
}

func (s *subSector) addFlats(md *goom.Map, gfx *goom.Graphics) {
	s.floors, s.ceilings = []*Mesh{}, []*Mesh{}
	var (
		fseg   = s.ref.Segments()[0]
		vfs    = md.Vert(fseg.GetStartVert())
		line   = md.LinesDefs[fseg.GetLineDef()]
		side   = md.SideDefs[line.Right]
		sector = md.Sectors[side.Sector]
	)

	if fseg.GetDirection() == 1 {
		side = md.SideDefs[line.Left]
		sector = md.Sectors[side.Sector]
	}

	floorData := []float32{}
	ceilData := []float32{}
	for _, seg := range s.ref.Segments() {
		var (
			s = md.Vert(seg.GetStartVert())
			e = md.Vert(seg.GetEndVert())
		)

		floorData = append(floorData,
			-vfs.X(), sector.FloorHeight(), vfs.Y(), -vfs.X()/64, vfs.Y()/64,
			-s.X(), sector.FloorHeight(), s.Y(), -s.X()/64, s.Y()/64,
			-e.X(), sector.FloorHeight(), e.Y(), -e.X()/64, e.Y()/64,
		)
		ceilData = append(ceilData,
			-vfs.X(), sector.CeilHeight(), vfs.Y(), -vfs.X()/64, vfs.Y()/64,
			-s.X(), sector.CeilHeight(), s.Y(), -s.X()/64, s.Y()/64,
			-e.X(), sector.CeilHeight(), e.Y(), -e.X()/64, e.Y()/64,
		)
	}

	s.sector = sector
	fm := NewMesh(floorData, sector.LightLevel(), gfx.GetFlat(sector.FloorTexture()), gfx.Palette(0))
	cm := NewMesh(ceilData, md.Sectors[side.Sector].LightLevel(), gfx.GetFlat(sector.CeilTexture()), gfx.Palette(0))
	s.floors = AddMesh(s.floors, fm)
	s.ceilings = AddMesh(s.ceilings, cm)
}

func (s *subSector) addWalls(md *goom.Map, gfx *goom.Graphics) {
	s.walls = []*Mesh{}
	for _, seg := range s.ref.Segments() {
		if seg.GetLineDef() == -1 {
			continue
		}
		line := md.LinesDefs[seg.GetLineDef()]
		side := md.SideDefs[line.Right]

		otherSide := md.OtherSide(&line, seg)
		sector := md.Sectors[side.Sector]

		var (
			start = md.Vert(uint32(line.Start))
			end   = md.Vert(uint32(line.End))
		)

		if side.MiddleName.ToString() == "-" &&
			side.UpperName.ToString() == "-" &&
			side.Lowername.ToString() == "-" {
			continue
		}

		if side.UpperName.ToString() != "-" && otherSide != nil {
			oppositeSector := md.Sectors[otherSide.Sector]
			wallData := []float32{
				-start.X(), sector.CeilHeight(), start.Y(), 0.0, 1.0,
				-start.X(), oppositeSector.CeilHeight(), start.Y(), 0.0, 0.0,
				-end.X(), oppositeSector.CeilHeight(), end.Y(), 1.0, 0.0,

				-end.X(), oppositeSector.CeilHeight(), end.Y(), 1.0, 0.0,
				-end.X(), sector.CeilHeight(), end.Y(), 1.0, 1.0,
				-start.X(), sector.CeilHeight(), start.Y(), 0.0, 1.0,
			}
			wm := NewMesh(wallData, sector.LightLevel(), gfx.GetTexture(side.UpperName.ToString()), gfx.Palette(0))
			s.walls = AddMesh(s.walls, wm)
		}

		if side.Lowername.ToString() != "-" && otherSide != nil {
			oppositeSector := md.Sectors[otherSide.Sector]
			wallData := []float32{
				-start.X(), sector.FloorHeight(), start.Y(), 0.0, 1.0,
				-start.X(), oppositeSector.FloorHeight(), start.Y(), 0.0, 0.0,
				-end.X(), oppositeSector.FloorHeight(), end.Y(), 1.0, 0.0,

				-end.X(), oppositeSector.FloorHeight(), end.Y(), 1.0, 0.0,
				-end.X(), sector.FloorHeight(), end.Y(), 1.0, 1.0,
				-start.X(), sector.FloorHeight(), start.Y(), 0.0, 1.0,
			}
			wm := NewMesh(wallData, sector.LightLevel(), gfx.GetTexture(side.Lowername.ToString()), gfx.Palette(0))
			s.walls = AddMesh(s.walls, wm)
		}

		if side.MiddleName.ToString() != "-" {
			wallData := []float32{
				-start.X(), sector.FloorHeight(), start.Y(), 0.0, 1.0,
				-start.X(), sector.CeilHeight(), start.Y(), 0.0, 0.0,
				-end.X(), sector.CeilHeight(), end.Y(), 1.0, 0.0,

				-end.X(), sector.CeilHeight(), end.Y(), 1.0, 0.0,
				-end.X(), sector.FloorHeight(), end.Y(), 1.0, 1.0,
				-start.X(), sector.FloorHeight(), start.Y(), 0.0, 1.0,
			}
			wm := NewMesh(wallData, sector.LightLevel(), gfx.GetTexture(side.MiddleName.ToString()), gfx.Palette(0))
			s.walls = AddMesh(s.walls, wm)
		}
	}
}

func (s *subSector) Draw() {
	for i := 0; i < len(s.floors); i++ {
		s.floors[i].DrawMesh(gl.TRIANGLE_FAN)
		s.ceilings[i].DrawMesh(gl.TRIANGLE_FAN)
	}
	for _, w := range s.walls {
		w.DrawMesh(gl.TRIANGLES)
	}
}
