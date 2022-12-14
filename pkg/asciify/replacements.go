package asciify

type Replacement struct {
	Glyph string
	Name  string
	Pct   float64
	Y     uint8
}

var AsciiReplacements []Replacement = []Replacement{
	{Glyph: "%", Name: "pct", Pct: 71.938776, Y: 184},
	{Glyph: "o", Name: "smol", Pct: 53.061224, Y: 135},
	{Glyph: "@", Name: "at", Pct: 90.306122, Y: 231},
	{Glyph: "!", Name: "bang", Pct: 25.000000, Y: 64},
	{Glyph: ":", Name: "colon", Pct: 14.285714, Y: 36},
	{Glyph: ".", Name: "dot", Pct: 7.142857, Y: 18},
	{Glyph: "=", Name: "eq", Pct: 36.734694, Y: 94},
	{Glyph: "#", Name: "hash", Pct: 62.244898, Y: 159},
	{Glyph: "0", Name: "zero", Pct: 70.408163, Y: 180},
	{Glyph: ",", Name: "comma", Pct: 13.775510, Y: 35},
	{Glyph: "O", Name: "ough", Pct: 66.836735, Y: 171},
	{Glyph: "|", Name: "pipe", Pct: 51.020408, Y: 130},
	{Glyph: ";", Name: "semi", Pct: 20.918367, Y: 53},
}
