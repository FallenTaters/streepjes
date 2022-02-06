package store

import (
	"fmt"

	"github.com/PotatoesFall/vecty-test/api"
	"github.com/PotatoesFall/vecty-test/domain"
)

type OrderEvent int

const (
	OrderEventItemsChanged OrderEvent = iota + 1
	OrderEventCategorySelected
	OrderEventClubChanged
)

var Order = OrderStore{
	Club: domain.ClubGladiators, // TODO
}

type Orderline struct {
	Item   domain.Item
	Amount int
}

type OrderStore struct {
	Catalog            api.Catalog
	Club               domain.Club
	SelectedCategoryID int
	Items              []Orderline
	ShownItems         []domain.Item

	OnChange func(OrderEvent)
}

func (os *OrderStore) Emit(event OrderEvent) {
	if os.OnChange == nil {
		return
	}

	os.OnChange(event)
}

func (os *OrderStore) CalculateTotal() domain.Price {
	var total domain.Price = 0
	for _, itemAmount := range os.Items {
		total += itemAmount.Item.Price(os.Club) * domain.Price(itemAmount.Amount)
	}

	return total
}

func (os *OrderStore) SelectCategory(id int) {
	os.SelectedCategoryID = id

	os.ShownItems = nil
	for _, item := range os.Catalog.Items {
		if item.CategoryID == id {
			os.ShownItems = append(os.ShownItems, item)
		}
	}

	os.Emit(OrderEventCategorySelected)
}

func (os *OrderStore) AddItem(item domain.Item) {
	defer os.Emit(OrderEventItemsChanged)

	for i, itemAmount := range os.Items {
		if itemAmount.Item.ID == item.ID {
			os.Items[i].Amount++
			return
		}
	}

	os.Items = append(os.Items, Orderline{
		Item:   item,
		Amount: 1,
	})
}

func (os *OrderStore) RemoveItem(item domain.Item) {
	defer os.Emit(OrderEventItemsChanged)

	for i, itemAmount := range os.Items {
		if itemAmount.Item.ID == item.ID {
			if itemAmount.Amount <= 1 {
				os.deleteAt(i)
				return
			}

			os.Items[i].Amount--
		}
	}
}

func (os *OrderStore) DeleteItem(item domain.Item) {
	fmt.Println(item)
	for i, itemAmount := range os.Items {
		if itemAmount.Item.ID == item.ID {
			os.deleteAt(i)
		}
	}

	os.Emit(OrderEventItemsChanged)
}

func (os *OrderStore) deleteAt(i int) {
	newItems := os.Items[:i]
	if len(os.Items) > i+1 {
		newItems = append(newItems, os.Items[i+1:]...)
	}

	os.Items = newItems
}

func (os *OrderStore) ToggleClub() {
	if os.Club == domain.ClubGladiators {
		os.Club = domain.ClubParabool
	} else {
		os.Club = domain.ClubGladiators
	}

	os.Emit(OrderEventClubChanged)
}
