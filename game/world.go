package game

import (
	"container/list"
	"errors"
	"fmt"
	"math"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/tinogoehlert/goom/drivers"
	"github.com/tinogoehlert/goom/goom"
	"github.com/tinogoehlert/goom/level"
	"github.com/tinogoehlert/goom/utils"
)

//World  holds the current world
type World struct {
	nodes       []level.Node
	things      []Thingable
	items       []*Item
	monsters    []*Monster
	players     []*Player
	definitions *DefStore
	projectiles *list.List
	me          *Player
	levelRef    *level.Level
	audioDriver drivers.Audio
	gameData    *goom.GameData
}

// NewWorld Creates a new world
func NewWorld(data *goom.GameData, defs *DefStore) *World {
	noopAudio, _ := drivers.AudioDrivers[drivers.NoopAudio](nil, "")

	return &World{
		definitions: defs,
		gameData:    data,
		audioDriver: noopAudio,
	}
}

func (w *World) Data() *goom.GameData {
	return w.gameData
}

// SetAudioDriver sets the audioDriver
func (w *World) SetAudioDriver(drv drivers.Audio) {
	w.audioDriver = drv
}

// AudioDriver returns the current audio driver.
func (w *World) AudioDriver() drivers.Audio {
	return w.audioDriver
}

// LoadLevel a specific level of the world
func (w *World) LoadLevel(lvl *level.Level) error {
	w.nodes = lvl.Nodes(level.GLNodesName)
	w.levelRef = lvl
	w.projectiles = list.New()

	for _, t := range w.levelRef.Things {
		if t.Type < 5 {
			player := NewPlayer(t.X, t.Y, 0, t.Angle, w)
			w.players = append(w.players, player)
			if t.Type == 1 {
				w.me = player
			}
			player.AddWeapon(w.definitions.GetWeapon("pistol"))
			player.SetCollision(w.doesCollide)
		}
		if obstacleDef := w.definitions.GetObstacleDef(int(t.Type)); obstacleDef != nil {
			obstacle := ThingFromDef(t.X, t.Y, 0, t.Angle, obstacleDef)
			w.things = appendDoomThing(w.things, obstacle, w.levelRef)
		}

		if itemDef := w.definitions.GetItemDef(int(t.Type)); itemDef != nil {
			item := ItemFromDef(t.X, t.Y, 0, t.Angle, itemDef)
			w.things = appendDoomThing(w.things, item, w.levelRef)
		}

		if obstacleDef := w.definitions.GetObstacleDef(int(t.Type)); obstacleDef != nil {
			obstacle := ThingFromDef(t.X, t.Y, 0, t.Angle, obstacleDef)
			w.things = appendDoomThing(w.things, obstacle, w.levelRef)
		}

		if monsterDef := w.definitions.GetMonsterDef(int(t.Type)); monsterDef != nil {
			sprite := w.gameData.Sprite(monsterDef.Sprite)
			img := sprite.FirstFrame().Angles()[1]

			monster := MonsterFromDef(
				t.X,
				t.Y,
				float32(img.Width())/2,
				float32(img.Width())/2,
				0,
				t.Angle,
				monsterDef,
			)
			w.monsters = append(w.monsters, monster)
			w.things = appendDoomThing(w.things, monster, w.levelRef)
		}
	}

	mus := w.levelRef.Name
	if w.levelRef.Name == "MAP01" {
		mus = "RUNNIN"
	}

	if err := w.audioDriver.PlayMusic(w.gameData.Music.Track(mus)); err != nil {
		fmt.Println("could not play music:", err.Error())
	}

	return nil
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
		return errors.New("out of range")
	}
	w.me = w.players[num]
	return nil
}

// Update the world (monster, thing and player position)
func (w *World) Update() {
	ppos := utils.V2(w.me.position[0], w.me.position[1])
	for _, m := range w.monsters {
		if !m.IsCorpse() {
			m.Update()
			m.Think(w.me)
		}
	}
	for e := w.projectiles.Front(); e != nil; e = e.Next() {
		var (
			p       = e.Value.(*Projectile)
			projPos = utils.V2(p.position[0], p.position[1])
		)

		for _, m := range w.monsters {
			if m.IsCorpse() {
				continue
			}
			if w.hitThing(m, p, m.sizeX, m.sizeY) {
				mpos := utils.V2(m.position[0], m.position[1])
				dist := ppos.DistanceTo(mpos)
				w.projectiles.Remove(e)
				state := m.Hit(p.damage, dist)
				id := m.sounds[state]
				sound := w.gameData.Sounds.GetByID(id)
				if sound == nil {
					fmt.Printf("bad monster state sound: %s = %d\n", id, int(state))
					break
				}
				if state > 0 {
					distP := mgl32.Vec2(m.position).Sub(w.me.position)
					angle := mgl32.RadToDeg(
						float32(math.Atan2(float64(distP.Y()),
							float64(distP.X()))),
					) - m.angle
					if angle < 0.0 {
						angle += 360
					}
					w.audioDriver.PlayAtPosition(sound.Name, dist/2.6, int16(angle))
				}
				break
			}
		}
		if int(ppos.DistanceTo(projPos)) > p.maxRange {
			w.projectiles.Remove(e)
		}
		p.Walk(20)
	}
	w.me.Update()
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
	w.audioDriver.Play("DS" + player.weapon.Sound)
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
					w.audioDriver.Play("DSWPNUP")
				}
			}
		}
	}

}

func (w *World) checkWallCollision(thing *DoomThing, to mgl32.Vec2) mgl32.Vec2 {
	var (
		collided = 0
		x        = to.X()
		y        = to.Y()
		radius   = float32(24)
		hitWall  level.Wall
		oldTo    = to
	)
	for _, wall := range w.levelRef.Walls {
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
				collided++
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
					collided++
				}
			}
		}
	}

	if collided > 0 {
		if hitWall.IsTwoSided {
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

// GetLevel return the currently loaded level.
func (w *World) GetLevel() *level.Level {
	return w.levelRef
}
