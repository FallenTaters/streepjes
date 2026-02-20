package order

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/FallenTaters/streepjes/domain/orderdomain"
)

func writeCSV(orders []orderdomain.Order, members []orderdomain.Member, timezone *time.Location) []byte {
	membersByID := make(map[int]orderdomain.Member)
	for _, m := range members {
		membersByID[m.ID] = m
	}

	filtered := make([]orderdomain.Order, 0, len(orders))
	for _, o := range orders {
		if o.MemberID != 0 {
			filtered = append(filtered, o)
		}
	}
	orders = filtered

	sort.Slice(orders, func(i, j int) bool {
		id1, id2 := orders[i].MemberID, orders[j].MemberID
		name1, name2 := strings.ToLower(membersByID[id1].Name), strings.ToLower(membersByID[id2].Name)
		if name1 == name2 {
			return orders[i].OrderTime.Before(orders[j].OrderTime)
		}
		return strings.Compare(name1, name2) < 0
	})

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	w.Write([]string{`Member`, `Price`, `Date`, `Order`})

	memberTotals := make(map[int]orderdomain.Price, len(members))
	for _, o := range orders {
		memberTotals[o.MemberID] += o.Price
		w.Write([]string{
			membersByID[o.MemberID].Name,
			o.Price.String(),
			o.OrderTime.In(timezone).Format(`2006-01-02 15:04`),
			csvOrderLines(o.Contents),
		})
	}

	w.Write(nil)
	w.Write([]string{`Member`, `Total`})
	for id, total := range memberTotals {
		w.Write([]string{membersByID[id].Name, total.String()})
	}

	w.Flush()

	return buf.Bytes()
}

func csvOrderLines(contents string) string {
	var lines []orderdomain.Line
	if err := json.Unmarshal([]byte(contents), &lines); err != nil {
		return `order data unreadable`
	}

	var out bytes.Buffer
	for _, line := range lines {
		out.WriteString(fmt.Sprintf("%dx %s, ", line.Amount, line.Item.Name))
	}

	return out.String()
}
