package pages

import (
	"errors"

	"github.com/FallenTaters/streepjes/api"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
	"github.com/FallenTaters/streepjes/frontend/backend/cache"
	"github.com/FallenTaters/streepjes/frontend/components/pages/order"
	"github.com/FallenTaters/streepjes/frontend/global"
	"github.com/FallenTaters/streepjes/frontend/jscall/window"
	"github.com/FallenTaters/streepjes/frontend/store"
	"github.com/vugu/vugu"
)

type Order struct {
	Catalog            api.Catalog            `vugu:"data"`
	Categories         []orderdomain.Category `vugu:"data"`
	Items              []orderdomain.Item     `vugu:"data"`
	SelectedCategoryID int                    `vugu:"data"`
	Large              bool                   `vugu:"data"`
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
			return // handle gracefully (unauthorized is already handled)
		}

		defer global.LockAndRender()()
		o.Catalog = catalog
	}()
}

func (o *Order) filterItems() {
	o.Items = []orderdomain.Item{}

	for _, item := range o.Catalog.Items {
		if item.CategoryID == o.SelectedCategoryID && item.Price(store.Order.Club) != 0 {
			o.Items = append(o.Items, item)
		}
	}
}

func (o *Order) filterCategories() {
	o.Categories = []orderdomain.Category{}

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

func (o *Order) selectCategory(category orderdomain.Category) {
	go func() {
		defer global.LockAndRender()()

		o.SelectedCategoryID = category.ID
		o.filterItems()
	}()
}

func (o *Order) selectItem(item orderdomain.Item) {
	store.Order.AddItem(item)
}
