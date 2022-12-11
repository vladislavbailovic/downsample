package pkg

import (
	"image/color"
	"sort"
)

type Normalizer interface {
	Normalize(color.Palette, uint8) color.Palette
}

type Quantizer interface {
	Quantize(color.Color) color.Color
}

type StraightNormalizer struct {
	Q Quantizer
}

func (x StraightNormalizer) Normalize(pxl color.Palette, size uint8) color.Palette {
	c := map[color.Color]int{}
	for _, p := range pxl {
		n := x.Q.Quantize(p)
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

type AverageNormalizer struct {
	Q Quantizer
}

func (x AverageNormalizer) Normalize(pxl color.Palette, size uint8) color.Palette {
	c := map[color.Color]internalPalette{}
	for _, p := range pxl {
		n := x.Q.Quantize(p)
		raw := c[n]
		raw.colors = append(raw.colors, p)
		raw.repr += 1
		c[n] = raw
	}

	keys := make(color.Palette, 0, len(c))
	for k := range c {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return c[keys[i]].repr > c[keys[j]].repr
	})

	palette := make(color.Palette, 0, size)
	var idx uint8 = 0
	for _, k := range keys {
		palette = append(palette, average(c[k].colors))
		idx++
		if idx >= size {
			break
		}
	}

	return palette
}

type RGBQuantizer struct {
	factor uint8
}

func (x RGBQuantizer) Quantize(c color.Color) color.Color {
	r, g, b, _ := c.RGBA()
	return color.RGBA{
		R: (uint8(r/256) / x.factor) * x.factor,
		G: (uint8(g/256) / x.factor) * x.factor,
		B: (uint8(b/256) / x.factor) * x.factor,
		A: 0xFF,
	}
}

type RGBShiftQuantizer struct {
	factor uint8
}

func (x RGBShiftQuantizer) Quantize(c color.Color) color.Color {
	r, g, b, _ := c.RGBA()
	return color.RGBA{
		R: (uint8(r/256) >> x.factor) << x.factor,
		G: (uint8(g/256) >> x.factor) << x.factor,
		B: (uint8(b/256) >> x.factor) << x.factor,
		A: 0xFF,
	}
}

type GrayQuantizer struct {
	factor uint8
}

func (x GrayQuantizer) Quantize(c color.Color) color.Color {
	rn := color.GrayModel.Convert(c).(color.Gray)
	return color.Gray{Y: rn.Y / x.factor * x.factor}
}

type GrayShiftQuantizer struct {
	factor uint8
}

func (x GrayShiftQuantizer) Quantize(c color.Color) color.Color {
	rn := color.GrayModel.Convert(c).(color.Gray)
	return color.Gray{Y: rn.Y << x.factor >> x.factor}
}

type internalPalette struct {
	colors color.Palette
	repr   int
}

func average(r color.Palette) color.Color {
	var r0, g0, b0, a0 int
	d := len(r)
	for _, c := range r {
		r, g, b, a := c.RGBA()
		r0 += int(r)
		g0 += int(g)
		b0 += int(b)
		a0 += int(a)
	}
	rd := uint8((r0 / 255) / d)
	gd := uint8((g0 / 255) / d)
	bd := uint8((b0 / 255) / d)
	ad := uint8((a0 / 255) / d)
	return color.RGBA{R: rd, G: gd, B: bd, A: ad}
}
