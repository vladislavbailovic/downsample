package main

import "testing"

func Test_Palette_ClosestTo(t *testing.T) {
	suite := map[int32]int32{
		0x880000: 0xB10000,
		0xFF0000: 0xF39300,
		0x008800: 0x006600,
		0x666666: 0x336633,
	}
	palette := Palette{
		PixelFromInt32(0xFF0000),
		PixelFromInt32(0x880000),
		PixelFromInt32(0x00FF00),
		PixelFromInt32(0x008800),
		PixelFromInt32(0x0000FF),
		PixelFromInt32(0x000088),
		PixelFromInt32(0x666666),
		PixelFromInt32(0xFFFFFF),
		PixelFromInt32(0x000000),
		PixelFromInt32(0xDE0A0D),
		PixelFromInt32(0xBADA55),
	}
	for expected, raw := range suite {
		t.Run("closest", func(t *testing.T) {
			test := PixelFromInt32(raw)
			actual := palette.ClosestTo(test)
			if expected != actual.Hex() {
				t.Errorf("wanted %06x, got %06x (%#v) for %x",
					expected, actual.Hex(), actual, raw)
			}
		})
	}
}
