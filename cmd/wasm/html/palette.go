package html

import (
	"fmt"
	"image/color"
	"strconv"
	"syscall/js"
)

type paletteElement struct {
	palette color.Palette
	newSize int
	wrapper htmlElement
	colors  []htmlElement
	add     htmlElement
	count   htmlElement
}

func NewPalette(palette color.Palette) *paletteElement {
	colors := make([]htmlElement, 0, len(palette))
	p := paletteElement{
		palette: palette,
		colors:  colors,
		newSize: len(palette),
		wrapper: htmlElement{
			tag:     tagName("div"),
			classes: []attributeValue{"palette"},
		},
		add: htmlElement{
			tag:     tagName("button"),
			classes: []attributeValue{"add"},
			text:    innerText("Add"),
		},
		count: htmlElement{
			tag: tagName("input"),
			params: map[attributeName]attributeValue{
				"type": "number",
				"min":  "2",
				"max":  "16",
			},
		},
	}
	return &p
}

func (x *paletteElement) Create(document js.Value) js.Value {
	w := x.wrapper.Create(document)
	for _, color := range x.palette {
		el := x.makeColorElement(color, document)
		w.Call("append", el)
	}
	x.add.Listen("click", func() bool {
		px := color.RGBA{R: 0x01, G: 0x31, B: 0x20, A: 0xFF}
		x.palette = append(x.palette, px)
		fireEvent("downsample:ui", document)
		return true
	})
	a := x.add.Create(document)
	w.Call("append", a)

	ops := x.createPaletteOpsElement(document)
	w.Call("append", ops)

	return w
}

func (x *paletteElement) makeColorElement(clr color.Color, document js.Value) js.Value {
	cr, cg, cb, _ := clr.RGBA()
	wrapper := htmlElement{
		tag:     tagName("div"),
		classes: []attributeValue{"color"},
	}
	control := htmlElement{
		tag:     tagName("div"),
		classes: []attributeValue{"control"},
	}
	input := htmlElement{
		tag:     tagName("input"),
		classes: []attributeValue{"color"},
		params: map[attributeName]attributeValue{
			"type":  "color",
			"value": attributeValue(fmt.Sprintf("#%02x%02x%02x", uint8(cr/256), uint8(cg/256), uint8(cb/256))),
		},
	}
	remove := htmlElement{
		tag:  tagName("button"),
		text: innerText("x"),
	}

	w := wrapper.Create(document)
	c := control.Create(document)

	input.Listen("change", func() bool {
		cs := input.ref.Get("value").String()[1:]
		rs, err := strconv.ParseInt(cs[0:2], 16, 16)
		if err != nil {
			fmt.Println(fmt.Sprintf("unable to parse int from %s: %s (%s)", cs[0:2], err, cs))
			return true
		}
		gs, err := strconv.ParseInt(cs[2:4], 16, 16)
		if err != nil {
			fmt.Println(fmt.Sprintf("unable to parse int from %s: %s (%s)", cs[2:4], err, cs))
			return true
		}
		bs, err := strconv.ParseInt(cs[4:6], 16, 16)
		if err != nil {
			fmt.Println(fmt.Sprintf("unable to parse int from %s: %s (%s)", cs[4:6], err, cs))
			return true
		}
		px := color.RGBA{R: uint8(rs), G: uint8(gs), B: uint8(bs), A: 0xFF}
		for idx, c := range x.palette {
			if c == clr {
				x.palette[idx] = px
				break // Just once
			}
		}
		fireEvent("downsample:ui", document)
		return true
	})
	i := input.Create(document)

	remove.Listen("click", func() bool {
		plt := make(color.Palette, 0, len(x.palette)-1)
		for _, px := range x.palette {
			if px == clr {
				continue
			}
			plt = append(plt, px)
		}
		x.palette = plt
		fireEvent("downsample:ui", document)
		return true
	})
	r := remove.Create(document)

	c.Call("append", i)
	c.Call("append", r)

	w.Call("append", c)

	return w
}

func (x *paletteElement) createPaletteOpsElement(document js.Value) js.Value {
	fmt.Println("creating ops element")
	wrapper := htmlElement{
		tag:     tagName("div"),
		classes: []attributeValue{"operations"},
	}
	w := wrapper.Create(document)

	x.count.params[attributeName("value")] = attributeValue(
		fmt.Sprintf("%d", x.newSize))
	x.count.Listen("change", func() bool {
		raw := x.count.ref.Get("value").String()
		if count, err := strconv.Atoi(raw); err == nil {
			x.newSize = count
			fireEvent("downsample:ui", document)
		}
		return true
	})
	w.Call("append", x.count.Create(document))

	load := htmlElement{
		tag:  tagName("button"),
		text: innerText("Load from image"),
	}
	load.Listen("click", func() bool {
		fireEvent("downsample:palette:image", document)
		return true
	})
	w.Call("append", load.Create(document))

	return w
}

func (x *paletteElement) Hide() {
	x.wrapper.Hide()
}
func (x *paletteElement) Show() {
	x.wrapper.ref.Get("style").Set("display", "flex")
}

func (x *paletteElement) ReplacePalette(palette color.Palette) {
	x.palette = palette
	x.newSize = len(palette)
}

func (x *paletteElement) GetPalette() color.Palette {
	return x.palette
}

func (x *paletteElement) GetNewSize() int {
	return x.newSize
}
