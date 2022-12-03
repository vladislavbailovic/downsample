package main

import (
	"path/filepath"
)

var squareSize int = 25

func main() {
	palette := []Pixel{
		PixelFromInt32(0xffb703),
		PixelFromInt32(0xfb8500),
		PixelFromInt32(0xd00000),
		PixelFromInt32(0x8ecae6),
		PixelFromInt32(0x023047),
		PixelFromInt32(0x124057),
		PixelFromInt32(0x225068),
		PixelFromInt32(0x219ebc),
		PixelFromInt32(0x2a9d8f),
		PixelFromInt32(0xccc5b9),
	}
	printAveragedImage(palette)
	printAveragedImage(Palette{})
	printHarsherPixelatedImage()
	printPixelatedImage()
}

func printPaletteImage(p Palette, paletteFname string) {
	if len(p) == 0 {
		file := filepath.Join("testdata", "sample.jpg")
		bfr := FromJPEG(file)
		p = bfr.Palette(8)
	}

	edit := p.ToImageBuffer(50)
	edit.ToJPEGFile(paletteFname)
}

func printAveragedImage(palette Palette) {
	file := filepath.Join("testdata", "sample.jpg")
	bfr := FromJPEG(file)
	outputFname := "average-with-palette.jpg"
	paletteFname := "supplied-palette.jpg"
	if len(palette) == 0 {
		palette = bfr.Palette(12) // from image itself
		outputFname = "average-image-palette.jpg"
		paletteFname = "image-palette.jpg"
	}
	square := squareSize
	edit := make([]*Pixel, len(bfr.pixels))

	for y := 0; y < bfr.height; y += square {
		for x := 0; x < bfr.width; x += square {
			p := make([]*Pixel, 0, square*square)

			for i := 0; i < square; i++ {
				for j := 0; j < square; j++ {
					dy := y + i
					if dy >= bfr.height {
						continue
					}
					dx := x + j
					if dx >= bfr.width {
						continue
					}
					px := bfr.pixels[dy*bfr.width+dx]
					p = append(p, px)
				}
			}

			normalized := normalizeColors_RGBA(p, 4)[0]
			closest := palette.ClosestTo(normalized)
			for i := 0; i < square; i++ {
				for j := 0; j < square; j++ {
					dy := y + i
					if dy >= bfr.height {
						continue
					}
					dx := x + j
					if dx >= bfr.width {
						continue
					}
					edit[dy*bfr.width+dx] = &closest
				}
			}
		}
	}
	b2 := &ImageBuffer{
		width:  bfr.width,
		height: bfr.height,
		pixels: edit}
	b2.ToJPEGFile(outputFname)
	printPaletteImage(palette, paletteFname)
}

func printHarsherPixelatedImage() {
	file := filepath.Join("testdata", "sample.jpg")
	bfr := FromJPEG(file)
	square := squareSize
	edit := make([]*Pixel, len(bfr.pixels))

	for y := 0; y < bfr.height; y += square {
		for x := 0; x < bfr.width; x += square {
			p := make([]*Pixel, 0, square*square)

			for i := 0; i < square; i++ {
				for j := 0; j < square; j++ {
					dy := y + i
					if dy >= bfr.height {
						continue
					}
					dx := x + j
					if dx >= bfr.width {
						continue
					}
					px := bfr.pixels[dy*bfr.width+dx]
					p = append(p, px)
				}
			}

			normalized := normalizeColors_RGBA(p, 4)[0]
			for i := 0; i < square; i++ {
				for j := 0; j < square; j++ {
					dy := y + i
					if dy >= bfr.height {
						continue
					}
					dx := x + j
					if dx >= bfr.width {
						continue
					}
					edit[dy*bfr.width+dx] = &normalized
				}
			}
		}
	}
	b2 := &ImageBuffer{
		width:  bfr.width,
		height: bfr.height,
		pixels: edit}
	b2.ToJPEGFile("harsher.jpg")
}

func printPixelatedImage() {
	file := filepath.Join("testdata", "sample.jpg")
	bfr := FromJPEG(file)
	square := squareSize
	edit := make([]*Pixel, len(bfr.pixels))

	for y := 0; y < bfr.height; y += square {
		for x := 0; x < bfr.width; x += square {
			normalized := bfr.pixels[y*bfr.width+x]
			for i := 0; i < square; i++ {
				for j := 0; j < square; j++ {
					dy := y + i
					if dy >= bfr.height {
						continue
					}
					dx := x + j
					if dx >= bfr.width {
						continue
					}
					edit[dy*bfr.width+dx] = normalized
				}
			}
		}
	}
	b2 := &ImageBuffer{
		width:  bfr.width,
		height: bfr.height,
		pixels: edit}
	b2.ToJPEGFile("pixelated.jpg")
}
