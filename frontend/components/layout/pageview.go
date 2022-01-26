package layout

import (
	"fmt"

	"github.com/PotatoesFall/vecty-test/frontend/components/pages"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

type PageView struct {
	vecty.Core

	Page Page
}

// Render implements the vecty.Component interface.
func (pv *PageView) Render() vecty.ComponentOrHTML {
	return elem.Div(
		&Header{
			Navigate: func(p Page) {
				pv.Page = p
				vecty.Rerender(pv)
			},
		},
		elem.Div(
			vecty.Markup(vecty.Class(`container`)),
			renderPage(pv.Page),
		),
	)
}

type Page int

const (
	PageOrder Page = iota + 1
	PageHistory

	PageCatalog
	PageMembers
	PageUsers
)

func renderPage(p Page) vecty.ComponentOrHTML {
	switch p {
	case PageOrder:
		return &pages.Order{}
	case PageHistory:
		return &pages.History{}
	}

	panic(fmt.Sprintf(`unknown page with value: %d`, p))
}
