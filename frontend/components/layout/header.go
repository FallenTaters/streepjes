package layout

import (
	"github.com/PotatoesFall/vecty-test/frontend/components/beercss"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
)

type Header struct {
	vecty.Core

	Navigate func(Page)
}

// Render implements the vecty.Component interface.
func (h *Header) Render() vecty.ComponentOrHTML {
	links := []vecty.MarkupOrChild{
		h.headerLink(beercss.IconTypeLocalBar, `Order`, PageOrder),
		h.headerLink(beercss.IconTypeHistory, `History`, PageHistory),
	}

	return elem.Div(vecty.List{
		elem.Div(append([]vecty.MarkupOrChild{vecty.Markup(vecty.Class(`menu`, `m`, `l`, `top`))}, links...)...),
		elem.Div(append([]vecty.MarkupOrChild{vecty.Markup(vecty.Class(`menu`, `s`, `bottom`))}, links...)...),
	})
}

func (h *Header) headerLink(icon beercss.IconType, text string, target Page) *vecty.HTML {
	return elem.Anchor(
		vecty.Markup(
			event.Click(func(*vecty.Event) {
				h.Navigate(target)
			}),
		),
		beercss.Icon(icon),
		elem.Div(
			vecty.Markup(vecty.Class(`m`, `l`)),
			vecty.Text(text),
		),
	)
}
