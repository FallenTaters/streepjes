package router

import (
	"net/http"
	"strings"

	"git.fuyu.moe/Fuyu/router"
	"github.com/PotatoesFall/vecty-test/api"
	"github.com/PotatoesFall/vecty-test/backend/application/auth"
)

type Static func(filename string) ([]byte, error)

func New(static Static, authService auth.Service) *router.Router {
	r := router.New()

	r.ErrorHandler = panicHandler

	r.GET(`/version`, getVersion)
	r.GET(`/`, getIndex(static))
	r.GET(`/static/*name`, getStatic(static))
	r.POST(`/login`, postLogin(authService))

	bar := r.Group(``) // TODO: add rolemiddleware
	bartenderRoutes(bar)

	admin := r.Group(``) // TODO: add rolemiddleware
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

func postLogin(authService auth.Service) func(*router.Context, api.Credentials) error {
	return func(c *router.Context, credentials api.Credentials) error {
		user, ok := authService.Login(credentials.Username, credentials.Password)
		if !ok {
			return c.NoContent(http.StatusUnauthorized)
		}

		return c.JSON(http.StatusOK, api.Token{Token: user.AuthToken})
	}
}
