package store

import (
	"github.com/PotatoesFall/vecty-test/domain"
	"github.com/PotatoesFall/vecty-test/domain/orderdomain"
)

type OrderEvent int

const (
	OrderEventItemsChanged OrderEvent = iota + 1
	OrderEventCategorySelected
	OrderEventClubChanged
)

var Order = OrderStore{
	Club:     domain.ClubGladiators, // TODO
	Lines:    nil,
	OnChange: nil,
}

type Orderline struct {
	Item   orderdomain.Item
	Amount int
}

func (ol Orderline) Price() orderdomain.Price {
	return ol.Item.Price(Order.Club).Times(ol.Amount)
}

type OrderStore struct {
	Club  domain.Club
	Lines []Orderline

	OnChange func(OrderEvent)
}

func (os *OrderStore) Emit(event OrderEvent) {
	if os.OnChange == nil {
		return
	}

	os.OnChange(event)
}

func (os *OrderStore) CalculateTotal() orderdomain.Price {
	var total orderdomain.Price = 0
	for _, itemAmount := range os.Lines {
		total += itemAmount.Item.Price(os.Club) * orderdomain.Price(itemAmount.Amount)
	}

	return total
}

func (os *OrderStore) AddItem(item orderdomain.Item) {
	defer os.Emit(OrderEventItemsChanged)

	for i, itemAmount := range os.Lines {
		if itemAmount.Item.ID == item.ID {
			os.Lines[i].Amount++
			return
		}
	}

	os.Lines = append(os.Lines, Orderline{
		Item:   item,
		Amount: 1,
	})
}

func (os *OrderStore) RemoveItem(item orderdomain.Item) {
	defer os.Emit(OrderEventItemsChanged)

	for i, itemAmount := range os.Lines {
		if itemAmount.Item.ID == item.ID {
			if itemAmount.Amount <= 1 {
				os.deleteAt(i)
				return
			}

			os.Lines[i].Amount--
		}
	}
}

func (os *OrderStore) DeleteItem(item orderdomain.Item) {
	for i, itemAmount := range os.Lines {
		if itemAmount.Item.ID == item.ID {
			os.deleteAt(i)
		}
	}

	os.Emit(OrderEventItemsChanged)
}

func (os *OrderStore) deleteAt(i int) {
	newItems := os.Lines[:i]
	if len(os.Lines) > i+1 {
		newItems = append(newItems, os.Lines[i+1:]...)
	}

	os.Lines = newItems
}

func (os *OrderStore) ToggleClub() {
	if os.Club == domain.ClubGladiators {
		os.Club = domain.ClubParabool
	} else {
		os.Club = domain.ClubGladiators
	}

	os.Emit(OrderEventClubChanged)
}
