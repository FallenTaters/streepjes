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

	Page Page `vecty:"prop"`

	OrderPage   *pages.OrderComponent `vecty:"prop"`
	HistoryPage *pages.History        `vecty:"prop"`
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
			vecty.Markup(vecty.Class(`full-height`)),
			getStyles(window.OnResize(func(s window.Size) { vecty.Rerender(pv) })),
			pv.renderPage(pv.Page),
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

func (pv *PageView) renderPage(p Page) vecty.ComponentOrHTML {
	switch p {
	case PageOrder:
		return pv.OrderPage
	case PageHistory:
		return pv.HistoryPage
	}

	panic(fmt.Sprintf(`unknown page with value: %d`, p))
}

func getStyles(screenSize window.Size) vecty.MarkupList {
	switch screenSize {
	case window.SizeL:
		return vecty.Markup(vecty.Style(`padding`, `70px 50px 50px 130px`))

	case window.SizeM:
		return vecty.Markup(vecty.Style(`padding`, `20px 20px 0px 100px`))

	case window.SizeS:
		return vecty.Markup(vecty.Style(`padding`, `20px 10px 100px 10px`))
	}

	panic(screenSize)
}
