package domain

import "fmt"

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Item struct {
	ID              int    `json:"id"`
	CategoryID      int    `json:"category_id"`
	Name            string `json:"name"`
	PriceGladiators int    `json:"price_gladiators"`
	PriceParabool   int    `json:"price_parabool"`
}

func (i Item) Price(c Club) int {
	switch c {
	case ClubGladiators:
		return i.PriceGladiators
	case ClubParabool:
		return i.PriceParabool
	}

	panic(c)
}

func (i Item) PriceString(c Club) string {
	return fmt.Sprintf(`€%.2f`, float64(i.Price(c))/100)
}
