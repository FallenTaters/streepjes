package club

import (
	"fmt"

	"github.com/PotatoesFall/vecty-test/domain"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

func Logo(club domain.Club, size int) *LogoComponent {
	l := &LogoComponent{
		size: size,
		club: club,
	}
	fmt.Println(`l upon creation:`, l)

	return l
}

type LogoComponent struct {
	vecty.Core

	size int
	club domain.Club
}

func (l *LogoComponent) Render() vecty.ComponentOrHTML {
	fmt.Println(`l at rendertime:`, l)

	return elem.Image(vecty.Markup(
		vecty.Style(`height`, fmt.Sprintf(`%dpx`, l.size)),
		vecty.Style(`width`, fmt.Sprintf(`%dpx`, l.size)),
		vecty.Attribute(`src`, path(l.club)),
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
