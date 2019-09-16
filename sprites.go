package goom

type SpriteFrames map[byte]SpriteFrame

// SpriteFrame DOOM sprite frame
type SpriteFrame struct {
	angle byte
	frame byte
	image *DoomImage
}

func (sf *SpriteFrame) Image() *DoomImage {
	return sf.image
}

// Sprite DOOM sprite
type Sprite struct {
	Name   string
	angles map[byte]SpriteFrames
}

// NewSprite Creates new sprite from lump
func NewSprite(name string) Sprite {
	return Sprite{
		Name:   name,
		angles: make(map[byte]SpriteFrames),
	}
}

// AddSpriteFrame Creates new sprite from lump
func (s *Sprite) AddSpriteFrame(lump *Lump) *SpriteFrame {
	var (
		angle = lump.Name[5] - 48
		frame = lump.Name[4] - 64
	)
	sf := SpriteFrame{
		angle: angle,
		frame: frame,
		image: &DoomImage{Name: lump.Name},
	}
	if _, ok := s.angles[angle]; !ok {
		s.angles[angle] = make(SpriteFrames)
	}
	s.angles[angle][frame] = sf
	return &sf
}

func (s *Sprite) GetHeadAngle() SpriteFrames {
	if _, ok := s.angles[1]; ok {
		return s.angles[1]
	}
	return s.angles[0]
}

func (s *Sprite) GetAngle(angle byte) SpriteFrames {
	if _, ok := s.angles[angle]; !ok {
		return s.angles[angle]
	}
	return s.GetHeadAngle()
}
