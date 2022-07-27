package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/FallenTaters/chio"
	"github.com/FallenTaters/streepjes/backend/application/auth"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func publicRoutes(r chi.Router, static Static, authService auth.Service) {
	r.Get(`/`, getIndex(static))
	r.Get(`/version`, getVersion)
	r.With(chiMiddleware.Compress(5)).Get(`/static/*`, getStatic(static))
	r.Post(`/login`, postLogin(authService))
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

func getIndex(assets Static) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		index, err := assets(`index.html`)
		if err != nil {
			panic(err)
		}

		setCacheHeader(w)
		chio.WriteBlob(w, http.StatusOK, `text/html`, index)
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

		setCacheHeader(w)
		chio.WriteBlob(w, http.StatusOK, http.DetectContentType(asset), asset)
	}
}

func setCacheHeader(w http.ResponseWriter) {
	w.Header().Set(`Cache-Control`, `max-age=86400, must-revalidate, private`)
}
