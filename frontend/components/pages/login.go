package pages

import (
	"github.com/PotatoesFall/vecty-test/frontend/backend"
	"github.com/PotatoesFall/vecty-test/frontend/components/beercss"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

func Login() (vecty.Component, error) {
	err := backend.Logout()
	if err != nil {
		return nil, err
	}

	return &loginComponent{}, nil
}

type loginComponent struct {
	vecty.Core
}

// Render implements the vecty.Component interface.
func (h *loginComponent) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(vecty.Class(`container`)),
		elem.Form(
			elem.Div(
				beercss.Input(`Username`),
			),
		),
	)
}
