package router

import (
	"net/http"

	"git.fuyu.moe/Fuyu/router"
	gorillaHandlers "github.com/gorilla/handlers"
)

func compressionMiddleware(next router.Handle) router.Handle {
	return func(c *router.Context) error {
		h := gorillaHandlers.CompressHandler(httpHandler{
			c: c,
			f: next,
		})

		h.ServeHTTP(c.Response, c.Request)

		return nil
	}
}

type httpHandler struct {
	c *router.Context
	f func(c *router.Context) error
}

func (h httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.c.Response = w
	h.c.Request = r

	err := h.f(h.c)
	if err != nil {
		panic(err)
	}
}
