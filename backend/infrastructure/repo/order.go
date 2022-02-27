package repo

import (
	"github.com/PotatoesFall/vecty-test/domain"
	"github.com/PotatoesFall/vecty-test/domain/orderdomain"
)

type Order interface {
	// Get a single order by ID
	Get(id int) (orderdomain.Order, bool)

	// Filter all orders
	Filter(filter OrderFilter) []orderdomain.Order

	// Create a new order
	Create(orderdomain.Order) error

	// Delete an order by ID
	Delete(id int) bool
}

type OrderFilter struct {
	Club        *domain.Club
	BartenderID *int
	MemberID    *int
	Status      []orderdomain.Status
	Month       *orderdomain.Month
}
