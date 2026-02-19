package router

import (
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/FallenTaters/streepjes/api"
	"github.com/FallenTaters/streepjes/backend/application/auth"
	"github.com/FallenTaters/streepjes/domain/authdomain"
	"go.uber.org/zap"
)

func publicRoutes(mux *http.ServeMux, static Static, authService auth.Service, secureCookies bool, logger *zap.Logger) {
	mux.HandleFunc("GET /{$}", getRoot(authService))
	mux.HandleFunc("GET /login", getLogin(logger))
	mux.HandleFunc("POST /login", postLoginPage(authService, secureCookies, logger))
	mux.HandleFunc("GET /version", getVersion)
	mux.HandleFunc("GET /static/", getStatic(static))
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
	io.WriteString(w, version())
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

func getLogin(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := struct{ Error bool }{Error: r.URL.Query().Get(`error`) == `1`}
		render(w, logger, "login.html", data)
	}
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
			Name:     api.AuthTokenCookieName,
			Value:    user.AuthToken,
			Path:     `/`,
			MaxAge:   int(authdomain.TokenDuration / time.Second),
			Secure:   secureCookies,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
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
			http.NotFound(w, r)
			return
		}

		contentType := mime.TypeByExtension(filepath.Ext(name))
		if contentType == "" {
			contentType = detectContentType(name, asset)
		}

		w.Header().Set(`Cache-Control`, `max-age=86400, must-revalidate, private`)
		w.Header().Set(`Content-Type`, contentType)
		w.Write(asset)
	}
}

func detectContentType(name string, data []byte) string {
	if strings.HasPrefix(string(data), "/*") || strings.HasPrefix(string(data), "@font-face") {
		return "text/css; charset=utf-8"
	}
	return http.DetectContentType(data)
}
