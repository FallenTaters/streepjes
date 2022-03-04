package users

import (
	"time"

	"github.com/FallenTaters/streepjes/migratelegacydb/shared"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Username string `json:"username"`

	Club shared.Club `json:"club"`
	Name string      `json:"name"`
	Role Role        `json:"role"`

	Password  []byte    `json:"password"`
	AuthToken string    `json:"authToken"`
	AuthTime  time.Time `json:"authDate"`
}

func (u User) Key() []byte {
	return []byte(u.Username)
}

func (u User) AsPayload() UserPayload {
	return UserPayload{
		Username: u.Username,
		Club:     u.Club,
		Name:     u.Name,
		Role:     u.Role,
	}
}

type UserPayload struct {
	Username string      `json:"username"`
	Club     shared.Club `json:"club"`
	Name     string      `json:"name"`
	Role     Role        `json:"role"`
}

func (u UserPayload) AsUser() User {
	return User{
		Username: u.Username,
		Club:     u.Club,
		Name:     u.Name,
		Role:     u.Role,
	}
}
