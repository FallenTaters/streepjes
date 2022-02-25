package layout

import (
	"github.com/PotatoesFall/vecty-test/domain/authdomain"
	"github.com/PotatoesFall/vecty-test/frontend/jscall/window"
	"github.com/PotatoesFall/vecty-test/frontend/store"
)

type Header struct {
	Navigate func(Page)

	Size window.Size
}

func (h *Header) Init() {
	h.Size = window.OnResize(`header`, func(s window.Size) {
		h.Size = s
	})
}

func (h *Header) menuClasses() string {
	side := `left`
	if h.Size == window.SizeS {
		side = `bottom`
	}

	return `menu ` + side
}

func (*Header) showAdminPages() bool {
	return store.Auth.User.Role == authdomain.RoleAdmin
}

func (*Header) username() string {
	username := store.Auth.User.Username
	if len(username) > 10 {
		username = username[:8] + `…`
	}

	return username
}
