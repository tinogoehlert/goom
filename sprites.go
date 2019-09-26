package goom

type SpriteFrames map[byte]SpriteFrame

// SpriteFrame DOOM sprite frame
type SpriteFrame struct {
	angle byte
	frame byte
	name  string
	image *DoomImage
}

func (sf *SpriteFrame) Image() *DoomImage {
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
func (s *Sprite) AddSpriteFrame(lump *Lump) *SpriteFrame {
	var (
		angle = lump.Name[5] - 48
		frame = lump.Name[4]
	)
	sf := &SpriteFrame{
		angle: angle,
		frame: frame,
		name:  lump.Name[:6],
		image: &DoomImage{Name: lump.Name[:6]},
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
