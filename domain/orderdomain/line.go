package orderdomain

import "github.com/FallenTaters/streepjes/domain"

// Line is usually saved as JSON and data integrity is not guaranteed.
type Line struct {
	Item   Item `json:"product"` // json tag is for backwards compatibility
	Amount int  `json:"amount"`
}

func (ol Line) Price(club domain.Club) Price {
	return ol.Item.Price(club).Times(ol.Amount)
}
