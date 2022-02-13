package beercss

import (
	"github.com/PotatoesFall/vecty-test/frontend/jscall/beercss"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
)

func Input(label string) vecty.Component {
	uiCalled := false
	listener := func(e *vecty.Event) {
		if uiCalled {
			return
		}

		go beercss.UI()

		uiCalled = true
	}

	return &input{
		label:    label,
		listener: listener,
	}
}

type input struct {
	vecty.Core

	label    string               `vecty:"prop"`
	listener func(r *vecty.Event) `vecty:"prop"`
}

func (i *input) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			vecty.Class(`field`, `label`, `border`),
			event.MouseMove(i.listener),
		),
		elem.Input(vecty.Markup(vecty.Attribute(`type`, `text`))),
		elem.Label(vecty.Text(i.label)),
	)
}
