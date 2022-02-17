package pages

import (
	"errors"
	"fmt"

	"github.com/PotatoesFall/vecty-test/api"
	"github.com/PotatoesFall/vecty-test/domain"
	"github.com/PotatoesFall/vecty-test/frontend/backend/cache"
	"github.com/PotatoesFall/vecty-test/frontend/global"
	"github.com/vugu/vugu"
)

type Order struct {
	Catalog            api.Catalog   `vugu:"data"`
	Items              []domain.Item `vugu:"data"`
	SelectedCategoryID int           `vugu:"data"`
}

var ErrLoadCatalog = errors.New(`unable to load catalog`)

func NewOrder() (vugu.Builder, error) {
	order := &Order{
		SelectedCategoryID: 1, // TODO remove this line
	}

	catalog, err := cache.Catalog()
	if err != nil {
		return nil, fmt.Errorf(`%w: %s`, ErrLoadCatalog, err.Error())
	}
	order.Catalog = catalog

	return order, nil
}

func (o *Order) filterItems() {
	o.Items = []domain.Item{}

	for _, item := range o.Catalog.Items {
		if item.CategoryID == o.SelectedCategoryID {
			o.Items = append(o.Items, item)
		}
	}
}

func (o *Order) selectCategory(category domain.Category) {
	fmt.Println(`selecting:`, category.ID)
	go func() {
		global.EventEnv.Lock()
		defer global.EventEnv.UnlockRender()

		o.SelectedCategoryID = category.ID
		o.filterItems()
		fmt.Println(`selected:`, category.ID)
	}()
}

// import (
// 	"errors"
// 	"fmt"

// 	"github.com/hexops/vecty"
// 	"github.com/hexops/vecty/elem"
// 	"github.com/hexops/vecty/event"

// 	"github.com/PotatoesFall/vecty-test/domain"
// 	"github.com/PotatoesFall/vecty-test/frontend/backend/cache"
// 	"github.com/PotatoesFall/vecty-test/frontend/components/catalog"
// 	"github.com/PotatoesFall/vecty-test/frontend/components/club"
// 	"github.com/PotatoesFall/vecty-test/frontend/components/pages/order"
// 	"github.com/PotatoesFall/vecty-test/frontend/jscall/window"
// 	"github.com/PotatoesFall/vecty-test/frontend/store"
// )

// func Order() (*OrderComponent, error) {
// 	catalog, err := cache.Catalog()
// 	if err != nil {
// 		return nil, fmt.Errorf(`%w: %s`, ErrLoadCatalog, err.Error())
// 	}

// 	store.Order.Catalog = catalog

// 	return &OrderComponent{}, err
// }

// type OrderComponent struct {
// 	vecty.Core
// }

// func (o *OrderComponent) Render() vecty.ComponentOrHTML {
// 	screenSize := window.GetSize()

// 	var child *vecty.HTML
// 	if screenSize == window.SizeL {
// 		child = o.grid()
// 	} else {
// 		child = o.reactive()
// 	}

// 	return elem.Div(
// 		vecty.Markup(vecty.Class(store.Order.Club.String())),
// 		child,
// 	)
// }

// func (o *OrderComponent) reactive() *vecty.HTML {
// 	return elem.Div(
// 		elem.Div(
// 			vecty.Markup(vecty.Class(`row`, `no-wrap`)),
// 			elem.Div(vecty.Markup(vecty.Class(`col`, `max`))),
// 			elem.Div(
// 				vecty.Markup(vecty.Class(`col`, `min`)),
// 				toggler(),
// 			),
// 			elem.Div(vecty.Markup(vecty.Class(`col`, `max`))),
// 		),
// 		elem.Div(
// 			vecty.Markup(vecty.Class(`row`)),
// 			elem.Div(
// 				vecty.Markup(vecty.Class(`col`, `s12`, `m6`)),
// 				elem.Heading5(vecty.Text("Categories")),
// 				categories(),
// 			),
// 			elem.Div(
// 				vecty.Markup(vecty.Class(`col`, `s12`, `m6`)),
// 				elem.Heading5(vecty.Text("Items")),
// 				items(),
// 			),
// 			elem.Div(
// 				vecty.Markup(vecty.Class(`col`, `s12`, `m12`)),
// 				elem.Heading5(vecty.Text("Overview")),
// 				&order.Overview{},
// 			),
// 			elem.Div(
// 				vecty.Markup(vecty.Class(`col`, `s12`, `m12`)),
// 				&order.Summary{},
// 			),
// 		),
// 	)
// }

// func categories() *catalog.CategoriesComponent {
// 	return catalog.Categories(store.Order.Categories(), store.Order.SelectedCategoryID, func(c domain.Category) {
// 		store.Order.SelectCategory(c.ID)
// 	})
// }

// func items() *catalog.ItemsComponent {
// 	return catalog.Items(store.Order.Items(), store.Order.AddItem)
// }

// func toggler() *vecty.HTML {
// 	return elem.Div(
// 		vecty.Markup(
// 			vecty.Class(`col`, `min`),
// 			event.Click(func(e *vecty.Event) {
// 				store.Order.ToggleClub()
// 			}),
// 		),
// 		club.Logo(store.Order.Club, 120),
// 	)
// }
