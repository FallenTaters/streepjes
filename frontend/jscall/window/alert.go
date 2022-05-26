package window

import "syscall/js"

func Alert(args ...any) {
	js.Global().Get(`window`).Call(`alert`, args...)
}
