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
	Get(id int) (orderdomain.Order, error)
	Filter(filter OrderFilter) ([]orderdomain.Order, error)
	Create(orderdomain.Order) (int, error)
	Delete(id int) error
}

type OrderFilter struct {
	Club        domain.Club
	BartenderID int
	MemberID    int
	StatusNot   []orderdomain.Status
	Start, End  time.Time
	Limit       int
}
