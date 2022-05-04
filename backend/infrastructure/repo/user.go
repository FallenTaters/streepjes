package repo

import (
	"errors"

	"github.com/FallenTaters/streepjes/domain/authdomain"
)

var (
	ErrNameTaken         = errors.New(`name taken`)
	ErrUserNotFound      = errors.New(`user not found`)
	ErrUsernameTaken     = errors.New(`username taken`)
	ErrUserMissingFields = errors.New(`user object is missing fields`)
	ErrUserHasOrders     = errors.New(`cannot delete user with orders`)
)

type User interface {
	// Get gets a single user by ID. It returns false if not found.
	Get(id int) (authdomain.User, bool)

	// GetAll gets all users
	GetAll() []authdomain.User

	// GetUser by token, false if not found
	GetByToken(token string) (authdomain.User, bool)

	// Get a specific user by username, returns false if not found
	GetByUsername(username string) (authdomain.User, bool)

	// Update a specific user. Returns ErrUserNotFound if the ID is not found
	Update(user authdomain.User) error

	// Create a new user. Returns ErrUsernameTaken if the username already exists.
	// if name is taken, it returns ErrNameTaken.
	// if mandatory fields are missing, it returns ErrUserMissingFields
	// it returns the id of the new user
	Create(user authdomain.User) (int, error)

	// Delete a user by id. Returns ErrUserHasOrders if there is a foreign key conflict, or ErrUserNotFound if id is unknown.
	Delete(id int) error
}
