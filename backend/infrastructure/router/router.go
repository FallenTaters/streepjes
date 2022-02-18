package router

import (
	"net/http"
	"strings"

	"git.fuyu.moe/Fuyu/router"
	"github.com/PotatoesFall/vecty-test/backend/application/auth"
	"github.com/PotatoesFall/vecty-test/domain"
)

type Static func(filename string) ([]byte, error)

func New(static Static, authService auth.Service) *router.Router {
	r := router.New()

	r.ErrorHandler = panicHandler

	publicRoutes(r, static, authService)

	auth := r.Group(`/`, authMiddleware(authService))
	authRoutes(auth, authService)

	bar := auth.Group(`/`, roleMiddleware(domain.RoleBartender))
	bartenderRoutes(bar)

	admin := auth.Group(`/`, roleMiddleware(domain.RoleAdmin))
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
