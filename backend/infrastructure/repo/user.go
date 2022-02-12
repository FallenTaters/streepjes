package repo

import "github.com/PotatoesFall/vecty-test/domain"

type User interface {
	// GetUser by token, false if not found
	GetByToken(token string) (domain.User, bool)

	// Get a specific user by username, returns false if not found
	GetByUsername(username string) (domain.User, bool)

	// Update a specific user. Returns ErrUserNotFound if the ID is not found
	Update(user domain.User) error

	// Delete a user by id. Return ErrUserHasOpenOrders if the month is not over, or ErrUserNotFound if id is unknown.
	Delete(id int) error
}
