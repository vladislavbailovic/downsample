package main

import (
	"downsample/pkg"
	"fmt"
	"syscall/js"
)

func jsToImageBuffer(args []js.Value, done chan bool) *pkg.ImageBuffer {
	if len(args) < 3 {
		fmt.Println("Missing expected argument(s): wanted 3, got", len(args))
		done <- true
		return nil
	}

	raw := args[0]
	width := args[1]
	height := args[2]
	fmt.Println(raw.Length() % 4)
	if raw.Length()%4 != 0 {
		fmt.Println("NOT divisible by 4")
		done <- true
		return nil
	}

	buffer := make([]*pkg.Pixel, 0, raw.Length()/4)
	idx := 0
	for idx < raw.Length() {
		r := uint8(raw.Index(idx).Int())
		idx++
		g := uint8(raw.Index(idx).Int())
		idx++
		b := uint8(raw.Index(idx).Int())
		idx++
		idx++ // A
		px := pkg.NewPixel(r, g, b)
		buffer = append(buffer, &px)
	}

	img := pkg.NewImageBuffer(width.Int(), height.Int(), buffer)
	return img
}

func imageBufferToJs(pixels []*pkg.Pixel) any {
	result := make([]byte, 0, len(pixels)*4)
	for _, p := range pixels {
		result = append(result, p.R)
		result = append(result, p.G)
		result = append(result, p.B)
		result = append(result, 0xFF)
	}

	data := js.Global().Get("Uint8ClampedArray").New(len(result))
	js.CopyBytesToJS(data, result)
	return data
}

func pixelateWrapper(done chan bool) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		img := jsToImageBuffer(args, done)
		b2 := pkg.PixelateImage(img, pkg.ModePixelate).Pixels()
		return imageBufferToJs(b2)
	})
}

func normalizeWrapper(done chan bool) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		img := jsToImageBuffer(args, done)
		b2 := pkg.PixelateImage(img, pkg.ModeAndNormalize).Pixels()
		return imageBufferToJs(b2)
	})
}

func averageWrapper(done chan bool) js.Func {
	palette := []pkg.Pixel{
		pkg.PixelFromInt32(0xffb703),
		pkg.PixelFromInt32(0xfb8500),
		pkg.PixelFromInt32(0xd00000),
		pkg.PixelFromInt32(0x8ecae6),
		pkg.PixelFromInt32(0x023047),
		pkg.PixelFromInt32(0x124057),
		pkg.PixelFromInt32(0x225068),
		pkg.PixelFromInt32(0x219ebc),
		pkg.PixelFromInt32(0x2a9d8f),
		pkg.PixelFromInt32(0xccc5b9),
	}
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		img := jsToImageBuffer(args, done)
		b2 := pkg.ConstrainImage(img, palette).Pixels()
		return imageBufferToJs(b2)
	})
}

func main() {
	done := make(chan bool)
	js.Global().Set("pixelate", pixelateWrapper(done))
	js.Global().Set("normalize", normalizeWrapper(done))
	js.Global().Set("average", averageWrapper(done))
	for {
		select {
		case <-done:
			return
		}
	}
}
