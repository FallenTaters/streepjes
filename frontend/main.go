//go:build js && wasm

package main

import (
	"github.com/PotatoesFall/vecty-test/frontend/backend"
	"github.com/PotatoesFall/vecty-test/frontend/components/layout"
	"github.com/PotatoesFall/vecty-test/frontend/components/pages"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

func main() {
	backend.Init(`http://localhost:8080`) // TODO: make setting or automatically get current location

	vecty.SetTitle("Streepjeslijst")
	vecty.RenderBody(&Body{})
}

// Body is our main page component.
type Body struct {
	vecty.Core
}

// Render implements the vecty.Component interface.
func (p *Body) Render() vecty.ComponentOrHTML {
	return elem.Body(
		vecty.Markup(vecty.Class(`is-dark`)),
		&layout.PageView{
			Page: layout.PageOrder,

			OrderPage:   pages.Order(),
			HistoryPage: &pages.History{},
		},
	)
}
