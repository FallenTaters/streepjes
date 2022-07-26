package router

import (
	"errors"
	"net/http"
	"os"

	"github.com/FallenTaters/chio"
	"github.com/FallenTaters/chio/middleware"
	"github.com/go-chi/chi/v5"

	"github.com/FallenTaters/streepjes/backend/application/auth"
	"github.com/FallenTaters/streepjes/backend/application/order"
	"github.com/FallenTaters/streepjes/domain/authdomain"
)

type Static func(filename string) ([]byte, error)

func New(static Static, authService auth.Service, orderService order.Service) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recover(middleware.DefaultPanicLogger(os.Stderr)))

	publicRoutes(r, static, authService)

	auth := r.With(authMiddleware(authService))
	authRoutes(auth, authService)

	bar := auth.With(permissionMiddleware(authdomain.PermissionBarStuff))
	bartenderRoutes(bar, orderService)

	admin := auth.Route(`/admin`, func(r chi.Router) {
		r.Use(permissionMiddleware(authdomain.PermissionAdminStuff))
	})
	adminRoutes(admin, authService, orderService)

	return r
}

func allowErrors(w http.ResponseWriter, err error, allowed ...error) {
	for _, er := range allowed {
		if errors.Is(err, er) {
			chio.WriteString(w, http.StatusBadRequest, er.Error())
			return
		}
	}

	chio.Empty(w, http.StatusInternalServerError)
}
