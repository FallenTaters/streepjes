//go:build wasm

package main

import (
	"github.com/PotatoesFall/vecty-test/frontend/components/layout"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

func main() {
	vecty.SetTitle("Streepjeslijst")
	vecty.AddStylesheet(`/css/order.css`)
	vecty.RenderBody(&Body{})
}

// Body is our main page component.
type Body struct {
	vecty.Core
}

// Render implements the vecty.Component interface.
func (p *Body) Render() vecty.ComponentOrHTML {
	return elem.Body(
		&layout.PageView{
			Page: layout.PageOrder,
		},
	)
}
