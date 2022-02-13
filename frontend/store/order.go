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
	Club:               domain.ClubGladiators, // TODO
	Catalog:            api.Catalog{},
	SelectedCategoryID: 0,
	Lines:              nil,
	OnChange:           nil,
}

type Orderline struct {
	Item   domain.Item
	Amount int
}

type OrderStore struct {
	Catalog            api.Catalog
	Club               domain.Club
	SelectedCategoryID int
	Lines              []Orderline

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
	for _, itemAmount := range os.Lines {
		total += itemAmount.Item.Price(os.Club) * domain.Price(itemAmount.Amount)
	}

	return total
}

func (os *OrderStore) SelectCategory(id int) {
	os.SelectedCategoryID = id
	os.Emit(OrderEventCategorySelected)
}

func (os *OrderStore) AddItem(item domain.Item) {
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

func (os *OrderStore) RemoveItem(item domain.Item) {
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

func (os *OrderStore) DeleteItem(item domain.Item) {
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

func (os *OrderStore) Categories() []domain.Category {
	hasItems := make(map[int]bool)
	for _, item := range os.Catalog.Items {
		if item.Price(os.Club) != 0 {
			hasItems[item.CategoryID] = true
		}
	}

	var categories []domain.Category
	for _, category := range os.Catalog.Categories {
		if hasItems[category.ID] {
			categories = append(categories, category)
		}
	}

	return categories
}

func (os *OrderStore) Items() []domain.Item {
	var items []domain.Item
	for _, item := range os.Catalog.Items {
		if item.CategoryID == os.SelectedCategoryID && item.Price(os.Club) != 0 {
			items = append(items, item)
		}
	}

	return items
}
