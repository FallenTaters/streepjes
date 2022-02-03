package order

import (
	"fmt"
	"sort"
	"strings"

	"github.com/PotatoesFall/vecty-test/api"
	"github.com/PotatoesFall/vecty-test/frontend/components/beercss"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
)

func Overview(items map[api.Item]int, onDelete func(api.Item)) *OverviewComponent {
	return &OverviewComponent{
		items:    items,
		onDelete: onDelete,
	}
}

type OverviewComponent struct {
	vecty.Core

	items map[api.Item]int

	onDelete func(api.Item)
}

func (o *OverviewComponent) Render() vecty.ComponentOrHTML {
	fmt.Println(`render`, o.items)
	markupAndChildren := []vecty.MarkupOrChild{
		vecty.Markup(vecty.Style(`overflow`, `auto`)),
		elem.Heading5(vecty.Text("Overview")),
	}

	markupAndChildren = append(markupAndChildren, makeCards(o.items, o.onDelete)...)

	return elem.Div(markupAndChildren...)
}

func makeCards(itemCounts map[api.Item]int, onDelete func(api.Item)) []vecty.MarkupOrChild {
	var itemsSorted []api.Item
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
		children = append(children, makeCard(item, itemCounts[item], func(e *vecty.Event) {
			onDelete(itm)
		}))
	}

	return children
}

func makeCard(item api.Item, count int, onClick func(e *vecty.Event)) vecty.MarkupOrChild {
	return elem.Article(
		elem.Div(
			vecty.Markup(vecty.Class(`row`, `no-wrap`)),
			elem.Div(
				vecty.Markup(vecty.Class(`col`, `max`)),
				elem.Heading6(
					vecty.Markup(vecty.Class(`no-margin`)),
					vecty.Text(fmt.Sprintf(`x%d %s`, count, item.Name)),
				),
			),
			elem.Div(
				vecty.Markup(vecty.Class(`col`, `min`)),
				elem.Button(
					vecty.Markup(
						vecty.Class(`round`, `error`),
						event.Click(onClick),
					),
					beercss.Icon(beercss.IconDelete),
				),
			),
		),
	)
}
