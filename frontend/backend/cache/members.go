package cache

import (
	"github.com/PotatoesFall/vecty-test/domain/orderdomain"
	"github.com/PotatoesFall/vecty-test/frontend/backend"
)

func Members() ([]orderdomain.Member, error) {
	data, err := getOrAdd(`members`, func() (interface{}, error) {
		return backend.GetMembers()
	})
	if err != nil {
		return []orderdomain.Member{}, err
	}

	return data.([]orderdomain.Member), nil
}

func FetchMembers() ([]orderdomain.Member, error) {
	members, err := backend.GetMembers()
	if err != nil {
		return members, err
	}

	add(`members`, members)

	return members, err
}
