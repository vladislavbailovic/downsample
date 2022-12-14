package main

import (
	"downsample/cmd/wasm/html"
	"downsample/pkg"
	"downsample/pkg/asciify"
	"fmt"
	"image"
	"image/color"
	"math/rand"
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
			r, g, b, _ := img.At(x, y).RGBA()
			pixels = append(pixels, uint8(r/256))
			pixels = append(pixels, uint8(g/256))
			pixels = append(pixels, uint8(b/256))
			pixels = append(pixels, 0xFF)
		}
	}
	source := js.Global().Get("Uint8ClampedArray").New(
		len(pixels))
	js.CopyBytesToJS(source, pixels)
	data.Get("data").Call("set", source)

	otx.Call("putImageData", data, 0, 0)
}

func renderAscii(output string, doc js.Value) {
	imgOut := doc.Call("getElementById", html.OutputElementID.String())
	imgOut.Set("width", 0)
	imgOut.Set("height", 0)

	asciiOut := doc.Call("getElementById", html.AsciiElementID.String())
	asciiOut.Set("innerHTML", output)
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
	var replacement []asciify.Replacement

	palette := color.Palette{
		color.RGBA{R: 0xba, G: 0xda, B: 0x55, A: 0xff},
		color.RGBA{R: 0xde, G: 0xad, B: 0x00, A: 0xff},
		color.RGBA{R: 0x00, G: 0xde, B: 0xaf, A: 0xff},
	}
	algorithm := pkg.Pixelate
	var factor byte = 5
	quantizer = pkg.RGBQuantizer{Factor: factor}
	normalizer = pkg.StraightNormalizer{Q: quantizer}
	replacement = asciify.AsciiReplacements
	newSize := len(palette)

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
	rpl := html.NewReplacements()
	elements := []struct {
		src  creatable
		el   js.Value
		kind uiKind
	}{
		{src: algo, el: algo.Create(doc), kind: uiControl},
		{src: plt, el: plt.Create(doc), kind: uiControl},
		{src: tile, el: tile.Create(doc), kind: uiControl},
		{src: norm, el: norm.Create(doc), kind: uiControl},
		{src: rpl, el: rpl.Create(doc), kind: uiControl},

		{src: &html.Input, el: html.Input.Create(doc), kind: uiInput},
		{src: &html.Output, el: html.Output.Create(doc), kind: uiOutput},
		{src: &html.Ascii, el: html.Ascii.Create(doc), kind: uiOutput},
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
		case pkg.Average:
			b2 := pkg.ConstrainImage(img, palette, normalizer)
			renderAscii("", doc)
			renderImageBuffer(b2, doc)
		case pkg.Normalize:
			b2 := pkg.PixelateImage(img, pkg.ModeAndNormalize, normalizer)
			renderAscii("", doc)
			renderImageBuffer(b2, doc)
		case pkg.Pixelate:
			b2 := pkg.PixelateImage(img, pkg.ModePixelate, normalizer)
			renderAscii("", doc)
			renderImageBuffer(b2, doc)
		case pkg.Asciify:
			a := asciify.Asciifier{
				Replacements: replacement,
				Replacer:     &asciify.HtmlReplacer{},
				TileWidth:    pkg.GetTileSize(),
			}
			renderAscii(a.Asciify(img), doc)
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
		newSize = plt.GetNewSize()
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

		switch rpl.GetReplacementType() {
		case asciify.ReplacementAscii:
			replacement = asciify.AsciiReplacements
		case asciify.ReplacementUnicode:
			replacement = asciify.UnicodeReplacements
		}

		if algorithm != pkg.Average {
			plt.Hide()
		} else {
			plt.Show()
			update()
		}

		if algorithm != pkg.Asciify {
			norm.Show()
			rpl.Hide()
		} else {
			norm.Hide()
			rpl.Show()
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
	doc.Call("addEventListener", "downsample:palette:image", js.FuncOf(
		func(this js.Value, args []js.Value) interface{} {
			if img == nil {
				return false
			}
			newPalette := pkg.ImagePalette(img, uint8(newSize), nil)
			plt.ReplacePalette(newPalette)
			renderUI()
			return true
		},
	))
	doc.Call("addEventListener", "downsample:palette:random", js.FuncOf(
		func(this js.Value, args []js.Value) interface{} {
			newPalette := make(color.Palette, newSize)
			for i := 0; i < newSize; i++ {
				newPalette[i] = color.RGBA{
					R: uint8(rand.Float32() * 0xFF),
					G: uint8(rand.Float32() * 0xFF),
					B: uint8(rand.Float32() * 0xFF),
					A: 0xff,
				}
			}
			plt.ReplacePalette(newPalette)
			renderUI()
			return true
		},
	))

	update()
	renderUI()
}

type creatable interface {
	Create(js.Value) js.Value
}
