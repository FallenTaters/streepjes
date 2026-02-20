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

type Server struct {
	auth          auth.Service
	order         order.Service
	static        Static
	secureCookies bool
	logger        *zap.Logger
}

func New(static Static, authService auth.Service, orderService order.Service, secureCookies bool, logger *zap.Logger) http.Handler {
	s := &Server{
		auth:          authService,
		order:         orderService,
		static:        static,
		secureCookies: secureCookies,
		logger:        logger,
	}

	mux := http.NewServeMux()
	s.routes(mux)
	return s.recoveryMiddleware()(mux)
}

func (s *Server) routes(mux *http.ServeMux) {
	s.publicRoutes(mux)

	authed := s.pageAuthMiddleware()
	bar := chain(authed, s.pagePermissionMiddleware(authdomain.PermissionBarStuff))
	admin := chain(authed, s.pagePermissionMiddleware(authdomain.PermissionAdminStuff))

	s.authRoutes(mux, authed)
	s.bartenderRoutes(mux, bar)
	s.adminRoutes(mux, admin)
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

func (s *Server) recoveryMiddleware() middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if v := recover(); v != nil {
					s.logger.Error("panic recovered",
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

func (s *Server) internalError(w http.ResponseWriter, msg string, err error) {
	s.logger.Error(msg, zap.Error(err))
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
