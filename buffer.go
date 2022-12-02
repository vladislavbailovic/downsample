package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
	"path/filepath"
	"sort"
)

type ImageBuffer struct {
	width, height int
	pixels        []*Pixel
}

func FromJPEG(imgfile string) *ImageBuffer {
	fp, err := os.Open(filepath.Join("testdata", imgfile))
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to open image: %v", err)
		return nil
	}
	defer fp.Close()

	img, err := jpeg.Decode(fp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to decode JPEG: %v", err)
	}

	bounds := img.Bounds()
	bfr := make([]*Pixel, bounds.Max.X*bounds.Max.Y)

	for y := 0; y < bounds.Max.Y; y++ {
		for x := 0; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			pxl := NewPixel(
				uint8(r>>8),
				uint8(g>>8),
				uint8(b>>8))
			bfr[y*bounds.Max.X+x] = &pxl
		}
	}

	return &ImageBuffer{
		width:  bounds.Max.X,
		height: bounds.Max.Y,
		pixels: bfr}
}

func (b *ImageBuffer) drawSquare(
	x, y, width, height int, color Pixel) error {

	px := color.Clone()
	for dy := y; dy < y+height; dy++ {
		if dy < 0 || dy >= b.height {
			continue
		}
		for dx := x; dx < x+width; dx++ {
			if dx < 0 || dx >= b.width {
				continue
			}
			b.pixels[dy*b.width+dx] = &px
		}
	}

	return nil
}

func (b *ImageBuffer) ToImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, b.width, b.height))

	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			c := b.pixels[y*b.width+x]
			img.Set(x, y, color.RGBA{c.R, c.G, c.B, 1})
		}
	}
	return img
}

type Palette []Pixel

func (b *ImageBuffer) Palette(size uint8) Palette {
	return normalizeColors_RGBA(b.pixels, size)
}

func (p Palette) ToImage(tileSize int) image.Image {
	img := image.NewRGBA(image.Rect(
		0, 0, tileSize*len(p), tileSize))

	for pos := 0; pos < len(p); pos++ {
		c := p[pos]
		for y := 0; y < tileSize; y++ {
			for x := 0; x < tileSize; x++ {
				img.Set(
					(pos*tileSize)+x,
					y,
					color.RGBA{c.R, c.G, c.B, 1})
			}
		}
	}
	return img
}

func getClosestColor(palette []Pixel, original Pixel) Pixel {
	min := math.MaxInt32
	result := original
	for _, px := range palette {

		tmp := math.Abs(float64(px.R) - float64(original.R))
		tmp += math.Abs(float64(px.G) - float64(original.G))
		tmp += math.Abs(float64(px.B) - float64(original.B))

		if int(tmp) < min {
			min = int(tmp)
			result = px
		}
	}
	return result
}

func normalizeColors_RGBA(pxl []*Pixel, size uint8) Palette {
	c := map[int32]int{}
	for _, p := range pxl {
		key := NewPixel(
			(p.R/size)*size,
			(p.G/size)*size,
			(p.B/size)*size,
		)
		c[key.Hex()] += 1
	}

	keys := make([]int32, 0, len(c))
	for k := range c {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return c[keys[i]] > c[keys[j]]
	})

	palette := make([]Pixel, 0, size)
	var idx uint8 = 0
	for _, k := range keys {
		palette = append(palette, PixelFromInt32(k))
		idx++
		if idx >= size {
			break
		}
	}

	return palette
}
