package pkg

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"sort"
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

func ImagePalette(src image.Image, size uint8) color.Palette {
	bounds := src.Bounds()
	all := make(color.Palette, 0, bounds.Max.X*bounds.Max.Y)
	for y := 0; y < bounds.Max.Y; y++ {
		for x := 0; x < bounds.Max.X; x++ {
			all = append(all, src.At(x, y))
		}
	}
	return normalizeColors_RGBA(all, size)
}

// TODO: improve this
func normalizeColors_RGBA(pxl color.Palette, size uint8) color.Palette {
	c := map[color.Color]int{}
	for _, p := range pxl {
		r, g, b, _ := p.RGBA()
		n := color.RGBA{
			R: (uint8(r/256) / size) * size,
			G: (uint8(g/256) / size) * size,
			B: (uint8(b/256) / size) * size,
			A: 0xFF,
		}
		c[n] += 1
	}

	keys := make(color.Palette, 0, len(c))
	for k := range c {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return c[keys[i]] > c[keys[j]]
	})

	palette := make(color.Palette, 0, size)
	var idx uint8 = 0
	for _, k := range keys {
		palette = append(palette, k)
		idx++
		if idx >= size {
			break
		}
	}

	return palette
}
