package pages

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/FallenTaters/streepjes/api"
	"github.com/FallenTaters/streepjes/domain"
	"github.com/FallenTaters/streepjes/domain/authdomain"
	"github.com/FallenTaters/streepjes/frontend/backend"
	"github.com/FallenTaters/streepjes/frontend/backend/cache"
	"github.com/FallenTaters/streepjes/frontend/global"
	"github.com/FallenTaters/streepjes/frontend/jscall/document"
	"github.com/vugu/vugu"
)

type Users struct {
	Loading bool `vugu:"data"`
	Error   bool `vugu:"data"`

	SubmitLoading bool   `vugu:"data"`
	SubmitError   string `vugu:"data"`
	DeleteConfirm bool   `vugu:"data"`

	Users        []authdomain.User `vugu:"data"`
	SelectedUser authdomain.User   `vugu:"data"`
	ShowForm     bool              `vugu:"data"`

	Username string          `vugu:"data"`
	Password string          `vugu:"data"`
	Name     string          `vugu:"data"`
	Role     authdomain.Role `vugu:"data"`
	Club     domain.Club     `vugu:"data"`
}

func (u *Users) Init() {
	u.Loading = true
	u.Error = false

	u.Username = ``
	u.Password = ``
	u.Name = ``
	u.Role = authdomain.RoleNotAuthorized
	u.Club = domain.ClubUnknown
	u.ShowForm = false
	u.DeleteConfirm = false

	go func() {
		users, err := cache.Users.Get()
		defer global.LockAndRender()()
		defer func() {
			u.Loading = false
		}()
		if err != nil {
			u.Error = true
			return
		}

		sort.Slice(users, func(i, j int) bool {
			n1, n2 := strings.ToLower(users[i].Name), strings.ToLower(users[j].Name)
			return strings.Compare(n1, n2) < 0
		})

		u.Users = users
	}()
}

func (u *Users) SelectUser(user authdomain.User) {
	u.SubmitError = ``
	u.DeleteConfirm = false

	u.SelectedUser = user

	u.Name = user.Name
	u.Username = user.Username
	u.Password = ``
	u.Role = user.Role
	u.Club = user.Club

	go u.selectClub(user.Club)
	go u.selectRole(user.Role)

	u.ShowForm = true
}

func (u *Users) NewUser() {
	u.SubmitError = ``
	u.DeleteConfirm = false

	u.SelectedUser = authdomain.User{}

	u.Name = ``
	u.Username = ``
	u.Password = ``
	u.Role = authdomain.RoleNotAuthorized
	u.Club = domain.ClubUnknown

	u.ShowForm = true

	go u.selectClub(u.Club)
	go u.selectRole(u.Role)
}

func (u *Users) FormTitle() string {
	if u.SelectedUser == (authdomain.User{}) {
		return `New User`
	}

	return `Edit User: ` + u.SelectedUser.Name
}

func (u *Users) SaveButtonText() string {
	if u.SelectedUser == (authdomain.User{}) {
		return `Add User`
	}

	return `Save Changes`
}

func (u *Users) Submit() {
	u.SubmitError = ``

	if u.SelectedUser == (authdomain.User{}) {
		u.submitNewUser()
		return
	}

	u.submitChanges()
}

func (u *Users) submitNewUser() {
	if u.Username == `` ||
		u.Password == `` ||
		u.Name == `` ||
		u.Role == authdomain.RoleNotAuthorized ||
		u.Club == domain.ClubUnknown {

		u.SubmitError = `All fields must be filled!`
		return
	}

	u.SubmitLoading = true
	go func() {
		defer func() {
			defer global.LockAndRender()()
			u.SubmitLoading = false
		}()

		err := backend.PostNewUser(api.UserWithPassword{
			Password: u.Password,
			User: authdomain.User{
				Username: u.Username,
				Club:     u.Club,
				Name:     u.Name,
				Role:     u.Role,
			},
		})
		if err != nil {
			u.SubmitError = `Unable to create new user. Maybe the username or name are already taken.`
			return
		}

		cache.Users.Invalidate()
		defer global.LockAndRender()()

		u.Init()
	}()
}

func (u *Users) submitChanges() {
	if u.Username == `` ||
		u.Name == `` ||
		u.Role == authdomain.RoleNotAuthorized ||
		u.Club == domain.ClubUnknown {

		u.SubmitError = `All fields must be filled!`
		return
	}

	u.SubmitLoading = true
	go func() {
		defer func() {
			defer global.LockAndRender()()
			u.SubmitLoading = false
		}()

		err := backend.PostEditUser(api.UserWithPassword{
			Password: u.Password,
			User: authdomain.User{
				ID:       u.SelectedUser.ID,
				Username: u.Username,
				Club:     u.Club,
				Name:     u.Name,
				Role:     u.Role,
			},
		})
		if err != nil {
			u.SubmitError = `Unable to update user. Maybe the username or name are already taken.`
			return
		}

		cache.Users.Invalidate()
		defer global.LockAndRender()()

		u.Init()
	}()
}

func (u *Users) Delete() {
	if u.DeleteConfirm {
		u.delete()
		return
	}

	u.DeleteConfirm = true
}

func (u *Users) delete() {
	u.SubmitLoading = true

	go func() {
		defer func() {
			defer global.LockAndRender()()
			u.SubmitLoading = false
		}()

		err := backend.PostDeleteUser(u.SelectedUser.ID)
		if err != nil {
			u.SubmitError = `Unable to delete user.`
			return
		}

		cache.Users.Invalidate()
		defer global.LockAndRender()()

		u.Init()
	}()
}

func (u *Users) DeleteText() string {
	if u.DeleteConfirm {
		return `Are you sure?`
	}

	return `Delete`
}

func (u *Users) SelectClub(event vugu.DOMEvent) {
	v, _ := strconv.Atoi(event.JSEventTarget().Get(`value`).String())

	u.Club = domain.Club(v)

	if !u.Club.IsAClub() {
		u.Club = domain.ClubUnknown
		return
	}
}

func (u *Users) selectClub(club domain.Club) {
	defer global.LockOnly()()

	elem, ok := document.GetElementById(`select-club`)
	if !ok {
		fmt.Fprintln(os.Stderr, `couldn't find select-club element`)
		return
	}

	elem.Set(`value`, int(club))
}

func (u *Users) SelectRole(event vugu.DOMEvent) {
	v, _ := strconv.Atoi(event.JSEventTarget().Get(`value`).String())

	u.Role = authdomain.Role(v)

	if !u.Role.IsARole() {
		u.Role = authdomain.RoleNotAuthorized
		return
	}
}

func (u *Users) selectRole(role authdomain.Role) {
	defer global.LockOnly()()

	elem, ok := document.GetElementById(`select-role`)
	if !ok {
		fmt.Fprintln(os.Stderr, `couldn't find select-role element`)
		return
	}

	elem.Set(`value`, int(role))
}

func (u *Users) PasswordLabel() string {
	if u.SelectedUser == (authdomain.User{}) {
		return `Password`
	}

	return `New Password (optional)`
}
