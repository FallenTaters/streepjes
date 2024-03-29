package authdomain

import (
	"time"

	"github.com/FallenTaters/streepjes/domain"
)

const (
	TokenDuration         = 24 * time.Hour
	LockScreenWarningTime = 24 * time.Hour
)

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

//go:generate enumer -type Role -linecomment -sql -json

type Role int

const (
	RoleNotAuthorized Role = iota // not_authorized
	RoleBartender                 // bartender
	RoleAdmin                     // admin
)

func (r Role) Has(p Permission) bool {
	for _, permission := range permissions[r] {
		if permission == p {
			return true
		}
	}

	return false
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
