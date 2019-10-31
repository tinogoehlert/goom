package graphics

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/tinogoehlert/goom/wad"
)

// FontName helps identifing the several embedded fonts in a WAD file
type FontName int

const (
	// font names
	FnNumGreySmall FontName = iota
	FnNumYellowSmall
	FnNumRedBig
	FnCompositeRed

	// WAD identifiers
	ptNumGreySmall   = `^STGNUM(\d)$`
	ptNumYellowSmall = `^STYSNUM(\d)$`
	ptNumRedBig      = `^STTNUM(\d)$`
	ptCompositeRed   = `^STCFN(\d{3})$`

	// Extras
	exNumRedBigMinus   = `STTMINUS`
	exNumRedBigPercent = `STTPRCNT`

	// monospace width's
	// TODO: check, how "small" and "big" actually look on screen. 15 for medium seemed to be OK
	spSmall  = 13
	spMedium = 15
	spBig    = 25
)

var (
	reNumGreySmall   = regexp.MustCompile(ptNumGreySmall)
	reNumYellowSmall = regexp.MustCompile(ptNumYellowSmall)
	reNumRedBig      = regexp.MustCompile(ptNumRedBig)
	reCompositeRed   = regexp.MustCompile(ptCompositeRed)
)

// Glyph is a DoomPicture with some font identifying meta data
type glyph struct {
	*DoomPicture
	name string
	char rune
}

func newGlyph(name string, char rune, buff []byte) glyph {
	return glyph{
		DoomPicture: NewDoomPicture(buff),
		name:        name,
		char:        char,
	}
}

func (g *glyph) GetName() string {
	return g.name
}

// font is a collection of glyphs
type font struct {
	name         FontName
	spacing      int
	fallback     rune
	glyphNameMap map[string]glyph
	glyphRuneMap map[rune]*glyph
}

func newFont(name FontName, spacing int, fallback rune) font {
	return font{
		name:         name,
		spacing:      spacing,
		fallback:     fallback,
		glyphNameMap: make(map[string]glyph),
		glyphRuneMap: make(map[rune]*glyph),
	}
}

func (f font) addGlyph(g glyph) {
	f.glyphNameMap[g.name] = g
	f.glyphRuneMap[g.char] = &g
}

func (f *font) GetSpacing() int {
	return f.spacing
}

func (f *font) GetGlyph(r rune) *glyph {
	g := f.glyphRuneMap[r]
	if g == nil {
		return f.glyphRuneMap[f.fallback]
	}

	return g
}

// FontBook is a collection of fonts
type FontBook map[FontName]font

// NewFontBook initializes a new fontBook with the defined fonts
func NewFontBook() FontBook {
	fb := make(FontBook)
	fb[FnNumGreySmall] = newFont(FnNumGreySmall, spSmall, rune('0'))
	fb[FnNumYellowSmall] = newFont(FnNumYellowSmall, spSmall, rune('0'))
	fb[FnNumRedBig] = newFont(FnNumRedBig, spBig, rune('-'))
	fb[FnCompositeRed] = newFont(FnCompositeRed, spMedium, rune('?'))

	return fb
}

// getFont returns a font from the fontBook
func (fb *FontBook) getFont(name FontName) (*font, error) {
	f, ok := (*fb)[name]
	if !ok {
		return &font{}, fmt.Errorf("font not found: %v", FontName(name))
	}

	return &f, nil
}

// GetAllGraphics return a map of all DoomGraphics contained in the Fontbook
func (fb *FontBook) GetAllGraphics() map[string]*DoomPicture {
	gMap := make(map[string]*DoomPicture)
	for _, font := range *fb {
		for name, glyph := range font.glyphNameMap {
			gMap[name] = glyph.DoomPicture
		}
	}

	return gMap
}

func (fb *FontBook) tryAddExtra(lump wad.Lump) (added bool) {
	added = true

	if lump.Name == exNumRedBigMinus {
		g := newGlyph(lump.Name, rune('-'), lump.Data)
		(*fb)[FnNumRedBig].addGlyph(g)
		return
	}

	if lump.Name == exNumRedBigPercent {
		g := newGlyph(lump.Name, rune('%'), lump.Data)
		(*fb)[FnNumRedBig].addGlyph(g)
		return
	}

	return false
}

func (fb *FontBook) tryAdd(lump wad.Lump) (added bool, err error) {
	added = true

	if m := reNumGreySmall.FindStringSubmatch(lump.Name); m != nil {
		// parse the actual digit-char from the first match group
		g := newGlyph(lump.Name, rune(m[1][0]), lump.Data)
		(*fb)[FnNumGreySmall].addGlyph(g)
		return
	}

	if m := reNumYellowSmall.FindStringSubmatch(lump.Name); m != nil {
		// parse the actual digit-char from the first match group
		g := newGlyph(lump.Name, rune(m[1][0]), lump.Data)
		(*fb)[FnNumYellowSmall].addGlyph(g)
		return
	}

	if m := reNumRedBig.FindStringSubmatch(lump.Name); m != nil {
		// parse the actual digit-char from the first match group
		g := newGlyph(lump.Name, rune(m[1][0]), lump.Data)
		(*fb)[FnNumRedBig].addGlyph(g)
		return
	}

	if m := reCompositeRed.FindStringSubmatch(lump.Name); m != nil {
		// here, the matched group actually is decimal ascii mapping of the char
		ascii, err := strconv.Atoi(m[1])
		if err != nil {
			return added, err
		}

		// fix for wrong mapping, see foot note in https://zdoom.org/wiki/Composite_font
		if ascii == 121 {
			ascii = 124
		}

		g := newGlyph(lump.Name, rune(ascii), lump.Data)
		(*fb)[FnCompositeRed].addGlyph(g)
		return added, err
	}

	return false, nil
}

// LoadWAD fills the font book with embedded fonts from the WAD file
func (fb *FontBook) LoadWAD(w *wad.WAD) error {
	for _, lump := range w.Lumps() {
		if lump.Size == 0 {
			// skip empty lumps
			continue
		}

		// cheap non-regex stuff
		if (*fb).tryAddExtra(lump) {
			continue
		}

		// match patterns of the glyp groups ("fonts")
		_, err := (*fb).tryAdd(lump)
		if err != nil {
			return err
		}
	}
	return nil
}
