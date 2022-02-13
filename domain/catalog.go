package domain

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

func (i Item) Price(c Club) Price {
	switch c {
	case ClubUnknown:
		return 0
	case ClubGladiators:
		return i.PriceGladiators
	case ClubParabool:
		return i.PriceParabool
	}

	panic(c)
}
