package order

import (
	"testing"

	"git.fuyu.moe/Fuyu/assert"
	"github.com/FallenTaters/streepjes/domain"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
)

func TestMakeLeaderboard(t *testing.T) {
	t.Parallel()

	members := []orderdomain.Member{
		{ID: 1, Name: "Alice", Club: domain.ClubGladiators},
		{ID: 2, Name: "Bob", Club: domain.ClubGladiators},
		{ID: 3, Name: "Carol", Club: domain.ClubGladiators},
	}

	t.Run("sorts members by total descending", func(t *testing.T) {
		assert := assert.New(t)

		orders := []orderdomain.Order{
			{MemberID: 1, Price: 100, Status: orderdomain.StatusOpen, Contents: `[]`},
			{MemberID: 2, Price: 300, Status: orderdomain.StatusOpen, Contents: `[]`},
			{MemberID: 3, Price: 200, Status: orderdomain.StatusOpen, Contents: `[]`},
		}

		lb := makeLeaderboard(members, orders)

		assert.Eq(3, len(lb.Members))
		assert.Eq("Bob", lb.Members[0].Member.Name)
		assert.Eq("Carol", lb.Members[1].Member.Name)
		assert.Eq("Alice", lb.Members[2].Member.Name)
	})

	t.Run("stable sort for equal totals", func(t *testing.T) {
		assert := assert.New(t)

		orders := []orderdomain.Order{
			{MemberID: 1, Price: 100, Status: orderdomain.StatusOpen, Contents: `[]`},
			{MemberID: 2, Price: 100, Status: orderdomain.StatusOpen, Contents: `[]`},
			{MemberID: 3, Price: 100, Status: orderdomain.StatusOpen, Contents: `[]`},
		}

		lb := makeLeaderboard(members, orders)

		assert.Eq(3, len(lb.Members))
		assert.Eq(orderdomain.Price(100), lb.Members[0].Total)
		assert.Eq(orderdomain.Price(100), lb.Members[1].Total)
		assert.Eq(orderdomain.Price(100), lb.Members[2].Total)
	})

	t.Run("excludes cancelled orders", func(t *testing.T) {
		assert := assert.New(t)

		orders := []orderdomain.Order{
			{MemberID: 1, Price: 100, Status: orderdomain.StatusOpen, Contents: `[]`},
			{MemberID: 1, Price: 200, Status: orderdomain.StatusCancelled, Contents: `[]`},
		}

		lb := makeLeaderboard(members, orders)

		assert.Eq(orderdomain.Price(100), lb.TotalPrice)
		assert.Eq(orderdomain.Price(100), lb.Members[0].Total)
	})

	t.Run("counts items from order contents", func(t *testing.T) {
		assert := assert.New(t)

		orders := []orderdomain.Order{
			{
				MemberID: 1,
				Price:    300,
				Status:   orderdomain.StatusOpen,
				Contents: `[{"product":{"name":"Beer"},"amount":2},{"product":{"name":"Wine"},"amount":1}]`,
			},
		}

		lb := makeLeaderboard(members, orders)

		assert.Eq(2, lb.Items["Beer"])
		assert.Eq(1, lb.Items["Wine"])
		assert.Eq(2, lb.Members[0].Amounts["Beer"])
		assert.Eq(1, lb.Members[0].Amounts["Wine"])
	})

	t.Run("handles malformed contents gracefully", func(t *testing.T) {
		assert := assert.New(t)

		orders := []orderdomain.Order{
			{MemberID: 1, Price: 100, Status: orderdomain.StatusOpen, Contents: `not json`},
		}

		lb := makeLeaderboard(members, orders)

		assert.Eq(orderdomain.Price(100), lb.TotalPrice)
		assert.Eq(orderdomain.Price(100), lb.Members[0].Total)
	})

	t.Run("empty orders", func(t *testing.T) {
		assert := assert.New(t)

		lb := makeLeaderboard(members, nil)

		assert.Eq(orderdomain.Price(0), lb.TotalPrice)
		assert.Eq(3, len(lb.Members))
	})
}
