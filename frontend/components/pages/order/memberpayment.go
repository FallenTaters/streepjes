package order

import (
	"github.com/FallenTaters/streepjes/api"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
	"github.com/FallenTaters/streepjes/frontend/backend"
	"github.com/FallenTaters/streepjes/frontend/backend/cache"
	"github.com/FallenTaters/streepjes/frontend/events"
	"github.com/FallenTaters/streepjes/frontend/global"
	"github.com/FallenTaters/streepjes/frontend/store"
)

type Memberpayment struct {
	Loading       bool              `vugu:"data"`
	Error         bool              `vugu:"data"`
	MemberDetails api.MemberDetails `vugu:"data"`

	LoadingPayment bool `vugu:"data"`
	ErrorPayment   bool `vugu:"data"`

	Close CloseHandler `vugu:"data"`
}

func (m *Memberpayment) Init() {
	m.Loading = true
	go func() {
		// do request before locking
		member, err := backend.GetMember(m.Member().ID)

		defer global.LockAndRender()()
		defer func() { m.Loading = false }()

		if err != nil {
			m.Error = true
			return
		}

		m.MemberDetails = member
	}()
}

func (m *Memberpayment) Destroy() {
	events.Listen(events.OrderPlaced, `memberlist-sorting`, nil)
}

func (m *Memberpayment) Member() orderdomain.Member {
	return store.Order.Member
}

func (m *Memberpayment) Price() orderdomain.Price {
	return store.Order.CalculateTotal()
}

func (m *Memberpayment) PlaceOrder() {
	m.LoadingPayment = true
	m.ErrorPayment = false
	go func() {
		// execute request before locking
		err := backend.PostOrder(store.Order.Make())

		defer global.LockAndRender()()
		defer func() { m.LoadingPayment = false }()

		if err != nil {
			m.ErrorPayment = true
			return
		}

		store.Order.Clear()
		events.Trigger(events.OrderPlaced)
		cache.InvalidateMembers()
		m.Close.CloseHandle(CloseEvent{})
	}()
}
