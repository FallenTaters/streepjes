package repo

import (
	"errors"

	"github.com/FallenTaters/streepjes/domain/authdomain"
)

var (
	ErrUserNotFound      = errors.New(`user not found`)
	ErrUsernameTaken     = errors.New(`username taken`)
	ErrUserMissingFields = errors.New(`user object is missing fields`)
	ErrUserHasOrders     = errors.New(`cannot delete user with orders`)
)

type User interface {
	Get(id int) (authdomain.User, error)
	GetAll() ([]authdomain.User, error)
	GetByToken(token string) (authdomain.User, error)
	GetByUsername(username string) (authdomain.User, error)
	Update(user authdomain.User) error
	UpdateActivity(user authdomain.User) error
	Create(user authdomain.User) (int, error)
	Delete(id int) error
}
