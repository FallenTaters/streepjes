package catalog

import (
	"github.com/PotatoesFall/vecty-test/api"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

func Categories(categories []api.Category) *CategoriesComponent {
	// todo add listeners
	return &CategoriesComponent{
		categories: categories,
	}
}

type CategoriesComponent struct {
	vecty.Core

	categories []api.Category
}

func (c *CategoriesComponent) Render() vecty.ComponentOrHTML {
	markupAndChildren := []vecty.MarkupOrChild{
		vecty.Markup(vecty.Style(`overflow`, `auto`)),
		elem.Heading2(vecty.Text("Categories")),
	}

	for _, category := range c.categories {
		markupAndChildren = append(markupAndChildren, elem.Div(vecty.Text(category.Name)))
	}

	return elem.Div(markupAndChildren...)
}
