package pages

import (
	"sort"
	"strings"

	"github.com/FallenTaters/streepjes/domain"
	"github.com/FallenTaters/streepjes/domain/authdomain"
	"github.com/FallenTaters/streepjes/frontend/backend/cache"
	"github.com/FallenTaters/streepjes/frontend/global"
)

type Users struct {
	Loading bool `vugu:"data"`
	Error   bool `vugu:"data"`

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
			return strings.Compare(users[i].Name, users[j].Name) < 0
		})

		u.Users = users
	}()
}

func (u *Users) Select(user authdomain.User) {
	u.SelectedUser = user

	u.Name = user.Name
	u.Username = user.Username
	u.Password = ``
	u.Role = user.Role
	u.Club = user.Club

	u.ShowForm = true
}

func (u *Users) NewUser() {
	u.SelectedUser = authdomain.User{}

	u.Name = ``
	u.Username = ``
	u.Password = ``
	u.Role = authdomain.RoleNotAuthorized
	u.Club = domain.ClubUnknown

	u.ShowForm = true
}

// TODO use this to prevent accidentally clearing unsaved changes when selecting another member
func (u *Users) unsavedChanges() bool {
	if !u.ShowForm {
		return false
	}

	return u.Name != u.SelectedUser.Name ||
		u.Username != u.SelectedUser.Username ||
		u.Password != `` ||
		u.Role != u.SelectedUser.Role ||
		u.Club != u.SelectedUser.Club
}

func (u *Users) Confirm(callback func()) {
	// TODO
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
	// TODO
}
