package catalog

import (
	"fmt"

	"github.com/PotatoesFall/vecty-test/domain"
)

type Categories struct {
	Categories         []domain.Category
	SelectedCategoryID int
	OnClick            func(domain.Category)
}

func (c *Categories) classes(category domain.Category) string {
	fmt.Println(`pointer in Categories:`, c)

	fmt.Println(c.SelectedCategoryID, category.ID)
	classes := `responsive extra small-margin`

	if c.SelectedCategoryID == category.ID {
		return classes + ` secondary`
	}

	return classes
}

// import (
// 	"github.com/PotatoesFall/vecty-test/domain"
// 	"github.com/hexops/vecty"
// 	"github.com/hexops/vecty/elem"
// 	"github.com/hexops/vecty/event"
// )

// func Categories(categories []domain.Category, selectedCategoryID int, onChange func(domain.Category)) *CategoriesComponent {
// 	return &CategoriesComponent{
// 		Categories:         categories,
// 		SelectedCategoryID: selectedCategoryID,
// 		OnChange:           onChange,
// 	}
// }

// type CategoriesComponent struct {
// 	vecty.Core

// 	Categories         []domain.Category `vecty:"prop"`
// 	SelectedCategoryID int               `vecty:"prop"`

// 	OnChange func(domain.Category) `vecty:"prop"`
// }

// func (c *CategoriesComponent) Render() vecty.ComponentOrHTML {
// 	markupAndChildren := []vecty.MarkupOrChild{}

// 	for _, category := range c.Categories {
// 		btn := categoryButton(category, category.ID == c.SelectedCategoryID, c.OnChange)

// 		markupAndChildren = append(markupAndChildren, btn)
// 	}

// 	return elem.Div(markupAndChildren...)
// }

// func categoryButton(category domain.Category, selected bool, onClick func(c domain.Category)) *vecty.HTML {
// 	classList := []string{`responsive`, `extra`, `small-margin`}
// 	if selected {
// 		classList = append(classList, `secondary`)
// 	}

// 	return elem.Button(
// 		vecty.Markup(
// 			vecty.Class(classList...),
// 			event.Click(func(*vecty.Event) { onClick(category) }),
// 		),
// 		vecty.Text(category.Name),
// 	)
// }
