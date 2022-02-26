package repo

import (
	"errors"

	"github.com/PotatoesFall/vecty-test/domain/authdomain"
)

var (
	ErrUserNotFound      = errors.New(`user not found`)
	ErrUsernameTaken     = errors.New(`username taken`)
	ErrUserMissingFields = errors.New(`user object is missing fields`)
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
	// if mandatory fields are missing, it returns ErrUserMissingFields
	Create(user authdomain.User) error

	// Delete a user by id. Return ErrUserHasOpenOrders if the month is not over, or ErrUserNotFound if id is unknown.
	Delete(id int) error
}
