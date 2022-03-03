package pages

import (
	"fmt"

	"github.com/PotatoesFall/vecty-test/api"
	"github.com/PotatoesFall/vecty-test/domain/authdomain"
	"github.com/PotatoesFall/vecty-test/frontend/backend"
	"github.com/PotatoesFall/vecty-test/frontend/events"
	"github.com/PotatoesFall/vecty-test/frontend/global"
	"github.com/PotatoesFall/vecty-test/frontend/store"
)

// TODO make forms and buttons work
type Profile struct {
	User authdomain.User `vugu:"data"`

	CurrentPassword string
	NewPassword     string
	PasswordError   string
	PasswordSuccess string
	PasswordLoading bool

	NewName     string
	NameError   string
	NameSuccess string
	NameLoading bool
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

func (p *Profile) ChangePassword() {
	p.PasswordError = ``
	p.PasswordSuccess = ``
	p.PasswordLoading = true

	go func() {
		defer func() {
			global.EventEnv.Lock()
			p.PasswordLoading = false
			global.EventEnv.UnlockRender()
		}()

		err := backend.PostChangePassword(api.ChangePassword{
			Original: p.CurrentPassword,
			New:      p.NewPassword,
		})
		if err != nil {
			defer global.LockOnly()()
			p.PasswordError = `That didn't work.`
			return
		}

		defer global.LockOnly()()
		p.NewPassword = ``
		p.CurrentPassword = ``
		p.PasswordSuccess = `Password changed successfully.`
	}()
}

func (p *Profile) ChangeName() {
	p.NameError = ``
	p.NameSuccess = ``
	p.NameLoading = true

	go func() {
		defer func() {
			global.EventEnv.Lock()
			p.NameLoading = false
			global.EventEnv.UnlockRender()
		}()

		err := backend.PostChangeName(p.NewName)
		fmt.Println(err)
		if err != nil {
			defer global.LockOnly()()
			p.NameError = `That didn't work.`
			return
		}

		defer global.LockOnly()()
		p.NewName = ``
		p.NameSuccess = `Name changed successfully.`
	}()
}
