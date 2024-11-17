package club

import (
	"fmt"

	"github.com/FallenTaters/streepjes/domain"
)

type Logo struct {
	Size int         `vugu:"data"`
	Club domain.Club `vugu:"data"`
}

func (l *Logo) style() string {
	return `background-color: white;` +
		`background-position: center;` +
		`background-repeat: no-repeat;` +
		`background-image: ` + l.path() + `;` +
		`background-size: ` + px(l.Size, 1) + `;` +
		`height: ` + px(l.Size, 1) + `;` +
		`width: ` + px(l.Size, 1) + `;` +
		`padding: ` + px(l.Size, 0.207) + `;` +
		`border-radius:` + px(l.Size, 0.707) + `;` +
		`border: none;`
}

func (l *Logo) path() string {
	switch l.Club {
	case domain.ClubUnknown:
		return ``
	case domain.ClubParabool:
		return `url("/static/logos/parabool.jpg")`
	case domain.ClubGladiators:
		return `url("/static/logos/gladiators.jpg")`
	case domain.ClubCalamari:
		return `url("/static/logos/calamari.jpg")`
	}

	return ``
}

func px(n int, factor float64) string {
	n = int(float64(n) * factor)
	return fmt.Sprintf(`%dpx`, n)
}
