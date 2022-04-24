package pages

import (
	"sort"
	"strings"

	"github.com/FallenTaters/streepjes/domain/authdomain"
	"github.com/FallenTaters/streepjes/frontend/backend/cache"
	"github.com/FallenTaters/streepjes/frontend/global"
)

type Users struct {
	Loading bool `vugu:"data"`
	Error   bool `vugu:"data"`

	Users        []authdomain.User `vugu:"data"`
	SelectedUser authdomain.User   `vugu:"data"`

	Username string          `vugu:"data"`
	Password string          `vugu:"data"`
	Name     string          `vugu:"data"`
	Role     authdomain.Role `vugu:"data"`
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
