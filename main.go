package main

import (
	"github.com/PotatoesFall/vecty-test/backend/application/auth"
	"github.com/PotatoesFall/vecty-test/backend/infrastructure/repo/sqlite"
	"github.com/PotatoesFall/vecty-test/backend/infrastructure/router"
	"github.com/PotatoesFall/vecty-test/static"
)

func main() {
	db, err := sqlite.OpenDB(`streepjes.db`)
	if err != nil {
		panic(err)
	}

	sqlite.Migrate(db)

	userRepo := sqlite.NewUserRepo(db)

	r := router.New(static.Get, auth.New(userRepo))

	panic(r.Start(`:8080`))
}
