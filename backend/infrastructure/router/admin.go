package router

import (
	"net/http"

	"github.com/FallenTaters/streepjes/backend/application/auth"
	"github.com/labstack/echo/v4"
)

func adminRoutes(r *echo.Group, authService auth.Service) {
	r.GET(`/users`, getUsers(authService))
}

func getUsers(authService auth.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		users := authService.GetUsers()
		return c.JSON(http.StatusOK, users)
	}
}
