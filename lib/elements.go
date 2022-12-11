package main

import (
	"fmt"
	"strconv"
	"syscall/js"
)

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
