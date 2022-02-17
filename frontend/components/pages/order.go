package pages

import (
	"errors"

	"github.com/PotatoesFall/vecty-test/api"
	"github.com/PotatoesFall/vecty-test/domain"
	"github.com/PotatoesFall/vecty-test/frontend/backend/cache"
	"github.com/PotatoesFall/vecty-test/frontend/components/pages/order"
	"github.com/PotatoesFall/vecty-test/frontend/global"
	"github.com/PotatoesFall/vecty-test/frontend/jscall/window"
	"github.com/PotatoesFall/vecty-test/frontend/store"
	"github.com/vugu/vugu"
)

type Order struct {
	Catalog            api.Catalog       `vugu:"data"`
	Categories         []domain.Category `vugu:"data"`
	Items              []domain.Item     `vugu:"data"`
	SelectedCategoryID int               `vugu:"data"`
	Large              bool              `vugu:"data"`
}

var ErrLoadCatalog = errors.New(`unable to load catalog`)

func (o *Order) Compute(vugu.ComputeCtx) {
	o.filterCategories()
	o.filterItems()
	o.Large = window.GetSize() == window.SizeL
}

func (o *Order) component() vugu.Builder {
	if window.GetSize() == window.SizeL {
		return &order.Grid{}
	}

	return &order.Reactive{}
}

func (o *Order) club() string {
	return store.Order.Club.String()
}

func (o *Order) Init(vugu.InitCtx) {
	go func() {
		catalog, err := cache.Catalog()
		if err != nil {
			panic(err) // TODO
		}

		global.EventEnv.Lock()
		defer global.EventEnv.UnlockRender()
		o.Catalog = catalog
	}()
}

func (o *Order) filterItems() {
	o.Items = []domain.Item{}

	for _, item := range o.Catalog.Items {
		if item.CategoryID == o.SelectedCategoryID {
			o.Items = append(o.Items, item)
		}
	}
}

func (o *Order) filterCategories() {
	o.Categories = []domain.Category{}

	seenCategoryIDs := map[int]bool{}
	for _, item := range o.Catalog.Items {
		seenCategoryIDs[item.CategoryID] = seenCategoryIDs[item.CategoryID] || item.Price(store.Order.Club) != 0
	}

	for _, category := range o.Catalog.Categories {
		if seenCategoryIDs[category.ID] {
			o.Categories = append(o.Categories, category)
		}
	}
}

func (o *Order) selectCategory(category domain.Category) {
	go func() {
		global.EventEnv.Lock()
		defer global.EventEnv.UnlockRender()

		o.SelectedCategoryID = category.ID
		o.filterItems()
	}()
}

func (o *Order) selectItem(item domain.Item) {
	store.Order.AddItem(item)
}
