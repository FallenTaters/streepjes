package club

import (
	"fmt"

	"github.com/PotatoesFall/vecty-test/domain"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

func Logo(club domain.Club, size int) *LogoComponent {
	return &LogoComponent{
		Club: club,
		Size: size,
	}
}

type LogoComponent struct {
	vecty.Core

	Size int         `vecty:"prop"`
	Club domain.Club `vecty:"prop"`
}

func (l *LogoComponent) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			vecty.Style(`background-color`, `white`),
			vecty.Style(`background-position`, `center`),
			vecty.Style(`background-repeat`, `no-repeat`),
			vecty.Style(`background-image`, path(l.Club)),
			vecty.Style(`background-size`, px(l.Size, 1)),
			vecty.Style(`height`, px(l.Size, 1)),
			vecty.Style(`width`, px(l.Size, 1)),
			vecty.Style(`padding`, px(l.Size, 0.207)),
			vecty.Style(`border-radius`, px(l.Size, 0.707)),
			vecty.Style(`border`, `none`),
		),
	)
}

func path(club domain.Club) string {
	switch club {
	case domain.ClubParabool:
		return `url("/static/logos/parabool.jpg")`
	case domain.ClubGladiators:
		return `url("/static/logos/gladiators.jpg")`
	}

	panic(club)
}

func px(n int, factor float64) string {
	n = int(float64(n) * factor)
	return fmt.Sprintf(`%dpx`, n)
}
