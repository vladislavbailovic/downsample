package main

import (
	"downsample/pkg"
<<<<<<< HEAD
	"image/color"
=======
>>>>>>> main
	"path/filepath"
	"testing"
)

func Benchmark_Pixelated(b *testing.B) {
	file := filepath.Join("testdata", "sample.jpg")
	bfr := pkg.FromJPEG(file)
	for i := 0; i < b.N; i++ {
		pkg.PixelateImage(bfr, pkg.ModePixelate)
	}
}

func Benchmark_Pixelated_Write(b *testing.B) {
	for i := 0; i < b.N; i++ {
		printPixelatedImage()
	}
}

func Benchmark_Averaged(b *testing.B) {
	file := filepath.Join("testdata", "sample.jpg")
	bfr := pkg.FromJPEG(file)
<<<<<<< HEAD
	palette := pkg.ImagePalette(bfr, 12)
=======
	palette := bfr.Palette(12) // from image itself
>>>>>>> main
	for i := 0; i < b.N; i++ {
		pkg.ConstrainImage(bfr, palette)
	}
}

func Benchmark_Averaged_Write(b *testing.B) {
	for i := 0; i < b.N; i++ {
<<<<<<< HEAD
		printAveragedImage([]color.Color{})
=======
		printAveragedImage(pkg.Palette{})
>>>>>>> main
	}
}
