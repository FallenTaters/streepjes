package pages

import (
	"fmt"

	"github.com/PotatoesFall/vecty-test/domain/authdomain"
	"github.com/PotatoesFall/vecty-test/frontend/events"
	"github.com/PotatoesFall/vecty-test/frontend/store"
)

type Profile struct {
	User authdomain.User `vugu:"data"`
}

func (p *Profile) Init() {
	fmt.Println(`init`)
	if !store.Auth.LoggedIn {
		events.Trigger(events.Unauthorized)
		return
	}

	p.User = store.Auth.User
}

func (p *Profile) Compute() {
	fmt.Println(`compute`)
}
