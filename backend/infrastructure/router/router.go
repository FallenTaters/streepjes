package router

import (
	"errors"
	"net/http"
	"os"

	"github.com/FallenTaters/chio"
	"github.com/FallenTaters/chio/middleware"
	"github.com/FallenTaters/streepjes/backend/application/auth"
	"github.com/FallenTaters/streepjes/backend/application/order"
	"github.com/FallenTaters/streepjes/domain/authdomain"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Static func(filename string) ([]byte, error)

func New(static Static, authService auth.Service, orderService order.Service, secureCookies bool, logger *zap.Logger) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recover(middleware.DefaultPanicLogger(os.Stderr)))
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	publicRoutes(r, static, authService, secureCookies, logger)

	authed := r.With(authMiddleware(authService, logger))
	authRoutes(authed, authService, logger)

	bar := authed.With(permissionMiddleware(authdomain.PermissionBarStuff, logger))
	bartenderRoutes(bar, orderService)

	admin := authed.Route(`/admin`, func(r chi.Router) {
		r.Use(permissionMiddleware(authdomain.PermissionAdminStuff, logger))
	})
	adminRoutes(admin, authService, orderService, logger)

	return r
}

func allowErrors(w http.ResponseWriter, logger *zap.Logger, err error, allowed ...error) {
	for _, er := range allowed {
		if errors.Is(err, er) {
			chio.WriteString(w, http.StatusBadRequest, er.Error())
			return
		}
	}

	logger.Error("unexpected error", zap.Error(err))
	chio.Empty(w, http.StatusInternalServerError)
}
