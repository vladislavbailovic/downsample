package html

import (
	"downsample/pkg/asciify"
	"fmt"
	"strconv"
	"syscall/js"
)

var replacements map[asciify.ReplacementsType]string = map[asciify.ReplacementsType]string{
	asciify.ReplacementAscii:   "ASCII",
	asciify.ReplacementUnicode: "Unicode",
}

type Replacements struct {
	replacement asciify.ReplacementsType

	wrapper htmlElement
	rpl     htmlElement
}

func NewReplacements() *Replacements {
	return &Replacements{
		wrapper: htmlElement{
			tag:     tagName("div"),
			classes: []attributeValue{"replacements"},
		},
		rpl: htmlElement{tag: tagName("select")},
	}
}

func (x *Replacements) Create(document js.Value) js.Value {
	w := x.wrapper.Create(document)

	rpl := x.createReplacements(document)
	w.Call("append", rpl)

	return w
}

func (x *Replacements) createReplacements(document js.Value) js.Value {
	x.rpl.Listen("change", func() bool {
		raw := x.rpl.ref.Get("value").String()
		if kind, err := strconv.Atoi(raw); err == nil {
			x.replacement = asciify.ReplacementsType(kind)
			fireEvent("downsample:ui", document)
		}
		return true
	})
	n := x.rpl.Create(document)
	for kind, name := range replacements {
		opts := map[attributeName]attributeValue{
			"value": attributeValue(fmt.Sprintf("%d", kind)),
		}
		if kind == x.replacement {
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

func (x *Replacements) GetReplacementType() asciify.ReplacementsType {
	return x.replacement
}
