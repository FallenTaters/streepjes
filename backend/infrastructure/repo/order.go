package repo

import (
	"errors"

	"github.com/PotatoesFall/vecty-test/domain"
	"github.com/PotatoesFall/vecty-test/domain/orderdomain"
)

var ErrOrderFieldsNotFilled = errors.New(`order fields not filled`)

type Order interface {
	// // Get a single order by ID
	// Get(id int) (orderdomain.Order, bool)

	// // Filter all orders
	// Filter(filter OrderFilter) []orderdomain.Order

	// Create a new order and return the id
	// if member id is unknown, it returns repo.ErrMemberNotFound
	// if bartender id is unknown, it returns repo.ErrUserNotFound
	Create(orderdomain.Order) (int, error)

	// // Delete an order by ID
	// Delete(id int) bool
}

type OrderFilter struct {
	Club        *domain.Club
	BartenderID *int
	MemberID    *int
	Status      []orderdomain.Status
	Month       *orderdomain.Month
}
