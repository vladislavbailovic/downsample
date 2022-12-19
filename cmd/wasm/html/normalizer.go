package html

import (
	"downsample/pkg"
	"fmt"
	"strconv"
	"syscall/js"
)

var normalizers map[pkg.NormalizerType]string = map[pkg.NormalizerType]string{
	pkg.NormalizerNormal:  "Normal",
	pkg.NormalizerAverage: "Average",
}

var quantizers map[pkg.QuantizerType]string = map[pkg.QuantizerType]string{
	pkg.QuantizerRGB:       "RGB",
	pkg.QuantizerShiftRGB:  "RGB (Shift)",
	pkg.QuantizerGray:      "Grayscale",
	pkg.QuantizerShiftGray: "Grayscale (Shift)",
}

type Normalizer struct {
	normKind  pkg.NormalizerType
	quantKind pkg.QuantizerType
	factor    int

	wrapper htmlElement
	norm    htmlElement
	quant   htmlElement
	fct     htmlElement
}

func NewNormalizer() *Normalizer {
	return &Normalizer{
		factor: 5,
		wrapper: htmlElement{
			tag:     tagName("div"),
			classes: []attributeValue{"normalizer"},
		},
		norm:  htmlElement{tag: tagName("select")},
		quant: htmlElement{tag: tagName("select")},
		fct: htmlElement{
			tag: tagName("input"),
			params: map[attributeName]attributeValue{
				"type": "number",
				"min":  "1",
				"max":  "255",
			},
		},
	}
}

func (x *Normalizer) Create(document js.Value) js.Value {
	w := x.wrapper.Create(document)

	norm := x.createNormalizer(document)
	w.Call("append", norm)

	quant := x.createQuantizer(document)
	w.Call("append", quant)

	fct := x.createFactor(document)
	w.Call("append", fct)

	return w
}

func (x *Normalizer) createNormalizer(document js.Value) js.Value {
	x.norm.Listen("change", func() bool {
		raw := x.norm.ref.Get("value").String()
		if kind, err := strconv.Atoi(raw); err == nil {
			x.normKind = pkg.NormalizerType(kind)
			fireEvent("downsample:ui", document)
		} else {
			fmt.Println(fmt.Sprintf("could not covert to kind: %s (%q)", err, raw))
		}
		return true
	})
	n := x.norm.Create(document)
	for kind, name := range normalizers {
		opts := map[attributeName]attributeValue{
			"value": attributeValue(fmt.Sprintf("%d", kind)),
		}
		if kind == x.normKind {
			opts["selected"] = "selected"
		}
		el := htmlElement{
			tag:    tagName("option"),
			params: opts,
			text:   innerText(name),
		}
		n.Call("append", el.Create(document))
	}
	return n
}

func (x *Normalizer) createQuantizer(document js.Value) js.Value {
	x.quant.Listen("change", func() bool {
		raw := x.quant.ref.Get("value").String()
		if kind, err := strconv.Atoi(raw); err == nil {
			x.quantKind = pkg.QuantizerType(kind)
			fireEvent("downsample:ui", document)
		} else {
			fmt.Println(fmt.Sprintf("could not covert to kind: %s (%q)", err, raw))
		}
		return true
	})
	n := x.quant.Create(document)
	for kind, name := range quantizers {
		opts := map[attributeName]attributeValue{
			"value": attributeValue(fmt.Sprintf("%d", kind)),
		}
		if kind == x.quantKind {
			opts["selected"] = "selected"
		}
		el := htmlElement{
			tag:    tagName("option"),
			params: opts,
			text:   innerText(name),
		}
		n.Call("append", el.Create(document))
	}
	return n
}

func (x *Normalizer) createFactor(document js.Value) js.Value {
	x.fct.params[attributeName("value")] = attributeValue(
		fmt.Sprintf("%d", x.factor))
	x.fct.Listen("change", func() bool {
		raw := x.fct.ref.Get("value").String()
		if fct, err := strconv.Atoi(raw); err == nil {
			x.factor = fct
			fireEvent("downsample:ui", document)
		} else {
			fmt.Println(fmt.Sprintf("could not covert to factor: %s (%q)", err, raw))
		}
		return true
	})
	return x.fct.Create(document)
}

func (x *Normalizer) GetNormalizerType() pkg.NormalizerType {
	return x.normKind
}

func (x *Normalizer) GetQuantizerType() pkg.QuantizerType {
	return x.quantKind
}

func (x *Normalizer) GetFactor() byte {
	return byte(x.factor)
}

func (x *Normalizer) Show() {
	x.wrapper.Show()
}

func (x *Normalizer) Hide() {
	x.wrapper.Hide()
}
