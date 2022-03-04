package order

import (
	"github.com/FallenTaters/streepjes/domain/orderdomain"
	"github.com/FallenTaters/streepjes/frontend/store"
	"github.com/vugu/vugu"
)

type Overview struct {
	Lines []store.Orderline
}

func (o *Overview) Compute(vugu.ComputeCtx) {
	o.Lines = store.Order.Lines
}

func (o *Overview) classes(ol store.Orderline) string {
	if ol.Item.Price(store.Order.Club) == 0 {
		return `error`
	}

	return ``
}

func (o *Overview) removeItem(item orderdomain.Item) {
	store.Order.RemoveItem(item)
}

func (o *Overview) addItem(item orderdomain.Item) {
	store.Order.AddItem(item)
}

func (o *Overview) delete(item orderdomain.Item) {
	store.Order.DeleteItem(item)
}
