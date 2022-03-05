package cache

import (
	"time"

	"github.com/FallenTaters/streepjes/domain/orderdomain"
	"github.com/FallenTaters/streepjes/frontend/backend"
)

func Members() ([]orderdomain.Member, error) {
	data, err := getOrAdd(`members`, time.Minute, func() (interface{}, error) {
		return backend.GetMembers()
	})
	if err != nil {
		return []orderdomain.Member{}, err
	}

	return data.([]orderdomain.Member), nil
}

// func Member(id int) ([]orderdomain.Member, error) {
// 	data, err := getOrAdd(`member`+strconv., func() (interface{}, error) {
// 		return backend.GetMember(id)
// 	})
// 	if err != nil {
// 		return []orderdomain.Member{}, err
// 	}

// 	return data.([]orderdomain.Member), nil
// }
