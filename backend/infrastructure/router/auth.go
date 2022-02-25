package router

import (
	"net/http"
	"time"

	"git.fuyu.moe/Fuyu/router"
	"github.com/PotatoesFall/vecty-test/api"
	"github.com/PotatoesFall/vecty-test/backend/application/auth"
	"github.com/PotatoesFall/vecty-test/backend/global/settings"
	"github.com/PotatoesFall/vecty-test/domain/authdomain"
)

func userFromContext(c *router.Context) authdomain.User {
	return c.Get(`user`).(authdomain.User)
}

func authRoutes(r *router.Group, authService auth.Service) {
	r.POST(`/logout`, postLogout(authService))
}

func postLogin(authService auth.Service) func(*router.Context, api.Credentials) error {
	return func(c *router.Context, credentials api.Credentials) error {
		user, ok := authService.Login(credentials.Username, credentials.Password)
		if !ok {
			return c.NoContent(http.StatusUnauthorized)
		}

		http.SetCookie(c.Response, &http.Cookie{ //nolint:exhaustivestruct
			Name:   `auth_token`,
			Value:  user.AuthToken,
			Path:   ``,
			Domain: ``,
			MaxAge: 24 * int(time.Hour/time.Second),
			Secure: !settings.DisableSecure,
		})

		return c.JSON(http.StatusOK, api.LoginResponse{Role: user.Role, Token: user.AuthToken})
	}
}

func postLogout(authService auth.Service) func(c *router.Context) error {
	return func(c *router.Context) error {
		user := userFromContext(c)

		authService.Logout(user.ID)

		return c.NoContent(http.StatusOK)
	}
}

func authMiddleware(authService auth.Service) router.Middleware {
	return func(next router.Handle) router.Handle {
		return func(c *router.Context) error {
			token, err := c.Request.Cookie(`auth_token`)
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

func permissionMiddleware(permission authdomain.Permission) router.Middleware {
	return func(next router.Handle) router.Handle {
		return func(c *router.Context) error {
			user := userFromContext(c)

			if !user.Role.Has(permission) {
				return c.NoContent(http.StatusForbidden)
			}

			return next(c)
		}
	}
}
