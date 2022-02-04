package order

import (
	"fmt"
	"sort"
	"strings"

	"github.com/PotatoesFall/vecty-test/domain"
	"github.com/PotatoesFall/vecty-test/frontend/components/beercss"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
)

func Overview(items map[domain.Item]int, onDelete func(domain.Item), club domain.Club) *OverviewComponent {
	return &OverviewComponent{
		club:     club,
		items:    items,
		onDelete: onDelete,
	}
}

type OverviewComponent struct {
	vecty.Core

	club     domain.Club
	items    Items
	onDelete func(domain.Item)
}

func (o *OverviewComponent) Render() vecty.ComponentOrHTML {
	fmt.Println(`render`, o.items)
	markupAndChildren := []vecty.MarkupOrChild{
		vecty.Markup(vecty.Style(`overflow`, `auto`), vecty.Style(`margin-top`, `7px`), vecty.Style(`padding-bottom`, `3px`)),
		elem.Heading5(vecty.Text("Overview")),
	}

	markupAndChildren = append(markupAndChildren, o.makeCards(o.items, o.onDelete)...)

	return elem.Div(markupAndChildren...)
}

func (o *OverviewComponent) makeCards(itemCounts map[domain.Item]int, onDelete func(domain.Item)) []vecty.MarkupOrChild {
	var itemsSorted []domain.Item
	for item := range itemCounts {
		itemsSorted = append(itemsSorted, item)
	}

	sort.Slice(itemsSorted, func(i, j int) bool {
		return strings.Compare(
			strings.ToLower(itemsSorted[i].Name),
			strings.ToLower(itemsSorted[j].Name),
		) < 0
	})

	var children []vecty.MarkupOrChild
	for _, item := range itemsSorted {
		itm := item
		children = append(children, o.makeCard(item, itemCounts[item], func(e *vecty.Event) {
			onDelete(itm)
		}))
	}

	return children
}

func (o *OverviewComponent) makeCard(item domain.Item, count int, onClick func(e *vecty.Event)) vecty.MarkupOrChild {
	return elem.Article(
		vecty.Markup(vecty.Class(`small-padding`)),
		elem.Div(
			vecty.Markup(vecty.Class(`row`, `no-wrap`, `large-text`)),
			elem.Div(
				vecty.Markup(vecty.Class(`col`, `max`, `middle-align`)),
				elem.Span(
					vecty.Markup(
						vecty.Class(`bold`),
						vecty.UnsafeHTML(fmt.Sprintf(`&nbsp;×%d&nbsp;&nbsp;`, count))),
				),
				elem.Span(vecty.Text(item.Name)),
			),
			elem.Div(
				vecty.Markup(vecty.Class(`col`, `min`, `middle-align`)),
				vecty.Text(item.Price(o.club).String()),
			),
			elem.Div(
				vecty.Markup(vecty.Class(`col`, `min`)),
				elem.Button(
					vecty.Markup(
						vecty.Class(`circle`, `left-round`),
						event.Click(func(e *vecty.Event) {
						}),
					),
					beercss.Icon(beercss.IconAdd),
				),
				elem.Button(),
			),
			// elem.Div(
			// 	vecty.Markup(vecty.Class(`col`, `min`)),
			// 	elem.Button(
			// 		vecty.Markup(
			// 			vecty.Class(`circle`, `error`),
			// 			event.Click(onClick),
			// 		),
			// 		beercss.Icon(beercss.IconDelete),
			// 	),
			// ),
		),
	)
}

type Items map[domain.Item]int

func (oi Items) Add(item domain.Item) {
	oi[item]++
}

func (oi Items) DeleteItem(item domain.Item) {
	n, ok := oi[item]
	if !ok {
		return
	}

	if n <= 1 {
		delete(oi, item)
		return
	}

	oi[item]--
}
