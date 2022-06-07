package window

import "syscall/js"

func Alert(args ...any) {
	js.Global().Get(`window`).Call(`alert`, args...)
}

func Confirm(msg string) bool {
	return js.Global().Get(`window`).Call(`confirm`, msg).Bool()
}
