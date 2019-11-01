package game

import (
	"container/list"
	"fmt"
	"math"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/tinogoehlert/goom"
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
	nodes       []level.Node
	things      []Thingable
	items       []*Item
	monsters    []*Monster
	players     []*Player
	definitions *DefStore
	walls       []Wall
	projectiles *list.List
	me          *Player
	levelRef    *level.Level
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
func NewWorld(doomLevel *level.Level, defs *DefStore, data *goom.GameData) *World {
	var w = &World{
		nodes:       doomLevel.Nodes(level.GLNodesName),
		walls:       make([]Wall, 0, len(doomLevel.LinesDefs)),
		levelRef:    doomLevel,
		definitions: defs,
		projectiles: list.New(),
	}

	for _, line := range doomLevel.LinesDefs {
		w.walls = append(w.walls, newWall(line, doomLevel))
	}

	for _, t := range doomLevel.Things {
		if t.Type < 5 {
			player := NewPlayer(t.X, t.Y, 0, t.Angle, w)
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

		if itemDef := defs.GetItemDef(int(t.Type)); itemDef != nil {
			item := ItemFromDef(t.X, t.Y, 0, t.Angle, itemDef)
			w.things = appendDoomThing(w.things, item, doomLevel)
		}

		if obstacleDef := defs.GetObstacleDef(int(t.Type)); obstacleDef != nil {
			obstacle := ThingFromDef(t.X, t.Y, 0, t.Angle, obstacleDef)
			w.things = appendDoomThing(w.things, obstacle, doomLevel)
		}

		if monsterDef := defs.GetMonsterDef(int(t.Type)); monsterDef != nil {
			sprite := data.Sprite(monsterDef.Sprite)
			img := sprite.FirstFrame().Angles()[1]
			fmt.Println(img.Width(), img.Height())

			monster := MonsterFromDef(
				t.X,
				t.Y,
				float32(img.Width()),
				float32(img.Height()),
				0,
				t.Angle,
				monsterDef,
			)
			w.monsters = append(w.monsters, monster)
			w.things = appendDoomThing(w.things, monster, doomLevel)
		}
	}
	return w
}

// Me returns current player
func (w *World) Me() *Player {
	return w.me
}

// Things returns things
func (w *World) Things() []Thingable {
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
func (w *World) Update(t float32) {
	ppos := geometry.V2(w.me.position[0], w.me.position[1])
	for _, m := range w.monsters {
		if !m.IsCorpse() {
			m.Update()
			m.Think(w.me, t)
		}
	}

	for e := w.projectiles.Front(); e != nil; e = e.Next() {
		var (
			p       = e.Value.(*Projectile)
			projPos = geometry.V2(p.position[0], p.position[1])
		)
		for _, m := range w.monsters {
			if m.IsCorpse() {
				continue
			}
			if w.hitThing(m, p, m.sizeX, m.sizeY) {
				mpos := geometry.V2(m.position[0], m.position[1])
				dist := ppos.DistanceTo(mpos)
				w.projectiles.Remove(e)
				m.Hit(p.damage, dist)
				break
			}
		}
		if int(ppos.DistanceTo(projPos)) > p.maxRange {
			w.projectiles.Remove(e)
		}
		p.Walk(1000, t)
	}
}

func (w *World) doesCollide(thing *DoomThing, to mgl32.Vec2) mgl32.Vec2 {
	w.checkThingCollision(thing, to)
	return w.checkWallCollision(thing, to)
}

func (w *World) spawnShot(player *Player) {
	w.projectiles.PushBack(NewProjectile(
		player.Position(),
		player.Direction(),
		player.weapon.Damage,
		player.weapon.Range,
	))
}

func (w *World) hitThing(t1, t2 Thingable, sx, sy float32) bool {
	var (
		x  = t1.Position()[0]
		y  = t1.Position()[1]
		x1 = t2.Position()[0] - (sx)
		x2 = t2.Position()[0] + (sx)
		y1 = t2.Position()[1] - (sy)
		y2 = t2.Position()[1] + (sy)
	)

	if x > x1 && x < x2 && y > y1 && y < y2 {
		return true
	}
	return false
}

func (w *World) checkThingCollision(thing *DoomThing, to mgl32.Vec2) {
	for _, thing2 := range w.things {
		if !thing2.IsShown() {
			continue
		}

		if w.hitThing(thing, thing2, 24, 24) {
			switch t := thing2.(type) {
			case *Item:
				if t.category == "weapon" {
					w.me.AddWeapon(w.definitions.GetWeapon(t.ref))
					t.consumed = true
				}
			}
		}
	}

}

func (w *World) checkWallCollision(thing *DoomThing, to mgl32.Vec2) mgl32.Vec2 {
	var (
		collided = false
		x        = to.X()
		y        = to.Y()
		radius   = float32(24)
		hitWall  Wall
		oldTo    = to
	)
	for _, wall := range w.walls {
		var (
			d   = wall.Start.Dot(wall.Normal)
			sd  = wall.Start.Dot(wall.Tangent)
			pd  = x*wall.Normal.X() + y*wall.Normal.Y() - d
			mul = float32(1.0)
		)
		if pd >= -radius && pd <= radius {
			if pd < 0 {
				pd = -pd
				mul = -1.0
			}
			psd := x*wall.Tangent.X() + y*wall.Tangent.Y() - sd
			if psd >= 0.0 && psd <= wall.Length {
				toPushOut := radius - pd + 0.001
				to[0] += wall.Normal.X() * toPushOut * mul
				to[1] += wall.Normal.Y() * toPushOut * mul
				hitWall = wall
				collided = true

			} else {
				var (
					tmpxd float32
					tmpyd float32
				)
				tmpxd = x - wall.Start.X()
				tmpyd = y - wall.Start.Y()
				if psd > 0.0 {
					tmpxd = x - wall.End.X()
					tmpyd = y - wall.End.Y()
				}

				distSqr := tmpxd*tmpxd + tmpyd*tmpyd
				if distSqr < radius*radius {
					dist := float32(math.Sqrt(float64(distSqr)))
					toPushOut := radius - dist + 0.001
					to[0] += tmpxd / dist * toPushOut
					to[1] += tmpyd / dist * toPushOut
					hitWall = wall
					collided = true
				}
			}
		}
	}

	if collided {
		if hitWall.lineDef.Left != -1 {
			var (
				lSector   = w.levelRef.Sectors[hitWall.Sides.Left.Sector]
				chkHeight = thing.currentSector.FloorHeight() + 32
			)
			if lSector.FloorHeight() < chkHeight {
				return oldTo
			}
		}
	}
	return to
}
