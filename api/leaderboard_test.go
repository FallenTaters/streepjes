package api_test

import (
	"testing"

	"git.fuyu.moe/Fuyu/assert"
	"github.com/FallenTaters/streepjes/api"
	"github.com/FallenTaters/streepjes/domain"
	"github.com/FallenTaters/streepjes/domain/orderdomain"
)

func TestMoneyRanking(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	lb := api.Leaderboard{
		TotalPrice: 600,
		Members: []api.LeaderboardMember{
			{
				Member:  orderdomain.Member{ID: 1, Name: "Alice", Club: domain.ClubGladiators},
				Total:   300,
				Amounts: map[string]int{"Beer": 3},
			},
			{
				Member:  orderdomain.Member{ID: 2, Name: "Bob", Club: domain.ClubGladiators},
				Total:   200,
				Amounts: map[string]int{"Wine": 2},
			},
			{
				Member: orderdomain.Member{ID: 3, Name: "Carol", Club: domain.ClubGladiators},
				Total:  100,
			},
		},
	}

	total, ranking := lb.MoneyRanking()

	assert.Eq(orderdomain.Price(600), total)
	assert.Eq(3, len(ranking))
	assert.Eq("Alice", ranking[0].Name)
	assert.Eq("€3.00", ranking[0].Total)
	assert.Eq("Bob", ranking[1].Name)
}

func TestItemRanking(t *testing.T) {
	t.Parallel()

	t.Run("ranks by weighted item count", func(t *testing.T) {
		assert := assert.New(t)

		lb := api.Leaderboard{
			Members: []api.LeaderboardMember{
				{
					Member:  orderdomain.Member{ID: 1, Name: "Alice"},
					Amounts: map[string]int{"Beer": 1, "Wine": 1},
				},
				{
					Member:  orderdomain.Member{ID: 2, Name: "Bob"},
					Amounts: map[string]int{"Beer": 5},
				},
			},
		}

		items := map[string]int{"Beer": 1}
		total, ranking := lb.ItemRanking(items)

		assert.Eq(6, total)
		assert.Eq(2, len(ranking))
		assert.Eq("Bob", ranking[0].Name)
		assert.Eq("5", ranking[0].Total)
		assert.Eq("Alice", ranking[1].Name)
		assert.Eq("1", ranking[1].Total)
	})

	t.Run("excludes members with zero items", func(t *testing.T) {
		assert := assert.New(t)

		lb := api.Leaderboard{
			Members: []api.LeaderboardMember{
				{
					Member:  orderdomain.Member{ID: 1, Name: "Alice"},
					Amounts: map[string]int{"Beer": 3},
				},
				{
					Member:  orderdomain.Member{ID: 2, Name: "Bob"},
					Amounts: map[string]int{"Wine": 5},
				},
			},
		}

		items := map[string]int{"Beer": 1}
		_, ranking := lb.ItemRanking(items)

		assert.Eq(1, len(ranking))
		assert.Eq("Alice", ranking[0].Name)
	})

	t.Run("weights items", func(t *testing.T) {
		assert := assert.New(t)

		lb := api.Leaderboard{
			Members: []api.LeaderboardMember{
				{
					Member:  orderdomain.Member{ID: 1, Name: "Alice"},
					Amounts: map[string]int{"Beer": 2, "Shot": 1},
				},
			},
		}

		items := map[string]int{"Beer": 1, "Shot": 3}
		total, ranking := lb.ItemRanking(items)

		assert.Eq(5, total) // 2*1 + 1*3
		assert.Eq(1, len(ranking))
		assert.Eq("5", ranking[0].Total)
	})

	t.Run("empty members", func(t *testing.T) {
		assert := assert.New(t)

		lb := api.Leaderboard{}
		_, ranking := lb.ItemRanking(map[string]int{"Beer": 1})
		assert.Eq(0, len(ranking))
	})
}
