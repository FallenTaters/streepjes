package order

import (
	"github.com/FallenTaters/streepjes/frontend/store"
)

type Summary struct{}

func (s *Summary) total() string {
	return store.Order.CalculateTotal().String()
}
