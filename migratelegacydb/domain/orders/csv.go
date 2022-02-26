package orders

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/PotatoesFall/vecty-test/migratelegacydb/domain/members"
	"github.com/PotatoesFall/vecty-test/migratelegacydb/shared"
)

func MakeCSV(orders []Order) ([]byte, error) {
	csvFile := new(bytes.Buffer)
	w := csv.NewWriter(csvFile)

	err := w.Write([]string{"member identifier", "member name", "month total", "orders"})
	if err != nil {
		panic(err)
	}

	ordersByMember := map[int][]Order{}
	for _, order := range orders {
		if order.Status != OrderStatusOpen && order.Status != OrderStatusBilled {
			continue
		}

		_, ok := ordersByMember[order.MemberID]
		if !ok {
			ordersByMember[order.MemberID] = []Order{order}
			continue
		}

		ordersByMember[order.MemberID] = append(ordersByMember[order.MemberID], order)
	}

	for _, orders := range ordersByMember {
		sort.Slice(orders, func(i, j int) bool {
			return orders[i].OrderTime.Sub(orders[j].OrderTime) < 0
		})

		member, err := members.Get(orders[0].MemberID)
		if err == members.ErrMemberNotFound {
			member.Name = "unknown"
		} else if err != nil {
			panic(err)
		}

		memberTotal := 0
		var orderSummary strings.Builder
		for _, o := range orders {
			if o.Status != OrderStatusOpen {
				continue
			}
			memberTotal += o.Price

			orderSummary.WriteString(fmt.Sprintf("%s %.2fe, placed by %s: %s\n",
				o.OrderTime.Format(`2006-01-02 03:04`),
				float64(o.Price)/100,
				o.Bartender,
				parseContents(o.Contents, o.Club),
			))
		}

		err = w.Write([]string{
			"", // placeholder for later
			member.Name,
			fmt.Sprint(memberTotal),
			orderSummary.String(),
		})
		if err != nil {
			panic(err)
		}
	}

	w.Flush()
	return csvFile.Bytes(), nil
}

type Product struct {
	Name            string `json:"name"`
	PriceGladiators int    `json:"priceGladiators"`
	PriceParabool   int    `json:"priceParabool"`
}

func (p Product) clubPrice(club shared.Club) int {
	switch club {
	case shared.ClubGladiators:
		return p.PriceGladiators
	case shared.ClubParabool:
		return p.PriceParabool
	default:
		return 0
	}
}

type Orderline struct {
	Product Product `json:"Product"`
	Amount  int     `json:"amount"`
}

func parseContents(data string, club shared.Club) string {
	contents := []Orderline{}
	err := json.Unmarshal([]byte(data), &contents)
	if err != nil {
		return "unable to read order contents"
	}

	var out string
	for _, ol := range contents {
		out += fmt.Sprintf(`%dx %s %.2fe; `, ol.Amount, ol.Product.Name, float64(ol.Product.clubPrice(club)*ol.Amount)/100)
	}

	return out
}
