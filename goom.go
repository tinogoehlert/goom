package goom

import (
	"sync"

	"github.com/tinogoehlert/goom/audio"
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
	Music    audio.MusicSuite
	Fonts    graphics.FontBook
}

var (
	gdCache = make(map[string]*GameData)
	gdLock  = sync.RWMutex{}
)

// LoadWAD loads engine data from a wad and pwad file.
// Use this function to load a specific level.
func LoadWAD(iwad, pwad string) (*GameData, error) {
	var wads = []string{}
	wads = append(wads, iwad+".WAD", iwad+".gwa")
	if pwad != "" {
		wads = append(wads, pwad+".WAD", pwad+".gwa")
	}
	return LoadGameData(wads...)
}

// GetWAD returns cached GameData, loading WADs to cache if required.
func GetWAD(iwadfile, pwadfile string) (*GameData, error) {
	gdLock.Lock()
	defer gdLock.Unlock()
	key := iwadfile + "_" + pwadfile
	if gd, ok := gdCache[key]; ok {
		return gd, nil
	}
	wad, err := LoadWAD(iwadfile, pwadfile)
	if err == nil {
		gdCache[key] = wad
	}
	return wad, err
}

// LoadGameData loads engine data from WAD files.
func LoadGameData(files ...string) (*GameData, error) {
	gd := &GameData{
		Levels:   level.NewStore(),
		Textures: graphics.NewTextureStore(),
		Flats:    graphics.NewFlatStore(),
		Sprites:  graphics.NewSpriteStore(),
		Music:    audio.NewMusicSuite(),
		Fonts:    graphics.NewFontBook(),
	}
	for _, file := range files {
		wad, err := wad.NewWADFromFile(file)
		if err != nil {
			return nil, err
		}
		if err := gd.Levels.LoadWAD(wad); err != nil {
			return nil, err
		}
		if err := gd.Music.LoadWAD(wad); err != nil {
			return nil, err
		}
		if p, _ := graphics.NewPalettes(wad); p != nil {
			gd.Palettes = p
		}
		if err := gd.Fonts.LoadWAD(wad); err != nil {
			return nil, err
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
