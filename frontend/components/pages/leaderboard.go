package pages

import (
	"sort"
	"strconv"
	"time"

	"github.com/FallenTaters/streepjes/api"
	"github.com/FallenTaters/streepjes/frontend/backend"
)

type Leaderboard struct {
	Loading bool
	Error   bool

	ItemWeights map[string]int

	Leaderboard api.Leaderboard

	ShowExpansion map[string]bool

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

	l.ShowExpansion = make(map[string]bool)

	l.Refresh()
}

func (l *Leaderboard) Refresh() {
	if len(l.ItemWeights) == 0 {
		total, ranking := l.Leaderboard.MoneyRanking()
		l.Total = total.String()
		l.Ranking = ranking
		return
	}

	total, ranking := l.Leaderboard.ItemRanking(l.ItemWeights)
	l.Total = strconv.Itoa(total)
	l.Ranking = ranking
	return
}

func (l *Leaderboard) calcTotal() string {
	if len(l.ItemWeights) == 0 {
		return l.Leaderboard.TotalPrice.String()
	}

	total := 0
	for item, weight := range l.ItemWeights {
		total += l.Leaderboard.Items[item] * weight
	}

	return strconv.Itoa(total)
}

func (l *Leaderboard) calcRanking() []api.LeaderboardRank {
	ranking := make([]api.LeaderboardRank, 0, len(l.Leaderboard.Members))

	for _, member := range l.Leaderboard.Members {
		// money only
		if len(l.ItemWeights) == 0 {
			ranking = append(ranking, api.LeaderboardRank{
				Total: member.Total.String(),
			})
			continue
		}

		// items/weights
		l.Leaderboard.ItemRanking(l.ItemWeights)
	}

	return ranking
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
		return out[i].count > out[j].count
	})

	return out
}
