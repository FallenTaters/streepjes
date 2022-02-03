package layout

import (
	"fmt"

	"github.com/PotatoesFall/vecty-test/frontend/components/pages"
	"github.com/PotatoesFall/vecty-test/frontend/jscall/window"
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
			getStyles(window.OnResize(func() {
				vecty.Rerender(pv)
			})),
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
		return pages.Order()
	case PageHistory:
		return &pages.History{}
	}

	panic(fmt.Sprintf(`unknown page with value: %d`, p))
}

func getStyles(largeScreen bool) vecty.MarkupList {
	if largeScreen {
		return vecty.Markup(vecty.Style(`padding`, `20px 20px 20px 100px`))
	}

	return vecty.Markup(vecty.Style(`padding`, `10px 10px 100px 10px`))
}
