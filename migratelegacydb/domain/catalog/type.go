package catalog

import "github.com/PotatoesFall/vecty-test/migratelegacydb/shared/buckets"

type Catalog struct {
	Categories []Category `json:"categories"`
	Products   []Product  `json:"products"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (c Category) Key() []byte {
	return buckets.Itob(c.ID)
}

type Product struct {
	ID              int    `json:"id"`
	CategoryID      int    `json:"category"`
	Name            string `json:"name"`
	PriceParabool   int    `json:"priceParabool"`
	PriceGladiators int    `json:"priceGladiators"`
}

func (p Product) Key() []byte {
	return buckets.Itob(p.ID)
}
