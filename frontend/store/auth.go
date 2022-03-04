package store

import (
	"github.com/FallenTaters/streepjes/domain/authdomain"
	"github.com/FallenTaters/streepjes/frontend/events"
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
