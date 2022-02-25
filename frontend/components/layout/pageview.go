package layout

import (
	"github.com/PotatoesFall/vecty-test/domain/authdomain"
	"github.com/PotatoesFall/vecty-test/frontend/events"
	"github.com/PotatoesFall/vecty-test/frontend/global"
	"github.com/PotatoesFall/vecty-test/frontend/jscall/window"
	"github.com/PotatoesFall/vecty-test/frontend/store"
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
	Page       Page `vugu:"data"`
	ShowHeader bool `vugu:"data"`
}

func (pv *Pageview) Init(vugu.InitCtx) {
	events.Listen(events.Unauthorized, `pageview`, func() {
		defer global.LockAndRender()()

		pv.Page = PageLogin
	})

	events.Listen(events.Login, `pageview`, func() {
		defer global.LockAndRender()()

		switch store.Auth.Role {
		case authdomain.RoleAdmin:
			pv.Page = PageMembers
		case authdomain.RoleBartender:
			pv.Page = PageOrder
		}
	})

	pv.Page = PageOrder
}

func (pv *Pageview) GetStyles() string {
	if pv.Page == PageLogin {
		return ``
	}

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
