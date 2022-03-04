package order

import (
	"fmt"

	"github.com/FallenTaters/streepjes/domain/orderdomain"
	"github.com/FallenTaters/streepjes/frontend/backend/cache"
	"github.com/FallenTaters/streepjes/frontend/global"
	"github.com/FallenTaters/streepjes/frontend/store"
)

type Summary struct {
	ShowMemberModal bool                 `vugu:"data"`
	Members         []orderdomain.Member `vugu:"data"`
	Loading         bool                 `vugu:"data"`
	Error           bool                 `vugu:"data"`
}

func (s *Summary) total() string {
	return store.Order.CalculateTotal().String()
}

func (s *Summary) GetMembers() []orderdomain.Member {
	return s.Members // TODO filter and fix search
}

func (s *Summary) Init() {
	s.Loading = true
	s.Error = false

	go func() {
		defer func() {
			defer global.LockAndRender()()
			s.Loading = false
		}()

		members, err := cache.Members()
		if err != nil {
			defer global.LockOnly()()
			s.Error = true
			return
		}

		fmt.Println(members)

		defer global.LockOnly()()
		s.Members = members
	}()
}
