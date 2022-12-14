package asciify

import (
	"downsample/pkg"
	"fmt"
	"image/color"
	"image/color/palette"
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

const (
	_DefaultTileWidth = 7
	// _DefaultTileHeight = 14
)

type Asciifier struct {
	Replacements []Replacement
	Replacer     replacer
	TileWidth    int
	TileHeight   int
}

type replacer interface {
	initialize([]Replacement)
	replace(color.Palette) (string, color.Color)
	wrap(string) string
}

func (a *Asciifier) getTileWidth() int {
	if a.TileWidth != 0 {
		return a.TileWidth
	}
	return _DefaultTileWidth
}

func (a *Asciifier) getTileHeight() int {
	if a.TileHeight != 0 {
		return a.TileHeight
	}
	return 2 * a.getTileWidth()
}

func (a *Asciifier) Asciify(imagePath string) string {
	a.Replacer.initialize(a.Replacements)

	bfr := pkg.FromJPEG(imagePath)

	b := bfr.Bounds()
	xincr := a.getTileWidth()
	yincr := a.getTileHeight()
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
			rpl, _ := a.Replacer.replace(p)
			cols = append(cols, rpl)
		}
		rows = append(rows, strings.Join(cols, ""))
	}
	return a.Replacer.wrap(strings.Join(rows, "\n"))
}

type PlainReplacer struct {
	norm       pkg.Normalizer
	rplPalette color.Palette
	rplMap     map[color.Color]string
}

func (x *PlainReplacer) wrap(res string) string {
	return res
}

func (x *PlainReplacer) initialize(rpl []Replacement) {
	x.rplPalette, x.rplMap = makePaletteReplacementTable(rpl)
	x.norm = pkg.DistributionNormalizer{Q: pkg.RGBQuantizer{Factor: 12}, Spread: 12.0}
}

func (x *PlainReplacer) replace(p color.Palette) (string, color.Color) {
	normalized := x.norm.Normalize(p, 4)[0]
	closest := x.rplPalette.Convert(normalized)

	return x.rplMap[closest], normalized
}

type ConsoleReplacer struct {
	PlainReplacer
	palette color.Palette
	colors  map[color.Color]int
}

const (
	AnsiBlack = iota + 30
	AnsiRed
	AnsiGreen
	AnsiYellow
	AnsiBlue
	AnsiMagenta
	AnsiCyan
	AnsiWhite
)
const (
	AnsiBrightBlack = iota + 90
	AnsiBrightRed
	AnsiBrightGreen
	AnsiBrightYellow
	AnsiBrightBlue
	AnsiBrightMagenta
	AnsiBrightCyan
	AnsiBrightWhite
)

func (x *ConsoleReplacer) initialize(rpl []Replacement) {
	x.PlainReplacer.initialize(rpl)
	x.colors = map[color.Color]int{
		color.RGBA{A: 0xFF, R: 0, G: 0, B: 0}:       AnsiBlack,
		color.RGBA{A: 0xFF, R: 170, G: 0, B: 0}:     AnsiRed,
		color.RGBA{A: 0xFF, R: 0, G: 170, B: 0}:     AnsiGreen,
		color.RGBA{A: 0xFF, R: 170, G: 85, B: 0}:    AnsiYellow,
		color.RGBA{A: 0xFF, R: 0, G: 0, B: 170}:     AnsiBlue,
		color.RGBA{A: 0xFF, R: 170, G: 0, B: 170}:   AnsiMagenta,
		color.RGBA{A: 0xFF, R: 0, G: 170, B: 170}:   AnsiCyan,
		color.RGBA{A: 0xFF, R: 170, G: 170, B: 170}: AnsiMagenta,

		color.RGBA{A: 0xFF, R: 85, G: 85, B: 85}:    AnsiBrightBlack,
		color.RGBA{A: 0xFF, R: 255, G: 85, B: 85}:   AnsiBrightRed,
		color.RGBA{A: 0xFF, R: 85, G: 255, B: 85}:   AnsiBrightGreen,
		color.RGBA{A: 0xFF, R: 255, G: 255, B: 85}:  AnsiBrightYellow,
		color.RGBA{A: 0xFF, R: 85, G: 85, B: 255}:   AnsiBrightBlue,
		color.RGBA{A: 0xFF, R: 255, G: 85, B: 255}:  AnsiBrightMagenta,
		color.RGBA{A: 0xFF, R: 85, G: 255, B: 255}:  AnsiBrightCyan,
		color.RGBA{A: 0xFF, R: 255, G: 255, B: 255}: AnsiBrightBlue,
	}
	for c, _ := range x.colors {
		x.palette = append(x.palette, c)
	}
}

func (x *ConsoleReplacer) replace(p color.Palette) (string, color.Color) {
	str, col := x.PlainReplacer.replace(p)
	c := x.palette.Convert(col)
	code := x.colors[c]
	return fmt.Sprintf("\u001B[%dm%s\u001B[0m", code, strings.Replace(str, "%", "%%", -1)), col
}

type HtmlReplacer struct {
	PlainReplacer
}

func (x *HtmlReplacer) replace(p color.Palette) (string, color.Color) {
	str, col := x.PlainReplacer.replace(p)
	r, g, b, _ := color.Palette(palette.WebSafe).Convert(col).RGBA()
	return fmt.Sprintf(
		`<span style="color: #%02x%02x%02x">%s</span>`,
		r/256, g/256, b/256,
		strings.Replace(str, "%", "%%", -1)), col
}
func (x *HtmlReplacer) wrap(s string) string {
	return `<pre style="background: #000">` + s + "</pre>"
}
