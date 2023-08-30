package router

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/FallenTaters/chio"
	"github.com/FallenTaters/chio/middleware"
	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"

	"github.com/FallenTaters/streepjes/backend/application/auth"
	"github.com/FallenTaters/streepjes/backend/application/order"
	"github.com/FallenTaters/streepjes/domain/authdomain"
)

type Static func(filename string) ([]byte, error)

func New(static Static, authService auth.Service, orderService order.Service) http.Handler {
	r := chi.NewRouter()

	r.Use(logMiddleware)
	r.Use(middleware.Recover(middleware.DefaultPanicLogger(os.Stderr)))
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

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

	log.Error("error not allowed for this call", "error", err.Error())
	chio.Empty(w, http.StatusInternalServerError)
}

func logMiddleware(next http.Handler) http.Handler {
	logger := log.New()
	logger.SetLevel(log.DebugLevel)
	logger.SetReportTimestamp(true)
	logger.SetPrefix("HTTP")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &statusRecorder{w, 0}
		next.ServeHTTP(rec, r)
		logger.Debug(r.Method + " " + fmt.Sprint(rec.status) + " " + r.URL.Path)
	})
}

type statusRecorder struct {
	http.ResponseWriter

	status int
}

func (sr *statusRecorder) Write(data []byte) (int, error) {
	if sr.status == 0 {
		sr.status = http.StatusOK
	}

	return sr.ResponseWriter.Write(data)
}

func (sr *statusRecorder) WriteHeader(status int) {
	if sr.status == 0 {
		sr.status = status
	}

	sr.ResponseWriter.WriteHeader(status)
}
