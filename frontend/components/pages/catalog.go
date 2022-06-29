package pages

import (
	"sort"
	"strings"

	"github.com/FallenTaters/streepjes/domain/orderdomain"
	"github.com/FallenTaters/streepjes/frontend/backend"
	"github.com/FallenTaters/streepjes/frontend/backend/cache"
	"github.com/FallenTaters/streepjes/frontend/components/beercss"
	"github.com/FallenTaters/streepjes/frontend/global"
	"github.com/FallenTaters/streepjes/frontend/jscall/window"
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

	LoadingForm bool
	FormError   bool

	// category form values
	CategoryName string

	// item form values
	CategoryID      int
	ItemName        string
	PriceGladiators orderdomain.Price
	PriceParabool   orderdomain.Price
}

func (c *Catalog) Init() {
	c.Loading = true
	c.Error = false

	c.reset()

	go func() {
		catalog, err := cache.Catalog.Get()
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
	c.FormError = false

	c.NewCategory = false
	c.NewItem = false
	c.SelectedItem = orderdomain.Item{}
	c.SelectedCategory = orderdomain.Category{}

	c.CategoryName = ``
	c.CategoryID = 0
	c.ItemName = ``
	c.PriceGladiators = 0
	c.PriceParabool = 0
}

func (c *Catalog) OnCategoryClick(category orderdomain.Category) {
	c.reset()

	c.SelectedCategory = category

	c.CategoryName = category.Name
}

func (c *Catalog) OnCategoryClickNew() {
	c.reset()

	c.NewCategory = true
}

func (c *Catalog) OnItemClick(item orderdomain.Item) {
	selCat := c.SelectedCategory

	c.reset()

	c.SelectedItem = item
	c.SelectedCategory = selCat

	c.ItemName = item.Name
	c.CategoryID = item.CategoryID
	c.PriceGladiators = item.PriceGladiators
	c.PriceParabool = item.PriceParabool
}

func (c *Catalog) OnItemClickNew() {
	selCat := c.SelectedCategory

	c.reset()

	c.SelectedCategory = selCat
	c.NewItem = true
}

func (c *Catalog) ShowCategoryForm() bool {
	return !c.ShowItemForm() && (c.NewCategory || c.SelectedCategory != (orderdomain.Category{}))
}

func (c *Catalog) ShowItemForm() bool {
	return c.NewItem || c.SelectedItem != (orderdomain.Item{})
}

func (c *Catalog) FormTitle() string {
	switch {
	case c.NewCategory:
		return `New Category`

	case c.NewItem:
		return `New Item`

	case c.SelectedItem != orderdomain.Item{}:
		return `Edit Item - ` + c.SelectedItem.Name

	case c.SelectedCategory != orderdomain.Category{}:
		return `Edit Category - ` + c.SelectedCategory.Name

	}

	return ``
}

func (c *Catalog) SubmitCategoryForm() {
	c.FormError = false
	c.LoadingForm = true

	go func() {
		var err error

		if c.NewCategory {
			err = backend.PostNewCategory(orderdomain.Category{
				Name: c.CategoryName,
			})
		} else {
			err = backend.PostUpdateCategory(orderdomain.Category{
				ID:   c.SelectedCategory.ID,
				Name: c.CategoryName,
			})
		}
		defer global.LockAndRender()()
		c.LoadingForm = false
		if err != nil {
			c.FormError = true
			return
		}
		cache.Catalog.Invalidate()

		c.Init()
	}()
}

func (c *Catalog) SubmitItemForm() {
	c.FormError = false
	c.LoadingForm = true

	go func() {
		var err error

		if c.NewItem {
			err = backend.PostNewItem(orderdomain.Item{
				CategoryID:      c.CategoryID,
				Name:            c.ItemName,
				PriceGladiators: c.PriceGladiators,
				PriceParabool:   c.PriceParabool,
			})
		} else {
			err = backend.PostUpdateItem(orderdomain.Item{
				ID:              c.SelectedItem.ID,
				CategoryID:      c.CategoryID,
				Name:            c.ItemName,
				PriceGladiators: c.PriceGladiators,
				PriceParabool:   c.PriceParabool,
			})
		}
		defer global.LockAndRender()()
		c.LoadingForm = false
		if err != nil {
			c.FormError = true
			return
		}
		cache.Catalog.Invalidate()

		c.Init()
	}()
}

func (c *Catalog) CategoryOptions() []beercss.Option {
	options := make([]beercss.Option, len(c.Categories))

	for i, cat := range c.Categories {
		options[i] = beercss.Option{
			Label: cat.Name,
			Value: cat.ID,
		}
	}

	return options
}

func (c *Catalog) ChooseCategory(id int) {
	c.CategoryID = id
}

func (c *Catalog) DeleteCategory() {
	if !window.Confirm(`Are you sure you want to delete this category?`) {
		return
	}

	c.FormError = false
	c.LoadingForm = true

	go func() {
		err := backend.PostDeleteCategory(c.SelectedCategory.ID)
		c.LoadingForm = false
		defer global.LockAndRender()()
		if err != nil {
			c.FormError = true
			return
		}
		cache.Catalog.Invalidate()

		c.Init()
	}()
}

func (c *Catalog) DeleteItem() {
	if !window.Confirm(`Are you sure you want to delete this item?`) {
		return
	}

	c.FormError = false
	c.LoadingForm = true

	go func() {
		err := backend.PostDeleteItem(c.SelectedItem.ID)
		c.LoadingForm = false
		defer global.LockAndRender()()
		if err != nil {
			c.FormError = true
			return
		}
		cache.Catalog.Invalidate()

		c.Init()
	}()
}
