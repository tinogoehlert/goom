package opengl

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/tinogoehlert/goom"
	"github.com/tinogoehlert/goom/level"
)

type doomLevel struct {
	name       string
	walls      []*Mesh
	subSectors []*subSector
	mapRef     *level.Level
}

type subSector struct {
	floors   []*Mesh
	ceilings []*Mesh
	walls    []*Mesh
	sector   level.Sector
	ref      level.SubSector
}

func RegisterMap(m *level.Level, gd *goom.GameData) *doomLevel {
	l := doomLevel{
		name:       m.Name,
		mapRef:     m,
		subSectors: make([]*subSector, 0, len(m.SubSectors("GL_SSECT"))),
	}

	for _, ssect := range m.SubSectors("GL_SSECT") {
		var s = &subSector{ref: ssect}
		s.addFlats(m, gd)
		s.addWalls(m, gd)
		l.subSectors = append(l.subSectors, s)
	}
	return &l
}

func (s *subSector) addFlats(md *level.Level, gd *goom.GameData) {
	s.floors, s.ceilings = []*Mesh{}, []*Mesh{}
	var (
		fseg   = s.ref.Segments()[0]
		vfs    = md.Vert(fseg.StartVert())
		line   = md.LinesDefs[fseg.LineDef()]
		side   = md.SideDefs[line.Right]
		sector = md.Sectors[side.Sector]
	)

	if fseg.Direction() == 1 {
		side = md.SideDefs[line.Left]
		sector = md.Sectors[side.Sector]
	}

	floorData := []float32{}
	ceilData := []float32{}
	for _, seg := range s.ref.Segments() {
		var (
			s = md.Vert(seg.StartVert())
			e = md.Vert(seg.EndVert())
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
	if len(gd.Flat(sector.FloorTexture())) > 0 {
		fm := NewMesh(floorData, sector.LightLevel(), gd.Flat(sector.FloorTexture())[0], gd.DefaultPalette(), sector.FloorTexture())
		s.floors = AddMesh(s.floors, fm)
	}
	if len(gd.Flat(sector.CeilTexture())) > 0 {
		cm := NewMesh(ceilData, md.Sectors[side.Sector].LightLevel(), gd.Flat(sector.CeilTexture())[0], gd.DefaultPalette(), sector.FloorTexture())
		s.ceilings = AddMesh(s.ceilings, cm)
	}
}

func (s *subSector) addWalls(md *level.Level, gd *goom.GameData) {
	s.walls = []*Mesh{}
	for _, seg := range s.ref.Segments() {
		if seg.LineDef() == -1 {
			continue
		}
		line := md.LinesDefs[seg.LineDef()]
		side := md.SideDefs[line.Right]

		otherSide := md.OtherSide(&line, seg)
		sector := md.Sectors[side.Sector]

		var (
			start = md.Vert(uint32(line.Start))
			end   = md.Vert(uint32(line.End))
		)

		if side.Middle() == "-" &&
			side.Upper() == "-" &&
			side.Lower() == "-" {
			continue
		}

		if side.Upper() != "-" && otherSide != nil {
			oppositeSector := md.Sectors[otherSide.Sector]
			wallData := []float32{
				-start.X(), sector.CeilHeight(), start.Y(), 0.0, 1.0,
				-start.X(), oppositeSector.CeilHeight(), start.Y(), 0.0, 0.0,
				-end.X(), oppositeSector.CeilHeight(), end.Y(), 1.0, 0.0,

				-end.X(), oppositeSector.CeilHeight(), end.Y(), 1.0, 0.0,
				-end.X(), sector.CeilHeight(), end.Y(), 1.0, 1.0,
				-start.X(), sector.CeilHeight(), start.Y(), 0.0, 1.0,
			}
			if gd.Texture(side.Upper()) != nil {
				wm := NewMesh(wallData, sector.LightLevel(), gd.Texture(side.Upper()), gd.DefaultPalette(), side.Upper())
				s.walls = AddMesh(s.walls, wm)
			}
		}

		if side.Lower() != "-" && otherSide != nil {
			oppositeSector := md.Sectors[otherSide.Sector]
			wallData := []float32{
				-start.X(), sector.FloorHeight(), start.Y(), 0.0, 1.0,
				-start.X(), oppositeSector.FloorHeight(), start.Y(), 0.0, 0.0,
				-end.X(), oppositeSector.FloorHeight(), end.Y(), 1.0, 0.0,

				-end.X(), oppositeSector.FloorHeight(), end.Y(), 1.0, 0.0,
				-end.X(), sector.FloorHeight(), end.Y(), 1.0, 1.0,
				-start.X(), sector.FloorHeight(), start.Y(), 0.0, 1.0,
			}
			if gd.Texture(side.Lower()) != nil {
				wm := NewMesh(wallData, sector.LightLevel(), gd.Texture(side.Lower()), gd.DefaultPalette(), side.Lower())
				s.walls = AddMesh(s.walls, wm)
			}
		}

		if side.Middle() != "-" {
			wallData := []float32{
				-start.X(), sector.FloorHeight(), start.Y(), 0.0, 1.0,
				-start.X(), sector.CeilHeight(), start.Y(), 0.0, 0.0,
				-end.X(), sector.CeilHeight(), end.Y(), 1.0, 0.0,

				-end.X(), sector.CeilHeight(), end.Y(), 1.0, 0.0,
				-end.X(), sector.FloorHeight(), end.Y(), 1.0, 1.0,
				-start.X(), sector.FloorHeight(), start.Y(), 0.0, 1.0,
			}
			wm := NewMesh(wallData, sector.LightLevel(), gd.Texture(side.Middle()), gd.DefaultPalette(), side.Middle())
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
