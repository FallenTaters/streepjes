package order

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
)

type Small struct {
	vecty.Core

	resize func(*vecty.Event)
}

func (s *Small) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			vecty.Class(`full-height`, `order-grid-sdfsdfsdfsdf`),
			event.Resize(s.resize),
		),
		elem.Div(vecty.Text(`SMALL`)),
	)
}
