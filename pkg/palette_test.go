package pkg

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

func Test_Palette_ToImage(t *testing.T) {
	tileSize := 5
	palette := Palette{
		PixelFromInt32(0xFF0000),
		PixelFromInt32(0x00FF00),
		PixelFromInt32(0x0000FF),
	}
	img := palette.ToImage(tileSize)
	bfr := FromImage(img)

	if bfr.height != tileSize {
		t.Errorf("unexpected palette image height: %d",
			bfr.height)
	}

	if bfr.width != len(palette)*tileSize {
		t.Errorf("unexpected palette image width: %d",
			bfr.width)
	}

	for idx, px := range bfr.pixels {
		pos := (idx / tileSize) % len(palette)
		expected := palette[pos].Hex()
		actual := px.Hex()
		if expected != actual {
			t.Errorf("at %d, wanted %06x, got %06x",
				idx, expected, actual)
		}
	}
}
