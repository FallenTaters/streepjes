package layout

import (
	"github.com/PotatoesFall/vecty-test/frontend/components/beercss"
	"github.com/PotatoesFall/vecty-test/frontend/jscall/window"
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

	largeScreen := window.OnResize(func() {
		vecty.Rerender(h)
	})

	var side string
	if largeScreen {
		side = `left`
	} else {
		side = `top`
	}

	return elem.Div(
		elem.Div(append([]vecty.MarkupOrChild{vecty.Markup(vecty.Class(`menu`, side))}, links...)...),
	)
}

func (h *Header) headerLink(icon beercss.IconType, text string, target Page) *vecty.HTML {
	return elem.Anchor(
		vecty.Markup(
			event.Click(func(*vecty.Event) {
				h.Navigate(target)
			}),
		),
		beercss.Icon(icon),
	)
}
