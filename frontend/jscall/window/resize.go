package window

import (
	"syscall/js"

	"github.com/PotatoesFall/vecty-test/frontend/global"
)

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

var listeners map[string]func(Size)

var size Size

// Listen is an init function that should be called as an Init step
func Listen() {
	listeners = make(map[string]func(Size))
	size = getSize()

	// wrap callback as js.Func
	jsFunc := js.FuncOf(func(js.Value, []js.Value) interface{} {
		newSize := getSize()
		if size == newSize {
			return nil
		}

		size = newSize

		global.EventEnv.Lock()
		defer global.EventEnv.UnlockRender()

		for _, listener := range listeners {
			listener(size)
		}

		return nil
	})

	js.Global().Call(`addEventListener`, `resize`, jsFunc)
}

// OnResize returns the current window size and adds a listener that goes off whenever the screen size changes.
// passing nil is not allowed, just use window.Size()
// OnResize will call a lock on the global EventEnv() and trigger a Rerender
func OnResize(key string, f func(Size)) Size {
	if f == nil {
		panic(`OnResize called with nil listener. Did you mean to use window.GetSize()?`)
	}

	listeners[key] = f

	return getSize()
}

// GetSize returns the current size
func GetSize() Size {
	return size
}
