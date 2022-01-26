package window

import "syscall/js"

func LargeScreen() bool {
	return js.Global().Call(`matchMedia`, `only screen and (min-width: 1200px)`).Get(`matches`).Bool()
}
