package main

import "testing"

func Test_Pixel_Int32(t *testing.T) {
	suite := map[int32][]uint8{
		0xFFFFFF: []uint8{255, 255, 255},
		0xFF0000: []uint8{255, 0, 0},
		0xFF0101: []uint8{255, 1, 1},
		0x00FF00: []uint8{0, 255, 0},
	}
	for code, test := range suite {
		t.Run("to hex", func(t *testing.T) {
			pixel := NewPixel(test[0], test[1], test[2])
			expected := code
			actual := pixel.Hex()
			if expected != actual {
				t.Errorf("wanted %x, got %x for %#v",
					expected, actual, pixel)
			}
		})
		t.Run("from hex", func(t *testing.T) {
			expected := NewPixel(test[0], test[1], test[2])
			actual := PixelFromInt32(code)
			if expected.R != actual.R ||
				expected.G != actual.G ||
				expected.B != actual.B {
				t.Errorf("wanted %#v (%x), got %#v (%x)",
					expected, expected.Hex(),
					actual, actual.Hex())
			}
		})
	}
}
