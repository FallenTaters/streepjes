package layout

import (
	"fmt"

	"github.com/PotatoesFall/vecty-test/frontend/components/pages"
	"github.com/PotatoesFall/vecty-test/frontend/jscall/window"
	"github.com/PotatoesFall/vecty-test/frontend/store"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

type PageView struct {
	vecty.Core

	Page Page `vecty:"prop"`
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
			getStyles(window.OnResize(func(s window.Size) { go vecty.Rerender(pv) })),
			pv.renderPage(pv.Page),
		),
	)
}

type Page int

const (
	PageLogin Page = iota + 1

	PageOrder
	PageHistory

	PageCatalog
	PageMembers
	PageUsers
)

func (pv *PageView) renderPage(p Page) vecty.ComponentOrHTML {
	switch p {
	case PageLogin:
		return login()
	case PageOrder:
		return order()
	case PageHistory:
		return &pages.History{}
	}

	fmt.Printf(`unknown page with value: %d`, p)
	return nil
}

func getStyles(screenSize window.Size) vecty.MarkupList {
	switch screenSize {
	case window.SizeL:
		return vecty.Markup(vecty.Style(`padding`, `20px 40px 20px 120px`))

	case window.SizeM:
		return vecty.Markup(vecty.Style(`padding`, `20px 20px 0px 100px`))

	case window.SizeS:
		return vecty.Markup(vecty.Style(`padding`, `20px 10px 50px 10px`))
	}

	return vecty.Markup()
}

var orderComponent vecty.Component

func order() vecty.Component {
	if orderComponent != nil {
		return orderComponent
	}

	component, err := pages.Order()
	if err != nil {
		return pages.Error(err.Error())
	}
	orderComponent = component

	store.Order.OnChange = func(oe store.OrderEvent) {
		vecty.Rerender(orderComponent)
	}

	return orderComponent
}

func login() vecty.Component {
	component, err := pages.Login()
	if err != nil {
		return pages.Error(err.Error())
	}

	return component
}
