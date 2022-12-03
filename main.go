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

	b2 := ConstrainImage(bfr, palette)
	b2.ToJPEGFile(outputFname)
	printPaletteImage(palette, paletteFname)
}

func printHarsherPixelatedImage() {
	file := filepath.Join("testdata", "sample.jpg")
	bfr := FromJPEG(file)
	b2 := PixelateImage(bfr, ModeAndNormalize)
	b2.ToJPEGFile("harsher.jpg")
}

func printPixelatedImage() {
	file := filepath.Join("testdata", "sample.jpg")
	bfr := FromJPEG(file)
	b2 := PixelateImage(bfr, ModePixelate)
	b2.ToJPEGFile("pixelated.jpg")
}
