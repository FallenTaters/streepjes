package jscall

import (
	"fmt"
	"syscall/js"
)

func Focus(id string) {
	fmt.Println(`focus1`)
	elem := js.Global().Get(`document`).Call(`getElementById`, id)
	if !elem.Truthy() || !elem.Get(`focus`).Truthy() {
		return
	}

	fmt.Println(`focus2`)
	elem.Call(`focus`)
}
