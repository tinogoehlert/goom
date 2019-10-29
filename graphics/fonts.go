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
	fnNumGreySmall FontName = iota
	fnNumYellowSmall
	fnNumRedBig
	fnCompositeRed

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
	spBig    = 32
)

var (
	numGreySmallRe   = regexp.MustCompile(ptNumGreySmall)
	numYellowSmallRe = regexp.MustCompile(ptNumYellowSmall)
	numRedBigRe      = regexp.MustCompile(ptNumRedBig)
	compositeRedRe   = regexp.MustCompile(ptCompositeRed)
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

// Font is a collection of glyphs
type Font struct {
	name         FontName
	spacing      int
	glyphNameMap map[string]glyph
	glyphRuneMap map[rune]*glyph
}

func newFont(name FontName, spacing int) Font {
	return Font{
		name:         name,
		spacing:      spacing,
		glyphNameMap: make(map[string]glyph),
		glyphRuneMap: make(map[rune]*glyph),
	}
}

func (f Font) addGlyph(g glyph) {
	f.glyphNameMap[g.name] = g
	f.glyphRuneMap[g.char] = &g
}

// FontBook is a collection of fonts
type FontBook map[FontName]Font

// NewFontBook initializes a new Fontbook with the defined fonts
func NewFontBook() FontBook {
	fb := make(FontBook)
	fb[fnNumGreySmall] = newFont(fnNumGreySmall, spSmall)
	fb[fnNumYellowSmall] = newFont(fnNumYellowSmall, spSmall)
	fb[fnNumRedBig] = newFont(fnNumRedBig, spBig)
	fb[fnCompositeRed] = newFont(fnCompositeRed, spMedium)

	return fb
}

// GetFont returns a Font from the fontbook
func (fb *FontBook) GetFont(name FontName) (*Font, error) {
	f, ok := (*fb)[name]
	if !ok {
		return &Font{}, fmt.Errorf("font not found: %v", FontName(name))
	}

	return &f, nil
}

func (fb *FontBook) tryAddExtra(lump wad.Lump) (added bool) {
	added = true

	if lump.Name == exNumRedBigMinus {
		g := newGlyph(lump.Name, rune('-'), lump.Data)
		(*fb)[fnNumRedBig].addGlyph(g)
		return
	}

	if lump.Name == exNumRedBigPercent {
		g := newGlyph(lump.Name, rune('%'), lump.Data)
		(*fb)[fnNumRedBig].addGlyph(g)
		return
	}

	return false
}

func (fb *FontBook) tryAdd(lump wad.Lump) (added bool, err error) {
	added = true

	if m := numGreySmallRe.FindStringSubmatch(lump.Name); m != nil {
		// parse the actual digit-char from the first match group
		g := newGlyph(lump.Name, rune(m[1][0]), lump.Data)
		(*fb)[fnNumGreySmall].addGlyph(g)
		return
	}

	if m := numYellowSmallRe.FindStringSubmatch(lump.Name); m != nil {
		// parse the actual digit-char from the first match group
		g := newGlyph(lump.Name, rune(m[1][0]), lump.Data)
		(*fb)[fnNumYellowSmall].addGlyph(g)
		return
	}

	if m := numRedBigRe.FindStringSubmatch(lump.Name); m != nil {
		// parse the actual digit-char from the first match group
		g := newGlyph(lump.Name, rune(m[1][0]), lump.Data)
		(*fb)[fnNumRedBig].addGlyph(g)
		return
	}

	if m := compositeRedRe.FindStringSubmatch(lump.Name); m != nil {
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
		(*fb)[fnNumRedBig].addGlyph(g)
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

		// >>> match patterns of the glyp groups ("fonts")
		_, err := (*fb).tryAdd(lump)
		if err != nil {
			return err
		}
	}
	return nil
}
