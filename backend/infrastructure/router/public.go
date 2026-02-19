package router

import (
	"fmt"
	"mime"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/FallenTaters/chio"
	"github.com/FallenTaters/streepjes/api"
	"github.com/FallenTaters/streepjes/backend/application/auth"
	"github.com/FallenTaters/streepjes/domain/authdomain"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func publicRoutes(r chi.Router, static Static, authService auth.Service, secureCookies bool, logger *zap.Logger) {
	r.Get(`/`, getRoot(authService))
	r.Get(`/login`, getLogin)
	r.Post(`/login`, postLoginPage(authService, secureCookies, logger))
	r.Get(`/version`, getVersion)
	r.With(chiMiddleware.Compress(5)).Get(`/static/*`, getStatic(static))
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

func getVersion(w http.ResponseWriter, r *http.Request) {
	chio.WriteString(w, http.StatusOK, version())
}

func getRoot(authService auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie(api.AuthTokenCookieName)
		if err != nil {
			http.Redirect(w, r, `/login`, http.StatusSeeOther)
			return
		}

		user, ok := authService.Check(token.Value)
		if !ok {
			http.Redirect(w, r, `/login`, http.StatusSeeOther)
			return
		}

		if user.Role.Has(authdomain.PermissionAdminStuff) {
			http.Redirect(w, r, `/admin/billing`, http.StatusSeeOther)
		} else {
			http.Redirect(w, r, `/order`, http.StatusSeeOther)
		}
	}
}

func getLogin(w http.ResponseWriter, r *http.Request) {
	data := struct{ Error bool }{Error: r.URL.Query().Get(`error`) == `1`}
	render(w, "login.html", data)
}

func postLoginPage(authService auth.Service, secureCookies bool, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Redirect(w, r, `/login?error=1`, http.StatusSeeOther)
			return
		}

		username := r.FormValue(`username`)
		password := r.FormValue(`password`)

		user, ok := authService.Login(username, password)
		if !ok {
			host, _, _ := net.SplitHostPort(r.RemoteAddr)
			logger.Warn("authentication failure", zap.String("ip", host))
			http.Redirect(w, r, `/login?error=1`, http.StatusSeeOther)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:   api.AuthTokenCookieName,
			Value:  user.AuthToken,
			Path:   `/`,
			MaxAge: 24 * int(time.Hour/time.Second),
			Secure: secureCookies,
		})

		if user.Role.Has(authdomain.PermissionAdminStuff) {
			http.Redirect(w, r, `/admin/billing`, http.StatusSeeOther)
		} else {
			http.Redirect(w, r, `/order`, http.StatusSeeOther)
		}
	}
}

func getStatic(assets Static) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, `/`), `static/`)

		asset, err := assets(name)
		if err != nil {
			chio.Empty(w, http.StatusNotFound)
			return
		}

		contentType := mime.TypeByExtension(filepath.Ext(name))
		if contentType == "" {
			contentType = detectContentType(name, asset)
		}

		setCacheHeader(w)
		chio.WriteBlob(w, http.StatusOK, contentType, asset)
	}
}

func detectContentType(name string, data []byte) string {
	if strings.HasPrefix(string(data), "/*") || strings.HasPrefix(string(data), "@font-face") {
		return "text/css; charset=utf-8"
	}
	return http.DetectContentType(data)
}

func setCacheHeader(w http.ResponseWriter) {
	w.Header().Set(`Cache-Control`, `max-age=86400, must-revalidate, private`)
}
