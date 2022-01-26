package order

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
)

type Large struct {
	vecty.Core

	resize func(*vecty.Event)
}

func (l *Large) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			vecty.Class(`full-height`, `order-grid`),
			event.Resize(l.resize),
		),
		elem.Div(vecty.Markup(vecty.Style(`overflow`, `auto`)), vecty.Text(longText)),
		elem.Div(vecty.Markup(vecty.Style(`overflow`, `auto`)), vecty.Text(`This is the order page!`)),
		elem.Div(vecty.Markup(vecty.Style(`overflow`, `auto`)), vecty.Text(`This is the order page!`)),
		elem.Div(vecty.Text(`This is the order page!`)),
		elem.Div(vecty.Text(`This is the order page!`)),
		elem.Div(vecty.Text(`This is the order page!`)),
	)
}

var longText = `This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>This is the order page!<br>`
