package club

import (
	"fmt"

	"github.com/PotatoesFall/vecty-test/domain"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

type Logo struct {
	vecty.Core

	Size   int         `vecty:"prop"`
	Margin int         `vecty:"prop"`
	Club   domain.Club `vecty:"prop"`
}

func (l *Logo) Render() vecty.ComponentOrHTML {
	return elem.Image(vecty.Markup(
		vecty.Style(`height`, fmt.Sprintf(`%dpx`, l.Size)),
		vecty.Style(`width`, fmt.Sprintf(`%dpx`, l.Size)),
		// vecty.Style(`margin`, fmt.Sprintf(`%dpx`, l.Margin)),
		// vecty.Style(`padding`, fmt.Sprintf(`%dpx`, l.Margin)),
		vecty.Style(`border-radius`, `100px`),
		vecty.Attribute(`src`, path(l.Club)),
	))
}

func path(club domain.Club) string {
	switch club {
	case domain.ClubParabool:
		return `/static/logos/parabool.jpg`
	case domain.ClubGladiators:
		return `/static/logos/gladiators.jpg`
	}

	panic(club)
}
