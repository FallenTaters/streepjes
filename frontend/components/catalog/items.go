package catalog

import (
	"github.com/PotatoesFall/vecty-test/api"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
)

func Items(items []api.Item, onChange func(api.Item)) *ItemsComponent {
	return &ItemsComponent{
		items:          items,
		selectedItemID: -1,
		onChange:       onChange,
	}
}

type ItemsComponent struct {
	vecty.Core

	items          []api.Item
	selectedItemID int

	onChange func(api.Item)
}

func (i *ItemsComponent) Render() vecty.ComponentOrHTML {
	markupAndChildren := []vecty.MarkupOrChild{
		vecty.Markup(vecty.Style(`overflow`, `auto`)),
		elem.Heading5(vecty.Text("Items")),
	}

	for _, item := range i.items {
		itm := item
		btn := itemButton(item, item.ID == i.selectedItemID, func() {
			i.onChange(itm)
		})

		markupAndChildren = append(markupAndChildren, btn)
	}

	return elem.Div(markupAndChildren...)
}

func itemButton(item api.Item, selected bool, onClick func()) vecty.ComponentOrHTML {
	classList := []string{`responsive`, `extra`, `small-margin`}
	if selected {
		classList = append(classList, `secondary`)
	}

	return elem.Button(
		vecty.Markup(
			vecty.Class(classList...),
			event.Click(func(*vecty.Event) { onClick() }),
		),
		vecty.Text(item.Name),
	)
}

func (i *ItemsComponent) SetItems(items []api.Item) {
	i.items = items
	vecty.Rerender(i)
}

func (i *ItemsComponent) SetSelected(id int) {
	i.selectedItemID = id
	vecty.Rerender(i)
}
