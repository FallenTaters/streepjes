package order

import (
	"github.com/PotatoesFall/vecty-test/frontend/store"
)

type Summary struct{}

func (s *Summary) total() string {
	return store.Order.CalculateTotal().String()
}
