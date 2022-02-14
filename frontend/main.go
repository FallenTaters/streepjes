//go:build js && wasm

package main

import (
	"github.com/PotatoesFall/vecty-test/frontend/backend"
	"github.com/PotatoesFall/vecty-test/frontend/jscall/window"
	"github.com/vugu/vugu"
	"github.com/vugu/vugu/domrender"
)

func main() {
	initPackages()

	startVugu()
}

func initPackages() {
	backend.Init(window.Location())
}

func startVugu() {
	renderer, err := domrender.New("#vugu_mount_point")
	if err != nil {
		panic(err)
	}
	defer renderer.Release()

	buildEnv, err := vugu.NewBuildEnv(renderer.EventEnv())
	if err != nil {
		panic(err)
	}

	root := &Root{}

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
