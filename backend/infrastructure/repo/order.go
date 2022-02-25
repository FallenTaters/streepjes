package repo

import (
	"time"

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
	DeleteByID(id int) bool
}

type OrderFilter struct {
	Club      *int
	Bartender *string
	Member    *int
	Status    []int
	Month     *time.Time
}
