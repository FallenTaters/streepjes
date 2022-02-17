package layout

import (
	"github.com/PotatoesFall/vecty-test/frontend/jscall/window"
	"github.com/vugu/vugu"
)

type Page int

const (
	PageLogin Page = iota + 1

	PageOrder
	PageHistory

	PageCatalog
	PageMembers
	PageUsers
)

type Pageview struct {
	Page Page
}

func (pv *Pageview) Init(vugu.InitCtx) {
	pv.Page = PageOrder
}

func getStyles() string {
	switch window.GetSize() {
	case window.SizeL:
		return `padding: 20px 40px 20px 120px;`

	case window.SizeM:
		return `padding: 20px 20px 0px 100px;`

	case window.SizeS:
		return `padding: 20px 10px 50px 10px;`
	}

	return ``
}

// type PageView struct {
// 	vecty.Core

// 	Page Page `vecty:"prop"`
// }

// // Render implements the vecty.Component interface.
// func (pv *PageView) Render() vecty.ComponentOrHTML {
// 	return elem.Div(
// 		&Header{
// 			Navigate: func(p Page) {
// 				pv.Page = p
// 				vecty.Rerender(pv)
// 			},
// 		},
// 		elem.Div(
// 			vecty.Markup(vecty.Class(`full-height`)),
// 			getStyles(window.OnResize(func(s window.Size) { go vecty.Rerender(pv) })),
// 			pv.renderPage(pv.Page),
// 		),
// 	)
// }

// type Page int

// const (
// 	PageLogin Page = iota + 1

// 	PageOrder
// 	PageHistory

// 	PageCatalog
// 	PageMembers
// 	PageUsers
// )

// func (pv *PageView) renderPage(p Page) vecty.ComponentOrHTML {
// 	switch p {
// 	case PageLogin:
// 		return login()
// 	case PageOrder:
// 		return order()
// 	case PageHistory:
// 		return &pages.History{}
// 	}

// 	fmt.Printf(`unknown page with value: %d`, p)
// 	return nil
// }

// var orderComponent vecty.Component

// func order() vecty.Component {
// 	if orderComponent != nil {
// 		return orderComponent
// 	}

// 	component, err := pages.Order()
// 	if err != nil {
// 		return pages.Error(err.Error())
// 	}
// 	orderComponent = component

// 	store.Order.OnChange = func(oe store.OrderEvent) {
// 		vecty.Rerender(orderComponent)
// 	}

// 	return orderComponent
// }

// func login() vecty.Component {
// 	component, err := pages.Login()
// 	if err != nil {
// 		return pages.Error(err.Error())
// 	}

// 	return component
// }
