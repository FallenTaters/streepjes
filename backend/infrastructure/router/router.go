package router

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/FallenTaters/streepjes/backend/application/auth"
	"github.com/FallenTaters/streepjes/backend/application/order"
	"github.com/FallenTaters/streepjes/domain/authdomain"
)

type Static func(filename string) ([]byte, error)

func New(static Static, authService auth.Service, orderService order.Service) http.Handler {
	r := echo.New()

	r.Use(middleware.Recover())
	r.HTTPErrorHandler = func(err error, ctx echo.Context) {
		r.DefaultHTTPErrorHandler(err, ctx)
	}

	auth := r.Group(``, authMiddleware(authService))
	authRoutes(auth, authService)

	bar := auth.Group(``, permissionMiddleware(authdomain.PermissionBarStuff))
	bartenderRoutes(bar, orderService)

	admin := auth.Group(``, permissionMiddleware(authdomain.PermissionAdminStuff))
	adminRoutes(admin)

	// must go last because of https://github.com/labstack/echo/issues/2141
	publicRoutes(r, static, authService)

	return r
}

func readJSON[T any](c echo.Context) (T, bool) {
	var t T
	err := json.NewDecoder(c.Request().Body).Decode(&t)
	if err != nil {
		c.NoContent(http.StatusBadRequest)
		return t, false
	}

	return t, true
}
