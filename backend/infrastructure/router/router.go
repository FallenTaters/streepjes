package router

import (
	"net/http"
	"os"

	"github.com/FallenTaters/chio/middleware"
	"github.com/FallenTaters/streepjes/backend/application/auth"
	"github.com/FallenTaters/streepjes/backend/application/order"
	"github.com/FallenTaters/streepjes/domain/authdomain"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Static func(filename string) ([]byte, error)

func New(static Static, authService auth.Service, orderService order.Service, secureCookies bool, logger *zap.Logger) http.Handler {
	pageLogger = logger
	r := chi.NewRouter()

	r.Use(middleware.Recover(middleware.DefaultPanicLogger(os.Stderr)))
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	publicRoutes(r, static, authService, secureCookies, logger)

	authed := r.With(pageAuthMiddleware(authService, logger))
	authRoutes(authed, authService, logger)

	bar := authed.With(pagePermissionMiddleware(authdomain.PermissionBarStuff, logger))
	bartenderPageRoutes(bar, orderService, logger)

	admin := authed.Route(`/admin`, func(r chi.Router) {
		r.Use(pagePermissionMiddleware(authdomain.PermissionAdminStuff, logger))
	})
	adminPageRoutes(admin, authService, orderService, logger)

	return r
}

