package graphics

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"regexp"
	"strings"

	"github.com/tinogoehlert/goom/internal/utils"
	"github.com/tinogoehlert/goom/wad"
)

const (
	textureSize = 22
)

// Texture - DOOM texture definition
type Texture struct {
	name       string
	isMasked   bool
	width      int
	height     int
	patchCount int
	patches    []*Patch
}

var pnameStore = []string{}

var picStore map[string]*DoomPicture

// NewTexture create a new DOOM texture
func NewTexture(buff []byte) (*Texture, error) {
	tex := &Texture{
		name:       strings.TrimRight(string(buff[0:8]), "\x00"),
		isMasked:   !(binary.LittleEndian.Uint32(buff[8:12]) == 0),
		width:      int(int16(binary.LittleEndian.Uint16(buff[12:14]))),
		height:     int(int16(binary.LittleEndian.Uint16(buff[14:16]))),
		patchCount: int(int16(binary.LittleEndian.Uint16(buff[20:22]))),
	}
	fmt.Println("new tex:", tex.name)
	tex.patches = make([]*Patch, tex.patchCount)
	pbuff := buff[22:]
	for i := 0; i < int(tex.patchCount); i++ {
		tex.patches[i] = NewPatch(pbuff[i*patchSize : (i*patchSize)+patchSize])
	}
	return tex, nil
}

// Width return width of image
func (t *Texture) Width() int { return int(t.width) }

// Height return height of image
func (t *Texture) Height() int { return int(t.height) }

// ToRGBA generates a go image from all patches
func (t *Texture) ToRGBA(palette [256]color.RGBA) *image.RGBA {
	var (
		bounds = image.Rect(0, 0, int(t.width), int(t.height))
		img    = image.NewRGBA(bounds)
	)
	for _, patch := range t.patches {
		if patch.DoomPicture == nil {
			continue
		}
		for y := 0; y < patch.height; y++ {
			for x := 0; x < patch.width; x++ {
				pixel := patch.data[y*int(patch.width)+x]
				var alpha uint8
				if pixel == transparentColor {
					alpha = 0
				} else {
					alpha = 255
				}
				rgb := palette[pixel]
				img.Set(int(patch.originX)+x, int(patch.originY)+y, color.RGBA{rgb.R, rgb.G, rgb.B, alpha})
			}
		}
	}
	return img
}

// ToPng exports picture to PNG
func (t *Texture) ToPng(out string, palette [256]color.RGBA) error {
	img := t.ToRGBA(palette)
	f, _ := os.Create(out)
	return png.Encode(f, img)
}

// TextureStore string map of textures
type TextureStore map[string]*Texture

// NewTextureStore creates new texture store
func NewTextureStore() TextureStore {
	return make(TextureStore)
}

func (ts TextureStore) LoadWAD(w *wad.WAD) {
	var (
		lumps           = w.Lumps()
		patchStartRegex = regexp.MustCompile(`^P?_START`)
		patchEndRegex   = regexp.MustCompile(`^P?_END`)
	)

	if picStore == nil {
		picStore = make(map[string]*DoomPicture)
	}

	for i := 0; i < len(lumps); i++ {
		var lump = &lumps[i]
		switch {
		case lump.Name == "PNAMES":
			loadPNAMES(lump)
		case lump.Name == "TEXTURE1" || lump.Name == "TEXTURE2":
			ts.loadTextures(lump)
		case patchStartRegex.Match([]byte(lump.Name)):
			for {
				lump := &lumps[i]
				if lump.Size > 0 {
					picStore[lump.Name] = NewDoomPicture(lump.Data)
				}
				if patchEndRegex.Match([]byte(lump.Name)) {
					break
				}
				i++
			}
		}
	}
}

func loadPNAMES(lump *wad.Lump) {
	numPnames := int(binary.LittleEndian.Uint32(lump.Data[0:4]))
	data := lump.Data[4:]
	for i := 0; i < numPnames; i++ {
		pname := strings.ToUpper(utils.WadString(data[8*i : (8*i)+8]))
		pnameStore = append(pnameStore, pname)
	}
}

func (ts TextureStore) loadTextures(lump *wad.Lump) {
	var (
		texCount = int(binary.LittleEndian.Uint32(lump.Data[0:4]))
		offsets  = make([]int, texCount)
	)

	odata := lump.Data[4:]
	for i := 0; i < texCount; i++ {
		offsets[i] = int(binary.LittleEndian.Uint32(odata[4*i : (i*4)+4]))
	}

	for _, offset := range offsets {
		tbuff := lump.Data[offset : offset+textureSize]
		tex, err := NewTexture(tbuff)
		if err != nil {
			fmt.Println(err)
			continue
		}
		ts[tex.name] = tex
	}
}

func (ts TextureStore) InitPatches() {
	for _, t := range ts {
		for _, patch := range t.patches {
			pic, ok := picStore[pnameStore[patch.pictureID]]
			if ok {
				patch.DoomPicture = pic
			}
		}
	}
}
