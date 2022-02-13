//go:build js && wasm

package main

import (
	"time"

	"github.com/PotatoesFall/vecty-test/frontend/backend"
	"github.com/PotatoesFall/vecty-test/frontend/components/layout"
	"github.com/PotatoesFall/vecty-test/frontend/jscall/window"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

func main() {
	initPackages()

	vecty.SetTitle("Streepjeslijst")
	vecty.RenderBody(&Body{})

	go func() {
		for {
			time.Sleep(time.Second)
		}
	}()
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
			Page: layout.PageLogin,
		},
	)
}

func initPackages() {
	backend.Init(window.Location())
}
