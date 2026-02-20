package mockdb

import (
	"github.com/FallenTaters/streepjes/backend/infrastructure/repo"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
)

type Order struct {
	GetFunc    func(id int) (orderdomain.Order, error)
	FilterFunc func(filter repo.OrderFilter) ([]orderdomain.Order, error)
	CreateFunc func(order orderdomain.Order) (int, error)
	DeleteFunc func(id int) error
}

func (o Order) Get(id int) (orderdomain.Order, error) {
	return o.GetFunc(id)
}

func (o Order) Filter(filter repo.OrderFilter) ([]orderdomain.Order, error) {
	return o.FilterFunc(filter)
}

func (o Order) Create(order orderdomain.Order) (int, error) {
	return o.CreateFunc(order)
}

func (o Order) Delete(id int) error {
	return o.DeleteFunc(id)
}
