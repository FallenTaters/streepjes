package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/FallenTaters/streepjes/backend/application/auth"
	"github.com/labstack/echo/v4"
)

func publicRoutes(r *echo.Echo, static Static, authService auth.Service) {
	r.GET(``, getIndex(static))
	r.GET(`/version`, getVersion)
	r.GET(`/static/*`, getStatic(static))
	r.POST(`/login`, postLogin(authService))
}

var (
	buildVersion string
	buildCommit  string
	buildTime    string
)

func version() string {
	return fmt.Sprintf("Version: %s\nDate: %s\nCommit: %s",
		buildVersion, buildTime, buildCommit)
}

func getVersion(c echo.Context) error {
	return c.String(http.StatusOK, version())
}

func getIndex(assets Static) echo.HandlerFunc {
	return func(c echo.Context) error {
		index, err := assets(`index.html`)
		if err != nil {
			panic(err)
		}

		return c.Blob(http.StatusOK, `text/html`, index)
	}
}

func getStatic(assets Static) echo.HandlerFunc {
	return func(c echo.Context) error {
		name := strings.TrimPrefix(strings.TrimPrefix(c.Request().URL.Path, `/`), `static/`)

		asset, err := assets(name)
		if err != nil {
			return c.NoContent(http.StatusNotFound)
		}

		return c.Blob(http.StatusOK, http.DetectContentType(asset), asset)
	}
}
