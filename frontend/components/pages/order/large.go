package order

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

type Large struct {
	vecty.Core
}

func (l *Large) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			vecty.Class(`full-height`, `order-grid`),
		),
		elem.Div(vecty.Markup(vecty.Style(`overflow`, `auto`)), elem.Heading2(vecty.Text("Categories"))),
		elem.Div(vecty.Markup(vecty.Style(`overflow`, `auto`)), elem.Heading2(vecty.Text("Products"))),
		elem.Div(vecty.Markup(vecty.Style(`overflow`, `auto`)), elem.Heading2(vecty.Text("Order"))),
		elem.Div(elem.Heading2(vecty.Text("Club"))),
		elem.Div(elem.Heading2(vecty.Text("Member"))),
		elem.Div(elem.Heading2(vecty.Text("Payment"))),
	)
}
