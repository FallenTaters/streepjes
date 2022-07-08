package window

import (
	"net/url"
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

func NewTab(url string) {
	js.Global().Get(`window`).Call(`open`, url, `_blank`).Call(`focus`)
}
