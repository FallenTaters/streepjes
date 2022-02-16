package layout

import (
	"github.com/PotatoesFall/vecty-test/frontend/components/beercss"
	"github.com/PotatoesFall/vecty-test/frontend/jscall/window"
)

type Header struct {
	Navigate func(Page)

	Size window.Size
}

func (h *Header) Init() {
	h.Size = window.OnResize(`header`, func(s window.Size) {
		h.Size = s
	})
}

func (h *Header) menuClasses() string {
	side := `left`
	if h.Size == window.SizeS {
		side = `bottom`
	}

	return `menu ` + side
}

var headerLinks = []struct {
	text    string
	icon    beercss.IconType
	onclick func(*Header)
}{
	{
		text:    `Order`,
		icon:    beercss.IconFastfood,
		onclick: func(h *Header) { h.Navigate(PageOrder) },
	},
	{
		text:    `History`,
		icon:    beercss.IconHistory,
		onclick: func(h *Header) { h.Navigate(PageHistory) },
	},
	{
		text: `Profile`,
		icon: beercss.IconPerson,
		onclick: func(h *Header) {
			// TODO show profile options (dropdown ?)
		},
	},
}

// import (
// 	"github.com/PotatoesFall/vecty-test/frontend/components/beercss"
// 	"github.com/PotatoesFall/vecty-test/frontend/jscall/window"
// 	"github.com/hexops/vecty"
// 	"github.com/hexops/vecty/elem"
// 	"github.com/hexops/vecty/event"
// )

// type Header struct {
// 	vecty.Core

// 	Navigate func(Page)
// }

// // Render implements the vecty.Component interface.
// func (h *Header) Render() vecty.ComponentOrHTML {
// 	links := []vecty.MarkupOrChild{
// 		h.headerLink(beercss.IconFastfood, `Order`, h.navigate(PageOrder)),
// 		h.headerLink(beercss.IconHistory, `History`, h.navigate(PageHistory)),
// 		h.headerLink(beercss.IconPerson, `Profile`, func(e *vecty.Event) {}), // TODO show profile options (dropdown?)
// 	}

// 	size := window.GetSize()
// 	side := `left`
// 	if size == window.SizeS {
// 		side = `bottom`
// 	}

// 	return elem.Div(
// 		elem.Div(append([]vecty.MarkupOrChild{vecty.Markup(vecty.Class(`menu`, side))}, links...)...),
// 	)
// }

// func (h *Header) headerLink(icon beercss.IconType, text string, onClick func(*vecty.Event)) *vecty.HTML {
// 	return elem.Anchor(
// 		vecty.Markup(event.Click(onClick)),
// 		beercss.Icon(icon),
// 	)
// }

// func (h *Header) navigate(p Page) func(*vecty.Event) {
// 	return func(e *vecty.Event) {
// 		go h.Navigate(p)
// 	}
// }
