package main

import (
	"github.com/PotatoesFall/vecty-test/backend/infrastructure/router"
	"github.com/PotatoesFall/vecty-test/static"
)

func main() {
	r := router.New(static.Get)

	panic(r.Start(`:8080`))
}
