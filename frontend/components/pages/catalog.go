package pages

import (
	"sort"
	"strings"

	"github.com/FallenTaters/streepjes/domain/orderdomain"
	"github.com/FallenTaters/streepjes/frontend/backend"
	"github.com/FallenTaters/streepjes/frontend/global"
)

type Catalog struct {
	Loading bool
	Error   bool

	Categories   []orderdomain.Category
	Items        []orderdomain.Item
	DisplayItems []orderdomain.Item

	NewCategory      bool
	NewItem          bool
	SelectedCategory orderdomain.Category
	SelectedItem     orderdomain.Item

	// form values
	CategoryID      int
	Name            string
	PriceGladiators orderdomain.Price
	PriceParabool   orderdomain.Price
}

func (c *Catalog) Init() {
	c.Loading = true
	c.Error = false

	go func() {
		catalog, err := backend.GetCatalog()
		defer global.LockAndRender()()
		c.Loading = false
		if err != nil {
			c.Error = true
			return
		}

		c.Categories = catalog.Categories
		c.Items = catalog.Items

		sort.Slice(c.Categories, func(i, j int) bool {
			return strings.Compare(c.Categories[i].Name, c.Categories[j].Name) < 0
		})

		sort.Slice(c.Items, func(i, j int) bool {
			return strings.Compare(c.Items[i].Name, c.Items[j].Name) < 0
		})
	}()
}

func (c *Catalog) Compute() {
	c.DisplayItems = []orderdomain.Item{}
	for _, item := range c.Items {
		if item.CategoryID == c.SelectedCategory.ID {
			c.DisplayItems = append(c.DisplayItems, item)
		}
	}
}

func (c *Catalog) reset() {
	c.NewCategory = false
	c.NewItem = false
	c.SelectedItem = orderdomain.Item{}
	c.SelectedCategory = orderdomain.Category{}
}

func (c *Catalog) OnCategoryClick(category orderdomain.Category) {
	if c.SelectedCategory == category {
		return
	}

	c.reset()

	c.SelectedCategory = category
}

func (c *Catalog) OnCategoryClickNew() {
	c.reset()

	c.NewCategory = true
}

func (c *Catalog) OnItemClick(item orderdomain.Item) {
	c.NewItem = false
	c.SelectedItem = item
}

func (c *Catalog) OnItemClickNew() {
	c.SelectedItem = orderdomain.Item{}
	c.NewItem = true
}

func (c *Catalog) ShowCategoryForm() bool {
	return !c.ShowItemForm() && (c.NewCategory || c.SelectedCategory != (orderdomain.Category{}))
}

func (c *Catalog) ShowItemForm() bool {
	return c.NewItem || c.SelectedItem != (orderdomain.Item{})
}
