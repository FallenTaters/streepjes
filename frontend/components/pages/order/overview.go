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

func Overview(items map[domain.Item]int, club *domain.Club) *OverviewComponent {
	return &OverviewComponent{
		club:  club,
		items: items,
	}
}

type OverviewComponent struct {
	vecty.Core

	club  *domain.Club `vecty:"prop"`
	items Items        `vecty:"prop"`
}

func (o *OverviewComponent) Render() vecty.ComponentOrHTML {
	markupAndChildren := []vecty.MarkupOrChild{
		vecty.Markup(vecty.Style(`margin-top`, `7px`), vecty.Style(`padding-bottom`, `3px`)),
		elem.Heading5(vecty.Text("Overview")),
	}

	markupAndChildren = append(markupAndChildren, o.makeCards()...)

	return elem.Div(markupAndChildren...)
}

func (o *OverviewComponent) makeCards() []vecty.MarkupOrChild {
	var itemsSorted []domain.Item
	for item := range o.items {
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
		children = append(children, o.makeCard(item, o.items))
	}

	return children
}

// TODO make smaller, remove duplication
func (o *OverviewComponent) makeCard(item domain.Item, items Items) vecty.MarkupOrChild {
	amount := items[item]

	return elem.Article(
		vecty.Markup(vecty.Class(`small-padding`)),
		elem.Div(
			vecty.Markup(vecty.Class(`row`, `no-wrap`, `large-text`, `no-club`)),
			elem.Div(
				vecty.Markup(vecty.Class(`col`, `min`, `middle-align`)),
				elem.Button(
					vecty.Markup(
						vecty.Class(`circle`, `left-round`, `no-margin`),
						event.Click(func(e *vecty.Event) {
							items.Remove(item)
							vecty.Rerender(o)
						}),
					),
					beercss.Icon(beercss.IconRemove),
				),
			), elem.Div(
				vecty.Markup(vecty.Class(`col`, `min`, `middle-align`)),
				elem.Span(
					vecty.Markup(vecty.Class(`bold`)),
					vecty.Text(fmt.Sprint(amount)),
				),
			), elem.Div(
				vecty.Markup(vecty.Class(`col`, `min`, `middle-align`)),
				elem.Button(
					vecty.Markup(
						vecty.Class(`circle`, `right-round`, `no-margin`),
						event.Click(func(e *vecty.Event) {
							o.items.Add(item)
							vecty.Rerender(o)
						}),
					),
					beercss.Icon(beercss.IconAdd),
				),
			), elem.Div(
				vecty.Markup(
					vecty.Class(`col`, `max`, `middle-align`),
					vecty.Style(`text-overflow`, `ellipsis`),
				),
				elem.Span(vecty.Text(item.Name)),
			), elem.Div(
				vecty.Markup(vecty.Class(`col`, `min`, `middle-align`)),
				vecty.Text(
					item.Price(*o.club).Times(amount).String(),
				),
			), elem.Div(
				vecty.Markup(vecty.Class(`col`, `min`, `middle-align`)),
				elem.Button(
					vecty.Markup(
						vecty.Class(`circle`, `error`),
						event.Click(func(e *vecty.Event) {
							o.items.Delete(item)
							vecty.Rerender(o)
						}),
					),
					beercss.Icon(beercss.IconDelete),
				),
			),
		),
	)
}

type Items map[domain.Item]int

func (oi Items) Add(item domain.Item) {
	oi[item]++
}

func (oi Items) Remove(item domain.Item) {
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

func (oi Items) Delete(item domain.Item) {
	delete(oi, item)
}
