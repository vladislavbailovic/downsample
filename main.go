package main

import (
	"downsample/pkg"
	"path/filepath"
)

func main() {
	// palette := []pkg.Pixel{
	// 	pkg.PixelFromInt32(0xffb703),
	// 	pkg.PixelFromInt32(0xfb8500),
	// 	pkg.PixelFromInt32(0xd00000),
	// 	pkg.PixelFromInt32(0x8ecae6),
	// 	pkg.PixelFromInt32(0x023047),
	// 	pkg.PixelFromInt32(0x124057),
	// 	pkg.PixelFromInt32(0x225068),
	// 	pkg.PixelFromInt32(0x219ebc),
	// 	pkg.PixelFromInt32(0x2a9d8f),
	// 	pkg.PixelFromInt32(0xccc5b9),
	// }
	// printAveragedImage(palette)
	printAveragedImage(pkg.Palette{})
	// printHarsherPixelatedImage()
	// printPixelatedImage()
}

func printPaletteImage(p pkg.Palette, paletteFname string) {
	if len(p) == 0 {
		file := filepath.Join("testdata", "sample.jpg")
		bfr := pkg.FromJPEG(file)
		p = bfr.Palette(8)
	}

	edit := p.ToImageBuffer(50)
	edit.ToJPEGFile(paletteFname)
}

func printAveragedImage(palette pkg.Palette) {
	file := filepath.Join("testdata", "sample.jpg")
	bfr := pkg.FromJPEG(file)
	outputFname := "average-with-palette.jpg"
	paletteFname := "supplied-palette.jpg"
	if len(palette) == 0 {
		palette = bfr.Palette(12) // from image itself
		outputFname = "average-image-palette.jpg"
		paletteFname = "image-palette.jpg"
	}

	b2 := pkg.ConstrainImage(bfr, palette)
	b2.ToJPEGFile(outputFname)
	printPaletteImage(palette, paletteFname)
}

func printHarsherPixelatedImage() {
	file := filepath.Join("testdata", "sample.jpg")
	bfr := pkg.FromJPEG(file)
	b2 := pkg.PixelateImage(bfr, pkg.ModeAndNormalize)
	b2.ToJPEGFile("harsher.jpg")
}

func printPixelatedImage() {
	file := filepath.Join("testdata", "sample.jpg")
	bfr := pkg.FromJPEG(file)
	b2 := pkg.PixelateImage(bfr, pkg.ModePixelate)
	b2.ToJPEGFile("pixelated.jpg")
}
