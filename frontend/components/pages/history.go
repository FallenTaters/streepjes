package pages

import (
	"sort"
	"time"

	"github.com/FallenTaters/streepjes/domain/orderdomain"
	"github.com/FallenTaters/streepjes/frontend/backend/cache"
	"github.com/FallenTaters/streepjes/frontend/components/pages/history"
	"github.com/FallenTaters/streepjes/frontend/events"
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
		orders, err1 := cache.Orders.Get()
		members, err2 := cache.Members.Get()
		defer global.LockAndRender()()
		defer func() {
			h.Loading = false
		}()
		if err1 != nil || err2 != nil {
			h.Error = true
			return
		}

		sort.Slice(orders, func(i, j int) bool {
			return orders[i].OrderTime.After(orders[j].OrderTime)
		})
		h.Orders = orders

		membersByID := make(map[int]orderdomain.Member)
		for _, member := range members {
			membersByID[member.ID] = member
		}
		h.MembersByID = membersByID
	}()

	events.Listen(events.OrderDeleted, `history-reload`, func() {
		defer global.LockAndRender()()
		h.Init()
	})
}

func (h *History) Click(order orderdomain.Order) {
	h.SelectedOrder = history.MemberOrder{
		Order:  order,
		Member: h.MembersByID[order.MemberID],
	}
	h.ShowOrderModal = true
}

func (h *History) formatDate(t time.Time) string {
	return shared.PrettyDatetime(t)
}

func (h *History) classes(order orderdomain.Order) string {
	classes := `responsive extra small-margin `

	if order.Status == orderdomain.StatusCancelled {
		classes += `grey`
	} else {
		classes += order.Club.String()
	}

	return classes
}
