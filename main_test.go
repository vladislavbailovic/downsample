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
