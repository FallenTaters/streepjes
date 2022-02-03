package catalog

import (
	"github.com/PotatoesFall/vecty-test/api"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
)

func Categories(categories []api.Category, onChange func(api.Category)) *CategoriesComponent {
	return &CategoriesComponent{
		categories:         categories,
		selectedCategoryID: -1,
		onChange:           onChange,
	}
}

type CategoriesComponent struct {
	vecty.Core

	categories         []api.Category
	selectedCategoryID int

	onChange func(api.Category)
}

func (c *CategoriesComponent) Render() vecty.ComponentOrHTML {
	markupAndChildren := []vecty.MarkupOrChild{
		vecty.Markup(vecty.Style(`overflow`, `auto`)),
		elem.Heading5(vecty.Text("Categories")),
	}

	for _, category := range c.categories {
		cat := category
		btn := categoryButton(category, category.ID == c.selectedCategoryID, func() {
			c.SetSelected(cat.ID)
			c.onChange(cat)
		})

		markupAndChildren = append(markupAndChildren, btn)
	}

	return elem.Div(markupAndChildren...)
}

func categoryButton(category api.Category, selected bool, onClick func()) vecty.ComponentOrHTML {
	classList := []string{`responsive`, `left-round`, `extra`, `small-margin`}
	if selected {
		classList = append(classList, `secondary`)
	}

	return elem.Button(
		vecty.Markup(
			vecty.Class(classList...),
			event.Click(func(*vecty.Event) { onClick() }),
		),
		vecty.Text(category.Name),
	)
}

func (c *CategoriesComponent) SetSelected(id int) {
	c.selectedCategoryID = id
	vecty.Rerender(c)
}
