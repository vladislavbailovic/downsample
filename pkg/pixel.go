package pkg

type Pixel struct {
	R, G, B uint8
}

func NewPixel(r, g, b uint8) Pixel {
	return Pixel{R: r, G: g, B: b}
}

func PixelFromInt32(code int32) Pixel {
	b := uint8(code & 0xFF)
	g := uint8((code >> 8) & 0xFF)
	r := uint8((code >> 16) & 0xFF)
	return Pixel{R: r, G: g, B: b}
}

func (x Pixel) Hex() int32 {
	b := (int32(x.B) & 0xFF)
	g := (int32(x.G) & 0xFF) << 8
	r := int32(x.R) & 0xFF << 16
	return r | g | b
}

func (x Pixel) Clone() Pixel {
	return NewPixel(x.R, x.G, x.B)
}
