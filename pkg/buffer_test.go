package pkg

/*
func Test_drawSquare_Center(t *testing.T) {
	raw := []int32{
		0, 0, 0, 0, 0,
		0, 1, 1, 1, 0,
		0, 1, 1, 1, 0,
		0, 1, 1, 1, 0,
		0, 0, 0, 0, 0,
	}
	buffer := make([]*Pixel, len(raw))
	for i, r := range raw {
		p := PixelFromInt32(r)
		buffer[i] = &p
	}
	ib := ImageBuffer{pixels: buffer, width: 5, height: 5}

	ib.drawSquare(1, 1, 3, 3, PixelFromInt32(0xFF0000))

	for idx, r := range raw {
		var expected int32
		if r == 1 {
			expected = 0xFF0000
		} else {
			expected = 0x000000
		}
		actual := ib.pixels[idx]
		if actual.Hex() != expected {
			t.Errorf("at %d, wanted %x, got %x (%#v)",
				idx, expected, actual.Hex(), actual)
		}
	}
}

func Test_drawSquare_Partial(t *testing.T) {
	raw := []int32{
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 0, 0,
		0, 0, 0, 1, 1,
		0, 0, 0, 1, 1,
	}
	buffer := make([]*Pixel, len(raw))
	for i, r := range raw {
		p := PixelFromInt32(r)
		buffer[i] = &p
	}
	ib := ImageBuffer{pixels: buffer, width: 5, height: 5}

	ib.drawSquare(3, 3, 3, 3, PixelFromInt32(0xFF0000))

	for idx, r := range raw {
		var expected int32
		if r == 1 {
			expected = 0xFF0000
		} else {
			expected = 0x000000
		}
		actual := ib.pixels[idx]
		if actual.Hex() != expected {
			t.Errorf("at %d, wanted %x, got %x (%#v)",
				idx, expected, actual.Hex(), actual)
		}
	}
}

func Test_FromJPEG(t *testing.T) {
	bfr := FromJPEG(filepath.Join("..", "testdata", "red.jpg"))

	if bfr.width != 50 {
		t.Errorf("unexpected width: %d", bfr.width)
	}
	if bfr.height != 50 {
		t.Errorf("unexpected height: %d", bfr.height)
	}

	if len(bfr.pixels) != 2500 {
		t.Errorf("unexpected number of pixels: %d",
			len(bfr.pixels))
	}

	for _, p := range bfr.pixels {
		if p.Hex() != 0xFE0000 {
			t.Errorf("unexpected pixel color: %x", p.Hex())
		}
	}
}

func Test_Palette(t *testing.T) {
	bfr := FromJPEG(filepath.Join("..", "testdata", "sample.jpg"))
	p := bfr.Palette(8)
	if len(p) != 8 {
		t.Errorf("unexpected palette length: %d", len(p))
	}
}
*/
