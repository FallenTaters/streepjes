package api

import (
	"sort"
	"time"

	"github.com/FallenTaters/streepjes/domain/orderdomain"
)

type LeaderboardFilter struct {
	Start      time.Time `json:"start"`
	End        time.Time `json:"end"`
	Gladiators bool      `json:"gladiators"`
	Parabool   bool      `json:"parabool"`
}

type Leaderboard struct {
	TotalPrice orderdomain.Price `json:"price"`

	Members []LeaderboardMember `json:"members"`

	// contains all items that have been purchased by name and total amount
	Items map[string]int `json:"items"`
}

// ForItems makes a leaderboard ranking for multiple item names
func (l Leaderboard) ForItems(items map[string]int) (int, []LeaderboardRank) {
	total := 0
	members := make([]LeaderboardRank, len(l.Members))

	for i, member := range l.Members {
		members[i] = LeaderboardRank{
			Member: member.Member,
			Amount: 0,
		}

		for name, weight := range items {
			total += member.Amounts[name] * weight
			members[i].Amount += member.Amounts[name] * weight
		}
	}

	sort.Slice(members, func(i, j int) bool {
		return members[i].Amount < members[j].Amount
	})

	return total, members
}

type LeaderboardMember struct {
	orderdomain.Member

	Total   orderdomain.Price `json:"amount"`
	Amounts map[string]int    `json:"items"`
}

type LeaderboardRank struct {
	orderdomain.Member

	Amount int `json:"amount"`
}
