package asciify

import (
	"downsample/pkg"
	"fmt"
	"image/color"
	"strings"
)

func makePaletteReplacementTable(replacements []Replacement) (color.Palette, map[color.Color]string) {
	palette := make(color.Palette, 0, len(replacements))
	rplMap := make(map[color.Color]string, len(replacements))
	for _, rpl := range replacements {
		c := color.Gray{Y: rpl.Y}
		palette = append(palette, c)
		rplMap[c] = rpl.Glyph
	}
	return palette, rplMap
}

func Asciify(imagePath string, rpl []Replacement) {
	palette, rplMap := makePaletteReplacementTable(rpl)

	bfr := pkg.FromJPEG(imagePath)
	norm := pkg.DistributionNormalizer{Q: pkg.RGBQuantizer{Factor: 12}, Spread: 12.0}

	b := bfr.Bounds()
	xincr := 7
	yincr := 14
	rows := make([]string, 0, b.Max.Y/yincr)
	for y := 0; y < b.Max.Y; y += yincr {
		cols := make([]string, 0, b.Max.X/xincr)
		for x := 0; x < b.Max.X; x += xincr {
			p := make(color.Palette, 0, yincr*xincr)

			for i := 0; i < yincr; i++ {
				for j := 0; j < xincr; j++ {
					dy := y + i
					if dy >= b.Max.Y {
						continue
					}
					dx := x + j
					if dx >= b.Max.X {
						continue
					}
					px := bfr.At(dx, dy)
					p = append(p, px)
				}
			}

			normalized := norm.Normalize(p, 4)[0]
			closest := palette.Convert(normalized)

			cols = append(cols, rplMap[closest])
		}
		rows = append(rows, strings.Join(cols, ""))
	}
	fmt.Println(strings.Join(rows, "\n"))
}
