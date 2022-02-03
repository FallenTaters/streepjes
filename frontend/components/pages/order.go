package pages

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"

	"github.com/PotatoesFall/vecty-test/api"
	"github.com/PotatoesFall/vecty-test/frontend/components/catalog"
	"github.com/PotatoesFall/vecty-test/frontend/state/cache"
)

func Order() *OrderComponent {
	return &OrderComponent{
		items: make(map[api.Item]int),
	}
}

type OrderComponent struct {
	vecty.Core

	// concept api.Order
	items map[api.Item]int

	selectedCategoryID int
}

func (o *OrderComponent) Render() vecty.ComponentOrHTML {
	cat, err := cache.Catalog()
	if err != nil {
		// TODO handle gracefully
		panic(err)
	}

	items := catalog.Items([]api.Item{}, func(i api.Item) {
		o.AddItem(i)
	})

	categories := catalog.Categories(cat.Categories, func(c api.Category) {
		o.selectedCategoryID = c.ID

		var newItems []api.Item
		for _, item := range cat.Items {
			if item.CategoryID == c.ID {
				newItems = append(newItems, item)
			}
		}

		items.SetItems(newItems)
	})

	return elem.Div(
		vecty.Markup(
			vecty.Class(`full-height`, `row`),
		),
		elem.Div(
			vecty.Markup(vecty.Class(`col`, `s12`, `m6`, `l4`)),
			categories,
		),
		elem.Div(
			vecty.Markup(vecty.Class(`col`, `s12`, `m6`, `l4`)),
			items,
		),
		elem.Div(
			vecty.Markup(vecty.Class(`col`, `s12`, `m6`, `l4`)),
			elem.Heading5(vecty.Text("Order")),
		),
		elem.Div(elem.Heading2(vecty.Text("Club"))),
		elem.Div(elem.Heading2(vecty.Text("Member"))),
		elem.Div(elem.Heading2(vecty.Text("Payment"))),
	)
}

func (o *OrderComponent) AddItem(item api.Item) {
	o.items[item]++
	// TODO: rerender correct components
}

func (o *OrderComponent) DeleteItem(item api.Item) {
	n, ok := o.items[item]
	if !ok || n < 1 {
		return
	}

	if n == 1 {
		delete(o.items, item)
		return
	}

	o.items[item]--
	// TODO: rerender correct components
}
