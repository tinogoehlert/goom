package goom

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"regexp"
	"strings"
)

const (
	textureSize             = 22
	patchSize               = 10
	numPalettes             = 14
	transparentPaletteIndex = 255
)

var (
	patchStartRegex  = regexp.MustCompile(`^P?_START`)
	patchEndRegex    = regexp.MustCompile(`^P?_END`)
	flatStartRegex   = regexp.MustCompile(`^F?_START`)
	flatEndRegex     = regexp.MustCompile(`^F?_END`)
	spriteStartRegex = regexp.MustCompile(`^S?_START`)
	spriteEndRegex   = regexp.MustCompile(`^S?_END`)
)

type Patch struct {
	OriginX  int16
	OriginY  int16
	ImageID  int16
	StepDir  int16
	ColorMap int16
}

type Palette struct {
	Colors [256]color.RGBA
}

type DoomImage struct {
	Name    string
	WidthI  int16
	HeightI int16
	XOffset int16
	YOffset int16
	pixels  []byte
}

func (img *DoomImage) Height() float32 {
	return float32(img.HeightI)
}

func (img *DoomImage) Width() float32 {
	return float32(img.WidthI)
}

func (img *DoomImage) ToRGBA(palette *Palette) *image.RGBA {
	bounds := image.Rect(0, 0, int(img.WidthI), int(img.HeightI))
	tex := image.NewRGBA(bounds)
	for y := 0; y < int(img.HeightI); y++ {
		for x := 0; x < int(img.WidthI); x++ {
			pixel := img.pixels[y*int(img.WidthI)+x]
			var alpha uint8
			if pixel == transparentPaletteIndex {
				alpha = 0
			} else {
				alpha = 255
			}
			rgb := palette.Colors[pixel]
			tex.Set(x, y, color.RGBA{rgb.R, rgb.G, rgb.B, alpha})
		}
	}
	return tex
}

type DoomTex interface {
	ToRGBA(palette *Palette) *image.RGBA
	Width() float32
	Height() float32
}

// Flat DOOM floor and ceiling images, always 64x64 (4096 bytes)
type Flat struct {
	pixels []byte
}

func (f *Flat) Width() float32  { return float32(64) }
func (f *Flat) Height() float32 { return float32(64) }

func (di *DoomImage) FromBuff(buff []byte) error {
	di.WidthI = int16(binary.LittleEndian.Uint16(buff[0:2]))
	di.HeightI = int16(binary.LittleEndian.Uint16(buff[2:4]))
	di.XOffset = int16(binary.LittleEndian.Uint16(buff[4:6]))
	di.YOffset = int16(binary.LittleEndian.Uint16(buff[6:8]))
	offsets := make([]int32, di.WidthI, di.WidthI)
	r := bytes.NewBuffer(buff[8 : 8+(di.WidthI*4)])

	if err := binary.Read(r, binary.LittleEndian, offsets); err != nil {
		return err
	}

	size := int(di.WidthI) * int(di.HeightI)
	di.pixels = make([]byte, size, size)
	for y := 0; y < int(di.HeightI); y++ {
		for x := 0; x < int(di.WidthI); x++ {
			di.pixels[y*int(di.WidthI)+x] = transparentPaletteIndex
		}
	}
	for columnIndex, offset := range offsets {
		for {
			rowStart := buff[offset]
			offset += 1
			if rowStart == 255 {
				break
			}
			numPixels := buff[offset]
			offset += 1
			offset += 1 /* Padding */
			for i := 0; i < int(numPixels); i++ {
				pixelOffset := (int(rowStart)+i)*int(di.WidthI) + columnIndex
				di.pixels[pixelOffset] = buff[offset]
				offset += 1
			}
			offset += 1 /* Padding */
		}
	}
	return nil
}

// Texture - DOOM image data
type Texture struct {
	Name       string
	IsMasked   bool
	WidthI     int16
	HeightI    int16
	patchCount int16
	patches    []Patch
	image      *image.RGBA
}

func newTexture(tbuff []byte) (*Texture, error) {
	tex := &Texture{
		Name:       strings.TrimRight(string(tbuff[0:8]), "\x00"),
		IsMasked:   !(binary.LittleEndian.Uint32(tbuff[8:12]) == 0),
		WidthI:     int16(binary.LittleEndian.Uint16(tbuff[12:14])),
		HeightI:    int16(binary.LittleEndian.Uint16(tbuff[14:16])),
		patchCount: int16(binary.LittleEndian.Uint16(tbuff[20:22])),
	}
	tex.patches = make([]Patch, tex.patchCount)
	r := bytes.NewBuffer(tbuff[22 : 22+(patchSize*tex.patchCount)])
	for i := 0; i < int(tex.patchCount); i++ {
		var patch = &Patch{}
		if err := binary.Read(r, binary.LittleEndian, patch); err != nil {
			fmt.Println("could not read patch for", tex.Name)
			return nil, err
		}
		tex.patches[i] = *patch
	}
	return tex, nil
}

func (t *Texture) ToRGBA(palette *Palette) *image.RGBA {
	if t == nil {
		return nil
	}
	return t.image
}

func (t *Texture) Width() float32  { return float32(t.WidthI) }
func (t *Texture) Height() float32 { return float32(t.HeightI) }

type Graphics struct {
	textures map[string]*Texture
	images   map[string]DoomImage
	flats    map[string]Flat
	sprites  map[string]Sprite
	palettes []Palette
}

func (g *Graphics) Palette(num int) *Palette {
	if g == nil {
		return nil
	}
	return &g.palettes[0]
}

func (wm *WadManager) LoadGraphics() (*Graphics, error) {
	gfx := Graphics{
		textures: make(map[string]*Texture),
		images:   make(map[string]DoomImage),
		flats:    make(map[string]Flat),
		sprites:  make(map[string]Sprite),
		palettes: []Palette{},
	}
	for _, wad := range wm.wads[:1] {
		var (
			lumps  = wad.GetLumps()
			pnames = make(map[int16]string)
			err    error
		)
		for i := 0; i < wad.NumLumps; i++ {
			lump := &lumps[i]

			switch {
			case lump.Name == "PLAYPAL":
				if err := gfx.loadPalettesFromBuff(lump.Data); err != nil {
					fmt.Println("could not load palette", err)
				}
			case lump.Name == "PNAMES":
				numPnames := int(binary.LittleEndian.Uint32(lump.Data[0:4]))
				data := lump.Data[4:]
				for i := 0; i < numPnames; i++ {
					pname := strings.ToUpper(wadString(data[8*i : (8*i)+8]))
					pnames[int16(i)] = pname
				}
			case lump.Name == "TEXTURE1" || lump.Name == "TEXTURE2":
				err := gfx.loadTexturesFromBuff(lump.Data)
				if err != nil {
					fmt.Println("could not load texture", lump.Name, err)
				}
			case patchStartRegex.Match([]byte(lump.Name)):
				for {
					lump := &lumps[i]
					if lump.Size > 0 {
						di := DoomImage{Name: lump.Name}
						err := di.FromBuff(lump.Data)
						if err != nil {
							fmt.Println("could not load patch", err)
						}
						gfx.images[lump.Name] = di
					}
					if patchEndRegex.Match([]byte(lump.Name)) {
						break
					}
					i++
				}
			case flatStartRegex.Match([]byte(lump.Name)):
				for {
					lump := &lumps[i]
					if lump.Size > 0 {
						flat := Flat{pixels: make([]byte, lump.Size)}
						copy(flat.pixels, lump.Data)
						gfx.flats[lump.Name] = flat
					}
					if flatEndRegex.Match([]byte(lump.Name)) {
						break
					}
					i++
				}
			case spriteStartRegex.Match([]byte(lump.Name)):
				for {
					lump := &lumps[i]
					if lump.Size > 0 {
						s, ok := gfx.sprites[lump.Name[:4]]
						if !ok {
							s = NewSprite(lump.Name[:4])
							s.first = lump.Name[:6]
							gfx.sprites[lump.Name[:4]] = s
						}
						sf := s.AddSpriteFrame(lump)
						sf.image.FromBuff(lump.Data)
						if len(lump.Name) == 8 {
							s.frames[lump.Name[:4]+lump.Name[6:8]] = sf
						}
					}
					if spriteEndRegex.Match([]byte(lump.Name)) {
						break
					}
					i++
				}
			}
		}

		for _, tex := range gfx.textures {
			bounds := image.Rect(0, 0, int(tex.WidthI), int(tex.HeightI))
			tex.image = image.NewRGBA(bounds)
			if tex.image.Stride != tex.image.Rect.Size().X*4 {
				return nil, fmt.Errorf("unsupported stride")
			}

			for _, patch := range tex.patches {
				image := gfx.images[pnames[patch.ImageID]]
				if err != nil {
					return nil, err
				}
				for y := 0; y < int(image.HeightI); y++ {
					for x := 0; x < int(image.WidthI); x++ {
						pixel := image.pixels[y*int(image.WidthI)+x]
						var alpha uint8
						if pixel == transparentPaletteIndex {
							alpha = 0
						} else {
							alpha = 255
						}
						rgb := gfx.palettes[0].Colors[pixel]
						tex.image.Set(int(patch.OriginX)+x, int(patch.OriginY)+y, color.RGBA{rgb.R, rgb.G, rgb.B, alpha})
					}
				}
			}
		}
	}
	return &gfx, nil
}

func (gfx *Graphics) loadPalettesFromBuff(data []byte) error {
	for i := 0; i < numPalettes; i++ {
		p := &Palette{}
		for ci := 0; ci < 256*3; ci += 3 {
			p.Colors[ci/3].R = data[ci]
			p.Colors[ci/3].G = data[ci+1]
			p.Colors[ci/3].B = data[ci+2]
			p.Colors[ci/3].A = 255
		}
		gfx.palettes = append(gfx.palettes, *p)
	}
	return nil
}

func (gfx *Graphics) loadTexturesFromBuff(data []byte) error {
	texCount := int(binary.LittleEndian.Uint32(data[0:4]))
	for i := 0; i < texCount; i++ {
		offset := int(binary.LittleEndian.Uint32(data[4*i : (4*i)+4]))
		tbuff := data[offset : offset+textureSize]
		tex, err := newTexture(tbuff)
		if err != nil {
			continue
		}
		gfx.textures[tex.Name] = tex
	}
	return nil
}

func (gfx *Graphics) GetTexture(name string) DoomTex {
	if gfx == nil {
		return nil
	}
	return gfx.textures[name]
}

func (gfx *Graphics) GetFlat(name string) DoomTex {
	if gfx == nil {
		return nil
	}
	flat := gfx.flats[name]
	return &flat
}

func (gfx *Graphics) GetSprites() map[string]Sprite {
	if gfx == nil {
		return nil
	}
	return gfx.sprites
}

func (f *Flat) ToRGBA(palette *Palette) *image.RGBA {
	if f == nil {
		return nil
	}
	rgba := image.NewRGBA(image.Rect(0, 0, 64, 64))

	for x := 0; x < 64; x++ {
		for y := 0; y < 64; y++ {
			ci := f.pixels[f.pixels[y+(64*x)]]
			rgba.Set(x, y, palette.Colors[ci])
		}
	}
	return rgba
}
