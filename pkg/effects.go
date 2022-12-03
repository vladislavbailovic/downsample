package pkg

var squareSize int = 25

type pixelateMode byte

const (
	ModePixelate     pixelateMode = iota
	ModeAndNormalize pixelateMode = iota
)

func PixelateImage(bfr *ImageBuffer, mode pixelateMode) *ImageBuffer {
	square := squareSize
	edit := make([]*Pixel, len(bfr.pixels))

	for y := 0; y < bfr.height; y += square {
		for x := 0; x < bfr.width; x += square {
			var normalized Pixel

			if ModeAndNormalize == mode {
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
				normalized = normalizeColors_RGBA(p, 4)[0]
			} else {
				normalized = bfr.pixels[y*bfr.width+x].Clone()
			}

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
	return &ImageBuffer{
		width:  bfr.width,
		height: bfr.height,
		pixels: edit}
}

func ConstrainImage(bfr *ImageBuffer, palette Palette) *ImageBuffer {
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
	return &ImageBuffer{
		width:  bfr.width,
		height: bfr.height,
		pixels: edit}
}
