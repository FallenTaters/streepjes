package order

import (
	"github.com/FallenTaters/streepjes/domain/orderdomain"
	"github.com/FallenTaters/streepjes/frontend/store"
)

type Memberpayment struct{}

func (m *Memberpayment) Member() orderdomain.Member {
	return store.Order.Member
}
