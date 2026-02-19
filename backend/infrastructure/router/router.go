package router

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/FallenTaters/streepjes/backend/application/auth"
	"github.com/FallenTaters/streepjes/backend/application/order"
	"github.com/FallenTaters/streepjes/domain/authdomain"
	"go.uber.org/zap"
)

type Static func(filename string) ([]byte, error)

type contextKey string

const userContextKey contextKey = "user"

func New(static Static, authService auth.Service, orderService order.Service, secureCookies bool, logger *zap.Logger) http.Handler {
	mux := http.NewServeMux()

	publicRoutes(mux, static, authService, secureCookies, logger)

	authed := pageAuthMiddleware(authService, logger)
	bar := chain(authed, pagePermissionMiddleware(authdomain.PermissionBarStuff, logger))
	admin := chain(authed, pagePermissionMiddleware(authdomain.PermissionAdminStuff, logger))

	authRoutes(mux, authed, authService, logger)
	bartenderPageRoutes(mux, bar, orderService, logger)
	adminPageRoutes(mux, admin, authService, orderService, logger)

	return recoveryMiddleware(logger)(mux)
}

type middleware func(http.Handler) http.Handler

func chain(middlewares ...middleware) middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

func handle(mux *http.ServeMux, pattern string, mw middleware, handler http.HandlerFunc) {
	if mw != nil {
		mux.Handle(pattern, mw(handler))
	} else {
		mux.HandleFunc(pattern, handler)
	}
}

func recoveryMiddleware(logger *zap.Logger) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if v := recover(); v != nil {
					logger.Error("panic recovered",
						zap.Any("value", v),
						zap.String("stack", string(debug.Stack())),
					)
					fmt.Fprintf(os.Stderr, "PANIC: %v\n%s", v, debug.Stack())
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
