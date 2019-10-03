package opengl

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/tinogoehlert/goom"
	"github.com/tinogoehlert/goom/level"
)

type doomLevel struct {
	name       string
	subSectors []*subSector
	mapRef     *level.Level
}

type subSector struct {
	floors   []*glWorldGeometry
	ceilings []*glWorldGeometry
	walls    []*glWorldGeometry
	sector   level.Sector
	ref      level.SubSector
}

func RegisterMap(m *level.Level, gd *goom.GameData, ts glTextureStore) *doomLevel {
	l := doomLevel{
		name:       m.Name,
		mapRef:     m,
		subSectors: make([]*subSector, 0, len(m.SubSectors("GL_SSECT"))),
	}

	for _, ssect := range m.SubSectors("GL_SSECT") {
		var s = &subSector{ref: ssect}
		s.addFlats(m, gd, ts)
		s.addWalls(m, gd, ts)
		l.subSectors = append(l.subSectors, s)
	}
	return &l
}

func (s *subSector) addFlats(md *level.Level, gd *goom.GameData, ts glTextureStore) {
	s.floors, s.ceilings = []*glWorldGeometry{}, []*glWorldGeometry{}
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
		fm := newGlWorldGeometry(floorData, sector.LightLevel(), ts[sector.FloorTexture()])
		s.floors = addGlWorldGeometry(s.floors, fm)
	}
	if len(gd.Flat(sector.CeilTexture())) > 0 {
		cm := newGlWorldGeometry(ceilData, md.Sectors[side.Sector].LightLevel(), ts[sector.CeilTexture()])
		s.ceilings = addGlWorldGeometry(s.ceilings, cm)
	}
}

func (s *subSector) addWalls(md *level.Level, gd *goom.GameData, ts glTextureStore) {
	s.walls = []*glWorldGeometry{}
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
				wm := newGlWorldGeometry(wallData, sector.LightLevel(), ts[side.Upper()])
				s.walls = addGlWorldGeometry(s.walls, wm)
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
				wm := newGlWorldGeometry(wallData, sector.LightLevel(), ts[side.Lower()])
				s.walls = addGlWorldGeometry(s.walls, wm)
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
			wm := newGlWorldGeometry(wallData, sector.LightLevel(), ts[side.Middle()])
			s.walls = addGlWorldGeometry(s.walls, wm)
		}
	}
}

func (s *subSector) Draw(ts glTextureStore) {
	for i := 0; i < len(s.floors); i++ {
		s.floors[i].Draw(gl.TRIANGLE_FAN)
		s.ceilings[i].Draw(gl.TRIANGLE_FAN)
	}
	for _, w := range s.walls {
		w.Draw(gl.TRIANGLES)
	}
}
