package store

import "github.com/PotatoesFall/vecty-test/domain/authdomain"

type AuthStore struct {
	LoggedIn bool
	Role     authdomain.Role
}

var Auth AuthStore

func (as *AuthStore) SetLoggedIn(role authdomain.Role) {
	as.LoggedIn = true
	as.Role = role
}
