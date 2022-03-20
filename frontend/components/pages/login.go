package pages

import (
	"github.com/FallenTaters/streepjes/api"
	"github.com/FallenTaters/streepjes/frontend/backend"
	"github.com/FallenTaters/streepjes/frontend/events"
	"github.com/FallenTaters/streepjes/frontend/global"
	"github.com/FallenTaters/streepjes/frontend/store"
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
		// execute request before locking
		user, err := backend.PostLogin(api.Credentials{
			Username: l.Username,
			Password: l.Password,
		})

		defer global.LockAndRender()()
		defer func() { l.Loading = false }()

		if err != nil {
			l.Error = true
			return
		}

		store.Auth.LogIn(user)
		store.Order.Club = user.Club
		store.Order.Lines = []store.Orderline{}

		events.Trigger(events.Login)
	}()
}
