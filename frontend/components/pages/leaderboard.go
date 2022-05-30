package pages

import (
	"time"

	"github.com/FallenTaters/streepjes/api"
	"github.com/FallenTaters/streepjes/frontend/backend"
)

type Leaderboard struct {
	Leaderboard api.Leaderboard
}

func (l *Leaderboard) Init() {
	// TODO: Loading, error handling, etc
	leaderboard, _ := backend.GetLeaderboard(api.LeaderboardFilter{
		Start: time.Now().AddDate(-10, 0, 0),
		End:   time.Now().AddDate(10, 0, 0),
	})

	l.Leaderboard = leaderboard
}
