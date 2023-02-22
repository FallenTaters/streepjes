package router

import (
	"net/http"
	"time"

	"github.com/FallenTaters/chio"
	"github.com/FallenTaters/chio/middleware"
	"github.com/FallenTaters/streepjes/api"
	"github.com/FallenTaters/streepjes/backend/application/auth"
	"github.com/FallenTaters/streepjes/backend/global/settings"
	"github.com/FallenTaters/streepjes/domain/authdomain"
	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
)

func userFromContext(r *http.Request) authdomain.User {
	return middleware.GetValue[authdomain.User](r, `user`)
}

func authRoutes(r chi.Router, authService auth.Service) {
	r.Post(`/logout`, postLogout(authService))
	r.Post(`/active`, postActive)

	r.Post(`/me/name`, postMeName(authService))
	r.Post(`/me/password`, postMePassword(authService))
}

func postLogin(authService auth.Service) http.HandlerFunc {
	return chio.JSON(func(w http.ResponseWriter, r *http.Request, credentials api.Credentials) {
		log.Debug("received login request", "username", credentials.Username)

		user, ok := authService.Login(credentials.Username, credentials.Password)
		if !ok {
			chio.Empty(w, http.StatusUnauthorized)
			return
		}

		http.SetCookie(w, &http.Cookie{ //nolint:exhaustivestruct
			Name:   api.AuthTokenCookieName,
			Value:  user.AuthToken,
			Path:   ``,
			Domain: ``,
			MaxAge: 24 * int(time.Hour/time.Second),
			Secure: !settings.DisableSecure,
		})

		chio.WriteJSON(w, http.StatusOK, user)
	})
}

func postLogout(authService auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := userFromContext(r)

		log.Debug("received logout request", "user", user.Username)

		authService.Logout(user.ID)

		chio.Empty(w, http.StatusOK)
	}
}

func postActive(w http.ResponseWriter, r *http.Request) {
	user := userFromContext(r)

	log.Debug("received activity refresh", "user", user.Username)

	chio.WriteJSON(w, http.StatusOK, user)
}

func authMiddleware(authService auth.Service) func(http.Handler) http.Handler {
	return middleware.SetValue(`user`, func(w http.ResponseWriter, r *http.Request) (authdomain.User, bool) {
		token, err := r.Cookie(`auth_token`)
		if err != nil {
			log.Debug("auth aborted - no token cookie")
			chio.Empty(w, http.StatusUnauthorized)
			return authdomain.User{}, false
		}

		user, ok := authService.Check(token.Value)
		if !ok {
			log.Debug("auth aborted - token not valid")
			chio.Empty(w, http.StatusUnauthorized)
		}

		return user, ok
	})
}

func permissionMiddleware(permission authdomain.Permission) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := userFromContext(r)

			if !user.Role.Has(permission) {
				log.Debug("auth aborted - no permission", "role", user.Role, "permission", permission)
				chio.Empty(w, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func postMeName(authService auth.Service) http.HandlerFunc {
	return chio.JSON(func(w http.ResponseWriter, r *http.Request, name string) {
		user := userFromContext(r)

		log.Debug("received change name request", "user", user, "name", name)

		if !authService.ChangeName(user, name) {
			chio.Empty(w, http.StatusBadRequest)
			return
		}

		chio.Empty(w, http.StatusOK)
	})
}

func postMePassword(authService auth.Service) http.HandlerFunc {
	return chio.JSON(func(w http.ResponseWriter, r *http.Request, changePassword api.ChangePassword) {
		user := userFromContext(r)

		log.Debug("received password change request", "user", user.Username)

		if !authService.ChangePassword(user, changePassword) {
			chio.Empty(w, http.StatusBadRequest)
			return
		}

		chio.Empty(w, http.StatusOK)
	})
}
