package pages

import (
	"sort"
	"time"

	"github.com/FallenTaters/streepjes/domain/orderdomain"
	"github.com/FallenTaters/streepjes/frontend/backend/cache"
	"github.com/FallenTaters/streepjes/frontend/components/pages/history"
	"github.com/FallenTaters/streepjes/frontend/global"
	"github.com/FallenTaters/streepjes/shared"
)

type History struct {
	Loading bool `vugu:"data"`
	Error   bool `vugu:"data"`

	SelectedOrder  history.MemberOrder `vugu:"data"`
	ShowOrderModal bool                `vugu:"data"`

	Orders      []orderdomain.Order        `vugu:"data"`
	MembersByID map[int]orderdomain.Member `vugu:"data"`
}

func (h *History) Init() {
	h.Error = false
	h.Loading = true

	go func() {
		defer func() {
			defer global.LockAndRender()()
			h.Loading = false
		}()

		orders, err := cache.Orders.Get()
		if err != nil {
			defer global.LockAndRender()()
			h.Error = true
			return
		}

		sort.Slice(orders, func(i, j int) bool {
			return orders[i].OrderTime.After(orders[j].OrderTime)
		})
		h.Orders = orders

		members, err := cache.Members.Get()
		if err != nil {
			h.Error = true
			return
		}

		defer global.LockAndRender()()

		membersByID := make(map[int]orderdomain.Member)
		for _, member := range members {
			membersByID[member.ID] = member
		}
		h.MembersByID = membersByID
	}()
}

func (h *History) Click(order orderdomain.Order) {
	h.SelectedOrder = history.MemberOrder{
		Order:  order,
		Member: h.MembersByID[order.ID],
	}
	h.ShowOrderModal = true
}

func (h *History) formatDate(t time.Time) string {
	return shared.PrettyDatetime(t)
}
