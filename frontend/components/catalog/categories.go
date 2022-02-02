package catalog

import (
	"github.com/PotatoesFall/vecty-test/api"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
)

func Categories(categories []api.Category, onChange func(api.Category)) *CategoriesComponent {
	return &CategoriesComponent{
		onChange:   onChange,
		categories: categories,
	}
}

type CategoriesComponent struct {
	vecty.Core

	onChange   func(api.Category)
	categories []api.Category
}

func (c *CategoriesComponent) Render() vecty.ComponentOrHTML {
	markupAndChildren := []vecty.MarkupOrChild{
		vecty.Markup(vecty.Style(`overflow`, `auto`)),
		elem.Heading2(vecty.Text("Categories")),
	}

	for _, category := range c.categories {
		categoryButton := elem.Button(
			vecty.Markup(
				vecty.Class(`responsive`), // TOOD add padding or margin
				event.Click(func(*vecty.Event) {
					c.onChange(category)
				}),
			),
			vecty.Text(category.Name),
		)

		markupAndChildren = append(markupAndChildren, categoryButton)
	}

	return elem.Div(markupAndChildren...)
}
