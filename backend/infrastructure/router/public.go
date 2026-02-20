package router

import (
	"errors"
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

func (s *Server) publicRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /{$}", s.getRoot)
	mux.HandleFunc("GET /login", s.getLogin)
	mux.HandleFunc("POST /login", s.postLogin)
	mux.HandleFunc("GET /version", getVersion)
	mux.HandleFunc("GET /static/", s.getStatic)
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

func getVersion(w http.ResponseWriter, _ *http.Request) {
	_, _ = io.WriteString(w, version())
}

func (s *Server) getRoot(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie(api.AuthTokenCookieName)
	if err != nil {
		http.Redirect(w, r, `/login`, http.StatusSeeOther)
		return
	}

	user, err := s.auth.Check(token.Value)
	if err != nil {
		if !errors.Is(err, auth.ErrInvalidToken) {
			s.logger.Error("root auth check error", zap.Error(err))
		}
		http.Redirect(w, r, `/login`, http.StatusSeeOther)
		return
	}

	if user.Role.Has(authdomain.PermissionAdminStuff) {
		http.Redirect(w, r, `/admin/billing`, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, `/order`, http.StatusSeeOther)
	}
}

func (s *Server) getLogin(w http.ResponseWriter, r *http.Request) {
	data := struct{ Error bool }{Error: r.URL.Query().Get(`error`) == `1`}
	s.render(w, "login.html", data)
}

func (s *Server) postLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, `/login?error=1`, http.StatusSeeOther)
		return
	}

	username := r.FormValue(`username`)
	password := r.FormValue(`password`)

	user, err := s.auth.Login(username, password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			host, _, _ := net.SplitHostPort(r.RemoteAddr)
			s.logger.Warn("authentication failure", zap.String("ip", host))
		} else {
			s.logger.Error("login error", zap.Error(err))
		}
		http.Redirect(w, r, `/login?error=1`, http.StatusSeeOther)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     api.AuthTokenCookieName,
		Value:    user.AuthToken,
		Path:     `/`,
		MaxAge:   int(authdomain.TokenDuration / time.Second),
		Secure:   s.secureCookies,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	if user.Role.Has(authdomain.PermissionAdminStuff) {
		http.Redirect(w, r, `/admin/billing`, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, `/order`, http.StatusSeeOther)
	}
}

func (s *Server) getStatic(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, `/`), `static/`)

	asset, err := s.static(name)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	contentType := mime.TypeByExtension(filepath.Ext(name))
	if contentType == "" {
		contentType = detectContentType(asset)
	}

	w.Header().Set(`Cache-Control`, `max-age=86400, must-revalidate, private`)
	w.Header().Set(`Content-Type`, contentType)
	_, _ = w.Write(asset)
}

func detectContentType(data []byte) string {
	if strings.HasPrefix(string(data), "/*") || strings.HasPrefix(string(data), "@font-face") {
		return "text/css; charset=utf-8"
	}
	return http.DetectContentType(data)
}
