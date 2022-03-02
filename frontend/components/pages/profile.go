package pages

import (
	"github.com/PotatoesFall/vecty-test/domain/authdomain"
	"github.com/PotatoesFall/vecty-test/frontend/backend"
	"github.com/PotatoesFall/vecty-test/frontend/events"
	"github.com/PotatoesFall/vecty-test/frontend/store"
)

// TODO make forms and buttons work
type Profile struct {
	User authdomain.User `vugu:"data"`

	CurrentPassword string
	NewPassword     string

	NewName string
}

func (p *Profile) Init() {
	if !store.Auth.LoggedIn {
		events.Trigger(events.Unauthorized)
		return
	}

	p.User = store.Auth.User
}

func (p *Profile) Logout() {
	go backend.PostLogout()
	events.Trigger(events.Unauthorized)
}

func (p *Profile) ChangeName() {
	// TODO implement and handle errors
}

func (p *Profile) ChangePassword() {
	// TODO implement and handle errors
}
