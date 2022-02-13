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
	children := make([]vecty.MarkupOrChild, 0, len(store.Order.Lines))
	for _, item := range store.Order.Lines {
		children = append(children, makeCard(item))
	}

	return children
}

func makeCard(item store.Orderline) *vecty.HTML { //nolint:funlen
	classList := []string{`small-padding`}
	if item.Item.Price(store.Order.Club) == 0 {
		classList = append(classList, `error`)
	}

	return elem.Article(
		vecty.Markup(vecty.Class(classList...)),
		elem.Div(
			vecty.Markup(vecty.Class(`row`, `no-wrap`, `large-text`, `no-club`)),
			minCol(
				elem.Button(
					vecty.Markup(
						vecty.Class(`circle`, `left-round`, `no-margin`),
						event.Click(func(*vecty.Event) { store.Order.RemoveItem(item.Item) }),
					),
					beercss.Icon(beercss.IconRemove),
				),
			),
			minCol(
				vecty.Markup(vecty.Style(`text-align`, `center`)),
				elem.Span(
					vecty.Markup(
						vecty.Class(`bold`),
						vecty.Style(`width`, `20px`),
					),
					vecty.Text(fmt.Sprint(item.Amount)),
				),
			),
			minCol(
				elem.Button(
					vecty.Markup(
						vecty.Class(`circle`, `right-round`, `no-margin`),
						event.Click(func(*vecty.Event) { store.Order.AddItem(item.Item) }),
					),
					beercss.Icon(beercss.IconAdd),
				),
			),

			maxCol(
				elem.Span(vecty.Text(item.Item.Name)),
			),

			minCol(
				vecty.Text(
					item.Item.Price(store.Order.Club).Times(item.Amount).String(),
				),
			),

			minCol(
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

func minCol(children ...vecty.MarkupOrChild) *vecty.HTML {
	children = append(children, vecty.Markup(vecty.Class(`col`, `min`, `middle-align`)))
	return elem.Div(children...)
}

func maxCol(children ...vecty.MarkupOrChild) *vecty.HTML {
	children = append(children, vecty.Markup(vecty.Class(`col`, `max`, `middle-align`)))
	return elem.Div(children...)
}
