package cache

import (
	"github.com/PotatoesFall/vecty-test/api"
	"github.com/PotatoesFall/vecty-test/frontend/backend"
)

func Catalog() (api.Catalog, error) {
	data, err := getOrAdd(`catalog`, func() (interface{}, error) {
		return backend.GetCatalog()
	})
	if err != nil {
		return api.Catalog{}, err
	}

	return data.(api.Catalog), nil
}

func FetchCatalog() (api.Catalog, error) {
	catalog, err := backend.GetCatalog()
	if err != nil {
		return catalog, err
	}

	add(`catalog`, catalog)

	return catalog, err
}
