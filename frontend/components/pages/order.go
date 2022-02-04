package pages

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"

	"github.com/PotatoesFall/vecty-test/api"
	"github.com/PotatoesFall/vecty-test/domain"
	"github.com/PotatoesFall/vecty-test/frontend/backend/cache"
	"github.com/PotatoesFall/vecty-test/frontend/components/catalog"
	"github.com/PotatoesFall/vecty-test/frontend/components/club"
	"github.com/PotatoesFall/vecty-test/frontend/components/pages/order"
)

func Order() *OrderComponent {
	o := &OrderComponent{
		club: domain.ClubGladiators, // TODO
	}

	o.toggler = &club.Toggler{
		Size:     70,
		Rerender: false,
		Club:     o.club,
		OnToggle: func(c domain.Club) {
			o.club = c
			vecty.Rerender(o)
		},
	}

	cat, err := cache.Catalog()
	if err != nil {
		// TODO handle gracefully
		panic(err)
	}
	o.catalog = cat

	o.items = catalog.Items([]domain.Item{}, func(i domain.Item) {
		o.AddItem(i)
	})

	o.categories = catalog.Categories(cat.Categories, func(c domain.Category) {
		o.selectedCategoryID = c.ID

		var newItems []domain.Item
		for _, item := range cat.Items {
			if item.CategoryID == c.ID {
				newItems = append(newItems, item)
			}
		}

		o.items.SetItems(newItems)
	})

	o.itemsInOrder = make(order.Items)
	o.overview = order.Overview(o.itemsInOrder, o.club)

	return o
}

type OrderComponent struct {
	vecty.Core

	club domain.Club

	catalog            api.Catalog
	selectedCategoryID int
	itemsInOrder       order.Items

	toggler    *club.Toggler
	items      *catalog.ItemsComponent
	categories *catalog.CategoriesComponent
	overview   *order.OverviewComponent
}

func (o *OrderComponent) Render() vecty.ComponentOrHTML {
	return elem.Div(
		elem.Div(
			vecty.Markup(vecty.Class(`row`, `s`, `no-wrap`)),
			elem.Div(vecty.Markup(vecty.Class(`col`, `max`))),
			elem.Div(
				vecty.Markup(vecty.Class(`col`, `min`)),
				o.toggler,
			),
			elem.Div(vecty.Markup(vecty.Class(`col`, `max`))),
		),
		elem.Div(
			vecty.Markup(vecty.Class(`row`, o.club.String())),
			elem.Div(
				vecty.Markup(vecty.Class(`col`, `s12`, `m6`, `l3`)),
				o.categories,
			),
			elem.Div(
				vecty.Markup(vecty.Class(`col`, `s12`, `m6`, `l4`)),
				o.items,
			),
			elem.Div(
				vecty.Markup(vecty.Class(`col`, `s12`, `m12`, `l5`)),
				o.overview,
			),
		),
		elem.Div(
			vecty.Markup(vecty.Class(`row`)),
			elem.Div(
				vecty.Markup(vecty.Class(`col`, `min`, `l`, `m`, `l2`)),
				o.toggler,
			),
			elem.Div(
				vecty.Markup(vecty.Class(`col`, `l10`, `s12`)),
				vecty.Text(`hello`),
			),
		),
	)
}

func (o *OrderComponent) AddItem(item domain.Item) {
	o.itemsInOrder.Add(item)
	vecty.Rerender(o.overview)
}

func (o *OrderComponent) DeleteItem(item domain.Item) {
	o.itemsInOrder.Delete(item)
	defer vecty.Rerender(o.overview)
}
