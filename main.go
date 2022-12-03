package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
)

var squareSize int = 40

func main() {
	palette := []Pixel{
		PixelFromInt32(0xffb703),
		PixelFromInt32(0xfb8500),
		PixelFromInt32(0xd00000),
		PixelFromInt32(0x8ecae6),
		PixelFromInt32(0x023047),
		PixelFromInt32(0x219ebc),
		PixelFromInt32(0x2a9d8f),
		PixelFromInt32(0xccc5b9),
	}
	// printAveragedImage(palette)
	// printAveragedImage(Palette{})

	printPaletteImage(palette)
	// printPaletteImage(Palette{})
	// printHarsherPixelatedImage()
	// printPixelatedImage()
}

func printPaletteImage(p Palette) {
	if len(p) == 0 {
		bfr := FromJPEG("sample.jpg")
		p = bfr.Palette(8)
	}

	edit := p.ToImage(50)
	printImage(edit)
}

func printAveragedImage(palette Palette) {
	bfr := FromJPEG("sample.jpg")
	if len(palette) == 0 {
		palette = bfr.Palette(24) // from image itself
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
	printImage(b2.ToImage())
}

func printHarsherPixelatedImage() {
	bfr := FromJPEG("sample.jpg")
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
	printImage(b2.ToImage())
}

func printPixelatedImage() {
	bfr := FromJPEG("sample.jpg")
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
	printImage(b2.ToImage())
}

func printImage(edit image.Image) {
	writer, err := os.Create("out.jpg")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open output file: %v", err)
		return
	}
	defer writer.Close()

	if err := jpeg.Encode(writer, edit, nil); err != nil {
		fmt.Fprintf(os.Stderr, "could not write output file: %v", err)
		return
	}
	fmt.Println("yo")
}
