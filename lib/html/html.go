package html

import (
	"fmt"
	"strings"
	"syscall/js"
	"unicode"
)

const (
	InputElementID  htmlAttributeValue = "input-file"
	OutputElementID htmlAttributeValue = "output"
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

func fireEvent(name eventType, document js.Value) {
	ev := js.Global().Get("Event").New(name.String())
	document.Call("dispatchEvent", ev)
}
