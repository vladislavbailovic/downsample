package main

import (
	"downsample/pkg"
	"fmt"
	"syscall/js"
)

const (
	elInput  htmlAttributeValue = "input-file"
	elOutput htmlAttributeValue = "output"
)

func getSource(doc js.Value) *pkg.ImageBuffer {
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
	buffer := make([]*pkg.Pixel, 0, raw.Length()/4)
	idx := 0
	for idx < source.Length() {
		r := uint8(source.Index(idx).Int())
		idx++
		g := uint8(source.Index(idx).Int())
		idx++
		b := uint8(source.Index(idx).Int())
		idx++
		idx++ // A
		px := pkg.NewPixel(r, g, b)
		buffer = append(buffer, &px)
	}
	img := pkg.NewImageBuffer(width, height, buffer)

	return img
}

func renderImageBuffer(img *pkg.ImageBuffer, doc js.Value) {
	out := doc.Call("getElementById", elOutput.String())
	out.Set("width", img.Width())
	out.Set("height", img.Height())

	otx := out.Call("getContext", "2d")
	data := otx.Call("createImageData", img.Width(), img.Height())

	pixels := make([]byte, 0, len(img.Pixels())*4)
	for _, px := range img.Pixels() {
		pixels = append(pixels, px.R)
		pixels = append(pixels, px.G)
		pixels = append(pixels, px.B)
		pixels = append(pixels, 0xff)
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
	var img *pkg.ImageBuffer

	palette := pkg.Palette{
		pkg.PixelFromInt32(0xbada55),
		pkg.PixelFromInt32(0x0dead0),
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
			b2 := pkg.ConstrainImage(img, palette)
			renderImageBuffer(b2, doc)
		case "normalize":
			b2 := pkg.PixelateImage(img, pkg.ModeAndNormalize)
			renderImageBuffer(b2, doc)
		case "pixelate":
			b2 := pkg.PixelateImage(img, pkg.ModePixelate)
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
