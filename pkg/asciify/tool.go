package asciify

import (
	"downsample/pkg"
	"fmt"
	"path/filepath"
)

type GlyphMap map[rune]string

const _THRESHOLD uint32 = 10

var DefaultGlyphMap GlyphMap = GlyphMap{
	'@': "at",
	'!': "bang",
	':': "colon",
	',': "comma",
	'.': "dot",
	'=': "eq",
	'#': "hash",
	'O': "ough",
	'%': "pct",
	'|': "pipe",
	';': "semi",
	'o': "smol",
	'0': "zero",
}

/// Inspects images within directory in `path` and
/// prepares the output for asciify replacements
func InspectGlyphsAt(path string, gmap GlyphMap) {
	for glyph, name := range gmap {
		fmt.Printf(`{ Glyph: '%c', Name: %q, `, glyph, name)
		img := pkg.FromPNG(filepath.Join(path, name) + ".png")
		b := img.Bounds()

		occupied := 0
		empty := 0
		for x := 0; x < b.Max.X; x++ {
			for y := 0; y < b.Max.Y; y++ {
				r, g, b, _ := img.At(x, y).RGBA()
				if r > _THRESHOLD && g > _THRESHOLD && b > _THRESHOLD {
					occupied += 1
				} else {
					empty += 1
				}
			}
		}

		total := b.Max.X * b.Max.Y
		op := (float64(occupied) / float64(total)) * 100.0
		fmt.Printf("Pct: %f,", op*2)
		fmt.Printf("Y: %d },\n", uint8(op*0.02*256))
	}
}
