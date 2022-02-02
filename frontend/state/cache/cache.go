package cache

import (
	"github.com/PotatoesFall/vecty-test/api"
	"github.com/PotatoesFall/vecty-test/frontend/backend"
)

var cache = map[string]interface{}{}

func getOrAdd(key string, f func() (interface{}, error)) (interface{}, error) {
	data, exists := cache[key]
	if exists {
		return data, nil
	}

	data, err := f()
	if err != nil {
		return nil, err
	}

	add(key, data)
	return data, nil
}

func add(key string, v interface{}) {
	cache[key] = v
}

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
