package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/FallenTaters/streepjes/domain/orderdomain"
)

type billingOrder struct {
	MemberName string
	Price      orderdomain.Price
	OrderTime  string
	Lines      []string
}

type billingData struct {
	pageData

	Month  string
	Orders []billingOrder
}

func parseOrderLines(contents string) []string {
	var lines []orderdomain.Line
	if err := json.Unmarshal([]byte(contents), &lines); err != nil {
		return []string{"order data unreadable"}
	}

	out := make([]string, 0, len(lines))
	for _, line := range lines {
		out = append(out, fmt.Sprintf("%dx %s", line.Amount, line.Item.Name))
	}
	return out
}

func (s *Server) getBillingPage(w http.ResponseWriter, r *http.Request) {
	user := userFromContext(r)

	data := billingData{
		pageData: newPageData(r, "billing"),
	}

	monthStr := r.URL.Query().Get("month")
	if monthStr != "" {
		month, err := orderdomain.ParseMonth(monthStr)
		if err == nil {
			data.Month = month.String()

			orders, err := s.order.GetOrdersByClub(user.Club, month)
			if err != nil {
				s.internalError(w, "get orders for billing", err)
				return
			}

			members, err := s.order.GetAllMembers()
			if err != nil {
				s.internalError(w, "get members for billing", err)
				return
			}

			membersByID := make(map[int]orderdomain.Member)
			for _, m := range members {
				membersByID[m.ID] = m
			}

			sort.Slice(orders, func(i, j int) bool {
				name1 := strings.ToLower(membersByID[orders[i].MemberID].Name)
				name2 := strings.ToLower(membersByID[orders[j].MemberID].Name)
				if name1 == name2 {
					return orders[i].OrderTime.Before(orders[j].OrderTime)
				}
				return name1 < name2
			})

			billingOrders := make([]billingOrder, 0, len(orders))
			for _, o := range orders {
				if o.MemberID == 0 {
					continue
				}
				billingOrders = append(billingOrders, billingOrder{
					MemberName: membersByID[o.MemberID].Name,
					Price:      o.Price,
					OrderTime:  o.OrderTime.Format("2006-01-02 15:04"),
					Lines:      parseOrderLines(o.Contents),
				})
			}
			data.Orders = billingOrders
		}
	}

	s.render(w, "admin/billing.html", data)
}
