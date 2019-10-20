package game

import (
	"fmt"
	"math"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/tinogoehlert/goom/geometry"
	"github.com/tinogoehlert/goom/level"
)

// Wall a DOOM Wall
type Wall struct {
	Start   geometry.Vec2
	End     geometry.Vec2
	Normal  geometry.Vec2
	Tangent geometry.Vec2
	Length  float32
	Sides   struct {
		Right *level.SideDef
		Left  *level.SideDef
	}
	lineDef level.LineDef
}

//World  holds the current world
type World struct {
	nodes    []level.Node
	things   []*DoomThing
	monsters []*Monster
	players  []*Player
	walls    []Wall
	me       *Player
	levelRef *level.Level
}

func newWall(line level.LineDef, lvl *level.Level) Wall {
	var w = Wall{
		lineDef: line,
		Start:   lvl.Vert(uint32(line.Start)),
		End:     lvl.Vert(uint32(line.End)),
		Sides: struct {
			Right *level.SideDef
			Left  *level.SideDef
		}{
			Right: &lvl.SideDefs[line.Right],
		},
	}

	w.Length = w.Start.DistanceTo(w.End)
	w.Tangent = w.End.Sub(w.Start).Normalize()
	w.Normal = w.Tangent.CrossVec2()

	if line.Left != -1 {
		w.Sides.Left = &lvl.SideDefs[line.Left]
	}
	return w
}

// NewWorld creates a new world
func NewWorld(doomLevel *level.Level, defs *DefStore) *World {
	var w = &World{
		nodes:    doomLevel.Nodes(level.GLNodesName),
		walls:    make([]Wall, 0, len(doomLevel.LinesDefs)),
		levelRef: doomLevel,
	}

	for _, line := range doomLevel.LinesDefs {
		w.walls = append(w.walls, newWall(line, doomLevel))
	}

	for _, t := range doomLevel.Things {
		if t.Type < 5 {
			player := NewPlayer(t.X, t.Y, 0, t.Angle)
			w.players = append(w.players, player)
			if t.Type == 1 {
				w.me = player
			}
			player.AddWeapon(defs.GetWeapon("pistol"))
			player.SetCollision(w.doesCollide)
		}
		if obstacleDef := defs.GetObstacleDef(int(t.Type)); obstacleDef != nil {
			obstacle := ThingFromDef(t.X, t.Y, 0, t.Angle, obstacleDef)
			w.things = appendDoomThing(w.things, obstacle, doomLevel)
		}

		if monsterDef := defs.GetMonsterDef(int(t.Type)); monsterDef != nil {
			monster := MonsterFromDef(t.X, t.Y, 0, t.Angle, monsterDef)
			w.monsters = append(w.monsters, monster)
			w.things = appendDoomThing(w.things, monster.DoomThing, doomLevel)
		}
	}
	return w
}

// Me returns current player
func (w *World) Me() *Player {
	return w.me
}

// Things returns things
func (w *World) Things() []*DoomThing {
	return w.things
}

// Monsters returns monsters
func (w *World) Monsters() []*Monster {
	return w.monsters
}

// SetPlayer (in case we will implement multiplayer LOL)
func (w *World) SetPlayer(num int) error {
	if num > 4 {
		return fmt.Errorf("out of range")
	}
	w.me = w.players[num]
	return nil
}

// Update updates the world (monster think and player position)
func (w *World) Update() {

}

func (w *World) doesCollide(thing *DoomThing, to mgl32.Vec2) mgl32.Vec2 {
	var (
		collided = false
		x        = thing.position[0]
		y        = thing.position[1]
		xNudge   = float32(0.0)
		yNudge   = float32(0.0)
		radius   = float32(23)
		hitWall  *Wall
	)
	for _, w := range w.walls {
		var (
			d   = w.Start.Dot(w.Normal)
			sd  = w.Start.Dot(w.Tangent)
			pd  = x*w.Normal.X() + y*w.Normal.Y() - d
			mul = float32(1.0)
		)
		hitWall = &w
		if pd >= -radius && pd <= radius {
			if pd < 0 {
				pd = -pd
				mul = -1.0
			}
			psd := x*w.Tangent.X() + y*w.Tangent.Y() - sd
			if psd >= 0.0 && psd <= w.Length {
				toPushOut := radius - pd + 0.001
				xNudge = w.Normal.X() * toPushOut * mul
				yNudge = w.Normal.Y() * toPushOut * mul
				collided = true
				break
			} else {
				var (
					tmpxd float32
					tmpyd float32
				)
				tmpxd = x - w.Start.X()
				tmpyd = y - w.Start.Y()
				if psd > 0.0 {
					tmpxd = x - w.End.X()
					tmpyd = y - w.End.Y()
				}

				distSqr := tmpxd*tmpxd + tmpyd*tmpyd
				if distSqr < radius*radius {
					dist := float32(math.Sqrt(float64(distSqr)))
					toPushOut := radius - dist + 0.001
					xNudge = tmpxd / dist * toPushOut
					yNudge = tmpyd / dist * toPushOut
					collided = true
					break
				}
			}
		}
	}

	if collided {
		if hitWall.lineDef.Left != -1 {
			var (
				lSector   = w.levelRef.Sectors[hitWall.Sides.Left.Sector]
				rSector   = w.levelRef.Sectors[hitWall.Sides.Left.Sector]
				chkHeight = thing.currentSector.FloorHeight() + 32
			)
			if rSector.FloorHeight() < chkHeight || lSector.FloorHeight() < chkHeight {
				return to
			}
		}
		return mgl32.Vec2{
			thing.position[0] + xNudge,
			thing.position[1] + yNudge,
		}
	}
	return to
}
