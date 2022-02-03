package pages

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"

	"github.com/PotatoesFall/vecty-test/api"
	"github.com/PotatoesFall/vecty-test/frontend/components/catalog"
	"github.com/PotatoesFall/vecty-test/frontend/components/pages/order"
	"github.com/PotatoesFall/vecty-test/frontend/state/cache"
)

func Order() *OrderComponent {
	o := &OrderComponent{
		items: make(map[api.Item]int),
	}

	o.overview = order.Overview(o.items, o.DeleteItem)

	return o
}

type OrderComponent struct {
	vecty.Core

	// concept api.Order
	items map[api.Item]int

	overview *order.OverviewComponent

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
		elem.Div(
			vecty.Markup(vecty.Class(`row`)),
			elem.Div(
				vecty.Markup(vecty.Class(`col`, `s12`, `m6`, `l3`)),
				categories,
			),
			elem.Div(
				vecty.Markup(vecty.Class(`col`, `s12`, `m6`, `l3`)),
				items,
			),
			elem.Div(
				vecty.Markup(vecty.Class(`col`, `s12`, `m12`, `l6`)),
				o.overview,
			),
		),
		elem.Div(
			vecty.Markup(vecty.Class(`row`)),
			elem.Div(
				vecty.Markup(vecty.Class(`col`, `s12`, `m6`, `l4`)),
				elem.Heading2(vecty.Text("Club"))),
			elem.Div(
				vecty.Markup(vecty.Class(`col`, `s12`, `m6`, `l4`)),
				elem.Heading2(vecty.Text("Member"))),
			elem.Div(
				vecty.Markup(vecty.Class(`col`, `s12`, `m6`, `l4`)),
				elem.Heading2(vecty.Text("Payment"))),
		),
	)
}

func (o *OrderComponent) AddItem(item api.Item) {
	o.items[item]++
	vecty.Rerender(o.overview)
}

func (o *OrderComponent) DeleteItem(item api.Item) {
	defer vecty.Rerender(o.overview)

	n, ok := o.items[item]
	if !ok {
		return
	}

	if n <= 1 {
		delete(o.items, item)
		return
	}

	o.items[item]--
}
