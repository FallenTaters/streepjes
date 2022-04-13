package order

import (
	"github.com/FallenTaters/streepjes/domain/orderdomain"
	"github.com/FallenTaters/streepjes/frontend/backend"
	"github.com/FallenTaters/streepjes/frontend/events"
	"github.com/FallenTaters/streepjes/frontend/global"
	"github.com/FallenTaters/streepjes/frontend/store"
)

type Anonymouspayment struct {
	LoadingPayment bool `vugu:"data"`
	ErrorPayment   bool `vugu:"data"`

	Close CloseHandler `vugu:"data"`
}

func (a *Anonymouspayment) Price() orderdomain.Price {
	return store.Order.CalculateTotal()
}

func (a *Anonymouspayment) PlaceOrder() {
	a.LoadingPayment = true
	a.ErrorPayment = false
	go func() {
		// execute request before locking
		err := backend.PostOrder(store.Order.Make())

		defer global.LockAndRender()()
		defer func() { a.LoadingPayment = false }()

		if err != nil {
			a.ErrorPayment = true
			return
		}

		store.Order.Clear()
		events.Trigger(events.OrderPlaced)
		a.Close.CloseHandle(CloseEvent{})
	}()
}
