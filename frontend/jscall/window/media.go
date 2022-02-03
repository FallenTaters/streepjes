package window

import "syscall/js"

func LargeScreen() bool {
	return js.Global().Call(`matchMedia`, `only screen and (min-width: 993px)`).Get(`matches`).Bool()
}
