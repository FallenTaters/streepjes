package window

import (
	"net/url"
	"strconv"
	"syscall/js"
)

func Location() *url.URL {
	href := js.Global().Get(`location`).Get(`href`).String()
	u, err := url.Parse(href)
	if err != nil {
		panic(err)
	}

	return u
}

func getSize() Size {
	if js.Global().Call(`matchMedia`, `only screen and (min-width: `+strconv.Itoa(minWidthL)+`px)`).Get(`matches`).Bool() {
		return SizeL
	}

	if js.Global().Call(`matchMedia`, `only screen and (min-width: `+strconv.Itoa(minWidthM)+`px)`).Get(`matches`).Bool() {
		return SizeM
	}

	return SizeS
}
