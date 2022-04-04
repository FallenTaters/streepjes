package cache

import (
	"time"

	"github.com/FallenTaters/streepjes/api"
	"github.com/FallenTaters/streepjes/frontend/backend"
)

func Catalog() (api.Catalog, error) {
	data, err := getOrAdd(`catalog`, time.Minute, func() (interface{}, error) {
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

	add(`catalog`, value{
		data:    catalog,
		expires: time.Now().Add(time.Hour),
	})

	return catalog, err
}

func InvalidateCatalog() {
	remove(`catalog`)
}
