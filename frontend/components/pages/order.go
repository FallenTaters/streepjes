package pages

import (
	"fmt"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"

	"github.com/PotatoesFall/vecty-test/api"
	"github.com/PotatoesFall/vecty-test/domain"
	"github.com/PotatoesFall/vecty-test/frontend/backend/cache"
	"github.com/PotatoesFall/vecty-test/frontend/components/catalog"
	"github.com/PotatoesFall/vecty-test/frontend/components/club"
	"github.com/PotatoesFall/vecty-test/frontend/components/pages/order"
	"github.com/PotatoesFall/vecty-test/frontend/jscall/window"
)

func Order() *OrderComponent {
	o := new(OrderComponent)
	clb := new(domain.Club)      // must be passed in callback
	*clb = domain.ClubGladiators // TODO
	o.Club = clb

	o.Toggler = &club.Toggler{
		Size:     70,
		Rerender: false,
		Club:     *clb,
		OnToggle: func(c domain.Club) {
			fmt.Println(c.String())
			*clb = c
			vecty.Rerender(o)
			vecty.Rerender(o.Overview)
		},
	}

	itemsInOrder := make(order.Items) // must be passed in callback
	o.ItemsInOrder = itemsInOrder

	cat, err := cache.Catalog()
	if err != nil {
		// TODO handle gracefully
		panic(err)
	}
	o.Catalog = cat

	overview := order.Overview(itemsInOrder, clb) // must be passed in callback
	o.Overview = overview

	items := catalog.Items([]domain.Item{}, func(i domain.Item) {
		itemsInOrder.Add(i)
		vecty.Rerender(overview)
	})
	o.Items = items

	var selectedCategoryID int
	o.SelectedCategoryID = &selectedCategoryID

	o.Categories = catalog.Categories(cat.Categories, func(c domain.Category) {
		selectedCategoryID = c.ID

		var newItems []domain.Item
		for _, item := range cat.Items {
			if item.CategoryID == c.ID {
				newItems = append(newItems, item)
			}
		}

		items.SetItems(newItems)
	})

	return o
}

type OrderComponent struct {
	vecty.Core

	Club *domain.Club `vecty:"prop"`

	Catalog            api.Catalog `vecty:"prop"`
	SelectedCategoryID *int        `vecty:"prop"`
	ItemsInOrder       order.Items `vecty:"prop"`

	Toggler    *club.Toggler                `vecty:"prop"`
	Items      *catalog.ItemsComponent      `vecty:"prop"`
	Categories *catalog.CategoriesComponent `vecty:"prop"`
	Overview   *order.OverviewComponent     `vecty:"prop"`
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
		vecty.Markup(vecty.Class(`row`, o.Club.String())),
		child,
	)
}

func (o *OrderComponent) DeleteItem(item domain.Item) {
	o.ItemsInOrder.Delete(item)
	defer vecty.Rerender(o.Overview)
}

// TODO THIS SPLITTING IS BROKEN AS FUCK, ON RESIZE THE WHOLE THING BREAKS
func (o *OrderComponent) grid() *vecty.HTML {
	return elem.Div(
		vecty.Markup(
			vecty.Class(`l`),
			vecty.Style(`display`, `grid`),
			vecty.Style(`grid-gap`, `5px`),
			vecty.Style(`grid-template-columns`, `30% 30% 40%`),
			vecty.Style(`grid-template-rows`, `1fr 200px`),
		),
		o.Categories,
		o.Items,
		o.Overview,
		o.Toggler,
	)
}

func (o *OrderComponent) reactive() *vecty.HTML {
	return elem.Div(
		elem.Div(
			vecty.Markup(vecty.Class(`row`, `no-wrap`)),
			elem.Div(vecty.Markup(vecty.Class(`col`, `max`))),
			elem.Div(
				vecty.Markup(vecty.Class(`col`, `min`)),
				o.Toggler,
			),
			elem.Div(vecty.Markup(vecty.Class(`col`, `max`))),
		),
		elem.Div(
			elem.Div(
				vecty.Markup(vecty.Class(`col`, `s12`, `m6`)),
				o.Categories,
			),
			elem.Div(
				vecty.Markup(vecty.Class(`col`, `s12`, `m6`)),
				o.Items,
			),
			elem.Div(
				vecty.Markup(vecty.Class(`col`, `s12`, `m12`)),
				o.Overview,
			),
		),
	)
}
