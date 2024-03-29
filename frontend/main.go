//go:build js && wasm

package main

import (
	"github.com/FallenTaters/streepjes/frontend/authroutine"
	"github.com/FallenTaters/streepjes/frontend/backend"
	"github.com/FallenTaters/streepjes/frontend/global"
	"github.com/FallenTaters/streepjes/frontend/jscall/window"
	"github.com/FallenTaters/streepjes/frontend/store"
	"github.com/vugu/vugu"
	"github.com/vugu/vugu/domrender"
)

func main() {
	initPackages()

	startVugu()
}

func initPackages() {
	u := window.Location()
	u.Path = ``
	u.RawQuery = ``
	u.ForceQuery = false

	backend.Init(u)
	store.Init()
}

func startVugu() {
	renderer, err := domrender.New("#vugu_mount_point")
	if err != nil {
		panic(err)
	}
	defer renderer.Release()

	global.EventEnv = renderer.EventEnv()

	buildEnv, err := vugu.NewBuildEnv(global.EventEnv)
	if err != nil {
		panic(err)
	}

	root := &Root{}

	window.Listen() // for resize
	authroutine.Start()

	render(renderer, buildEnv, root)
	for renderer.EventWait() {
		render(renderer, buildEnv, root)
	}
}

func render(renderer *domrender.JSRenderer, buildEnv *vugu.BuildEnv, builder vugu.Builder) {
	buildResults := buildEnv.RunBuild(builder)
	err := renderer.Render(buildResults)
	if err != nil {
		panic(err)
	}
}
