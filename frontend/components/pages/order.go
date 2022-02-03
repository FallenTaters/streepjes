package pages

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"

	"github.com/PotatoesFall/vecty-test/domain"
	"github.com/PotatoesFall/vecty-test/frontend/components/catalog"
	"github.com/PotatoesFall/vecty-test/frontend/components/pages/order"
	"github.com/PotatoesFall/vecty-test/frontend/state/cache"
)

func Order() *OrderComponent {
	o := &OrderComponent{
		items: make(map[domain.Item]int),
		club:  domain.ClubGladiators, // TODO
		// club: domain.ClubParabool, // TODO
	}

	o.overview = order.Overview(o.items, o.DeleteItem, o.club)

	return o
}

type OrderComponent struct {
	vecty.Core

	club               domain.Club
	items              map[domain.Item]int
	overview           *order.OverviewComponent
	selectedCategoryID int
}

func (o *OrderComponent) Render() vecty.ComponentOrHTML {
	cat, err := cache.Catalog()
	if err != nil {
		// TODO handle gracefully
		panic(err)
	}

	items := catalog.Items([]domain.Item{}, func(i domain.Item) {
		o.AddItem(i)
	})

	categories := catalog.Categories(cat.Categories, func(c domain.Category) {
		o.selectedCategoryID = c.ID

		var newItems []domain.Item
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
				vecty.Markup(vecty.Class(`col`, `s12`, `m6`, `l4`)),
				items,
			),
			elem.Div(
				vecty.Markup(vecty.Class(`col`, `s12`, `m12`, `l5`)),
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

func (o *OrderComponent) AddItem(item domain.Item) {
	o.items[item]++
	vecty.Rerender(o.overview)
}

func (o *OrderComponent) DeleteItem(item domain.Item) {
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
