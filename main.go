//go:build !dev

package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/FallenTaters/streepjes/backend/application/auth"
	"github.com/FallenTaters/streepjes/backend/application/order"
	"github.com/FallenTaters/streepjes/backend/global/settings"
	"github.com/FallenTaters/streepjes/backend/infrastructure/repo"
	"github.com/FallenTaters/streepjes/backend/infrastructure/repo/sqlite"
	"github.com/FallenTaters/streepjes/backend/infrastructure/router"
	"github.com/FallenTaters/streepjes/domain"
	"github.com/FallenTaters/streepjes/domain/authdomain"
	"github.com/FallenTaters/streepjes/static"
	"github.com/charmbracelet/log"
)

func main() {
	os.Exit(run())
}

func run() int {
	log.Default().SetLevel(log.DebugLevel)

	dbPath := os.Getenv("STREEPJES_DB_PATH")
	if dbPath == "" {
		dbPath = "streepjes.db"
	}

	db, err := sqlite.OpenDB(dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sigChan := make(chan os.Signal)
	shutdown := make(chan int)
	signal.Notify(sigChan, os.Interrupt, os.Kill)
	go func() {
		<-sigChan
		shutdown <- 0
	}()

	readSettings()

	sqlite.Migrate(db)

	userRepo := sqlite.NewUserRepo(db)
	memberRepo := sqlite.NewMemberRepo(db)
	orderRepo := sqlite.NewOrderRepo(db)
	catalogRepo := sqlite.NewCatalogRepo(db)

	authService := auth.New(userRepo, orderRepo)
	checkNoUsers(userRepo, authService)

	orderService := order.New(memberRepo, orderRepo, catalogRepo)

	handler := router.New(static.Get, authService, orderService)

	fmt.Printf("Starting server on port %d\n", settings.Port)
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(`:%d`, settings.Port), handler)
		log.Fatal("server exited", "error", err)
		shutdown <- 1
	}()

	return <-shutdown
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
