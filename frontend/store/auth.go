package store

import (
	"github.com/PotatoesFall/vecty-test/domain/authdomain"
	"github.com/PotatoesFall/vecty-test/frontend/events"
)

type AuthStore struct {
	LoggedIn bool
	User     authdomain.User
}

var Auth AuthStore

func Init() {
	events.Listen(events.Unauthorized, `auth-store`, func() {
		Auth.LoggedIn = false
		Auth.User = authdomain.User{}
	})
}

func (auth *AuthStore) LogIn(user authdomain.User) {
	auth.LoggedIn = true
	auth.User = user
}
