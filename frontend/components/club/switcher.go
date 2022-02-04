package club

import (
	"github.com/PotatoesFall/vecty-test/domain"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
)

func Switcher(rerender bool, club domain.Club, onToggle func(domain.Club)) *SwitcherComponent {
	return &SwitcherComponent{
		rerender: rerender,
		onToggle: onToggle,
		club:     club,
	}
}

type SwitcherComponent struct {
	vecty.Core

	rerender bool
	onToggle func(domain.Club)
	club     domain.Club
}

func (s *SwitcherComponent) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			event.Click(func(e *vecty.Event) {
				s.club = otherClub(s.club)
				s.onToggle(s.club)
				if s.rerender {
					vecty.Rerender(s)
				}
			}),
		),
		Logo(s.club, 150),
	)
}

func otherClub(c domain.Club) domain.Club {
	if c == domain.ClubGladiators {
		return domain.ClubParabool
	}

	return domain.ClubGladiators
}
