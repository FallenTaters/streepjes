package order

import (
	"github.com/PotatoesFall/vecty-test/frontend/components/beercss"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
)

type Small struct {
	vecty.Core

	Page SubPage
}

func (s *Small) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(vecty.Class(`container`)),
		elem.Div(vecty.Markup(vecty.Class(`menu`, `bottom`)),
			s.footerLink(beercss.IconTypeSwapHoriz, func(s *Small) {}),
			s.footerLink(beercss.IconTypeAddCircle, func(s *Small) { s.Page = SubPageCategories }),
			s.footerLink(beercss.IconTypePayments, func(s *Small) { s.Page = SubPagePayment }),
		),
		elem.Div(vecty.Text(`TODO: make subComponents`)),
		elem.Button(vecty.Text(`TODO: make subComponents`)),
		elem.Button(
			vecty.Markup(vecty.Style(`background-color`, `var(--secondary)`)),
			vecty.Text(`TODO: make subComponents`)),
	)
}

type SubPage int

const (
	SubPageCategories SubPage = iota + 1
	SubPageProducts
	SubPageOrderOverview
	SubPagePayment
)

func (s *Small) footerLink(icon beercss.IconType, f func(*Small)) *vecty.HTML {
	return elem.Anchor(
		vecty.Markup(
			event.Click(func(*vecty.Event) {
				f(s)
				vecty.Rerender(s)
			}),
		),
		beercss.Icon(icon),
	)
}
