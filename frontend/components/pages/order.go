package pages

import (
	"github.com/hexops/vecty"

	"github.com/PotatoesFall/vecty-test/api"
	"github.com/PotatoesFall/vecty-test/frontend/components/catalog"
	"github.com/PotatoesFall/vecty-test/frontend/components/pages/order"
	"github.com/PotatoesFall/vecty-test/frontend/jscall/window"
	"github.com/PotatoesFall/vecty-test/frontend/state/cache"
)

type Order struct {
	vecty.Core

	catalog api.Catalog

	conceptOrder order.Concept
}

func (o *Order) Render() vecty.ComponentOrHTML {
	largeScreen := window.OnResize(func() {
		vecty.Rerender(o)
	})

	// TODO not all this in render maybe?

	cat, err := cache.Catalog()
	if err != nil {
		// TODO handle gracefully
		panic(err)
	}

	o.catalog = cat

	categories := catalog.Categories(o.catalog.Categories, func(c api.Category) {
		// TODO
	})

	if largeScreen {
		return order.Large(categories)
	}
	return &order.Small{
		Page: order.SubPageCategories,
	}
}
