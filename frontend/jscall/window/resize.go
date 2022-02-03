package window

import "syscall/js"

func OnResize(f func()) bool {
	largeScreen := LargeScreen()

	var destroy func()
	jsFunc := js.FuncOf(func(js.Value, []js.Value) interface{} {
		if largeScreen != LargeScreen() {
			destroy()
			if f != nil {
				f()
			}
		}

		return nil
	})

	destroy = func() { js.Global().Call(`removeEventListener`, `resize`, jsFunc) }

	js.Global().Call(`addEventListener`, `resize`, jsFunc)

	return largeScreen
}
