package history

import (
	"encoding/json"

	"github.com/FallenTaters/streepjes/domain/orderdomain"
	"github.com/FallenTaters/streepjes/frontend/backend"
	"github.com/FallenTaters/streepjes/frontend/backend/cache"
	"github.com/FallenTaters/streepjes/frontend/events"
	"github.com/FallenTaters/streepjes/frontend/global"
	"github.com/FallenTaters/streepjes/frontend/store"
)

type Ordermodal struct {
	Order      MemberOrder       `vugu:"data"`
	Contents   []store.Orderline `vugu:"data"`
	ParseError bool              `vugu:"data"`

	Loading     bool `vugu:"data"`
	DeleteError bool `vugu:"data"`

	Close CloseHandler `vugu:"data"`
}

type MemberOrder struct {
	orderdomain.Order

	Member orderdomain.Member
}

func (o *Ordermodal) Init() {
	o.ParseError = false

	err := json.Unmarshal([]byte(o.Order.Contents), &o.Contents)
	if err != nil {
		o.ParseError = true
	}
}

func (o *Ordermodal) Delete() {
	o.DeleteError = false
	o.Loading = true

	go func() {
		defer func() {
			defer global.LockAndRender()()
			o.Loading = false
		}()

		err := backend.PostDeleteOrder(o.Order.ID)
		if err != nil {
			o.DeleteError = true
			return
		}

		o.Close.CloseHandle(CloseEvent{})
		cache.Orders.Invalidate()
		events.Trigger(events.OrderDeleted)
	}()
}
