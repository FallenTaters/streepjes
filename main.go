//go:build !dev

package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/PotatoesFall/vecty-test/backend/application/auth"
	"github.com/PotatoesFall/vecty-test/backend/global/settings"
	"github.com/PotatoesFall/vecty-test/backend/infrastructure/repo"
	"github.com/PotatoesFall/vecty-test/backend/infrastructure/repo/sqlite"
	"github.com/PotatoesFall/vecty-test/backend/infrastructure/router"
	"github.com/PotatoesFall/vecty-test/domain"
	"github.com/PotatoesFall/vecty-test/domain/authdomain"
	"github.com/PotatoesFall/vecty-test/static"
)

func main() {
	db, err := sqlite.OpenDB(`streepjes.db`)
	if err != nil {
		panic(err)
	}

	readSettings()

	sqlite.Migrate(db)

	userRepo := sqlite.NewUserRepo(db)

	authService := auth.New(userRepo)
	checkNoUsers(userRepo, authService)

	r := router.New(static.Get, authService)

	fmt.Printf("Starting server on port %d\n", settings.Port)
	panic(r.Start(fmt.Sprintf(`:%d`, settings.Port)))
}

// check if there are no users in the database, if so, insert some
func checkNoUsers(userRepo repo.User, authService auth.Service) {
	if len(userRepo.GetAll()) == 0 {
		_ = authService.Register(authdomain.User{ //nolint:exhaustivestruct
			Username: `adminGladiators`,
			Club:     domain.ClubGladiators,
			Name:     `Gladiators Admin`,
			Role:     authdomain.RoleAdmin,
		}, `playlacrossebecauseitsfun`)
		_ = authService.Register(authdomain.User{ //nolint:exhaustivestruct
			Username: `adminParabool`,
			Club:     domain.ClubParabool,
			Name:     `Parabool Admin`,
			Role:     authdomain.RoleAdmin,
		}, `groningerstudentenkorfbalcommissie`)
	}
}

// read settings from environment
func readSettings() {
	disableSecure, ok := os.LookupEnv(`STREEPJES_DISABLE_SECURE`)
	settings.DisableSecure = ok && disableSecure == `true`

	port, ok := os.LookupEnv(`STREEPJES_PORT`)
	if ok {
		portN, err := strconv.Atoi(port)
		if err != nil {
			panic(err)
		}

		settings.Port = portN
	}
}
