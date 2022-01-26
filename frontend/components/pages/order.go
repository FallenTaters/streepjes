package pages

import (
	"github.com/hexops/vecty"

	"github.com/PotatoesFall/vecty-test/frontend/components/pages/order"
	"github.com/PotatoesFall/vecty-test/frontend/jscall/window"
)

type Order struct {
	vecty.Core
}

func (o *Order) Render() vecty.ComponentOrHTML {
	largeScreen := window.OnResize(func() {
		vecty.Rerender(o)
	})

	if largeScreen {
		return &order.Large{}
	}
	return &order.Small{}
}
