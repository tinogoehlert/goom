package graphics

import (
	"regexp"

	"github.com/tinogoehlert/goom/wad"
)

// SpriteFrames map of frames
type SpriteFrames map[byte]SpriteFrame

// SpriteFrame DOOM sprite frame
type SpriteFrame struct {
	angle byte
	frame byte
	name  string
	image Image
}

func (sf *SpriteFrame) Image() Image {
	return sf.image
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
		angle = lump.Name[5] - 48
		frame = lump.Name[4]
	)
	sf := &SpriteFrame{
		angle: angle,
		frame: frame,
		name:  lump.Name[:6],
		image: NewDoomPicture(lump.Data),
	}

	s.frames[lump.Name[:6]] = sf
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

func (ss SpriteStore) LoadWAD(w *wad.WAD) error {
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
						s.first = lump.Name[:6]
						ss[lump.Name[:4]] = s
					}
					sf := s.AddSpriteFrame(lump)
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

	return nil
}
