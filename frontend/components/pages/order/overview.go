package order

import (
	"fmt"

	"github.com/PotatoesFall/vecty-test/frontend/components/beercss"
	"github.com/PotatoesFall/vecty-test/frontend/store"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
)

type Overview struct {
	vecty.Core
}

func (o *Overview) Render() vecty.ComponentOrHTML {
	markupAndChildren := []vecty.MarkupOrChild{
		vecty.Markup(vecty.Style(`margin-top`, `7px`), vecty.Style(`padding-bottom`, `3px`)),
	}

	markupAndChildren = append(markupAndChildren, makeCards()...)

	return elem.Div(markupAndChildren...)
}

func makeCards() []vecty.MarkupOrChild {
	var children []vecty.MarkupOrChild
	for _, item := range store.Order.Items {
		children = append(children, makeCard(item))
	}

	return children
}

// TODO make smaller, remove duplication
func makeCard(item store.Orderline) vecty.MarkupOrChild {
	return elem.Article(
		vecty.Markup(vecty.Class(`small-padding`)),
		elem.Div(
			vecty.Markup(vecty.Class(`row`, `no-wrap`, `large-text`, `no-club`)),
			elem.Div(
				vecty.Markup(vecty.Class(`col`, `min`, `middle-align`)),
				elem.Button(
					vecty.Markup(
						vecty.Class(`circle`, `left-round`, `no-margin`),
						event.Click(func(*vecty.Event) { store.Order.RemoveItem(item.Item) }),
					),
					beercss.Icon(beercss.IconRemove),
				),
			), elem.Div(
				vecty.Markup(vecty.Class(`col`, `min`, `middle-align`)),
				elem.Span(
					vecty.Markup(vecty.Class(`bold`)),
					vecty.Text(fmt.Sprint(item.Amount)),
				),
			), elem.Div(
				vecty.Markup(vecty.Class(`col`, `min`, `middle-align`)),
				elem.Button(
					vecty.Markup(
						vecty.Class(`circle`, `right-round`, `no-margin`),
						event.Click(func(*vecty.Event) { store.Order.AddItem(item.Item) }),
					),
					beercss.Icon(beercss.IconAdd),
				),
			), elem.Div(
				vecty.Markup(
					vecty.Class(`col`, `max`, `middle-align`),
					vecty.Style(`text-overflow`, `ellipsis`),
				),
				elem.Span(vecty.Text(item.Item.Name)),
			), elem.Div(
				vecty.Markup(vecty.Class(`col`, `min`, `middle-align`)),
				vecty.Text(
					item.Item.Price(store.Order.Club).Times(item.Amount).String(),
				),
			), elem.Div(
				vecty.Markup(vecty.Class(`col`, `min`, `middle-align`)),
				elem.Button(
					vecty.Markup(
						vecty.Class(`circle`, `error`),
						event.Click(func(*vecty.Event) { store.Order.DeleteItem(item.Item) }),
					),
					beercss.Icon(beercss.IconDelete),
				),
			),
		),
	)
}
