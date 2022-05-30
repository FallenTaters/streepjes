package pages

import (
	"time"

	"github.com/FallenTaters/streepjes/api"
	"github.com/FallenTaters/streepjes/frontend/backend"
)

type Leaderboard struct {
	Loading bool
	Error   bool

	Items []string

	Leaderboard api.Leaderboard

	Total string
	Ranking []orderdomain.LeaderboardRank
}

func (l *Leaderboard) Init() {
	// TODO: Loading, error handling, etc
	leaderboard, _ := backend.GetLeaderboard(api.LeaderboardFilter{
		Start: time.Now().AddDate(-10, 0, 0),
		End:   time.Now().AddDate(10, 0, 0),
	})

	l.Leaderboard = leaderboard
}

func (l *Leaderboard) Refresh() {

}

func (l *Leaderboard) Total() string {
	if len(l.Items) == 0 {
		return l.Leaderboard.TotalPrice.String()
	}

	return l.
}
