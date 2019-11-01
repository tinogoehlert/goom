package graphics

import (
	"regexp"

	"github.com/tinogoehlert/goom/wad"
)

// SpriteFrames map of frames
type SpriteFrames map[byte]SpriteFrame

// SpriteFrame DOOM sprite frame
type SpriteFrame struct {
	angles [9]Image
	frame  byte
	name   string
}

func (sf *SpriteFrame) Image(angle int) Image {
	if sf.angles[angle] != nil && angle == 1 {
		return sf.angles[0]
	}
	return sf.angles[angle]
}

func (sf *SpriteFrame) Angles() [9]Image {
	return sf.angles
}

func (sf *SpriteFrame) FrameID() byte {
	return sf.frame
}

func (sf *SpriteFrame) Name() string {
	return sf.name
}

// Sprite DOOM sprite
type Sprite struct {
	Name   string
	frames map[string]*SpriteFrame
	first  string
}

// NewSprite Creates new sprite from lump
func NewSprite(name string) Sprite {
	return Sprite{
		Name:   name,
		frames: make(map[string]*SpriteFrame),
	}
}

// AddSpriteFrame Creates new sprite from lump
func (s *Sprite) AddSpriteFrame(lump *wad.Lump) *SpriteFrame {
	var (
		frame = lump.Name[4]
	)

	sf, ok := s.frames[lump.Name[:5]]
	if !ok {
		sf = &SpriteFrame{
			angles: [9]Image{},
			frame:  frame,
			name:   lump.Name[:5],
		}
		s.frames[lump.Name[:5]] = sf
	}
	sf.angles[lump.Name[5]-48] = NewDoomPicture(lump.Data)
	if len(lump.Name) == 8 {
		sf.angles[lump.Name[7]-48] = sf.angles[lump.Name[5]-48]
	}
	return sf
}

func (s *Sprite) GetFrame(angle, frame byte) *SpriteFrame {
	if frame, ok := s.frames[s.Name+string(angle)+string(frame)]; !ok && frame != nil {
		return frame
	}
	return s.frames[s.first]
}

func (s *Sprite) FirstFrame() *SpriteFrame {
	return s.frames[s.first]
}

func (s *Sprite) Frames(cb func(sf *SpriteFrame)) {
	for _, v := range s.frames {
		cb(v)
	}
}

// SpriteStore string map of sprites
type SpriteStore map[string]Sprite

// NewSpriteStore creates new level store
func NewSpriteStore() SpriteStore {
	return make(SpriteStore)
}

func (ss SpriteStore) LoadWAD(w *wad.WAD) {
	var (
		spriteStartRegex = regexp.MustCompile(`^S?_START`)
		spriteEndRegex   = regexp.MustCompile(`^S?_END`)
		lumps            = w.Lumps()
	)

	for i := 0; i < len(lumps); i++ {
		lump := &lumps[i]
		if spriteStartRegex.Match([]byte(lump.Name)) {
			for {
				lump := &lumps[i]
				if lump.Size > 0 {
					s, ok := ss[lump.Name[:4]]
					if !ok {
						s = NewSprite(lump.Name[:4])
						s.first = lump.Name[:5]
						ss[lump.Name[:4]] = s
					}
					s.AddSpriteFrame(lump)
				}
				if spriteEndRegex.Match([]byte(lump.Name)) {
					break
				}
				i++
			}
		}
	}
}
