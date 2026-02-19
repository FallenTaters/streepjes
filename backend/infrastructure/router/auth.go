package router

import (
	"context"
	"net/http"

	"github.com/FallenTaters/streepjes/api"
	"github.com/FallenTaters/streepjes/backend/application/auth"
	"github.com/FallenTaters/streepjes/domain/authdomain"
	"go.uber.org/zap"
)

func userFromContext(r *http.Request) authdomain.User {
	user, _ := r.Context().Value(userContextKey).(authdomain.User)
	return user
}

func authRoutes(mux *http.ServeMux, authed middleware, authService auth.Service, logger *zap.Logger) {
	handle(mux, "GET /logout", authed, getLogout(authService))
	handle(mux, "POST /active", authed, postActive(authService, logger))

	handle(mux, "GET /profile", authed, getProfilePage(logger))
	handle(mux, "POST /profile/password", authed, postProfilePasswordPage(authService, logger))
	handle(mux, "POST /profile/name", authed, postProfileNamePage(authService, logger))
}

func getLogout(authService auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := userFromContext(r)
		authService.Logout(user.ID)
		http.SetCookie(w, &http.Cookie{
			Name:     api.AuthTokenCookieName,
			Value:    ``,
			Path:     `/`,
			MaxAge:   -1,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})
		http.Redirect(w, r, `/login`, http.StatusSeeOther)
	}
}

func postActive(authService auth.Service, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := userFromContext(r)
		authService.Active(user.ID)
		logger.Debug("received activity refresh", zap.String("user", user.Username))
		w.WriteHeader(http.StatusNoContent)
	}
}

func pageAuthMiddleware(authService auth.Service, logger *zap.Logger) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := r.Cookie(api.AuthTokenCookieName)
			if err != nil {
				logger.Debug("page auth - no token cookie")
				http.Redirect(w, r, `/login`, http.StatusSeeOther)
				return
			}

			user, ok := authService.Check(token.Value)
			if !ok {
				logger.Debug("page auth - token not valid")
				http.Redirect(w, r, `/login`, http.StatusSeeOther)
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func pagePermissionMiddleware(permission authdomain.Permission, logger *zap.Logger) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := userFromContext(r)

			if !user.Role.Has(permission) {
				logger.Debug("page auth - no permission",
					zap.String("role", user.Role.String()),
					zap.Int("permission", int(permission)),
				)
				http.Redirect(w, r, `/`, http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
