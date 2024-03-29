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

	_ "time/tzdata"
)

func init() { //nolint:gochecknoinits
	tz, err := time.LoadLocation("Europe/Amsterdam")
	if err != nil {
		panic(err)
	}

	timezone = tz
}

var timezone *time.Location

func writeCSV(orders []orderdomain.Order, members []orderdomain.Member) []byte { //nolint:funlen,cyclop
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
	err := w.Write([]string{`Member`, `Price`, `Date`, `Order`})
	if err != nil {
		panic(err)
	}

	memberTotals := make(map[int]orderdomain.Price, len(members))
	for _, o := range orders {
		memberTotals[o.MemberID] += o.Price
		err = w.Write([]string{
			membersByID[o.MemberID].Name,
			o.Price.String(),
			o.OrderTime.In(timezone).Format(`2006-01-02 15:04`),
			parseOrderLines(o.Contents),
		})
		if err != nil {
			panic(err)
		}
	}
	err = w.Write(nil)
	if err != nil {
		panic(err)
	}
	err = w.Write([]string{`Member`, `Total`})
	if err != nil {
		panic(err)
	}
	for id, total := range memberTotals {
		err := w.Write([]string{membersByID[id].Name, total.String()})
		if err != nil {
			panic(err)
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func parseOrderLines(contents string) string {
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
