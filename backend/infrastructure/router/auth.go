package router

import (
	"context"
	"errors"
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

func (s *Server) authRoutes(mux *http.ServeMux, authed middleware) {
	handle(mux, "GET /logout", authed, s.getLogout)
	handle(mux, "POST /active", authed, s.postActive)
	handle(mux, "GET /profile", authed, s.getProfilePage)
	handle(mux, "POST /profile/password", authed, s.postProfilePasswordPage)
	handle(mux, "POST /profile/name", authed, s.postProfileNamePage)
}

func (s *Server) getLogout(w http.ResponseWriter, r *http.Request) {
	user := userFromContext(r)
	if err := s.auth.Logout(user.ID); err != nil {
		s.logger.Error("logout failed", zap.Int("user_id", user.ID), zap.Error(err))
	}
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

func (s *Server) postActive(w http.ResponseWriter, r *http.Request) {
	user := userFromContext(r)
	if err := s.auth.Active(user.ID); err != nil {
		s.logger.Error("activity refresh failed", zap.Int("user_id", user.ID), zap.Error(err))
	}
	s.logger.Debug("received activity refresh", zap.String("user", user.Username))
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) pageAuthMiddleware() middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := r.Cookie(api.AuthTokenCookieName)
			if err != nil {
				s.logger.Debug("page auth - no token cookie")
				http.Redirect(w, r, `/login`, http.StatusSeeOther)
				return
			}

			user, err := s.auth.Check(token.Value)
			if err != nil {
				if !errors.Is(err, auth.ErrInvalidToken) {
					s.logger.Error("auth check error", zap.Error(err))
				}
				s.logger.Debug("page auth - token not valid")
				http.Redirect(w, r, `/login`, http.StatusSeeOther)
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (s *Server) pagePermissionMiddleware(permission authdomain.Permission) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := userFromContext(r)

			if !user.Role.Has(permission) {
				s.logger.Debug("page auth - no permission",
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
