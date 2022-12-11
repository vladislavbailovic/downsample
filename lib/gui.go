package main

import (
	"downsample/pkg"
	"fmt"
	"image"
	"image/color"
	"syscall/js"
)

const (
	elInput  htmlAttributeValue = "input-file"
	elOutput htmlAttributeValue = "output"
)

func getSource(doc js.Value) image.Image {
	data := doc.Call("createElement", "canvas")
	ctx := data.Call("getContext", "2d")

	input := doc.Call("getElementById", elInput.String())
	width := input.Get("width").Int()
	height := input.Get("height").Int()

	data.Set("width", width)
	data.Set("height", height)
	ctx.Call("drawImage", input, 0, 0)

	raw := ctx.Call("getImageData", 0, 0, width, height)
	source := raw.Get("data")
	buffer := image.NewRGBA(image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{width, height},
	})
	idx := 0
	for idx < source.Length() {
		pos := idx
		r := uint8(source.Index(idx).Int())
		idx++
		g := uint8(source.Index(idx).Int())
		idx++
		b := uint8(source.Index(idx).Int())
		idx++
		a := uint8(source.Index(idx).Int())
		idx++ // A
		px := color.RGBA{R: r, G: g, B: b, A: a}

		y := (pos / 4) / width
		x := (pos / 4) % width
		buffer.Set(x, y, px)
	}

	return buffer
}

func renderImageBuffer(img image.Image, doc js.Value) {
	out := doc.Call("getElementById", elOutput.String())
	bounds := img.Bounds()
	out.Set("width", bounds.Max.X)
	out.Set("height", bounds.Max.Y)

	otx := out.Call("getContext", "2d")
	data := otx.Call("createImageData", bounds.Max.X, bounds.Max.Y)

	pixels := make([]byte, 0, bounds.Max.X*bounds.Max.Y*4)
	for y := 0; y < bounds.Max.Y; y++ {
		for x := 0; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixels = append(pixels, uint8(r/256))
			pixels = append(pixels, uint8(g/256))
			pixels = append(pixels, uint8(b/256))
			pixels = append(pixels, uint8(a/256))
		}
	}
	source := js.Global().Get("Uint8ClampedArray").New(
		len(pixels))
	js.CopyBytesToJS(source, pixels)
	data.Get("data").Call("set", source)

	otx.Call("putImageData", data, 0, 0)
}

type uiKind byte

const (
	uiStructure uiKind = iota
	uiControl   uiKind = iota
	uiInput     uiKind = iota
	uiOutput    uiKind = iota
)

func initGui() {
	doc := js.Global().Get("document")
	var img image.Image

	palette := color.Palette{
		color.RGBA{R: 0xba, G: 0xda, B: 0x55, A: 0xff},
		color.RGBA{R: 0x0d, G: 0xea, B: 0xd0, A: 0xff},
	}
	algorithm := "pixelate"

	root := Root.Create(doc)
	controls := Controls.Create(doc)
	io := Io.Create(doc)
	root.Call("append", controls)
	root.Call("append", io)
	doc.Call("querySelector", "body>div").Call("replaceWith", root)

	algo := NewAlgo(algorithm)
	plt := NewPalette(palette)
	tile := NewTileSize(pkg.GetTileSize())
	elements := []struct {
		src  creatable
		el   js.Value
		kind uiKind
	}{
		{src: algo, el: algo.Create(doc), kind: uiControl},
		{src: plt, el: plt.Create(doc), kind: uiControl},
		{src: tile, el: tile.Create(doc), kind: uiControl},

		{src: &Input, el: Input.Create(doc), kind: uiInput},
		{src: &Output, el: Output.Create(doc), kind: uiOutput},
	}

	update := func() {
		for _, item := range elements {
			destination := io
			if item.kind == uiControl {
				destination = controls
			} else if img != nil {
				continue
			}
			item.el = item.src.Create(doc)
			destination.Call("append", item.el)
		}
	}

	render := func() {
		if img == nil {
			return
		}
		switch algorithm {
		case "average":
			b2 := pkg.ConstrainImage(img, palette, nil)
			renderImageBuffer(b2, doc)
		case "normalize":
			b2 := pkg.PixelateImage(img, pkg.ModeAndNormalize, nil)
			renderImageBuffer(b2, doc)
		case "pixelate":
			b2 := pkg.PixelateImage(img, pkg.ModePixelate, nil)
			renderImageBuffer(b2, doc)
		default:
			fmt.Println("ignoring the unknown algo", algo)
			renderImageBuffer(img, doc)
		}
	}

	Input.Listen("load", func() bool {
		img = getSource(doc)
		render()
		return true
	})

	doc.Call("addEventListener", "downsample:ui", js.FuncOf(
		func(this js.Value, args []js.Value) interface{} {
			algorithm = algo.GetAlgorithm()
			pkg.SetTileSize(tile.size)
			palette = plt.GetPalette()

			if algorithm != "average" {
				plt.Hide()
			} else {
				plt.Show()
				update()
			}

			render()
			return true
		},
	))
	doc.Call("addEventListener", "downsample:render", js.FuncOf(
		func(this js.Value, args []js.Value) interface{} {
			render()
			return true
		},
	))

	update()
	fireEvent("downsample:ui", doc)
}

type creatable interface {
	Create(js.Value) js.Value
}
