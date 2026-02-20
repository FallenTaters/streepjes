package mockdb

import (
	"github.com/FallenTaters/streepjes/domain/orderdomain"
)

type Member struct {
	GetAllFunc func() ([]orderdomain.Member, error)
	GetFunc    func(id int) (orderdomain.Member, error)
	UpdateFunc func(member orderdomain.Member) error
	CreateFunc func(member orderdomain.Member) (int, error)
	DeleteFunc func(id int) error
}

func (m Member) GetAll() ([]orderdomain.Member, error) {
	return m.GetAllFunc()
}

func (m Member) Get(id int) (orderdomain.Member, error) {
	return m.GetFunc(id)
}

func (m Member) Update(member orderdomain.Member) error {
	return m.UpdateFunc(member)
}

func (m Member) Create(member orderdomain.Member) (int, error) {
	return m.CreateFunc(member)
}

func (m Member) Delete(id int) error {
	return m.DeleteFunc(id)
}
