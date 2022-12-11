package pkg

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
)

func FromJPEG(imgfile string) image.Image {
	fp, err := os.Open(imgfile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to open image: %v", err)
		return nil
	}
	defer fp.Close()

	img, err := jpeg.Decode(fp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to decode JPEG: %v", err)
	}

	return img
}

func ImagePalette(src image.Image, size uint8, norm Normalizer) color.Palette {
	if norm == nil {
		norm = StraightNormalizer{Q: RGBQuantizer{Factor: size}}
	}

	bounds := src.Bounds()
	all := make(color.Palette, 0, bounds.Max.X*bounds.Max.Y)
	for y := 0; y < bounds.Max.Y; y++ {
		for x := 0; x < bounds.Max.X; x++ {
			all = append(all, src.At(x, y))
		}
	}
	return norm.Normalize(all, size)
}

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
