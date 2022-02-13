package beercss

import "syscall/js"

func UI() {
	js.Global().Call(`ui`)
}
