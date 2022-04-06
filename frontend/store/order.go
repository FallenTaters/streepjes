package store

import (
	"encoding/json"

	"github.com/FallenTaters/streepjes/domain"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
)

type OrderEvent int

const (
	OrderEventItemsChanged OrderEvent = iota + 1
	OrderEventCategorySelected
	OrderEventClubChanged
)

type OrderStore struct {
	Club   domain.Club
	Lines  []Orderline
	Member orderdomain.Member
}

var Order = OrderStore{
	Club:   domain.ClubGladiators,
	Lines:  nil,
	Member: orderdomain.Member{},
}

type Orderline struct {
	Item   orderdomain.Item
	Amount int
}

func (ol Orderline) Price() orderdomain.Price {
	return ol.Item.Price(Order.Club).Times(ol.Amount)
}

func (os *OrderStore) CalculateTotal() orderdomain.Price {
	var total orderdomain.Price = 0
	for _, itemAmount := range os.Lines {
		total += itemAmount.Item.Price(os.Club) * orderdomain.Price(itemAmount.Amount)
	}

	return total
}

func (os *OrderStore) AddItem(item orderdomain.Item) {
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
}

func (os *OrderStore) Contents() string {
	data, err := json.Marshal(os.Lines)
	if err != nil {
		panic(err)
	}

	return string(data)
}

func (os *OrderStore) Make() orderdomain.Order {
	return orderdomain.Order{ //nolint:exhaustivestruct
		Club:     os.Club,
		MemberID: os.Member.ID,
		Contents: os.Contents(),
		Price:    os.CalculateTotal(),
		Status:   orderdomain.StatusOpen,
	}
}

func (os *OrderStore) Clear() {
	os.Lines = nil
	os.Member = orderdomain.Member{}
}
