package level

import (
	"strings"

	"github.com/tinogoehlert/goom/utils"
)

// Wall a DOOM Wall
type Wall struct {
	Start      utils.Vec2
	End        utils.Vec2
	Normal     utils.Vec2
	Tangent    utils.Vec2
	IsTwoSided bool
	IsSky      bool
	Length     float32
	Sectors    struct {
		Right *Sector
		Left  *Sector
	}
	Sides struct {
		Right *SideDef
		Left  *SideDef
	}
	lineDef *LineDef
}

// NewWall creates a new wall from linedef and level
func NewWall(line *LineDef, lvl *Level) Wall {
	var w = Wall{
		lineDef: line,
		Start:   lvl.Vert(uint32(line.Start)),
		End:     lvl.Vert(uint32(line.End)),
		Sides: struct {
			Right *SideDef
			Left  *SideDef
		}{
			Right: &lvl.SideDefs[line.Right],
		},
	}

	w.Sectors.Right = &lvl.Sectors[w.Sides.Right.Sector]
	w.Length = w.Start.DistanceTo(w.End)
	w.Tangent = w.End.Sub(w.Start).Normalize()
	w.Normal = w.Tangent.CrossVec2()

	if line.Left != -1 {
		w.IsTwoSided = true
		w.Sides.Left = &lvl.SideDefs[line.Left]
		w.Sectors.Left = &lvl.Sectors[w.Sides.Left.Sector]
	}

	if strings.Contains(w.Sectors.Right.CeilTexture(), "SKY") {
		w.IsSky = true
	}

	return w
}
