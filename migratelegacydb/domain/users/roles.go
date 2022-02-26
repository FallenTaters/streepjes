package users

//go:generate enumer -type Role -linecomment -sql -json

type Role int

const (
	RoleNotAuthorized Role = iota // not_authorized
	RoleBartender                 // bartender
	RoleAdmin                     // admin
)
