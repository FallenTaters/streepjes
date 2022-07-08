package pages

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/FallenTaters/streepjes/domain/orderdomain"
	"github.com/FallenTaters/streepjes/frontend/backend"
	"github.com/FallenTaters/streepjes/frontend/global"
	"github.com/FallenTaters/streepjes/frontend/jscall/window"
	"github.com/FallenTaters/timefmt"
	"github.com/vugu/vugu"
)

type monthFormat struct{}

func (monthFormat) DateFormat() string {
	return `2006-01`
}

type Billing struct {
	Loading, Error bool

	Month orderdomain.Month

	MembersByID map[int]orderdomain.Member
	Orders      []orderdomain.Order
}

func (b *Billing) SetMonth(event vugu.DOMEvent) {
	v := event.JSEventTarget().Get(`value`).String()

	var date timefmt.Date[monthFormat]
	err := date.UnmarshalText([]byte(v))
	if err != nil {
		panic(err) // TODO
	}
	b.Month = orderdomain.MonthOf(date.Time())

	b.Loading = true
	b.Error = false

	go func() {
		var members []orderdomain.Member
		var ordersErr, membersErr error

		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			var orders []orderdomain.Order
			orders, ordersErr = backend.GetBillingOrders(b.Month)
			b.Orders = make([]orderdomain.Order, 0, len(orders))
			for _, o := range orders {
				if o.MemberID != 0 {
					b.Orders = append(b.Orders, o)
				}
			}
			wg.Done()
		}()
		go func() {
			members, membersErr = backend.GetMembers()
			wg.Done()
		}()
		wg.Wait()

		defer global.LockAndRender()()
		b.Loading = false
		if ordersErr != nil || membersErr != nil {
			b.Error = true
			return
		}

		b.MembersByID = make(map[int]orderdomain.Member)
		for _, member := range members {
			b.MembersByID[member.ID] = member
		}

		sort.Slice(b.Orders, func(i, j int) bool {
			id1, id2 := b.Orders[i].MemberID, b.Orders[j].MemberID
			name1, name2 := strings.ToLower(b.MembersByID[id1].Name), strings.ToLower(b.MembersByID[id2].Name)
			if name1 == name2 {
				return b.Orders[i].OrderTime.Before(b.Orders[j].OrderTime)
			}
			return strings.Compare(name1, name2) < 0
		})
	}()
}

func (b *Billing) Parse(contents string) []string {
	var lines []orderdomain.Line
	if err := json.Unmarshal([]byte(contents), &lines); err != nil {
		return []string{`order data unreadable`}
	}

	var out []string
	for _, line := range lines {
		out = append(out, fmt.Sprintf("%dx %s\n", line.Amount, line.Item.Name))
	}

	return out
}

func (b *Billing) Download() {
	window.NewTab(`/admin/download?month=` + b.Month.String())
}
