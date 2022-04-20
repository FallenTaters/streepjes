package router

import (
	"net/http"
	"time"

	"github.com/FallenTaters/streepjes/api"
	"github.com/FallenTaters/streepjes/backend/application/auth"
	"github.com/FallenTaters/streepjes/backend/global/settings"
	"github.com/FallenTaters/streepjes/domain/authdomain"
	"github.com/labstack/echo/v4"
)

func userFromContext(c echo.Context) authdomain.User {
	return c.Get(`user`).(authdomain.User)
}

func authRoutes(r *echo.Group, authService auth.Service) {
	r.POST(`/logout`, postLogout(authService))
	r.POST(`/active`, postActive)

	r.POST(`/me/name`, postMeName(authService))
	r.POST(`/me/password`, postMePassword(authService))
}

func postLogin(authService auth.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		credentials, ok := readJSON[api.Credentials](c)
		if !ok {
			return nil
		}

		user, ok := authService.Login(credentials.Username, credentials.Password)
		if !ok {
			return c.NoContent(http.StatusUnauthorized)
		}

		http.SetCookie(c.Response(), &http.Cookie{ //nolint:exhaustivestruct
			Name:   api.AuthTokenCookieName,
			Value:  user.AuthToken,
			Path:   ``,
			Domain: ``,
			MaxAge: 24 * int(time.Hour/time.Second),
			Secure: !settings.DisableSecure,
		})

		return c.JSON(http.StatusOK, user)
	}
}

func postLogout(authService auth.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := userFromContext(c)

		authService.Logout(user.ID)

		return c.NoContent(http.StatusOK)
	}
}

func postActive(c echo.Context) error {
	return c.JSON(http.StatusOK, userFromContext(c))
}

func authMiddleware(authService auth.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := c.Request().Cookie(`auth_token`)
			if err != nil {
				return c.NoContent(http.StatusUnauthorized)
			}

			user, ok := authService.Check(token.Value)
			if !ok {
				return c.NoContent(http.StatusUnauthorized)
			}

			c.Set(`user`, user)

			return next(c)
		}
	}
}

func permissionMiddleware(permission authdomain.Permission) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := userFromContext(c)

			if !user.Role.Has(permission) {
				return c.NoContent(http.StatusForbidden)
			}

			return next(c)
		}
	}
}

func postMeName(authService auth.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		name, ok := readJSON[string](c)
		if !ok {
			return nil
		}

		if !authService.ChangeName(userFromContext(c), name) {
			return c.NoContent(http.StatusBadRequest)
		}

		return c.NoContent(http.StatusOK)
	}
}

func postMePassword(authService auth.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		changePassword, ok := readJSON[api.ChangePassword](c)
		if !ok {
			return nil
		}

		if !authService.ChangePassword(userFromContext(c), changePassword) {
			return c.NoContent(http.StatusBadRequest)
		}

		return c.NoContent(http.StatusOK)
	}
}
