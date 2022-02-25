package window

import "syscall/js"

func OnClick(f func()) {
	fn := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		f()
		return nil
	})

	js.Global().Get(`document`).Get(`body`).Call(`addEventListener`, `click`, fn)
}
