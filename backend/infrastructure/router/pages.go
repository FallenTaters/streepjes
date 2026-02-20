package router

import (
	"net/http"

	"github.com/FallenTaters/streepjes/domain/authdomain"
	"github.com/FallenTaters/streepjes/templates"
	"go.uber.org/zap"
)

type pageData struct {
	ActivePage  string
	User        authdomain.User
	UserDisplay string
	IsBartender bool
	IsAdmin     bool
}

func newPageData(r *http.Request, activePage string) pageData {
	user := userFromContext(r)
	display := user.Username
	if len(display) > 10 {
		display = display[:8] + "…"
	}
	return pageData{
		ActivePage:  activePage,
		User:        user,
		UserDisplay: display,
		IsBartender: user.Role.Has(authdomain.PermissionBarStuff),
		IsAdmin:     user.Role.Has(authdomain.PermissionAdminStuff),
	}
}

func (s *Server) render(w http.ResponseWriter, tmpl string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := templates.Render(w, tmpl, data); err != nil {
		s.logger.Error("template render failed", zap.String("template", tmpl), zap.Error(err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
