package pages

import (
	"github.com/FallenTaters/streepjes/api"
	"github.com/FallenTaters/streepjes/domain/authdomain"
	"github.com/FallenTaters/streepjes/frontend/backend"
	"github.com/FallenTaters/streepjes/frontend/events"
	"github.com/FallenTaters/streepjes/frontend/global"
	"github.com/FallenTaters/streepjes/frontend/store"
)

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
		// do request before locking
		err := backend.PostChangePassword(api.ChangePassword{
			Original: p.CurrentPassword,
			New:      p.NewPassword,
		})

		defer global.LockAndRender()()
		defer func() { p.PasswordLoading = false }()
		if err != nil {
			p.PasswordError = `That didn't work.`
			return
		}

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
		// do request before locking
		err := backend.PostChangeName(p.NewName)

		defer global.LockAndRender()()
		defer func() { p.NameLoading = false }()

		if err != nil {
			p.NameError = `That didn't work.`
			return
		}

		p.NewName = ``
		p.NameSuccess = `Name changed successfully.`
	}()
}
