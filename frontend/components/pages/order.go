package pages

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"

	"github.com/PotatoesFall/vecty-test/domain"
	"github.com/PotatoesFall/vecty-test/frontend/backend/cache"
	"github.com/PotatoesFall/vecty-test/frontend/components/catalog"
	"github.com/PotatoesFall/vecty-test/frontend/components/club"
	"github.com/PotatoesFall/vecty-test/frontend/components/pages/order"
	"github.com/PotatoesFall/vecty-test/frontend/jscall/window"
	"github.com/PotatoesFall/vecty-test/frontend/store"
)

func Order() *OrderComponent {
	catalog, err := cache.Catalog()
	if err != nil {
		panic(err) // TODO
	}

	store.Order.Catalog = catalog

	return &OrderComponent{}
}

type OrderComponent struct {
	vecty.Core
}

func (o *OrderComponent) Render() vecty.ComponentOrHTML {
	screenSize := window.GetSize()

	var child *vecty.HTML
	if screenSize == window.SizeL {
		child = o.grid()
	} else {
		child = o.reactive()
	}

	return elem.Div(
		vecty.Markup(vecty.Class(store.Order.Club.String())),
		child,
	)
}

func (o *OrderComponent) grid() *vecty.HTML {
	return elem.Div(
		vecty.Markup(
			vecty.Class(`l`),
			vecty.Style(`height`, `100%`),
			vecty.Style(`display`, `grid`),
			vecty.Style(`grid-gap`, `5px`),
			vecty.Style(`grid-template-columns`, `30% 30% 40%`),
			vecty.Style(`grid-template-rows`, `50px 1fr 200px`),
			vecty.Style(`grid-template-areas`, `"topleft topcenter topright" "midleft midcenter right" "bottom bottom right"`),
		),
		elem.Heading5(vecty.Text("Categories")),
		elem.Heading5(vecty.Text("Items")),
		elem.Heading5(vecty.Text("Overview")),
		elem.Div(
			vecty.Markup(vecty.Style(`overflow`, `auto`)),
			categories(),
		),
		elem.Div(
			vecty.Markup(vecty.Style(`overflow`, `auto`)),
			items(),
		),
		elem.Div(
			vecty.Markup(
				vecty.Style(`grid-area`, `right`),
			),
			vecty.Markup(
				vecty.Style(`display`, `grid`),
				vecty.Style(`grid-template-columns`, `1fr`),
				vecty.Style(`grid-template-rows`, `1fr 70px`),
			),
			elem.Div(
				vecty.Markup(vecty.Style(`overflow`, `auto`)),
				&order.Overview{},
			),
			&order.Summary{},
		),
		elem.Div(
			vecty.Markup(vecty.Style(`grid-area`, `bottom`)),
			toggler(),
		),
	)
}

func (o *OrderComponent) reactive() *vecty.HTML {
	return elem.Div(
		elem.Div(
			vecty.Markup(vecty.Class(`row`, `no-wrap`)),
			elem.Div(vecty.Markup(vecty.Class(`col`, `max`))),
			elem.Div(
				vecty.Markup(vecty.Class(`col`, `min`)),
				toggler(),
			),
			elem.Div(vecty.Markup(vecty.Class(`col`, `max`))),
		),
		elem.Div(
			vecty.Markup(vecty.Class(`row`)),
			elem.Div(
				vecty.Markup(vecty.Class(`col`, `s12`, `m6`)),
				elem.Heading5(vecty.Text("Categories")),
				categories(),
			),
			elem.Div(
				vecty.Markup(vecty.Class(`col`, `s12`, `m6`)),
				elem.Heading5(vecty.Text("Items")),
				items(),
			),
			elem.Div(
				vecty.Markup(vecty.Class(`col`, `s12`, `m12`)),
				elem.Heading5(vecty.Text("Overview")),
				&order.Overview{},
			),
			elem.Div(
				vecty.Markup(vecty.Class(`col`, `s12`, `m12`)),
				&order.Summary{},
			),
		),
	)
}

func categories() *catalog.CategoriesComponent {
	return catalog.Categories(store.Order.Catalog.Categories, store.Order.SelectedCategoryID, func(c domain.Category) {
		store.Order.SelectCategory(c.ID)
	})
}

func items() *catalog.ItemsComponent {
	return catalog.Items(store.Order.ShownItems, store.Order.AddItem)
}

func toggler() *vecty.HTML {
	return elem.Div(
		vecty.Markup(
			vecty.Class(`col`, `min`),
			event.Click(func(e *vecty.Event) {
				store.Order.ToggleClub()
			}),
		),
		club.Logo(store.Order.Club, 120),
	)
}
