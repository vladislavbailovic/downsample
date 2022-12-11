package main

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
	"syscall/js"
	"unicode"
)

func noSpecialChars(r rune) rune {
	if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) {
		return r
	}
	if '-' == r || '#' == r || ':' == r || '.' == r {
		return r
	}
	fmt.Println(fmt.Sprintf("invalid: %d %c", r, r))
	return rune(-1)
}

type htmlTag string

func (x htmlTag) String() string {
	switch x {
	case "button", "label", "input", "select", "option", "canvas", "img":
		return string(x)
	default:
		return "div"
	}
}

type htmlInnerText string

func (x htmlInnerText) String() string {
	return strings.Map(noSpecialChars, string(x))
}

type htmlAttributeName string

func (x htmlAttributeName) String() string {
	return strings.Map(noSpecialChars, string(x))
}

type htmlAttributeValue string

func (x htmlAttributeValue) String() string {
	return strings.Map(noSpecialChars, string(x))
}

type eventType string

func (x eventType) String() string {
	return strings.Map(noSpecialChars, string(x))
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
	palette color.Palette
	wrapper htmlElement
	colors  []htmlElement
	add     htmlElement
}

func NewPalette(palette color.Palette) *paletteElement {
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
		px := color.RGBA{R: 0x01, G: 0x31, B: 0x20, A: 0xFF}
		x.palette = append(x.palette, px)
		fireEvent("downsample:ui", document)
		return true
	})
	a := x.add.Create(document)

	w.Call("append", a)
	return w
}

func (x *paletteElement) makeColorElement(clr color.Color, document js.Value) js.Value {
	cr, cg, cb, _ := clr.RGBA()
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
			"value": htmlAttributeValue(fmt.Sprintf("#%02x%02x%02x", uint8(cr/256), uint8(cg/256), uint8(cb/256))),
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
		rs, err := strconv.ParseInt(cs[0:2], 16, 8)
		if err != nil {
			return true
		}
		gs, err := strconv.ParseInt(cs[2:4], 16, 8)
		if err != nil {
			return true
		}
		bs, err := strconv.ParseInt(cs[4:6], 16, 8)
		if err != nil {
			return true
		}
		px := color.RGBA{R: uint8(rs), G: uint8(gs), B: uint8(bs), A: 0xFF}
		for idx, c := range x.palette {
			if c == clr {
				x.palette[idx] = px
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

func (x *paletteElement) Hide() {
	x.wrapper.Hide()
}
func (x *paletteElement) Show() {
	x.wrapper.ref.Get("style").Set("display", "flex")
}

func (x *paletteElement) GetPalette() color.Palette {
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

var Root htmlElement = htmlElement{
	tag:     htmlTag("div"),
	classes: []htmlAttributeValue{"interface"},
}
var Controls htmlElement = htmlElement{
	tag:     htmlTag("div"),
	classes: []htmlAttributeValue{"controls"},
}
var Io htmlElement = htmlElement{
	tag:     htmlTag("div"),
	classes: []htmlAttributeValue{"io"},
}
