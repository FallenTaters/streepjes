package repo

import (
	"time"

	"github.com/PotatoesFall/vecty-test/src/domain"
)

type Order interface {
	// Get a single order by ID
	Get(id int) (domain.Order, bool)

	// Filter all orders
	Filter(filter OrderFilter) []domain.Order

	// Create a new order
	Create(domain.Order) error

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
