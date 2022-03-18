package jscall

import (
	"syscall/js"
)

func Focus(id string) {
	elem := js.Global().Get(`document`).Call(`getElementById`, id)
	if !elem.Truthy() || !elem.Get(`focus`).Truthy() {
		return
	}

	elem.Call(`focus`)
}
