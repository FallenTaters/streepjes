package domain

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
