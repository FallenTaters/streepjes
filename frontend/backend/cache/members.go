package cache

import (
	"github.com/PotatoesFall/vecty-test/domain"
	"github.com/PotatoesFall/vecty-test/frontend/backend"
)

func Members() ([]domain.Member, error) {
	data, err := getOrAdd(`members`, func() (interface{}, error) {
		return backend.GetMembers()
	})
	if err != nil {
		return []domain.Member{}, err
	}

	return data.([]domain.Member), nil
}

func FetchMembers() ([]domain.Member, error) {
	members, err := backend.GetMembers()
	if err != nil {
		return members, err
	}

	add(`members`, members)

	return members, err
}
