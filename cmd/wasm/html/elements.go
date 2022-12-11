package html

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
			tag:     tagName("label"),
			classes: []attributeValue{"tile-size"},
		},
		input: htmlElement{
			tag: tagName("input"),
			params: map[attributeName]attributeValue{
				"type": "number",
				"min":  "2",
				"max":  "100",
			},
		},
	}
}

func (x *tileElement) Create(document js.Value) js.Value {
	w := x.wrapper.Create(document)

	x.input.params[attributeName("value")] = attributeValue(
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

func (x *tileElement) GetSize() int {
	return x.size
}

type algoElement struct {
	algorithm string
	wrapper   htmlElement
}

func NewAlgo(algorithm string) *algoElement {
	return &algoElement{
		algorithm: algorithm,
		wrapper: htmlElement{
			tag: tagName("select"),
			id:  attributeValue("algo"),
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
		opts := map[attributeName]attributeValue{
			"value": attributeValue(a),
		}
		if a == x.algorithm {
			opts["selected"] = "selected"
		}
		el := htmlElement{
			tag:    tagName("option"),
			params: opts,
			text:   innerText(a),
		}
		w.Call("append", el.Create(document))
	}

	return w
}

func (x *algoElement) GetAlgorithm() string {
	return x.algorithm
}

var Input htmlElement = htmlElement{
	id:  attributeValue(InputElementID),
	tag: tagName("img"),
	params: map[attributeName]attributeValue{
		"src": "sample.jpg",
	},
}
var Output htmlElement = htmlElement{
	id:  attributeValue(OutputElementID),
	tag: tagName("canvas"),
}

var Root htmlElement = htmlElement{
	tag:     tagName("div"),
	classes: []attributeValue{"interface"},
}
var Controls htmlElement = htmlElement{
	tag:     tagName("div"),
	classes: []attributeValue{"controls"},
}
var Io htmlElement = htmlElement{
	tag:     tagName("div"),
	classes: []attributeValue{"io"},
}
