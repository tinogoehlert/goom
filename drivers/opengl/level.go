package opengl

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/tinogoehlert/goom/goom"
	"github.com/tinogoehlert/goom/level"
)

type doomLevel struct {
	name       string
	subSectors []*subSector
	mapRef     *level.Level
	skyName    string
	sky        *glTexture
}

type subSector struct {
	floors   []*glWorldGeometry
	ceilings []*glWorldGeometry
	walls    []*glWorldGeometry
	sector   level.Sector
	ref      level.SubSector
}

func RegisterMap(m *level.Level, gd *goom.GameData, ts glTextureStore, skyName string) *doomLevel {
	l := doomLevel{
		name:       m.Name,
		mapRef:     m,
		subSectors: make([]*subSector, 0, len(m.SubSectors(level.GLNodesName))),
	}
	for _, ssect := range m.SubSectors(level.GLSsectsName) {
		var s = &subSector{ref: ssect}
		s.addFlats(m, gd, ts)
		s.addWalls(m, gd, ts)
		l.subSectors = append(l.subSectors, s)
	}

	l.sky = makeGLTexture(gd.Texture(skyName))

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
		fm := newGlWorldutils(floorData, sector.LightLevel(), ts[sector.FloorTexture()])
		s.floors = addGlWorldutils(s.floors, fm)
	}
	if len(gd.Flat(sector.CeilTexture())) > 0 {
		var (
			tex   = sector.CeilTexture()
			isSky = false
		)
		if sector.CeilTexture() == "F_SKY1" {
			isSky = true
		}
		cm := newGlWorldutils(ceilData, md.Sectors[side.Sector].LightLevel(), ts[tex])
		cm.isSky = isSky
		s.ceilings = addGlWorldutils(s.ceilings, cm)
	}
}

func (s *subSector) addWalls(md *level.Level, gd *goom.GameData, ts glTextureStore) {
	s.walls = []*glWorldGeometry{}
	for _, seg := range s.ref.Segments() {
		if seg.LineDef() == -1 {
			continue
		}
		line := md.LinesDefs[seg.LineDef()]
		side := &md.SideDefs[line.Right]

		otherSide := md.OtherSide(&line, seg)
		sector := md.Sectors[side.Sector]
		//isSky := false

		var (
			start  = md.Vert(uint32(line.Start))
			end    = md.Vert(uint32(line.End))
			upTex  = side.Upper()
			midTex = side.Middle()
			lowTex = side.Lower()
		)

		if side.Middle() == "-" &&
			side.Upper() == "-" &&
			side.Lower() == "-" {
			continue
		}

		var lside *level.SideDef

		if line.Left > 0 {
			lside = &md.SideDefs[line.Left]
			if side.Upper() == "-" {
				upTex = lside.Upper()
			}
			if side.Middle() == "-" {
				midTex = lside.Middle()
			}
			if side.Lower() == "-" {
				lowTex = lside.Lower()
			}
		}

		dist := start.DistanceTo(end)

		if upTex != "-" && otherSide != nil {
			oppositeSector := md.Sectors[otherSide.Sector]

			tex := ts[upTex][0]
			var (
				tw     = float32(tex.image.Width()) + float32(tex.image.Left())
				th     = float32(tex.image.Height()) + float32(tex.image.Top())
				height = oppositeSector.CeilHeight() - sector.CeilHeight()
				el     = dist / tw
			)

			wallData := []float32{
				-start.X(), sector.CeilHeight(), start.Y(), 0, height / th,
				-start.X(), oppositeSector.CeilHeight(), start.Y(), 0, 0.0,
				-end.X(), oppositeSector.CeilHeight(), end.Y(), el, 0.0,

				-end.X(), oppositeSector.CeilHeight(), end.Y(), el, 0.0,
				-end.X(), sector.CeilHeight(), end.Y(), el, height / th,
				-start.X(), sector.CeilHeight(), start.Y(), 0, height / th,
			}
			if gd.Texture(upTex) != nil {
				wm := newGlWorldutils(wallData, sector.LightLevel(), ts[upTex])
				if oppositeSector.CeilTexture() == "F_SKY1" {
					wm.isSky = true
				}
				s.walls = addGlWorldutils(s.walls, wm)
			}
		}

		if lowTex != "-" && otherSide != nil {
			oppositeSector := md.Sectors[otherSide.Sector]

			tex := ts[lowTex][0]
			var (
				tw     = float32(tex.image.Width()) + float32(tex.image.Left())
				th     = float32(tex.image.Height()) + float32(tex.image.Top())
				height = sector.FloorHeight() - oppositeSector.FloorHeight()
				el     = dist / tw
			)

			wallData := []float32{
				-start.X(), sector.FloorHeight(), start.Y(), 0, height / th,
				-start.X(), oppositeSector.FloorHeight(), start.Y(), 0, 0.0,
				-end.X(), oppositeSector.FloorHeight(), end.Y(), el, 0.0,

				-end.X(), oppositeSector.FloorHeight(), end.Y(), el, 0.0,
				-end.X(), sector.FloorHeight(), end.Y(), el, height / th,
				-start.X(), sector.FloorHeight(), start.Y(), 0, height / th,
			}
			if gd.Texture(lowTex) != nil {
				wm := newGlWorldutils(wallData, sector.LightLevel(), ts[lowTex])
				if lowTex == "F_SKY1" {
					wm.isSky = true
				}
				s.walls = addGlWorldutils(s.walls, wm)
			}
		}

		if midTex != "-" {
			tex := ts[midTex][0]
			var (
				tw     = float32(tex.image.Width()) + float32(tex.image.Left())
				th     = float32(tex.image.Height()) + float32(tex.image.Top())
				height = sector.CeilHeight() - sector.FloorHeight()
				el     = dist / tw
			)

			wallData := []float32{
				-start.X(), sector.FloorHeight(), start.Y(), 0, height / th,
				-start.X(), sector.CeilHeight(), start.Y(), 0, 0,
				-end.X(), sector.CeilHeight(), end.Y(), el, 0,

				-end.X(), sector.CeilHeight(), end.Y(), el, 0.0,
				-end.X(), sector.FloorHeight(), end.Y(), el, height / th,
				-start.X(), sector.FloorHeight(), start.Y(), 0, height / th,
			}
			wm := newGlWorldutils(wallData, sector.LightLevel(), ts[midTex])
			if midTex == "F_SKY1" {
				wm.isSky = true
			}
			s.walls = addGlWorldutils(s.walls, wm)
		}
	}
}

func (s *subSector) Draw(ts glTextureStore) {
	for i := 0; i < len(s.floors); i++ {
		s.floors[i].Draw(gl.TRIANGLE_FAN)
		if !s.ceilings[i].isSky {
			s.ceilings[i].Draw(gl.TRIANGLE_FAN)
		}
	}
	for _, w := range s.walls {
		if !w.isSky {
			w.Draw(gl.TRIANGLES)
		}
	}
}

func (s *subSector) DrawSky(ts glTextureStore, sky *glTexture) {
	for i := 0; i < len(s.floors); i++ {
		if s.ceilings[i].isSky {
			s.ceilings[i].DrawWithTexture(gl.TRIANGLE_FAN, sky)
		}
	}
	for _, w := range s.walls {
		if w.isSky {
			w.DrawWithTexture(gl.TRIANGLES, sky)
		}
	}
}
