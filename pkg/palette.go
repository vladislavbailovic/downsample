package pkg

import (
	"image"
	"image/color"
)

func ToPaletteImage(p color.Palette, tileSize int) image.Image {
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
