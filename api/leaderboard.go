package api

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/FallenTaters/streepjes/domain/orderdomain"
)

type LeaderboardFilter struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type Leaderboard struct {
	TotalPrice orderdomain.Price `json:"price"`

	Members []LeaderboardMember `json:"members"`

	// contains all items that have been purchased by name and total amount
	Items map[string]int `json:"items"`
}

func (l Leaderboard) MoneyRanking() (orderdomain.Price, []LeaderboardRank) {
	ranking := make([]LeaderboardRank, len(l.Members))
	for i, member := range l.Members {
		ranking[i] = LeaderboardRank{
			Member:       member.Member,
			sortingValue: int(member.Total),
			Total:        member.Total.String(),
			ItemInfo:     member.Amounts,
		}
	}

	// sorting not needed?

	return l.TotalPrice, ranking
}

// ItemRanking makes a leaderboard ranking for multiple item names
func (l Leaderboard) ItemRanking(items map[string]int) (int, []LeaderboardRank) {
	total := 0
	ranking := make([]LeaderboardRank, 0, len(l.Members))

	for _, member := range l.Members {
		var memberTotal int
		for name, weight := range items {
			total += member.Amounts[name] * weight
			memberTotal += member.Amounts[name] * weight
		}

		if memberTotal == 0 {
			continue
		}

		ranking = append(ranking, LeaderboardRank{
			Member:       member.Member,
			sortingValue: memberTotal,
			Total:        strconv.Itoa(memberTotal),
			ItemInfo:     member.Amounts,
		})
	}

	sort.Slice(ranking, func(i, j int) bool {
		if ranking[i].sortingValue == ranking[j].sortingValue {
			return strings.Compare(ranking[i].Name, ranking[j].Name) < 0
		}

		return ranking[i].sortingValue > ranking[j].sortingValue
	})

	return total, ranking
}

type LeaderboardMember struct {
	orderdomain.Member

	Total   orderdomain.Price `json:"amount"`
	Amounts map[string]int    `json:"items"`
}

type LeaderboardRank struct {
	orderdomain.Member

	Total        string
	ItemInfo     map[string]int
	sortingValue int
}
