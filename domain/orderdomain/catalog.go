package orderdomain

import "github.com/FallenTaters/streepjes/domain"

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Item struct {
	ID              int    `json:"id"`
	CategoryID      int    `json:"category_id"`
	Name            string `json:"name"`
	PriceGladiators Price  `json:"price_gladiators"`
	PriceParabool   Price  `json:"price_parabool"`
}

func (i Item) Price(c domain.Club) Price {
	switch c {
	case domain.ClubUnknown:
		return 0
	case domain.ClubGladiators:
		return i.PriceGladiators
	case domain.ClubParabool:
		return i.PriceParabool
	}

	panic(c)
}
