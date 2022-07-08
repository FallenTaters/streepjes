package window

import (
	"fmt"
	"syscall/js"
)

func Alert(args ...any) {
	for i := range args {
		args[i] = fmt.Sprint(args[i])
	}

	js.Global().Get(`window`).Call(`alert`, args...)
}

func Confirm(msg string) bool {
	return js.Global().Get(`window`).Call(`confirm`, msg).Bool()
}
