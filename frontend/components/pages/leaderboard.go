package pages

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/FallenTaters/streepjes/api"
	"github.com/FallenTaters/streepjes/domain"
	"github.com/FallenTaters/streepjes/frontend/backend"
)

type Leaderboard struct {
	Loading bool
	Error   bool

	ItemWeights map[string]int

	Leaderboard api.Leaderboard

	ShowExpansion map[string]bool

	Gladiators, Parabool bool

	// Display state
	Total   string                `vugu:"data"`
	Ranking []api.LeaderboardRank `vugu:"data"`
}

func (l *Leaderboard) Init() {
	// TODO: Loading, error handling, etc
	leaderboard, _ := backend.GetLeaderboard(api.LeaderboardFilter{
		Start: time.Now().AddDate(-10, 0, 0),
		End:   time.Now().AddDate(10, 0, 0),
	})

	l.Leaderboard = leaderboard

	// TODO make this adjustable ?
	l.ItemWeights = map[string]int{
		`Bier`:           1,
		`Weizen glas`:    1,
		`Pitcher`:        5,
		`Weizen Pitcher`: 5,
		`Flugel`:         1,
		`Bier Barcie`:    1,
		`Bier BarCie`:    1,
		`Seltzer BarCie`: 1,
		`Seltzer`:        1,
		`Wine Bottle`:    5,
		`Wijn`:           1,
		`Radler`:         0,
	}

	l.Gladiators = true

	l.ShowExpansion = make(map[string]bool)

	l.Refresh()
}

func (l *Leaderboard) Refresh() {
	if len(l.ItemWeights) == 0 {
		total, ranking := l.Leaderboard.MoneyRanking()
		l.Total = total.String()
		l.Ranking = l.FilterRanking(ranking)
		return
	}

	total, ranking := l.Leaderboard.ItemRanking(l.ItemWeights)
	l.Total = strconv.Itoa(total)
	l.Ranking = l.FilterRanking(ranking)
	return
}

func (l *Leaderboard) FilterRanking(ranking []api.LeaderboardRank) []api.LeaderboardRank {
	newRanking := make([]api.LeaderboardRank, 0, len(ranking))

	for _, rank := range ranking {
		if (rank.Club == domain.ClubGladiators && l.Gladiators) ||
			(rank.Club == domain.ClubParabool && l.Parabool) {
			newRanking = append(newRanking, rank)
		}
	}

	return newRanking
}

type MsgCount struct {
	Msg   string
	count int
}

func (l *Leaderboard) SortItemInfo(itemInfo map[string]int) []MsgCount {
	out := make([]MsgCount, 0, len(itemInfo))
	for name, count := range itemInfo {
		if len(l.ItemWeights) == 0 || l.ItemWeights[name] > 0 {
			out = append(out, MsgCount{
				Msg:   strconv.Itoa(count) + ` ` + name,
				count: count,
			})
		}
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].count == out[j].count {
			return strings.Compare(out[i].Msg, out[j].Msg) < 0
		}

		return out[i].count > out[j].count
	})

	return out
}
