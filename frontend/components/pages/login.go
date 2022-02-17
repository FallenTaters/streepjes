package pages

import (
	"github.com/PotatoesFall/vecty-test/frontend/backend"
	"github.com/vugu/vugu"
)

type Login struct{}

func (l *Login) Init(vugu.InitCtx) {
	go backend.Logout()
}
