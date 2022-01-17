package main

import (
	_ "embed"
	"net/http"

	"git.fuyu.moe/Fuyu/router"
)

//go:embed frontend/assets/app.wasm
var wasmApp []byte

//go:embed frontend/assets/index.html
var indexPage []byte

//go:embed frontend/assets/wasm_exec.js
var wasmExec []byte

func main() {
	r := router.New()

	r.GET(`/`, handle)
	r.GET(`/wasm_exec.js`, serveFile(wasmExec))
	r.GET(`/app.wasm`, serveFile(wasmApp))

	panic(r.Start(`127.0.0.1:3000`))
}

func handle(c *router.Context) error {
	return c.Bytes(http.StatusOK, indexPage)
}

func serveFile(file []byte) router.Handle {
	return func(c *router.Context) error {
		return c.Bytes(http.StatusOK, file)
	}
}
