package main

import (
	"downsample/pkg"
	"fmt"
	"strconv"
	"syscall/js"
)

const (
	elInput        string = "input-file"
	elOutput       string = "output"
	elAlgo         string = "algo"
	elPalette      string = ".palette"
	elTileSize     string = ".tile-size input"
	elAddColor     string = ".palette .add"
	elsColor       string = ".color"
	elsColorValues string = ".palette input"
)

func getSource(doc js.Value) *pkg.ImageBuffer {
	data := doc.Call("createElement", "canvas")
	ctx := data.Call("getContext", "2d")

	input := doc.Call("getElementById", elInput)
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
	out := doc.Call("getElementById", elOutput)
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

func getPalette(doc js.Value) pkg.Palette {
	raw := doc.Call("querySelectorAll", elsColorValues)
	p := make([]pkg.Pixel, 0, raw.Length())
	for i := 0; i < raw.Length(); i++ {
		color := raw.Index(i).Get("value").String()[1:]
		if intval, err := strconv.ParseInt(color, 16, 32); err == nil {
			px := pkg.PixelFromInt32(int32(intval))
			p = append(p, px)

		}
	}
	return p
}

func render(algo string, doc js.Value) {
	img := getSource(doc)
	switch algo {
	case "average":
		palette := getPalette(doc)
		fmt.Println(palette)
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
	updateInterface(algo, doc)
}

func updateInterface(algo string, doc js.Value) {
	tile := doc.Call("querySelector", elTileSize)
	tile.Set("value", pkg.GetTileSize())

	palette := doc.Call("querySelector", elPalette)
	if algo != "average" {
		palette.Get("style").Set("display", "none")
		return
	}
	palette.Get("style").Set("display", "flex")
}

func rerender(doc js.Value) {
	algo := doc.Call("getElementById", elAlgo)
	render(algo.Get("value").String(), doc)
}

func initGui() {
	doc := js.Global().Get("document")

	doc.Call("addEventListener", "change", js.FuncOf(
		func(this js.Value, args []js.Value) any {
			tgt := args[0].Get("target")
			if tgt.Get("nodeName").String() != "INPUT" {
				return false
			}

			closest := tgt.Call("closest", elsColor)
			if !closest.Truthy() {
				return false
			}

			rerender(doc)

			return true
		}),
	)

	doc.Call("addEventListener", "click", js.FuncOf(
		func(this js.Value, args []js.Value) any {
			tgt := args[0].Get("target")
			if tgt.Get("nodeName").String() != "BUTTON" {
				return false
			}

			closest := tgt.Call("closest", elsColor)
			if !closest.Truthy() {
				return false
			}

			closest.Call("remove")
			rerender(doc)

			return true
		}),
	)

	algo := doc.Call("getElementById", elAlgo)
	algo.Call("addEventListener", "change", js.FuncOf(
		func(this js.Value, args []js.Value) any {
			rerender(doc)
			return true
		},
	))

	tile := doc.Call("querySelector", elTileSize)
	tile.Call("addEventListener", "change", js.FuncOf(
		func(this js.Value, args []js.Value) any {
			tileSize := 0
			raw := tile.Get("value").String()
			if ts, err := strconv.Atoi(raw); err == nil {
				tileSize = ts
			}
			if tileSize == 0 {
				return false
			}
			pkg.SetTileSize(tileSize)
			rerender(doc)
			return true
		},
	))

	add := doc.Call("querySelector", elAddColor)
	add.Call("addEventListener", "click", js.FuncOf(
		func(this js.Value, args []js.Value) any {
			palette := doc.Call("querySelector", elPalette)
			clr := palette.Call("querySelector", elsColor).
				Call("cloneNode", true)
			add.Call("before", clr)
			rerender(doc)
			return true
		},
	))
}
