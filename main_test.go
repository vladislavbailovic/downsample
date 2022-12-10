package main

import (
	"downsample/pkg"
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
	palette := bfr.Palette(12) // from image itself
	for i := 0; i < b.N; i++ {
		pkg.ConstrainImage(bfr, palette)
	}
}

func Benchmark_Averaged_Write(b *testing.B) {
	for i := 0; i < b.N; i++ {
		printAveragedImage(pkg.Palette{})
	}
}
