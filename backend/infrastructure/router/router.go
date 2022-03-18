package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/FallenTaters/streepjes/backend/application/auth"
	"github.com/FallenTaters/streepjes/backend/application/order"
	"github.com/FallenTaters/streepjes/domain/authdomain"
)

type Static func(filename string) ([]byte, error)

func New(static Static, authService auth.Service, orderService order.Service) http.Handler {
	r := echo.New()

	r.Use(recoverMiddleware)
	r.HTTPErrorHandler = func(err error, ctx echo.Context) {
		r.DefaultHTTPErrorHandler(err, ctx)
	}

	auth := r.Group(``, authMiddleware(authService))
	authRoutes(auth, authService)

	bar := auth.Group(``, permissionMiddleware(authdomain.PermissionBarStuff))
	bartenderRoutes(bar, orderService)

	admin := auth.Group(``, permissionMiddleware(authdomain.PermissionAdminStuff))
	adminRoutes(admin)

	// must go last because of https://github.com/labstack/echo/issues/2141
	publicRoutes(r, static, authService)

	return r
}

func recoverMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			v := recover()
			if v == nil {
				return
			}

			fmt.Fprintf(os.Stderr, "(%s %s) panic: %s\n", c.Request().Method, c.Request().URL.Path, fmt.Sprint(v))
			fmt.Fprintf(os.Stderr, getStack())
			fmt.Println()

			c.Error(nil)
		}()
		return next(c)
	}
}

func getStack() string {
	pc := make([]uintptr, 50)
	entries := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:entries])

	var stackTrace strings.Builder

	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		if frame.Function == `runtime.gopanic` { // If a panic occurred, start at the frame that called panic
			stackTrace.Reset()
			continue
		}

		stackTrace.WriteString(fmt.Sprintf("%s:%d\n", frame.File, frame.Line))
	}
	return stackTrace.String()
}

func readJSON[T any](c echo.Context) (T, bool) {
	var t T
	err := json.NewDecoder(c.Request().Body).Decode(&t)
	if err != nil {
		c.NoContent(http.StatusBadRequest)
		return t, false
	}

	return t, true
}
