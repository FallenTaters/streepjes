package pages

import (
	"github.com/PotatoesFall/vecty-test/api"
	"github.com/PotatoesFall/vecty-test/frontend/backend"
	"github.com/PotatoesFall/vecty-test/frontend/events"
	"github.com/PotatoesFall/vecty-test/frontend/global"
	"github.com/PotatoesFall/vecty-test/frontend/store"
	"github.com/vugu/vugu"
)

type Login struct {
	Username string `vugu:"data"`
	Password string `vugu:"data"`
	Error    bool   `vugu:"data"`
	Loading  bool   `vugu:"data"`
}

func (l *Login) Init(vugu.InitCtx) {
	go backend.PostLogout()
}

func (l *Login) Submit() {
	l.Error = false
	l.Loading = true

	go func() {
		defer func() {
			global.EventEnv.Lock()
			l.Loading = false
			global.EventEnv.UnlockRender()
		}()

		resp, err := backend.PostLogin(api.Credentials{
			Username: l.Username,
			Password: l.Password,
		})
		if err != nil {
			defer global.LockOnly()()
			l.Error = true
			return
		}

		store.Auth.SetLoggedIn(resp.Role)

		events.Trigger(events.Login)
	}()
}
