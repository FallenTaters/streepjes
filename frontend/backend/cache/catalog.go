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

// TODO: fix all this. make a map of enum to getFunc.
// enum is strings (keys) of type cache.Key
// then have: cache.Get(cache.Catalog) en zo
// remove the expiration thing
// instad also have a cache.Invalidate(key) thing
// for complex keys, like member-id, make builder functions that make values
// like cache.MemberKey(id) cache.Key
// don't cache non-generic
// remove the fetch thing
