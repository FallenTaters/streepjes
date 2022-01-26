package pages

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

type History struct {
	vecty.Core
}

// Render implements the vecty.Component interface.
func (h *History) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(vecty.Class(`container`)),
		elem.Heading1(vecty.Text(`This is the history page.`)),
	)
}
