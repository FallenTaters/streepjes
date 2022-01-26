package pages

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

type Order struct {
	vecty.Core
}

// Render implements the vecty.Component interface.
func (o *Order) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Text(`blablabla`),
	)
}
