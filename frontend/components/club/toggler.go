package club

import (
	"github.com/PotatoesFall/vecty-test/domain"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
)

type Toggler struct {
	vecty.Core

	Rerender bool              `vecty:"prop"`
	OnToggle func(domain.Club) `vecty:"prop"`
	Club     domain.Club       `vecty:"prop"`
}

func (s *Toggler) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			vecty.Class(`col`, `min`),
			event.Click(func(e *vecty.Event) {
				s.toggle()
			}),
		),
		&Logo{Size: 150, Margin: 100, Club: s.Club},
	)
}

func (s *Toggler) toggle() {
	s.Club = otherClub(s.Club)
	if s.Rerender {
		vecty.Rerender(s)
	}
	s.OnToggle(s.Club)
}

func otherClub(c domain.Club) domain.Club {
	if c == domain.ClubGladiators {
		return domain.ClubParabool
	}

	return domain.ClubGladiators
}
