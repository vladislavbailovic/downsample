package main

import (
	"downsample/lib/html"
	"downsample/pkg"
	"fmt"
	"image"
	"image/color"
	"syscall/js"
)

func getSource(doc js.Value) image.Image {
	data := doc.Call("createElement", "canvas")
	ctx := data.Call("getContext", "2d")

	input := doc.Call("getElementById", html.InputElementID.String())
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
	out := doc.Call("getElementById", html.OutputElementID.String())
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
	var quantizer pkg.Quantizer
	var normalizer pkg.Normalizer

	palette := color.Palette{
		color.RGBA{R: 0xba, G: 0xda, B: 0x55, A: 0xff},
		color.RGBA{R: 0x0d, G: 0xea, B: 0xd0, A: 0xff},
	}
	algorithm := "pixelate"
	var factor byte = 5
	quantizer = pkg.RGBQuantizer{Factor: factor}
	normalizer = pkg.StraightNormalizer{Q: quantizer}

	root := html.Root.Create(doc)
	controls := html.Controls.Create(doc)
	io := html.Io.Create(doc)
	root.Call("append", controls)
	root.Call("append", io)
	doc.Call("querySelector", "body>div").Call("replaceWith", root)

	algo := html.NewAlgo(algorithm)
	plt := html.NewPalette(palette)
	tile := html.NewTileSize(pkg.GetTileSize())
	norm := html.NewNormalizer()
	elements := []struct {
		src  creatable
		el   js.Value
		kind uiKind
	}{
		{src: algo, el: algo.Create(doc), kind: uiControl},
		{src: plt, el: plt.Create(doc), kind: uiControl},
		{src: tile, el: tile.Create(doc), kind: uiControl},
		{src: norm, el: norm.Create(doc), kind: uiControl},

		{src: &html.Input, el: html.Input.Create(doc), kind: uiInput},
		{src: &html.Output, el: html.Output.Create(doc), kind: uiOutput},
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
			b2 := pkg.ConstrainImage(img, palette, normalizer)
			renderImageBuffer(b2, doc)
		case "normalize":
			b2 := pkg.PixelateImage(img, pkg.ModeAndNormalize, normalizer)
			renderImageBuffer(b2, doc)
		case "pixelate":
			b2 := pkg.PixelateImage(img, pkg.ModePixelate, normalizer)
			renderImageBuffer(b2, doc)
		default:
			fmt.Println("ignoring the unknown algo", algo)
			renderImageBuffer(img, doc)
		}
	}

	html.Input.Listen("load", func() bool {
		img = getSource(doc)
		render()
		return true
	})

	renderUI := func() {
		algorithm = algo.GetAlgorithm()
		pkg.SetTileSize(tile.GetSize())
		palette = plt.GetPalette()
		factor = norm.GetFactor()

		switch norm.GetQuantizerType() {
		case pkg.QuantizerRGB:
			quantizer = pkg.RGBQuantizer{Factor: factor}
		case pkg.QuantizerShiftRGB:
			quantizer = pkg.RGBShiftQuantizer{Factor: factor}
		case pkg.QuantizerGray:
			quantizer = pkg.GrayQuantizer{Factor: factor}
		case pkg.QuantizerShiftGray:
			quantizer = pkg.GrayShiftQuantizer{Factor: factor}
		}

		switch norm.GetNormalizerType() {
		case pkg.NormalizerNormal:
			normalizer = pkg.StraightNormalizer{Q: quantizer}
		case pkg.NormalizerAverage:
			normalizer = pkg.AverageNormalizer{Q: quantizer}
		}

		if algorithm != "average" {
			plt.Hide()
		} else {
			plt.Show()
			update()
		}

		render()
	}
	doc.Call("addEventListener", "downsample:ui", js.FuncOf(
		func(this js.Value, args []js.Value) interface{} {
			renderUI()
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
	renderUI()
}

type creatable interface {
	Create(js.Value) js.Value
}
