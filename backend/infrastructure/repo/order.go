package repo

import (
	"errors"
	"time"

	"github.com/FallenTaters/streepjes/domain"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
)

var (
	ErrOrderFieldsNotFilled = errors.New(`order fields not filled`)
	ErrOrderNotFound        = errors.New(`order not found`)
)

type Order interface {
	// Get a single order by ID
	Get(id int) (orderdomain.Order, bool)

	// Filter all orders. The zero value for any filter is ignored.
	Filter(filter OrderFilter) []orderdomain.Order

	// Create a new order and return the id
	// if member id is unknown, it returns repo.ErrMemberNotFound
	// if bartender id is unknown, it returns repo.ErrUserNotFound
	Create(orderdomain.Order) (int, error)

	// Delete an order by ID
	Delete(id int) error
}

type OrderFilter struct {
	Club        domain.Club
	BartenderID int
	MemberID    int
	// Status      []orderdomain.Status
	StatusNot  []orderdomain.Status
	Start, End time.Time
	Limit      int
}
