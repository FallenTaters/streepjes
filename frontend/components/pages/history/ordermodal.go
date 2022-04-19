package history

import (
	"encoding/json"

	"github.com/FallenTaters/streepjes/domain/orderdomain"
	"github.com/FallenTaters/streepjes/frontend/store"
)

type Ordermodal struct {
	Order      MemberOrder       `vugu:"data"`
	Contents   []store.Orderline `vugu:"data"`
	ParseError bool              `vugu:"data"`

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
	// TODO
}
