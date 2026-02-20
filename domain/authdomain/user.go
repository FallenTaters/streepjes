package authdomain

import (
	"slices"
	"time"

	"github.com/FallenTaters/streepjes/domain"
)

const TokenDuration = 20 * time.Minute

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`

	Club domain.Club `json:"club"`
	Name string      `json:"name"`
	Role Role        `json:"role"`

	AuthToken string    `json:"-"`
	AuthTime  time.Time `json:"-"`
}

//go:generate go tool enumer -type Role -linecomment -sql -json

type Role int

const (
	RoleNotAuthorized Role = iota // not_authorized
	RoleBartender                 // bartender
	RoleAdmin                     // admin
)

func (r Role) Has(p Permission) bool {
	return slices.Contains(permissions[r], p)
}

type Permission int

const (
	PermissionBarStuff Permission = iota + 1
	PermissionAdminStuff
)

var permissions = map[Role][]Permission{
	RoleNotAuthorized: {},
	RoleBartender: {
		PermissionBarStuff,
	},
	RoleAdmin: {
		PermissionBarStuff,
		PermissionAdminStuff,
	},
}
