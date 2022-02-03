package window

import "syscall/js"

type Size int

// from beercss
const (
	minWidthL = 993
	minWidthM = 601
)

const (
	SizeS Size = iota
	SizeM
	SizeL
)

// OnResize returns the current window size and adds a listener that goes off whenever the screen size changes.
// passing nil is not allowed, just use window.GetSize()
func OnResize(f func(Size)) Size {
	size := GetSize()

	if f == nil {
		panic(`OnResize called with nil listener. Did you mean to use window.GetSize()?`)
	}

	var destroy func()
	jsFunc := js.FuncOf(func(js.Value, []js.Value) interface{} {
		if size != GetSize() {
			destroy()
			f(size)
		}

		return nil
	})

	destroy = func() { js.Global().Call(`removeEventListener`, `resize`, jsFunc) }

	js.Global().Call(`addEventListener`, `resize`, jsFunc)

	return size
}
