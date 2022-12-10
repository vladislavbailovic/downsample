package main

import (
	"downsample/pkg"
	"fmt"
	"strconv"
	"syscall/js"
)

type htmlTag string

func (x htmlTag) String() string {
	// TODO: validate
	return string(x)
}

type htmlInnerText string

func (x htmlInnerText) String() string {
	// TODO: validate
	return string(x)
}

type htmlAttributeName string

func (x htmlAttributeName) String() string {
	// TODO: validate
	return string(x)
}

type htmlAttributeValue string

func (x htmlAttributeValue) String() string {
	// TODO: validate
	return string(x)
}

type eventType string

func (x eventType) String() string {
	// TODO: validate
	return string(x)
}

type handlerCallback func() bool

type htmlElement struct {
	id       htmlAttributeValue
	classes  []htmlAttributeValue
	tag      htmlTag
	params   map[htmlAttributeName]htmlAttributeValue
	text     htmlInnerText
	handlers map[eventType]handlerCallback
	ref      js.Value
}

func (x *htmlElement) Create(document js.Value) js.Value {
	if !x.ref.IsUndefined() {
		x.ref.Call("remove")
	}
	el := document.Call("createElement", x.tag.String())

	if x.id != "" {
		el.Call("setAttribute", "id", x.id.String())
	}
	if len(x.classes) > 0 {
		for _, cls := range x.classes {
			el.Get("classList").Call("add", cls.String())
		}
	}
	if len(x.params) > 0 {
		for name, value := range x.params {
			el.Call("setAttribute", name.String(), value.String())
		}
	}
	if x.text != "" {
		el.Set("innerText", x.text.String())
	}

	if x.handlers != nil {
		for event, handler := range x.handlers {
			el.Call("addEventListener", event.String(), js.FuncOf(
				func(this js.Value, args []js.Value) interface{} {
					handler()
					return true
				},
			))
		}
	}
	x.ref = el

	return el
}

func (x *htmlElement) Listen(event eventType, handler handlerCallback) {
	if x.handlers == nil {
		x.handlers = map[eventType]handlerCallback{}
	}
	x.handlers[event] = handler
}

func (x *htmlElement) Trigger(event string) {
	if x.ref.IsUndefined() {
		return
	}
	fireEvent(eventType(event), x.ref)
}

func (x *htmlElement) Show() {
	if x.ref.IsUndefined() {
		return
	}
	x.ref.Get("style").Set("display", "block")
}

func (x *htmlElement) Hide() {
	if x.ref.IsUndefined() {
		return
	}
	x.ref.Get("style").Set("display", "none")
}

type paletteElement struct {
	palette pkg.Palette
	wrapper htmlElement
	colors  []htmlElement
	add     htmlElement
}

func NewPalette(palette pkg.Palette) *paletteElement {
	colors := make([]htmlElement, 0, len(palette))
	p := paletteElement{
		palette: palette,
		colors:  colors,
		wrapper: htmlElement{
			tag:     htmlTag("div"),
			classes: []htmlAttributeValue{"palette"},
		},
		add: htmlElement{
			tag:     htmlTag("button"),
			classes: []htmlAttributeValue{"add"},
			text:    htmlInnerText("Add"),
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
		px := pkg.PixelFromInt32(0x013120)
		x.palette = append(x.palette, px)
		fireEvent("downsample:ui", document)
		return true
	})
	a := x.add.Create(document)

	w.Call("append", a)
	return w
}

func (x *paletteElement) makeColorElement(color pkg.Pixel, document js.Value) js.Value {
	wrapper := htmlElement{
		tag:     htmlTag("div"),
		classes: []htmlAttributeValue{"color"},
	}
	control := htmlElement{
		tag:     htmlTag("div"),
		classes: []htmlAttributeValue{"control"},
	}
	input := htmlElement{
		tag:     htmlTag("input"),
		classes: []htmlAttributeValue{"color"},
		params: map[htmlAttributeName]htmlAttributeValue{
			"type":  "color",
			"value": htmlAttributeValue(fmt.Sprintf("#%06x", color.Hex())),
		},
	}
	remove := htmlElement{
		tag:  htmlTag("button"),
		text: htmlInnerText("x"),
	}

	w := wrapper.Create(document)
	c := control.Create(document)

	input.Listen("change", func() bool {
		cs := input.ref.Get("value").String()[1:]
		if clr, err := strconv.ParseInt(cs, 16, 32); err == nil {
			px := pkg.PixelFromInt32(int32(clr))
			for idx, clr := range x.palette {
				if clr.Hex() == color.Hex() {
					x.palette[idx] = px
				}
			}
			fireEvent("downsample:ui", document)
		} else {
			fmt.Println("error parsing new color", cs, err)
		}
		return true
	})
	i := input.Create(document)

	remove.Listen("click", func() bool {
		plt := make([]pkg.Pixel, 0, len(x.palette)-1)
		for _, px := range x.palette {
			if px.Hex() == color.Hex() {
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

func (x *paletteElement) Hide() {
	x.wrapper.Hide()
}
func (x *paletteElement) Show() {
	x.wrapper.ref.Get("style").Set("display", "flex")
}

func (x *paletteElement) GetPalette() pkg.Palette {
	return x.palette
}

type tileElement struct {
	size    int
	wrapper htmlElement
	input   htmlElement
}

func NewTileSize(size int) *tileElement {
	return &tileElement{
		size: size,
		wrapper: htmlElement{
			tag:     htmlTag("label"),
			classes: []htmlAttributeValue{"tile-size"},
		},
		input: htmlElement{
			tag: htmlTag("input"),
			params: map[htmlAttributeName]htmlAttributeValue{
				"type": "number",
				"min":  "2",
				"max":  "100",
			},
		},
	}
}

func (x *tileElement) Create(document js.Value) js.Value {
	w := x.wrapper.Create(document)

	x.input.params[htmlAttributeName("value")] = htmlAttributeValue(
		fmt.Sprintf("%d", x.size))
	x.input.Listen("change", func() bool {
		if s, err := strconv.Atoi(x.input.ref.Get("value").String()); err != nil {
			fmt.Println("unable to set new size!", x.input.ref.Get("value"), err)
		} else {
			x.size = s
			fireEvent("downsample:ui", document)
		}
		return true
	})
	t := x.input.Create(document)

	w.Call("append", t)
	return w
}

type algoElement struct {
	algorithm string
	wrapper   htmlElement
}

func NewAlgo(algorithm string) *algoElement {
	return &algoElement{
		algorithm: algorithm,
		wrapper: htmlElement{
			tag: htmlTag("select"),
			id:  htmlAttributeValue("algo"),
		},
	}
}

func (x *algoElement) Create(document js.Value) js.Value {
	x.wrapper.Listen("change", func() bool {
		x.algorithm = x.wrapper.ref.Get("value").String()
		fireEvent("downsample:ui", document)
		return true
	})
	w := x.wrapper.Create(document)

	algos := []string{
		"pixelate",
		"normalize",
		"average",
	}
	for _, a := range algos {
		opts := map[htmlAttributeName]htmlAttributeValue{
			"value": htmlAttributeValue(a),
		}
		if a == x.algorithm {
			opts["selected"] = "selected"
		}
		el := htmlElement{
			tag:    htmlTag("option"),
			params: opts,
			text:   htmlInnerText(a),
		}
		w.Call("append", el.Create(document))
	}

	return w
}

func (x *algoElement) GetAlgorithm() string {
	return x.algorithm
}

func fireEvent(name eventType, document js.Value) {
	ev := js.Global().Get("Event").New(name.String())
	document.Call("dispatchEvent", ev)
}

var Input htmlElement = htmlElement{
	id:  htmlAttributeValue(elInput),
	tag: htmlTag("img"),
	params: map[htmlAttributeName]htmlAttributeValue{
		"src": "sample.jpg",
	},
}
var Output htmlElement = htmlElement{
	id:  htmlAttributeValue(elOutput),
	tag: htmlTag("canvas"),
}
