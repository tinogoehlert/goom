package graphics

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
	"image/png"
	"os"
)

const (
	transparentColor = 255
)

// Image generic Image
type Image interface {
	ToRGBA(palette [256]color.RGBA) *image.RGBA
	ToPng(out string, palette [256]color.RGBA) error
	Width() int
	Height() int
	Top() int
	Left() int
}

// DoomPicture holds an image in the doom picture format
type DoomPicture struct {
	width  int
	height int
	left   int
	top    int
	data   []uint8
}

// NewDoomPicture gets picture from buffer
func NewDoomPicture(buff []byte) *DoomPicture {
	dp := &DoomPicture{
		width:  int(int16(binary.LittleEndian.Uint16(buff[0:2]))),
		height: int(int16(binary.LittleEndian.Uint16(buff[2:4]))),
		left:   int(int16(binary.LittleEndian.Uint16(buff[4:6]))),
		top:    int(int16(binary.LittleEndian.Uint16(buff[6:8]))),
	}
	offsets := make([]int32, dp.width, dp.width)
	r := bytes.NewBuffer(buff[8 : 8+(dp.width*4)])

	if err := binary.Read(r, binary.LittleEndian, offsets); err != nil {
		return nil
	}

	size := int(dp.width) * int(dp.height)
	dp.data = make([]byte, size, size)
	for y := 0; y < int(dp.height); y++ {
		for x := 0; x < int(dp.width); x++ {
			dp.data[y*int(dp.width)+x] = transparentColor
		}
	}
	for columnIndex, offset := range offsets {
		for {
			rowStart := buff[offset]
			offset++
			if rowStart == 255 {
				break
			}
			numPixels := buff[offset]
			offset++
			offset++ /* Padding */
			for i := 0; i < int(numPixels); i++ {
				pixelOffset := (int(rowStart)+i)*int(dp.width) + columnIndex
				dp.data[pixelOffset] = buff[offset]
				offset++
			}
			offset++ /* Padding */
		}
	}
	return dp
}

func newDummyPicture(width, height int) *DoomPicture {
	return &DoomPicture{
		width:  width,
		height: height,
		data:   make([]uint8, width*height),
	}
}

// Width return width of image
func (p *DoomPicture) Width() int { return int(p.width) }

// Height return height of image
func (p *DoomPicture) Height() int { return int(p.height) }

// Top offset. The number of pixels above the origin; where the top row is.
func (p *DoomPicture) Top() int { return int(p.top) }

// Left offset. The number of pixels to the left of the center; where the first column gets drawn.
func (p *DoomPicture) Left() int { return int(p.left) }

// ToRGBA converts picture to go image
func (p *DoomPicture) ToRGBA(palette [256]color.RGBA) *image.RGBA {
	bounds := image.Rect(0, 0, p.width, p.height)
	tex := image.NewRGBA(bounds)
	for y := 0; y < p.height; y++ {
		for x := 0; x < p.width; x++ {
			pixel := p.data[y*p.width+x]
			if pixel == transparentColor {
				tex.Set(x, y, color.RGBA{0, 0, 0, 0})
				continue
			}
			tex.Set(x, y, palette[pixel])
		}
	}
	return tex
}

// ToPng exports picture to PNG
func (p *DoomPicture) ToPng(out string, palette [256]color.RGBA) error {
	img := p.ToRGBA(palette)
	f, _ := os.Create(out)
	return png.Encode(f, img)
}
