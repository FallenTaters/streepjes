package order

import (
	"strings"
	"testing"
	"time"

	"git.fuyu.moe/Fuyu/assert"
	"github.com/FallenTaters/streepjes/domain"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
)

func TestWriteCSV(t *testing.T) {
	t.Parallel()

	tz := time.UTC
	members := []orderdomain.Member{
		{ID: 1, Name: "Alice", Club: domain.ClubGladiators},
		{ID: 2, Name: "Bob", Club: domain.ClubGladiators},
	}

	t.Run("basic output with header and totals", func(t *testing.T) {
		assert := assert.New(t)

		orders := []orderdomain.Order{
			{MemberID: 1, Price: 200, OrderTime: time.Date(2025, 6, 15, 10, 0, 0, 0, time.UTC), Contents: `[]`},
			{MemberID: 2, Price: 300, OrderTime: time.Date(2025, 6, 15, 11, 0, 0, 0, time.UTC), Contents: `[]`},
		}

		result := string(writeCSV(orders, members, tz))
		assert.True(strings.Contains(result, "Member,Price,Date,Order"))
		assert.True(strings.Contains(result, "Alice"))
		assert.True(strings.Contains(result, "Bob"))
		assert.True(strings.Contains(result, "Member,Total"))
	})

	t.Run("filters orders without member", func(t *testing.T) {
		assert := assert.New(t)

		orders := []orderdomain.Order{
			{MemberID: 0, Price: 100, OrderTime: time.Date(2025, 6, 15, 10, 0, 0, 0, time.UTC), Contents: `[]`},
			{MemberID: 1, Price: 200, OrderTime: time.Date(2025, 6, 15, 11, 0, 0, 0, time.UTC), Contents: `[]`},
		}

		result := string(writeCSV(orders, members, tz))
		assert.True(strings.Contains(result, "Alice"))
		lines := strings.Split(strings.TrimSpace(result), "\n")

		dataLines := 0
		for _, line := range lines {
			if strings.HasPrefix(line, "Alice") || strings.HasPrefix(line, "Bob") {
				dataLines++
			}
		}
		assert.Eq(2, dataLines) // 1 order line + 1 total line for Alice
	})

	t.Run("sorts by member name", func(t *testing.T) {
		assert := assert.New(t)

		orders := []orderdomain.Order{
			{MemberID: 2, Price: 300, OrderTime: time.Date(2025, 6, 15, 10, 0, 0, 0, time.UTC), Contents: `[]`},
			{MemberID: 1, Price: 200, OrderTime: time.Date(2025, 6, 15, 11, 0, 0, 0, time.UTC), Contents: `[]`},
		}

		result := string(writeCSV(orders, members, tz))
		aliceIdx := strings.Index(result, "Alice")
		bobIdx := strings.Index(result, "Bob")
		assert.True(aliceIdx < bobIdx)
	})

	t.Run("empty orders", func(t *testing.T) {
		assert := assert.New(t)

		result := string(writeCSV(nil, members, tz))
		assert.True(strings.Contains(result, "Member,Price,Date,Order"))
		assert.True(strings.Contains(result, "Member,Total"))
	})
}

func TestCsvOrderLines(t *testing.T) {
	t.Parallel()

	t.Run("valid contents", func(t *testing.T) {
		assert := assert.New(t)
		result := csvOrderLines(`[{"product":{"name":"Beer"},"amount":2},{"product":{"name":"Wine"},"amount":1}]`)
		assert.Eq("2x Beer, 1x Wine, ", result)
	})

	t.Run("empty list", func(t *testing.T) {
		assert := assert.New(t)
		result := csvOrderLines(`[]`)
		assert.Eq("", result)
	})

	t.Run("invalid json", func(t *testing.T) {
		assert := assert.New(t)
		result := csvOrderLines(`not json`)
		assert.Eq("order data unreadable", result)
	})
}
