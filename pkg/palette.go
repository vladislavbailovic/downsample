package pkg

import (
	"image"
	"image/color"
	"math"
)

type Palette []Pixel

func (p Palette) ToImage(tileSize int) image.Image {
	img := image.NewRGBA(image.Rect(
		0, 0, tileSize*len(p), tileSize))

	for pos := 0; pos < len(p); pos++ {
		c := p[pos]
		for y := 0; y < tileSize; y++ {
			for x := 0; x < tileSize; x++ {
				img.Set(
					(pos*tileSize)+x,
					y,
					color.RGBA{c.R, c.G, c.B, 1})
			}
		}
	}
	return img
}

// TODO: wow this is super inefficient
func (p Palette) ToImageBuffer(tileSize int) *ImageBuffer {
	return FromImage(p.ToImage(tileSize))
}

func (p Palette) ClosestTo(original Pixel) Pixel {
	min := math.MaxInt32
	result := original
	for _, px := range p {

		r := math.Abs(float64(px.R) - float64(original.R))
		g := math.Abs(float64(px.G) - float64(original.G))
		b := math.Abs(float64(px.B) - float64(original.B))
		tmp := int(r + g + b)

		if int(tmp) < min {
			min = int(tmp)
			result = px
		}
	}
	return result
}
