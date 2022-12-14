package main

import (
	"downsample/pkg"
	"downsample/pkg/asciify"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
)

func main() {
	asciifyImage()
	// palette := color.Palette{
	// 	color.RGBA{R: 0xFF, G: 0xB7, B: 0x03, A: 0xFF},
	// 	color.RGBA{R: 0xFB, G: 0x85, B: 0x00, A: 0xFF},
	// 	color.RGBA{R: 0xD0, G: 0x00, B: 0x00, A: 0xFF},
	// 	color.RGBA{R: 0x8E, G: 0xCA, B: 0xE6, A: 0xFF},
	// 	color.RGBA{R: 0x02, G: 0x30, B: 0x47, A: 0xFF},
	// 	color.RGBA{R: 0x12, G: 0x40, B: 0x57, A: 0xFF},
	// 	color.RGBA{R: 0x22, G: 0x50, B: 0x68, A: 0xFF},
	// 	color.RGBA{R: 0x21, G: 0x9E, B: 0xBC, A: 0xFF},
	// 	color.RGBA{R: 0x2A, G: 0x9D, B: 0x8F, A: 0xFF},
	// 	color.RGBA{R: 0xCC, G: 0xC5, B: 0xB9, A: 0xFF},
	// }
	// printAveragedImage(palette)
	// printAveragedImage(color.Palette{})
	// printHarsherPixelatedImage()
	// printPixelatedImage()
}

func asciifyImage() {
	file := filepath.Join("testdata", "sample.jpg")
	a := asciify.Asciifier{
		Replacements: asciify.UnicodeReplacements,
		Replacer:     &asciify.ConsoleReplacer{},
	}
	a.Asciify(file)
}

func printPaletteImage(p color.Palette, paletteFname string) {
	if len(p) == 0 {
		file := filepath.Join("testdata", "sample.jpg")
		bfr := pkg.FromJPEG(file)
		p = pkg.ImagePalette(bfr, 12, nil)
	}

	edit := pkg.ToPaletteImage(p, 50)
	ToJPEGFile(edit, paletteFname)
}

func printAveragedImage(palette color.Palette) {
	file := filepath.Join("testdata", "sample.jpg")
	bfr := pkg.FromJPEG(file)
	outputFname := "average-with-palette.jpg"
	paletteFname := "supplied-palette.jpg"
	if len(palette) == 0 {
		palette = pkg.ImagePalette(bfr, 12, nil)
		outputFname = "average-image-palette.jpg"
		paletteFname = "image-palette.jpg"
	}

	b2 := pkg.ConstrainImage(bfr, palette, nil)
	ToJPEGFile(b2, outputFname)
	printPaletteImage(palette, paletteFname)
}

func printHarsherPixelatedImage() {
	file := filepath.Join("testdata", "sample.jpg")
	bfr := pkg.FromJPEG(file)
	b2 := pkg.PixelateImage(bfr, pkg.ModeAndNormalize, nil)
	ToJPEGFile(b2, "harsher.jpg")
}

func printPixelatedImage() {
	file := filepath.Join("testdata", "sample.jpg")
	bfr := pkg.FromJPEG(file)
	b2 := pkg.PixelateImage(bfr, pkg.ModePixelate, nil)
	ToJPEGFile(b2, "pixelated.jpg")
}

func ToJPEGFile(edit image.Image, imgpath string) error {
	writer, err := os.Create(imgpath)
	if err != nil {
		return err
	}
	defer writer.Close()

	if err := jpeg.Encode(writer, edit, nil); err != nil {
		return err
	}

	return nil
}
