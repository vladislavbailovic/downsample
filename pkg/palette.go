package pkg

import (
	"image"
	"image/color"
	"math"
)

func ToPaletteImage(p []color.Color, tileSize int) image.Image {
	img := image.NewRGBA(image.Rect(
		0, 0, tileSize*len(p), tileSize))

	for pos := 0; pos < len(p); pos++ {
		c := p[pos]
		for y := 0; y < tileSize; y++ {
			for x := 0; x < tileSize; x++ {
				img.Set(
					(pos*tileSize)+x,
					y,
					c)
			}
		}
	}
	return img
}

func ClosestTo(original color.Color, palette []color.Color) color.Color {
	min := math.MaxInt32
	result := original

	r0, g0, b0, _ := original.RGBA()
	for _, px := range palette {
		rx, gx, bx, _ := px.RGBA()

		r := math.Abs(float64(rx) - float64(r0))
		g := math.Abs(float64(gx) - float64(g0))
		b := math.Abs(float64(bx) - float64(b0))
		tmp := int(r + g + b)

		if int(tmp) < min {
			min = int(tmp)
			result = px
		}
	}
	return result
}
