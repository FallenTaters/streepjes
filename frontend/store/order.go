package store

import (
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

type OrderStore struct {
	Catalog            api.Catalog
	Club               domain.Club
	SelectedCategoryID int
	Items              map[domain.Item]int
	ShownItems         []domain.Item

	OnChange func(OrderEvent)
}

func (os *OrderStore) Emit(event OrderEvent) {
	if os.OnChange == nil {
		return
	}

	os.OnChange(event)
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
	if os.Items == nil {
		os.Items = make(map[domain.Item]int)
	}

	os.Items[item]++
	os.Emit(OrderEventItemsChanged)
}

func (os *OrderStore) RemoveItem(item domain.Item) {
	defer os.Emit(OrderEventItemsChanged)

	if os.Items == nil {
		os.Items = make(map[domain.Item]int)
	}

	if os.Items[item] > 1 {
		os.Items[item]--
		return
	}

	delete(os.Items, item)
}

func (os *OrderStore) DeleteItem(item domain.Item) {
	if os.Items == nil {
		os.Items = make(map[domain.Item]int)
	}

	delete(os.Items, item)
	os.Emit(OrderEventItemsChanged)
}

func (os *OrderStore) ToggleClub() {
	if os.Club == domain.ClubGladiators {
		os.Club = domain.ClubParabool
	} else {
		os.Club = domain.ClubGladiators
	}

	os.Emit(OrderEventClubChanged)
}
