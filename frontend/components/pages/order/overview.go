package order

import (
	"github.com/PotatoesFall/vecty-test/domain"
	"github.com/PotatoesFall/vecty-test/frontend/store"
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

func (o *Overview) removeItem(item domain.Item) {
	store.Order.RemoveItem(item)
}

func (o *Overview) addItem(item domain.Item) {
	store.Order.AddItem(item)
}

func (o *Overview) delete(item domain.Item) {
	store.Order.DeleteItem(item)
}
