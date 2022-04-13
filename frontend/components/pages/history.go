package pages

import (
	"github.com/FallenTaters/streepjes/domain/orderdomain"
	"github.com/FallenTaters/streepjes/frontend/backend/cache"
)

type History struct {
	Orders      []orderdomain.Order        `vugu:"data"`
	MembersByID map[int]orderdomain.Member `vugu:"data"`
}

func (h *History) Init() {
	go func() {
		orders, err := cache.Orders.Get()
		if err != nil {
			return
		}
		h.Orders = orders

		members, err := cache.Members.Get()
		if err != nil {
			return
		}

		membersByID := make(map[int]orderdomain.Member)
		for _, member := range members {
			membersByID[member.ID] = member
		}
		h.MembersByID = membersByID
	}()
}
