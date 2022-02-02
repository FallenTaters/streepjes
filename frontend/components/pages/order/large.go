package order

import (
	"github.com/PotatoesFall/vecty-test/frontend/components/catalog"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

func Large(categories *catalog.CategoriesComponent) vecty.Component {
	return &LargeComponent{
		categories: categories,
	}
}

type LargeComponent struct {
	vecty.Core

	categories *catalog.CategoriesComponent
}

func (l *LargeComponent) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			vecty.Class(`full-height`, `order-grid`),
		),
		l.categories,
		elem.Div(vecty.Markup(vecty.Style(`overflow`, `auto`)), elem.Heading2(vecty.Text("Products"))),
		elem.Div(vecty.Markup(vecty.Style(`overflow`, `auto`)), elem.Heading2(vecty.Text("Order"))),
		elem.Div(elem.Heading2(vecty.Text("Club"))),
		elem.Div(elem.Heading2(vecty.Text("Member"))),
		elem.Div(elem.Heading2(vecty.Text("Payment"))),
	)
}
