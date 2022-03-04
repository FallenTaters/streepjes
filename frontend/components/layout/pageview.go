package layout

import (
	"github.com/FallenTaters/streepjes/domain/authdomain"
	"github.com/FallenTaters/streepjes/frontend/events"
	"github.com/FallenTaters/streepjes/frontend/global"
	"github.com/FallenTaters/streepjes/frontend/jscall/window"
	"github.com/FallenTaters/streepjes/frontend/store"
	"github.com/vugu/vugu"
)

type Page int

const (
	PageLogin Page = iota + 1
	PageProfile

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

		switch store.Auth.User.Role {
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
