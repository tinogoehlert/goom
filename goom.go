package goom

import (
	"github.com/tinogoehlert/goom/graphics"
	"github.com/tinogoehlert/goom/level"
	"github.com/tinogoehlert/goom/wad"
)

//GameData Game Data
type GameData struct {
	Levels   level.Store
	Textures graphics.TextureStore
	Flats    graphics.FlatStore
	Sprites  graphics.SpriteStore
	Palettes *graphics.Palettes
}

// LoadGameData Load Engine data from WAD files
func LoadGameData(files ...string) (*GameData, error) {
	gd := &GameData{
		Levels:   level.NewStore(),
		Textures: graphics.NewTextureStore(),
		Flats:    graphics.NewFlatStore(),
		Sprites:  graphics.NewSpriteStore(),
	}
	for _, file := range files {
		wad, err := wad.NewWADFromFile(file)
		if err != nil {
			return nil, err
		}
		if err := gd.Levels.LoadWAD(wad); err != nil {
			return nil, err
		}
		if p, _ := graphics.NewPalettes(wad); p != nil {
			gd.Palettes = p
		}
		gd.Sprites.LoadWAD(wad)
		gd.Flats.LoadWAD(wad)
		gd.Textures.LoadWAD(wad)

	}
	gd.Textures.InitPatches()
	return gd, nil
}

// Level return level by name
func (gd *GameData) Level(name string) *level.Level {
	return gd.Levels[name]
}

// Texture return texture by name
func (gd *GameData) Texture(name string) *graphics.Texture {
	return gd.Textures[name]
}

// Flat return flat(s) by name
func (gd *GameData) Flat(name string) []*graphics.Flat {
	return gd.Flats[name]
}

// Sprite return sprite by name
func (gd *GameData) Sprite(name string) graphics.Sprite {
	return gd.Sprites[name]
}

// DefaultPalette gets the default game palette
func (gd *GameData) DefaultPalette() graphics.Palette {
	return gd.Palettes[0]
}
