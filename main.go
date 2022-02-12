package main

import (
	"github.com/PotatoesFall/vecty-test/backend/application/auth"
	"github.com/PotatoesFall/vecty-test/backend/infrastructure/router"
	"github.com/PotatoesFall/vecty-test/static"
)

func main() {
	r := router.New(static.Get, auth.New())

	panic(r.Start(`:8080`))
}
