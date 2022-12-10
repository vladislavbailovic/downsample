package main

import (
	"downsample/pkg"
	"image/color"
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
	palette := pkg.ImagePalette(bfr, 12)
	for i := 0; i < b.N; i++ {
		pkg.ConstrainImage(bfr, palette)
	}
}

func Benchmark_Averaged_Write(b *testing.B) {
	for i := 0; i < b.N; i++ {
		printAveragedImage([]color.Color{})
	}
}
