package catalog

import (
	"github.com/PotatoesFall/vecty-test/domain"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
)

func Items(items []domain.Item, onChange func(domain.Item)) *ItemsComponent {
	return &ItemsComponent{
		Items:          items,
		SelectedItemID: -1,
		OnChange:       onChange,
	}
}

type ItemsComponent struct {
	vecty.Core

	Items          []domain.Item `vecty:"prop"`
	SelectedItemID int           `vecty:"prop"`

	OnChange func(domain.Item) `vecty:"prop"`
}

func (i *ItemsComponent) Render() vecty.ComponentOrHTML {
	markupAndChildren := []vecty.MarkupOrChild{}

	for _, item := range i.Items {
		btn := itemButton(item, item.ID == i.SelectedItemID, i.OnChange)

		markupAndChildren = append(markupAndChildren, btn)
	}

	return elem.Div(markupAndChildren...)
}

func itemButton(item domain.Item, selected bool, onClick func(i domain.Item)) vecty.ComponentOrHTML {
	classList := []string{`responsive`, `extra`, `small-margin`}
	if selected {
		classList = append(classList, `secondary`)
	}

	return elem.Button(
		vecty.Markup(
			vecty.Class(classList...),
			event.Click(func(*vecty.Event) { onClick(item) }),
		),
		vecty.Text(item.Name),
	)
}

func (i *ItemsComponent) SetItems(items []domain.Item) {
	i.Items = items
	vecty.Rerender(i)
}

func (i *ItemsComponent) SetSelected(id int) {
	i.SelectedItemID = id
	vecty.Rerender(i)
}
