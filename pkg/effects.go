package pkg

import (
	"image"
	"image/color"
)

var squareSize int = 25

func GetTileSize() int {
	return squareSize
}

func SetTileSize(newSize int) {
	squareSize = newSize
}

type pixelateMode byte

const (
	ModePixelate     pixelateMode = iota
	ModeAndNormalize pixelateMode = iota
)

func PixelateImage(src image.Image, mode pixelateMode, norm Normalizer) image.Image {
	square := squareSize
	bounds := src.Bounds()
	dest := image.NewRGBA(bounds)

	if norm == nil {
		norm = StraightNormalizer{Q: RGBQuantizer{Factor: 4}}
	}

	for y := 0; y < bounds.Max.Y; y += square {
		for x := 0; x < bounds.Max.X; x += square {
			var normalized color.Color

			if ModeAndNormalize == mode {
				p := make(color.Palette, 0, square*square)

				for i := 0; i < square; i++ {
					for j := 0; j < square; j++ {
						dy := y + i
						if dy >= bounds.Max.Y {
							continue
						}
						dx := x + j
						if dx >= bounds.Max.X {
							continue
						}
						px := src.At(dx, dy)
						p = append(p, px)
					}
				}
				normalized = norm.Normalize(p, 4)[0]
			} else {
				normalized = src.At(x, y)
			}

			for i := 0; i < square; i++ {
				for j := 0; j < square; j++ {
					dy := y + i
					if dy >= bounds.Max.Y {
						continue
					}
					dx := x + j
					if dx >= bounds.Max.X {
						continue
					}
					dest.Set(dx, dy, normalized)
				}
			}
		}
	}
	return dest
}

func ConstrainImage(src image.Image, palette color.Palette, norm Normalizer) image.Image {
	square := squareSize
	bounds := src.Bounds()
	dest := image.NewRGBA(bounds)

	if norm == nil {
		norm = StraightNormalizer{Q: RGBQuantizer{Factor: 4}}
	}

	for y := 0; y < bounds.Max.Y; y += square {
		for x := 0; x < bounds.Max.X; x += square {
			p := make(color.Palette, 0, square*square)

			for i := 0; i < square; i++ {
				for j := 0; j < square; j++ {
					dy := y + i
					if dy >= bounds.Max.Y {
						continue
					}
					dx := x + j
					if dx >= bounds.Max.X {
						continue
					}
					px := src.At(dx, dy)
					p = append(p, px)
				}
			}

			normalized := norm.Normalize(p, 4)[0]
			closest := palette.Convert(normalized)
			for i := 0; i < square; i++ {
				for j := 0; j < square; j++ {
					dy := y + i
					if dy >= bounds.Max.Y {
						continue
					}
					dx := x + j
					if dx >= bounds.Max.X {
						continue
					}
					dest.Set(dx, dy, closest)
				}
			}
		}
	}
	return dest
}
