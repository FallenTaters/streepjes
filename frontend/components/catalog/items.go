package catalog

import (
	"github.com/PotatoesFall/vecty-test/domain"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
)

func Items(items []domain.Item, onChange func(domain.Item)) *ItemsComponent {
	return &ItemsComponent{
		items:          items,
		selectedItemID: -1,
		onChange:       onChange,
	}
}

type ItemsComponent struct {
	vecty.Core

	items          []domain.Item
	selectedItemID int

	onChange func(domain.Item)
}

func (i *ItemsComponent) Render() vecty.ComponentOrHTML {
	markupAndChildren := []vecty.MarkupOrChild{
		elem.Heading5(vecty.Text("Items")),
	}

	for _, item := range i.items {
		btn := itemButton(item, item.ID == i.selectedItemID, i.onChange)

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
	i.items = items
	vecty.Rerender(i)
}

func (i *ItemsComponent) SetSelected(id int) {
	i.selectedItemID = id
	vecty.Rerender(i)
}
