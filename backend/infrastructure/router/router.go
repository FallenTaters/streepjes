package router

import (
	"net/http"
	"strings"

	"git.fuyu.moe/Fuyu/router"
	"github.com/PotatoesFall/vecty-test/backend/application/auth"
)

type Static func(filename string) ([]byte, error)

func New(static Static, authService auth.Service) *router.Router {
	r := router.New()

	r.ErrorHandler = panicHandler

	r.GET(`/version`, getVersion)
	r.GET(`/`, getIndex(static))
	r.GET(`/static/*name`, getStatic(static))

	bar := r.Group(``)
	bartenderRoutes(bar)

	admin := r.Group(``)
	adminRoutes(admin)

	return r
}

func getVersion(c *router.Context) error {
	return c.String(http.StatusOK, version())
}

func getIndex(assets Static) router.Handle {
	return func(c *router.Context) error {
		index, err := assets(`index.html`)
		if err != nil {
			panic(err)
		}

		c.Response.Header().Set(`Content-Type`, `text/html`)

		return c.Bytes(http.StatusOK, index)
	}
}

func getStatic(assets Static) router.Handle {
	return func(c *router.Context) error {
		name := strings.TrimPrefix(c.Param(`name`), `/`)

		asset, err := assets(name)
		if err != nil {
			return c.NoContent(http.StatusNotFound)
		}

		return c.Bytes(http.StatusOK, asset)
	}
}
