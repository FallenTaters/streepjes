package pages

import (
	"github.com/PotatoesFall/vecty-test/api"
	"github.com/PotatoesFall/vecty-test/frontend/backend"
	"github.com/PotatoesFall/vecty-test/frontend/events"
	"github.com/PotatoesFall/vecty-test/frontend/global"
	"github.com/vugu/vugu"
)

type Login struct {
	Username string `vugu:"data"`
	Password string `vugu:"data"`
	Error    bool   `vugu:"data"`
}

func (l *Login) Init(vugu.InitCtx) {
	go backend.PostLogout()
}

func (l *Login) Submit() {
	l.Error = false

	go func() {
		err := backend.PostLogin(api.Credentials{
			Username: l.Username,
			Password: l.Password,
		})
		if err != nil {
			defer global.LockAndRender()()
			l.Error = true
			return
		}

		events.Trigger(events.Login)
	}()
}
