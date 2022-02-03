package window

import (
	"strconv"
	"syscall/js"
)

func GetSize() Size {
	if js.Global().Call(`matchMedia`, `only screen and (min-width: `+strconv.Itoa(minWidthL)+`px)`).Get(`matches`).Bool() {
		return SizeL
	}

	if js.Global().Call(`matchMedia`, `only screen and (min-width: `+strconv.Itoa(minWidthM)+`px)`).Get(`matches`).Bool() {
		return SizeM
	}

	return SizeS
}
