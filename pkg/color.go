package pkg

import (
	"image/color"
	"sort"
)

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
