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
	fireEvent(event, x.ref)
}

type colorElement struct {
	color   pkg.Pixel
	wrapper htmlElement
	control htmlElement
	input   htmlElement
	remove  htmlElement
}

func NewColor(color pkg.Pixel) colorElement {
	return colorElement{
		color: color,
		wrapper: htmlElement{
			tag:     htmlTag("div"),
			classes: []htmlAttributeValue{"color"},
		},
		control: htmlElement{
			tag:     htmlTag("div"),
			classes: []htmlAttributeValue{"control"},
		},
		input: htmlElement{
			tag:     htmlTag("input"),
			classes: []htmlAttributeValue{"color"},
			params: map[htmlAttributeName]htmlAttributeValue{
				"type":  "color",
				"value": htmlAttributeValue(fmt.Sprintf("#%06x", color.Hex())),
			},
		},
		remove: htmlElement{
			tag:  htmlTag("button"),
			text: htmlInnerText("x"),
		},
	}
}

func (x colorElement) Create(document js.Value) js.Value {
	w := x.wrapper.Create(document)
	c := x.control.Create(document)

	x.input.Listen("change", func() bool {
		fireEvent("downsample:render", document)
		return true
	})
	i := x.input.Create(document)

	x.remove.Listen("click", func() bool {
		x.wrapper.Trigger("color:remove")
		fireEvent("downsample:render", document)
		return true
	})
	r := x.remove.Create(document)

	c.Call("append", i)
	c.Call("append", r)

	w.Call("append", c)

	return w
}

type paletteElement struct {
	palette pkg.Palette
	wrapper htmlElement
	add     htmlElement
}

func NewPalette(palette pkg.Palette) *paletteElement {
	p := paletteElement{
		palette: palette,
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
	for _, c := range x.palette {
		color := NewColor(c)
		color.wrapper.Listen("color:remove", func() bool {
			x.removeColor(color.color)
			fireEvent("downsample:ui", document)
			return true
		})
		el := color.Create(document)
		w.Call("append", el)
	}
	x.add.Listen("click", func() bool {
		x.addColor(0xbada55)
		fireEvent("downsample:ui", document)
		return true
	})
	a := x.add.Create(document)

	w.Call("append", a)
	return w
}

func (x *paletteElement) addColor(clr int32) {
	x.palette = append(x.palette, pkg.PixelFromInt32(clr))
}

func (x *paletteElement) removeColor(clr pkg.Pixel) {
	newPalette := make([]pkg.Pixel, 0, len(x.palette)-1)
	for _, c := range x.palette {
		if c.Hex() == clr.Hex() {
			continue
		}
		newPalette = append(newPalette, c)
	}
	x.palette = newPalette
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
			tag: htmlTag("label"),
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

	t := x.input.Create(document)

	t.Call("addEventListener", "change", js.FuncOf(
		func(this js.Value, args []js.Value) interface{} {
			if s, err := strconv.Atoi(t.Get("value").String()); err != nil {
				fmt.Println("unable to set new size!", t.Get("value"), err)
			} else {
				x.size = s
				fireEvent("downsample:ui", document)
				fireEvent("downsample:render", document)
			}
			return true
		},
	))

	w.Call("append", t)
	return w
}

func fireEvent(name string, document js.Value) {
	ev := js.Global().Get("Event").New(name)
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
