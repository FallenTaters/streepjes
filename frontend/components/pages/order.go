package pages

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"

	"github.com/PotatoesFall/vecty-test/domain"
	"github.com/PotatoesFall/vecty-test/frontend/backend/cache"
	"github.com/PotatoesFall/vecty-test/frontend/components/catalog"
	"github.com/PotatoesFall/vecty-test/frontend/components/club"
	"github.com/PotatoesFall/vecty-test/frontend/components/pages/order"
)

func Order() *OrderComponent {
	o := &OrderComponent{
		items: make(order.Items),
		club:  domain.ClubParabool, // TODO dont always default to gladiators
		// idea: before entering order screen, club must be selected?
		// the toggle then becomes unnecessary or could at least be removed on mobile
	}

	o.overview = order.Overview(o.items, o.club)

	return o
}

type OrderComponent struct {
	vecty.Core

	club               domain.Club
	items              order.Items
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
			vecty.Markup(vecty.Class(`row`, o.club.String())),
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
			vecty.Markup(vecty.Class(`row`, `no-wrap`)),
			elem.Div(
				vecty.Markup(vecty.Class(`col`, `min`)),
				&club.Toggler{
					Rerender: false,
					Club:     o.club,
					OnToggle: func(c domain.Club) {
						o.club = c
						vecty.Rerender(o)
					},
				},
			),
			// elem.Div(
			// 	vecty.Markup(vecty.Class(`col`, `s12`, `m6`, `l4`)),
			// 	// club.Logo(domain.ClubGladiators, 100),
			// 	// club.Logo(domain.ClubParabool, 100),
			// ), elem.Div(
			// 	vecty.Markup(vecty.Class(`col`, `s12`, `m6`, `l4`)),
			// 	elem.Heading2(vecty.Text("Payment"))
			// ),
		),
	)
}

func (o *OrderComponent) AddItem(item domain.Item) {
	o.items.Add(item)
	vecty.Rerender(o.overview)
}

func (o *OrderComponent) DeleteItem(item domain.Item) {
	o.items.Delete(item)
	defer vecty.Rerender(o.overview)
}
