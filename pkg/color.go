package pkg

import (
	"image/color"
	"math"
	"sort"
)

type NormalizerType byte

const (
	NormalizerNormal  NormalizerType = iota
	NormalizerAverage NormalizerType = iota
)

type Normalizer interface {
	Normalize(color.Palette, uint8) color.Palette
}

type QuantizerType byte

const (
	QuantizerRGB       QuantizerType = iota
	QuantizerShiftRGB  QuantizerType = iota
	QuantizerGray      QuantizerType = iota
	QuantizerShiftGray QuantizerType = iota
)

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

type DistributionNormalizer struct {
	Q      Quantizer
	Spread float64
}

func (x DistributionNormalizer) Normalize(pxl color.Palette, size uint8) color.Palette {
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
		if len(palette) > 0 {
			if x.isTooClose(k, palette) {
				continue
			}
		}
		palette = append(palette, k)
		idx++
		if idx >= size {
			break
		}
	}

	return palette
}

func (x DistributionNormalizer) distance(c1 color.Color, c2 color.Color) float64 { // in pct
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()

	rmean := (float64(r1/256) + float64(r2/256)) / 2
	r := float64(r1/256) - float64(r2/256)
	g := float64(g1/256) - float64(g2/256)
	b := float64(b1/256) - float64(b2/256)

	wr := 2.0 + rmean/256.0
	wg := 4.0
	wb := 2.0 + (255.0-rmean)/256.0

	dist := math.Sqrt(wr*r*r + wg*g*g + wb*b*b)
	max := 764.8339663572415
	return (dist / max) * 100
}

func (x DistributionNormalizer) isTooClose(subject color.Color, palette color.Palette) bool {
	return x.isTooClose_Diff(subject, palette)
}

func (x DistributionNormalizer) isTooClose_Diff(subject color.Color, palette color.Palette) bool {
	for _, to := range palette {
		dist := x.distance(subject, to)
		if dist < x.Spread {
			return true
		}
	}
	return false
}
func (x DistributionNormalizer) isTooClose_Rgb(subject color.Color, palette color.Palette) bool {
	rs, gs, bs, _ := subject.RGBA()
	for _, to := range palette {
		rt, gt, bt, _ := to.RGBA()

		rd := math.Abs(float64(rs/256) - float64(rt/256))
		if rd < x.Spread {
			return true
		}

		gd := math.Abs(float64(gs/256) - float64(gt/256))
		if gd < x.Spread {
			return true
		}

		bd := math.Abs(float64(bs/256) - float64(bt/256))
		if bd < x.Spread {
			return true
		}
	}
	return false
}

type RGBQuantizer struct {
	Factor uint8
}

func (x RGBQuantizer) Quantize(c color.Color) color.Color {
	r, g, b, _ := c.RGBA()
	return color.RGBA{
		R: (uint8(r/256) / x.Factor) * x.Factor,
		G: (uint8(g/256) / x.Factor) * x.Factor,
		B: (uint8(b/256) / x.Factor) * x.Factor,
		A: 0xFF,
	}
}

type RGBShiftQuantizer struct {
	Factor uint8
}

func (x RGBShiftQuantizer) Quantize(c color.Color) color.Color {
	r, g, b, _ := c.RGBA()
	return color.RGBA{
		R: (uint8(r/256) >> x.Factor) << x.Factor,
		G: (uint8(g/256) >> x.Factor) << x.Factor,
		B: (uint8(b/256) >> x.Factor) << x.Factor,
		A: 0xFF,
	}
}

type GrayQuantizer struct {
	Factor uint8
}

func (x GrayQuantizer) Quantize(c color.Color) color.Color {
	rn := color.GrayModel.Convert(c).(color.Gray)
	return color.Gray{Y: rn.Y / x.Factor * x.Factor}
}

type GrayShiftQuantizer struct {
	Factor uint8
}

func (x GrayShiftQuantizer) Quantize(c color.Color) color.Color {
	rn := color.GrayModel.Convert(c).(color.Gray)
	return color.Gray{Y: rn.Y << x.Factor >> x.Factor}
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
