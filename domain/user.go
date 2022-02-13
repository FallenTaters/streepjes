package domain

import "time"

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`

	Club Club   `json:"club"`
	Name string `json:"name"`
	Role Role   `json:"role"`

	AuthToken string    `json:"authToken"`
	AuthTime  time.Time `json:"authDate"`
}

//go:generate enumer -type Role -linecomment -sql -json

type Role int

const (
	RoleNotAuthorized Role = iota // not_authorized
	RoleBartender                 // bartender
	RoleAdmin                     // admin
)
